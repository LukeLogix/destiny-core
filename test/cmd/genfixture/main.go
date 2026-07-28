// genfixture 由 NASA JPL Horizons 的太陽視黃經反推節氣精確時刻，產生 golden fixture。
//
// 節氣定義：太陽視黃經達 15° 倍數之時刻（冬至 270°，每節氣 +15°）。
// 月柱只依「節」分界，「氣」不影響四柱，故僅取 12 個節。
//
// 用法：go run ./cmd/genfixture -from 1900 -to 2100 -out fixtures/solarterms_jpl.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/6tail/tyme4go/tyme"
)

// Horizons 的 TLIST 送超過 80 筆會「靜默截斷」——不報錯，只回傳部分資料；
// 送 150 筆則直接伺服器錯誤。故取 50 為安全批量，並在每批強制核對回傳筆數。
const batchSize = 50

// 12 個「節」的索引（tyme4go：0=冬至，奇數為節）
var jieIndices = []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}

var jieKeys = map[int]string{
	1: "xiaohan", 3: "lichun", 5: "jingzhe", 7: "qingming",
	9: "lixia", 11: "mangzhong", 13: "xiaoshu", 15: "liqiu",
	17: "bailu", 19: "hanlu", 21: "lidong", 23: "daxue",
}

// targetLongitude 節氣索引對應的太陽視黃經
func targetLongitude(idx int) float64 {
	return math.Mod(float64(270+idx*15), 360)
}

type Term struct {
	Year         int     `json:"year"`
	Index        int     `json:"index"`
	Key          string  `json:"key"`
	Longitude    float64 `json:"longitude_deg"`
	UTC          string  `json:"utc"`
	JD           float64 `json:"jd"`
	TymeDeltaSec float64 `json:"tyme_delta_sec"` // tyme4go 相對本基準的偏差（正=tyme 較晚）
}

type Fixture struct {
	Source      string  `json:"source"`
	Definition  string  `json:"definition"`
	Ephemeris   string  `json:"ephemeris"`
	TimeScale   string  `json:"time_scale"`
	YearFrom    int     `json:"year_from"`
	YearTo      int     `json:"year_to"`
	Count       int     `json:"count"`
	MaxDeltaSec float64 `json:"tyme_max_delta_sec"`
	Terms       []Term  `json:"terms"`
}

type termProbe struct {
	year, idx  int
	loJD, hiJD float64
	estUTC     time.Time
}

var monMap = map[string]time.Month{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

func utcToJD(t time.Time) float64 {
	return float64(t.UnixNano())/1e9/86400.0 + 2440587.5
}

func jdToUTC(jd float64) time.Time {
	sec := (jd - 2440587.5) * 86400
	whole := math.Floor(sec)
	return time.Unix(int64(whole), int64(math.Round((sec-whole)*1e9))).UTC()
}

// fetchBatch 查詢一批時刻的視黃經。
// jds 必須已排序遞增——Horizons 回傳一律按時間排序，故以順序對應，
// 不以浮點 JD 當 map key（送出值與回算值會有 1e-8 級差異，無法匹配）。
func fetchBatch(jds []float64) ([]float64, error) {
	parts := make([]string, len(jds))
	for i, jd := range jds {
		parts[i] = strconv.FormatFloat(jd, 'f', 8, 64)
	}
	q := url.Values{}
	q.Set("format", "text")
	q.Set("COMMAND", "'10'")
	q.Set("OBJ_DATA", "'NO'")
	q.Set("MAKE_EPHEM", "'YES'")
	q.Set("EPHEM_TYPE", "'OBSERVER'")
	q.Set("CENTER", "'500@399'")
	q.Set("QUANTITIES", "'31'")
	q.Set("TLIST", "'"+strings.Join(parts, ",")+"'")

	resp, err := http.Get("https://ssd.jpl.nasa.gov/api/horizons.api?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	i, j := strings.Index(body, "$$SOE"), strings.Index(body, "$$EOE")
	if i < 0 || j < 0 {
		return nil, fmt.Errorf("回應無星曆區塊（可能超量或參數錯誤，len=%d）", len(body))
	}

	var lons []float64
	for _, line := range strings.Split(body[i+5:j], "\n") {
		f := strings.Fields(line)
		if len(f) < 3 {
			continue
		}
		dp := strings.Split(f[0], "-")
		if len(dp) != 3 {
			continue
		}
		if _, ok := monMap[dp[1]]; !ok {
			continue
		}
		lon, err := strconv.ParseFloat(f[2], 64)
		if err != nil {
			continue
		}
		lons = append(lons, lon)
	}

	// 核對筆數：Horizons 超量時會靜默截斷，不核對就會產出缺漏的 fixture
	if len(lons) != len(jds) {
		return nil, fmt.Errorf("筆數不符：送出 %d，回傳 %d（疑似靜默截斷）", len(jds), len(lons))
	}
	return lons, nil
}

// solarTimeToUTC 將 tyme4go 的 SolarTime（北京時間 UTC+8）轉為 UTC
func solarTimeToUTC(st tyme.SolarTime) time.Time {
	parts := strings.FieldsFunc(st.String(), func(r rune) bool {
		return r == '年' || r == '月' || r == '日' || r == ':' || r == ' '
	})
	var v [6]int
	for i := 0; i < 6 && i < len(parts); i++ {
		v[i], _ = strconv.Atoi(parts[i])
	}
	cst := time.FixedZone("UTC+8", 8*3600)
	return time.Date(v[0], time.Month(v[1]), v[2], v[3], v[4], v[5], 0, cst).UTC()
}

func main() {
	from := flag.Int("from", 1900, "起始年")
	to := flag.Int("to", 2100, "結束年")
	out := flag.String("out", "fixtures/solarterms_jpl.json", "輸出路徑")
	window := flag.Float64("window", 120, "取樣半窗（秒）：需大於 tyme4go 的最大偏差")
	flag.Parse()

	// 1) 以 tyme4go 的估計時刻為中心，每個節氣取前後兩點包夾
	half := time.Duration(*window) * time.Second
	var probes []termProbe
	for y := *from; y <= *to; y++ {
		for _, idx := range jieIndices {
			est := solarTimeToUTC(tyme.SolarTerm{}.FromIndex(y, idx).GetJulianDay().GetSolarTime())
			probes = append(probes, termProbe{
				year: y, idx: idx, estUTC: est,
				loJD: utcToJD(est.Add(-half)),
				hiJD: utcToJD(est.Add(half)),
			})
		}
	}

	// 2) 彙整所有待查時刻並排序（Horizons 回傳按時間排序，需以此對應）
	jdSet := make([]float64, 0, len(probes)*2)
	for _, p := range probes {
		jdSet = append(jdSet, p.loJD, p.hiJD)
	}
	sort.Float64s(jdSet)

	batches := (len(jdSet) + batchSize - 1) / batchSize
	fmt.Fprintf(os.Stderr, "節氣 %d 個，取樣點 %d 個，分 %d 批（每批 %d）\n",
		len(probes), len(jdSet), batches, batchSize)

	lonByJD := make(map[float64]float64, len(jdSet))
	start := time.Now()
	for i := 0; i < len(jdSet); i += batchSize {
		end := i + batchSize
		if end > len(jdSet) {
			end = len(jdSet)
		}
		chunk := jdSet[i:end]

		var lons []float64
		var err error
		for attempt := 0; attempt < 4; attempt++ {
			lons, err = fetchBatch(chunk)
			if err == nil {
				break
			}
			fmt.Fprintf(os.Stderr, "\n  批 %d 第 %d 次重試：%v\n", i/batchSize+1, attempt+1, err)
			time.Sleep(time.Duration(2<<attempt) * time.Second)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n批 %d 連續失敗，中止（已完成 %d 點）\n", i/batchSize+1, len(lonByJD))
			os.Exit(1)
		}
		for k, jd := range chunk {
			lonByJD[jd] = lons[k]
		}

		done := i/batchSize + 1
		elapsed := time.Since(start)
		eta := time.Duration(float64(elapsed) / float64(done) * float64(batches-done))
		fmt.Fprintf(os.Stderr, "\r進度 %d/%d 批  已耗時 %s  預估剩餘 %s    ",
			done, batches, elapsed.Round(time.Second), eta.Round(time.Second))
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Fprintln(os.Stderr)

	// 3) 線性插值求跨越 15° 倍數的精確時刻
	var terms []Term
	var maxDelta float64
	var failed int
	for _, p := range probes {
		a, okA := lonByJD[p.loJD]
		b, okB := lonByJD[p.hiJD]
		if !okA || !okB {
			failed++
			continue
		}
		target := targetLongitude(p.idx)
		if b < a { // 360→0 環繞
			b += 360
		}
		t := target
		if t < a-180 {
			t += 360
		}
		if t < a || t > b {
			fmt.Fprintf(os.Stderr, "警告 %d/idx%d：目標 %.0f° 不在區間 [%.6f, %.6f]，請加大 -window\n",
				p.year, p.idx, target, a, b)
			failed++
			continue
		}
		frac := (t - a) / (b - a)
		jd := p.loJD + frac*(p.hiJD-p.loJD)
		utc := jdToUTC(jd)

		delta := p.estUTC.Sub(utc).Seconds()
		if math.Abs(delta) > math.Abs(maxDelta) {
			maxDelta = delta
		}
		terms = append(terms, Term{
			Year: p.year, Index: p.idx, Key: jieKeys[p.idx],
			Longitude:    target,
			UTC:          utc.Format("2006-01-02T15:04:05.000Z"),
			JD:           jd,
			TymeDeltaSec: math.Round(delta*10) / 10,
		})
	}

	sort.Slice(terms, func(i, j int) bool { return terms[i].JD < terms[j].JD })

	f := Fixture{
		Source:      "NASA JPL Horizons (ssd.jpl.nasa.gov/api/horizons.api)",
		Definition:  "solar term = instant when apparent geocentric solar ecliptic longitude of date reaches a multiple of 15 deg (winter solstice = 270)",
		Ephemeris:   "QUANTITIES=31 (observer ecliptic lon/lat of date), CENTER=500@399 (geocentric)",
		TimeScale:   "UTC",
		YearFrom:    *from,
		YearTo:      *to,
		Count:       len(terms),
		MaxDeltaSec: math.Round(maxDelta*10) / 10,
		Terms:       terms,
	}

	buf, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "序列化失敗：%v\n", err)
		os.Exit(1)
	}
	if dir := filepath.Dir(*out); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "建立目錄失敗：%v\n", err)
		}
	}
	if err := os.WriteFile(*out, buf, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "寫檔失敗：%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("已寫出 %s\n  節氣 %d 筆，失敗 %d 筆\n  tyme4go 相對 JPL 最大偏差 %.1f 秒\n",
		*out, len(terms), failed, maxDelta)
}
