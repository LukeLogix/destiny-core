package test

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"os"
	"sort"
	"testing"
)

// bin 與 JSON 的一致性驗證放在這裡，而非 internal/calendar 內。
//
// 理由：golden JSON 位於本巢狀 module 中，Go 會把整個 test/ 目錄排除在
// 主 module 的發布 zip 之外，下游對模組快取跑測試時讀不到它。
// 主 module 的測試因此改以內嵌表自足，轉換是否無損則由本檔把關——
// 這是 repo 內才需要跑的驗證，正好與 cmd/ 工具同層。

const (
	jsonPath = "fixtures/solarterms_jpl.json"
	binPath  = "../internal/calendar/solarterms.bin"

	binHeaderSize = 13
	binRecordSize = 12
	jiePerYear    = 12
)

type goldenTerm struct {
	Year         int     `json:"year"`
	Index        int     `json:"index"`
	Key          string  `json:"key"`
	JD           float64 `json:"jd"`
	TymeDeltaSec float64 `json:"tyme_delta_sec"`
}

type goldenFile struct {
	YearFrom int          `json:"year_from"`
	YearTo   int          `json:"year_to"`
	Count    int          `json:"count"`
	Terms    []goldenTerm `json:"terms"`
}

func loadGolden(t *testing.T) goldenFile {
	t.Helper()
	b, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("讀取 golden JSON 失敗: %v", err)
	}
	var f goldenFile
	if err := json.Unmarshal(b, &f); err != nil {
		t.Fatalf("解析 golden JSON 失敗: %v", err)
	}
	if len(f.Terms) != f.Count {
		t.Fatalf("JSON 自述 %d 筆，實際 %d 筆", f.Count, len(f.Terms))
	}
	return f
}

func loadBin(t *testing.T) (yearFrom, yearTo int, records []byte) {
	t.Helper()
	b, err := os.ReadFile(binPath)
	if err != nil {
		t.Fatalf("讀取內嵌表失敗: %v", err)
	}
	if len(b) < binHeaderSize {
		t.Fatal("內嵌表過短")
	}
	if string(b[0:4]) != "DCST" {
		t.Fatalf("magic 為 %q，應為 DCST", b[0:4])
	}
	if b[4] != 1 {
		t.Fatalf("版本為 %d，應為 1", b[4])
	}
	yearFrom = int(int16(binary.LittleEndian.Uint16(b[5:7])))
	yearTo = int(int16(binary.LittleEndian.Uint16(b[7:9])))
	count := int(binary.LittleEndian.Uint32(b[9:13]))

	if want := binHeaderSize + count*binRecordSize; len(b) != want {
		t.Fatalf("檔案長度 %d bytes，依 %d 筆應為 %d", len(b), count, want)
	}
	return yearFrom, yearTo, b[binHeaderSize:]
}

// TestBinMatchesGoldenJSON 內嵌表須與 JPL golden data 逐筆一致。
//
// 這是二進位轉換無損的唯一保證——若 genbin 有截斷、位移或精度損失，此處會發現。
func TestBinMatchesGoldenJSON(t *testing.T) {
	g := loadGolden(t)
	yearFrom, yearTo, rec := loadBin(t)

	if yearFrom != g.YearFrom || yearTo != g.YearTo {
		t.Fatalf("年份範圍不符：bin %d-%d，JSON %d-%d", yearFrom, yearTo, g.YearFrom, g.YearTo)
	}

	// JSON 依時間排序，bin 依 (年份, 節氣索引) 排序，故先建索引
	type key struct{ y, i int }
	byKey := make(map[key]goldenTerm, len(g.Terms))
	for _, term := range g.Terms {
		byKey[key{term.Year, term.Index}] = term
	}

	order := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}
	sort.Ints(order)

	checked := 0
	for y := yearFrom; y <= yearTo; y++ {
		for n, idx := range order {
			term, ok := byKey[key{y, idx}]
			if !ok {
				t.Fatalf("JSON 缺 %d 年索引 %d", y, idx)
			}
			off := ((y-yearFrom)*jiePerYear + n) * binRecordSize
			jd := math.Float64frombits(binary.LittleEndian.Uint64(rec[off : off+8]))
			unc := float64(math.Float32frombits(
				binary.LittleEndian.Uint32(rec[off+8 : off+12])))

			if d := math.Abs(jd - term.JD); d > 1e-9 {
				t.Errorf("%d %s: bin 的 JD %.10f 與 JSON %.10f 相差 %.2e",
					y, term.Key, jd, term.JD, d)
			}
			// 不確定度以 float32 儲存且取絕對值，容差放寬到 0.01 秒
			if d := math.Abs(unc - math.Abs(term.TymeDeltaSec)); d > 0.01 {
				t.Errorf("%d %s: bin 的不確定度 %.3f 與 JSON %.3f 相差 %.3f",
					y, term.Key, unc, math.Abs(term.TymeDeltaSec), d)
			}
			checked++
		}
	}
	if checked != g.Count {
		t.Fatalf("只驗證 %d 筆，JSON 共 %d 筆", checked, g.Count)
	}
	t.Logf("內嵌表與 JPL golden JSON 逐筆一致，共 %d 筆", checked)
}
