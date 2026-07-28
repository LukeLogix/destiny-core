package calendar

import (
	"errors"
	"math"
	"testing"
)

// TestLookupTermMatchesFixture 內嵌表必須與 JPL golden data 逐筆一致。
// 這是二進位轉換是否無損的唯一保證——轉換過程若有截斷或精度損失，此處會發現。
func TestLookupTermMatchesFixture(t *testing.T) {
	f := loadFixture(t)

	checked := 0
	for _, term := range f.Terms {
		jd, unc, err := LookupTerm(term.Year, term.Index)
		if err != nil {
			t.Fatalf("%d idx=%d 查表失敗: %v", term.Year, term.Index, err)
		}
		// float64 原樣存取，應完全相等；容許 1e-9 天（約 0.1 毫秒）的浮點誤差
		if d := math.Abs(jd - term.JD); d > 1e-9 {
			t.Errorf("%d %s: 查表 JD %.10f 與 fixture %.10f 相差 %.2e",
				term.Year, term.Key, jd, term.JD, d)
		}
		if unc < 0 {
			t.Errorf("%d %s: 不確定度為負 %.2f", term.Year, term.Key, unc)
		}
		checked++
	}
	if checked != f.Count {
		t.Fatalf("只驗證了 %d 筆，fixture 共 %d 筆", checked, f.Count)
	}
	t.Logf("內嵌表與 JPL golden data 逐筆一致，共 %d 筆", checked)
}

// TestLookupTermUncertaintyGrowsWithEra 不確定度須隨年代放大。
// 這是 BoundaryFlags 臨界判定的依據；若恆為定值，臨界標記就退化成固定閾值。
func TestLookupTermUncertaintyGrowsWithEra(t *testing.T) {
	// 立春（idx=3）取三個年代比較
	early, _, err := LookupTerm(2000, 3)
	if err != nil {
		t.Fatal(err)
	}
	_ = early

	var uncs [3]float64
	for i, y := range []int{2000, 2060, 2100} {
		_, u, err := LookupTerm(y, 3)
		if err != nil {
			t.Fatalf("%d: %v", y, err)
		}
		uncs[i] = u
	}
	t.Logf("立春不確定度: 2000=%.1fs  2060=%.1fs  2100=%.1fs", uncs[0], uncs[1], uncs[2])

	if !(uncs[0] < uncs[1] && uncs[1] < uncs[2]) {
		t.Errorf("不確定度未隨年代放大: %.1f, %.1f, %.1f", uncs[0], uncs[1], uncs[2])
	}
	if uncs[0] > 10 {
		t.Errorf("2000 年不確定度 %.1f 秒過大，ΔT 觀測可靠區間不該如此", uncs[0])
	}
	if uncs[2] < 30 {
		t.Errorf("2100 年不確定度 %.1f 秒過小，未反映 ΔT 外推分歧", uncs[2])
	}
}

// TestLookupTermOutOfRange 超出支援範圍須回 ErrYearOutOfRange，不得靜默回零值。
func TestLookupTermOutOfRange(t *testing.T) {
	for _, y := range []int{1899, 2101, 0, -1, 3000} {
		_, _, err := LookupTerm(y, 3)
		if !errors.Is(err, ErrYearOutOfRange) {
			t.Errorf("年份 %d 應回 ErrYearOutOfRange，實得 %v", y, err)
		}
	}
}

// TestLookupTermRejectsQi 表中只存 12 個「節」，「氣」不影響四柱故未收錄。
// 傳入偶數索引須明確報錯，不可回傳鄰近節氣的資料。
func TestLookupTermRejectsQi(t *testing.T) {
	for _, idx := range []int{0, 2, 4, 12, 22} {
		if _, _, err := LookupTerm(2000, idx); !errors.Is(err, ErrNotJie) {
			t.Errorf("索引 %d 為「氣」，應回 ErrNotJie，實得 %v", idx, err)
		}
	}
	for _, idx := range []int{-1, 24, 100} {
		if _, _, err := LookupTerm(2000, idx); err == nil {
			t.Errorf("索引 %d 超出範圍，應回 error", idx)
		}
	}
}

// TestEmbeddedTableAgainstComputed 雙軌交叉驗證：內嵌表 vs VSOP87 自算。
//
// 這是「資料源失真」的防線——表若被竄改或損壞，此測試會發現，
// 且不需要連上任何外部服務。
func TestEmbeddedTableAgainstComputed(t *testing.T) {
	f := loadFixture(t)

	var maxDiff float64
	var worst string
	for _, term := range f.Terms {
		jd, unc, err := LookupTerm(term.Year, term.Index)
		if err != nil {
			t.Fatalf("%d idx=%d: %v", term.Year, term.Index, err)
		}

		w := CumulativeLongitude(term.Longitude, jd)
		computed := SolarTermJD(w)
		diff := math.Abs((computed - jd) * secondsPerDay)

		// 自算與查表的差距，應落在該筆記錄的不確定度附近；
		// 放寬到 3 倍容忍模型差異，但仍能攔下被竄改的資料。
		tolerance := math.Max(unc*3, 10)
		if diff > tolerance {
			t.Errorf("%d %s: 查表與自算相差 %.1f 秒，超出容忍 %.1f 秒",
				term.Year, term.Key, diff, tolerance)
		}
		if diff > maxDiff {
			maxDiff, worst = diff, term.Key
		}
	}
	t.Logf("雙軌最大分歧 %.1f 秒（%s）", maxDiff, worst)
}
