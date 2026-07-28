package calendar

import (
	"math"
	"sort"
	"testing"
)

// fixtureTerm 一筆基準節氣。
//
// 基準取自 go:embed 的內嵌表，而非 test/ module 內的 JSON——
// test/ 有自己的 go.mod，Go 會將整個目錄排除在發布的 module zip 之外，
// 下游使用者對模組快取執行 go test 時根本讀不到那個檔案。
// 內嵌表本身即 JPL 值，且 bin 與 JSON 的逐筆一致由 test/fixture_test.go 把關。
type fixtureTerm struct {
	Year      int
	Index     int
	Key       string
	Longitude float64
	JD        float64
}

type fixtureFile struct {
	YearFrom int
	YearTo   int
	Count    int
	Terms    []fixtureTerm
}

// jieKeys 供錯誤訊息辨識用
var jieKeys = map[int]string{
	1: "xiaohan", 3: "lichun", 5: "jingzhe", 7: "qingming",
	9: "lixia", 11: "mangzhong", 13: "xiaoshu", 15: "liqiu",
	17: "bailu", 19: "hanlu", 21: "lidong", 23: "daxue",
}

func loadFixture(t *testing.T) fixtureFile {
	t.Helper()
	from, to := YearRange()
	f := fixtureFile{YearFrom: from, YearTo: to}

	for y := from; y <= to; y++ {
		for idx := 1; idx <= 23; idx += 2 {
			jd, _, err := LookupTerm(y, idx)
			if err != nil {
				t.Fatalf("內嵌表缺 %d 年索引 %d: %v", y, idx, err)
			}
			f.Terms = append(f.Terms, fixtureTerm{
				Year:      y,
				Index:     idx,
				Key:       jieKeys[idx],
				Longitude: math.Mod(float64(270+idx*15), 360),
				JD:        jd,
			})
		}
	}
	f.Count = len(f.Terms)
	if f.Count == 0 {
		t.Fatal("內嵌表無資料")
	}
	return f
}

// TestSolarTermAgainstJPL 以 JPL golden data 驗證 VSOP87 移植的正確性。
//
// 這是雙軌設計的核心保障：內嵌表與自行計算互為守護。若此測試失敗，
// 表示兩條路徑之一出了問題，而非單純的精度不足。
func TestSolarTermAgainstJPL(t *testing.T) {
	f := loadFixture(t)

	var errs []float64
	worst := struct {
		term fixtureTerm
		sec  float64
	}{}

	for _, term := range f.Terms {
		w := CumulativeLongitude(term.Longitude, term.JD)
		gotJD := SolarTermJD(w)
		diffSec := (gotJD - term.JD) * secondsPerDay

		errs = append(errs, math.Abs(diffSec))
		if math.Abs(diffSec) > math.Abs(worst.sec) {
			worst.term, worst.sec = term, diffSec
		}
	}

	sort.Float64s(errs)
	var sum float64
	for _, e := range errs {
		sum += e
	}
	mean := sum / float64(len(errs))
	median := errs[len(errs)/2]
	p99 := errs[len(errs)*99/100]
	max := errs[len(errs)-1]

	t.Logf("VSOP87 計算 vs JPL golden data（%d 筆，%d-%d）",
		len(errs), f.YearFrom, f.YearTo)
	t.Logf("  平均 %.2f 秒  中位 %.2f 秒  p99 %.2f 秒  最大 %.2f 秒",
		mean, median, p99, max)
	t.Logf("  最差: %d %s  偏差 %+.2f 秒", worst.term.Year, worst.term.Key, worst.sec)

	// 門檻依 2.5 節的 ΔT 分析設定：偏差主要來自 ΔT 模型差異而非黃經計算，
	// 2040 年後本就有分鐘級不可知性，故以「中位數」把關計算本身的正確性。
	if median > 5 {
		t.Errorf("中位偏差 %.2f 秒超過 5 秒，移植可能有誤", median)
	}
	if max > 300 {
		t.Errorf("最大偏差 %.2f 秒超過 300 秒，超出 ΔT 分歧可解釋的範圍", max)
	}
}

// TestSolarTermByEra 分年代檢視偏差，用以區分「移植錯誤」與「ΔT 模型分歧」。
//
// 若移植有誤，各年代會一致地偏；若只是 ΔT 分歧，偏差會隨年代單調放大。
func TestSolarTermByEra(t *testing.T) {
	f := loadFixture(t)

	type stat struct {
		n        int
		sum, max float64
	}
	eras := map[int]*stat{}

	for _, term := range f.Terms {
		w := CumulativeLongitude(term.Longitude, term.JD)
		diffSec := math.Abs((SolarTermJD(w) - term.JD) * secondsPerDay)

		era := term.Year / 20 * 20
		if eras[era] == nil {
			eras[era] = &stat{}
		}
		s := eras[era]
		s.n++
		s.sum += diffSec
		if diffSec > s.max {
			s.max = diffSec
		}
	}

	keys := make([]int, 0, len(eras))
	for k := range eras {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	t.Log("年代        筆數   平均偏差   最大偏差")
	for _, k := range keys {
		s := eras[k]
		t.Logf("%d-%d  %4d   %7.2fs   %7.2fs", k, k+19, s.n, s.sum/float64(s.n), s.max)
	}

	// 1940-2039 是 ΔT 觀測值可靠的區間，此處偏差應極小；
	// 若這段也偏得厲害，就不是 ΔT 的問題，而是移植錯了。
	for _, k := range []int{1940, 1960, 1980, 2000, 2020} {
		if s := eras[k]; s != nil {
			if avg := s.sum / float64(s.n); avg > 3 {
				t.Errorf("%d 年代平均偏差 %.2f 秒過大（ΔT 可靠區間），移植可能有誤", k, avg)
			}
		}
	}
}

// TestCumulativeLongitude 驗證圈數定位正確——這是最易出錯處：
// 同一個 15° 倍數每年重現，圈數算錯會定位到相鄰年份的同名節氣。
func TestCumulativeLongitude(t *testing.T) {
	f := loadFixture(t)

	for _, term := range f.Terms {
		w := CumulativeLongitude(term.Longitude, term.JD)
		// 還原為 0-360 度，應與目標黃經相符
		deg := math.Mod(math.Mod(w*180/math.Pi, 360)+360, 360)
		want := math.Mod(math.Mod(term.Longitude, 360)+360, 360)
		if d := math.Abs(deg - want); d > 1e-6 && math.Abs(d-360) > 1e-6 {
			t.Fatalf("%d %s: 累積黃經還原為 %.6f°，應為 %.6f°",
				term.Year, term.Key, deg, want)
		}
		// 求解結果須落在原時刻附近，否則是定位到別年了
		if gotJD := SolarTermJD(w); math.Abs(gotJD-term.JD) > 1 {
			t.Fatalf("%d %s: 求得 JD %.5f 與目標 %.5f 相差超過一天，圈數定位錯誤",
				term.Year, term.Key, gotJD, term.JD)
		}
	}
}
