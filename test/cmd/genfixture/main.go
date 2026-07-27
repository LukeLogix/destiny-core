// 用 NASA JPL Horizons 的太陽視黃經，反推節氣精確時刻，驗證 tyme4go。
// 這就是未來 golden fixture 的產生器。
package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/6tail/tyme4go/tyme"
)

var monMap = map[string]time.Month{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6,
	"Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

type sample struct {
	t   time.Time
	lon float64
}

func fetchLon(start, stop time.Time) ([]sample, error) {
	q := url.Values{}
	q.Set("format", "text")
	q.Set("COMMAND", "'10'")
	q.Set("OBJ_DATA", "'NO'")
	q.Set("MAKE_EPHEM", "'YES'")
	q.Set("EPHEM_TYPE", "'OBSERVER'")
	q.Set("CENTER", "'500@399'")
	q.Set("START_TIME", "'"+start.Format("2006-01-02 15:04")+"'")
	q.Set("STOP_TIME", "'"+stop.Format("2006-01-02 15:04")+"'")
	q.Set("STEP_SIZE", "'1m'")
	q.Set("QUANTITIES", "'31'")

	resp, err := http.Get("https://ssd.jpl.nasa.gov/api/horizons.api?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	body := string(b)

	i, j := strings.Index(body, "$$SOE"), strings.Index(body, "$$EOE")
	if i < 0 || j < 0 {
		return nil, fmt.Errorf("unexpected response")
	}
	var out []sample
	for _, line := range strings.Split(body[i+5:j], "\n") {
		f := strings.Fields(line)
		if len(f) < 3 {
			continue
		}
		dp := strings.Split(f[0], "-")
		if len(dp) != 3 {
			continue
		}
		yy, _ := strconv.Atoi(dp[0])
		dd, _ := strconv.Atoi(dp[2])
		mm, ok := monMap[dp[1]]
		if !ok {
			continue
		}
		tp := strings.Split(f[1], ":")
		hh, _ := strconv.Atoi(tp[0])
		mi, _ := strconv.Atoi(tp[1])
		lon, err := strconv.ParseFloat(f[2], 64)
		if err != nil {
			continue
		}
		out = append(out, sample{
			time.Date(yy, mm, dd, hh, mi, 0, 0, time.UTC), lon,
		})
	}
	return out, nil
}

// solve 在樣本中找視黃經跨越 target 的時刻（線性插值）
func solve(ss []sample, target float64) (time.Time, bool) {
	for k := 1; k < len(ss); k++ {
		a, b := ss[k-1].lon, ss[k].lon
		// 處理 360→0 環繞
		if b < a {
			b += 360
		}
		t := target
		if t < a-180 {
			t += 360
		}
		if a <= t && t <= b {
			frac := (t - a) / (b - a)
			d := ss[k].t.Sub(ss[k-1].t)
			return ss[k-1].t.Add(time.Duration(frac * float64(d))), true
		}
	}
	return time.Time{}, false
}

func main() {
	cst := time.FixedZone("UTC+8", 8*3600)

	cases := []struct {
		year, idx int
		name      string
	}{
		{1950, 6, "春分"},
		{1990, 9, "立夏"},
		{2024, 3, "立春"},
		{2050, 0, "冬至"},
	}

	fmt.Printf("%-6s %-6s  %-21s %-21s %s\n",
		"年份", "節氣", "tyme4go (UTC+8)", "JPL Horizons (UTC+8)", "誤差")
	fmt.Println(strings.Repeat("-", 78))

	var maxErr float64
	for _, c := range cases {
		// tyme4go 的答案（鐘面時間為 UTC+8）
		st := tyme.SolarTerm{}.FromIndex(c.year, c.idx)
		sd := st.GetJulianDay().GetSolarTime()
		parts := strings.FieldsFunc(sd.String(), func(r rune) bool {
			return r == '年' || r == '月' || r == '日' || r == ':' || r == ' '
		})
		var v [6]int
		for i := 0; i < 6 && i < len(parts); i++ {
			v[i], _ = strconv.Atoi(parts[i])
		}
		tymeT := time.Date(v[0], time.Month(v[1]), v[2], v[3], v[4], v[5], 0, cst)

		// 目標視黃經：冬至=270°，每 index +15°
		target := math.Mod(float64(270+c.idx*15), 360)

		center := tymeT.UTC()
		ss, err := fetchLon(center.Add(-4*time.Minute), center.Add(4*time.Minute))
		if err != nil {
			fmt.Printf("%-6d %-6s  查詢失敗: %v\n", c.year, c.name, err)
			continue
		}
		jplT, ok := solve(ss, target)
		if !ok {
			fmt.Printf("%-6d %-6s  區間內未跨越 %.0f°（樣本 %d 筆）\n", c.year, c.name, target, len(ss))
			continue
		}

		diff := tymeT.Sub(jplT).Seconds()
		if math.Abs(diff) > maxErr {
			maxErr = math.Abs(diff)
		}
		fmt.Printf("%-6d %-6s  %-21s %-21s %+.1f 秒\n",
			c.year, c.name,
			tymeT.Format("2006-01-02 15:04:05"),
			jplT.In(cst).Format("2006-01-02 15:04:05"),
			diff)

		time.Sleep(1200 * time.Millisecond) // 對 API 客氣一點
	}
	fmt.Printf("\n最大誤差: %.1f 秒\n", maxErr)
}
