package bazi

import (
	"errors"
	"testing"
	"time"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

var cst = time.FixedZone("UTC+8", 8*3600)

func mustCompute(t *testing.T, b Birth, opt Options) *Chart {
	t.Helper()
	c, err := Compute(b, opt)
	if err != nil {
		t.Fatalf("Compute 失敗: %v", err)
	}
	return c
}

func birthAt(y, mo, d, h, mi int) Birth {
	return Birth{
		Time:   time.Date(y, time.Month(mo), d, h, mi, 0, 0, cst),
		Gender: Male,
	}
}

// TestComputeKnownChart 已知命例：1990-05-20 10:30（UTC+8）為庚午 辛巳 乙酉 辛巳。
// 此盤經 tyme4go 與 lunar-go 雙庫交叉驗證。
func TestComputeKnownChart(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	cases := []struct {
		got    Pillar
		stem   ganzhi.StemIndex
		branch ganzhi.BranchIndex
		desc   string
	}{
		{c.Year, 6, 6, "年柱庚午"},
		{c.Month, 7, 5, "月柱辛巳"},
		{c.Day, 1, 9, "日柱乙酉"},
		{c.Hour, 7, 5, "時柱辛巳"},
	}
	for _, cs := range cases {
		if cs.got.Stem != cs.stem || cs.got.Branch != cs.branch {
			t.Errorf("%s: 得干 %d 支 %d，應為干 %d 支 %d",
				cs.desc, cs.got.Stem, cs.got.Branch, cs.stem, cs.branch)
		}
	}
}

// TestComputeDayMasterTenGods 日主乙木，年月時干的十神。
func TestComputeDayMasterTenGods(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	if c.DayMaster() != 1 {
		t.Fatalf("日主為 %d，應為乙(1)", c.DayMaster())
	}
	if c.Year.StemTenGod != DirectOfficer {
		t.Errorf("年干十神為 %d，應為正官", c.Year.StemTenGod)
	}
	if c.Month.StemTenGod != SevenKilling {
		t.Errorf("月干十神為 %d，應為七殺", c.Month.StemTenGod)
	}
	if c.Hour.StemTenGod != SevenKilling {
		t.Errorf("時干十神為 %d，應為七殺", c.Hour.StemTenGod)
	}
}

// TestComputeHiddenTenGods 月支巳藏丙庚戊，對日主乙木為傷官、正官、正財。
func TestComputeHiddenTenGods(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	want := []TenGod{Hurting, DirectOfficer, DirectWealth}
	if len(c.Month.Hidden) != len(want) {
		t.Fatalf("月支藏干 %d 個，應為 %d 個", len(c.Month.Hidden), len(want))
	}
	for i, w := range want {
		if c.Month.Hidden[i].TenGod != w {
			t.Errorf("月支第 %d 個藏干十神為 %d，應為 %d", i, c.Month.Hidden[i].TenGod, w)
		}
	}
}

// TestYearPillarSwitchesAtLichun 年柱以立春為界，非以元旦或農曆正月初一。
// 2024 立春為 UTC+8 16:27:08.694。
func TestYearPillarSwitchesAtLichun(t *testing.T) {
	before := mustCompute(t, birthAt(2024, 2, 4, 16, 26), Default())
	after := mustCompute(t, birthAt(2024, 2, 4, 16, 28), Default())

	// 立春前為癸卯年（癸=9，卯=3）
	if before.Year.Stem != 9 || before.Year.Branch != 3 {
		t.Errorf("立春前年柱為干 %d 支 %d，應為癸卯", before.Year.Stem, before.Year.Branch)
	}
	// 立春後為甲辰年（甲=0，辰=4）
	if after.Year.Stem != 0 || after.Year.Branch != 4 {
		t.Errorf("立春後年柱為干 %d 支 %d，應為甲辰", after.Year.Stem, after.Year.Branch)
	}
	// 月柱同時由乙丑轉丙寅
	if before.Month.Branch != 1 {
		t.Errorf("立春前月支為 %d，應為丑(1)", before.Month.Branch)
	}
	if after.Month.Branch != 2 {
		t.Errorf("立春後月支為 %d，應為寅(2)", after.Month.Branch)
	}
}

// TestYearPillarBeforeLichunBelongsToPreviousYear 元旦出生仍屬前一干支年。
func TestYearPillarBeforeLichunBelongsToPreviousYear(t *testing.T) {
	c := mustCompute(t, birthAt(2024, 1, 1, 12, 0), Default())
	// 2024-01-01 尚未立春，屬癸卯年
	if c.Year.Stem != 9 || c.Year.Branch != 3 {
		t.Errorf("2024 元旦年柱為干 %d 支 %d，應為癸卯", c.Year.Stem, c.Year.Branch)
	}
}

// TestLateZiKeepsDay 子時換日口徑：23 時出生，日柱算今日或明日。
// 兩庫預設相反，故必須顯式配置。
func TestLateZiKeepsDay(t *testing.T) {
	opt := Default()

	// 預設（早子時派）：23 時進次日
	opt.LateZiKeepsDay = false
	early := mustCompute(t, birthAt(1990, 5, 20, 23, 30), opt)

	// 晚子時派：23 時仍算當日
	opt.LateZiKeepsDay = true
	late := mustCompute(t, birthAt(1990, 5, 20, 23, 30), opt)

	if early.Day.Sexagenary == late.Day.Sexagenary {
		t.Fatal("兩種子時口徑的日柱應不同")
	}
	// 晚子時派日柱應與當日 10:30 相同（乙酉）
	sameDay := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)
	if late.Day.Sexagenary != sameDay.Day.Sexagenary {
		t.Errorf("晚子時派 23:30 的日柱應與當日相同")
	}
	// 早子時派日柱應為次日
	if early.Day.Sexagenary != sameDay.Day.Sexagenary.Next(1) {
		t.Errorf("早子時派 23:30 的日柱應為次日")
	}
	// 時支皆為子
	if early.Hour.Branch != 0 || late.Hour.Branch != 0 {
		t.Errorf("23:30 時支應為子，得 %d 與 %d", early.Hour.Branch, late.Hour.Branch)
	}
}

// TestComputeCarriesOptions 命盤須攜帶所用口徑，確保結果可完整重現。
func TestComputeCarriesOptions(t *testing.T) {
	opt := Default()
	opt.LateZiKeepsDay = true
	opt.HiddenStem = HiddenStemWithEarth

	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)
	if !c.Options.LateZiKeepsDay {
		t.Error("命盤未記錄 LateZiKeepsDay")
	}
	if c.Options.HiddenStem != HiddenStemWithEarth {
		t.Error("命盤未記錄 HiddenStem 口徑")
	}
}

// TestComputeRejectsYearOutOfRange 支援範圍外須回 ErrYearOutOfRange，
// 不得靜默回傳可能錯誤的結果。
func TestComputeRejectsYearOutOfRange(t *testing.T) {
	for _, y := range []int{1899, 2101} {
		if _, err := Compute(birthAt(y, 6, 1, 12, 0), Default()); !errors.Is(err, ErrYearOutOfRange) {
			t.Errorf("年份 %d 應回 ErrYearOutOfRange，實得 %v", y, err)
		}
	}
	// 邊界年份應可計算
	for _, y := range []int{1900, 2100} {
		if _, err := Compute(birthAt(y, 6, 1, 12, 0), Default()); err != nil {
			t.Errorf("年份 %d 應可計算，實得 %v", y, err)
		}
	}
}

// TestComputeRejectsMissingGender 大運順逆需要性別。
func TestComputeRejectsMissingGender(t *testing.T) {
	b := birthAt(1990, 5, 20, 10, 30)
	b.Gender = GenderUnset
	if _, err := Compute(b, Default()); !errors.Is(err, ErrGenderRequired) {
		t.Errorf("性別未設應回 ErrGenderRequired，實得 %v", err)
	}
}

// TestComputeRequiresLongitudeForSolarCorrection 真太陽時修正需要經度；
// 零值 0 是合法經度（格林威治），不可用來表達「未提供」。
func TestComputeRequiresLongitudeForSolarCorrection(t *testing.T) {
	b := birthAt(1990, 5, 20, 10, 30)

	opt := Default()
	opt.SolarTime = LongitudeOnly
	if _, err := Compute(b, opt); !errors.Is(err, ErrLongitudeMissing) {
		t.Errorf("未提供經度應回 ErrLongitudeMissing，實得 %v", err)
	}

	// 明確提供 0 度應可計算
	zero := 0.0
	b.Longitude = &zero
	opt.SolarTime = WallClock // 0 度與 UTC+8 差距過大，故改用鐘面時間驗證欄位本身
	if _, err := Compute(b, opt); err != nil {
		t.Errorf("經度 0 度為合法值，應可計算，實得 %v", err)
	}
}

// TestComputeRejectsInvalidLongitude 經度超出 ±180 須報錯。
func TestComputeRejectsInvalidLongitude(t *testing.T) {
	b := birthAt(1990, 5, 20, 10, 30)
	opt := Default()
	opt.SolarTime = LongitudeOnly

	for _, lon := range []float64{181, -181, 360} {
		l := lon
		b.Longitude = &l
		if _, err := Compute(b, opt); !errors.Is(err, ErrLongitudeInvalid) {
			t.Errorf("經度 %.0f 應回 ErrLongitudeInvalid，實得 %v", lon, err)
		}
	}
}

// TestComputeRejectsLongitudeTimezoneMismatch 經度與時區差距過大表示填錯。
//
// 門檻取 4 小時：中國全境單一時區（新疆達 2.9 小時）、西班牙夏令（2.6 小時）
// 等全球合法情況皆須通過，而「台北時間配紐約經度」（12.9 小時）須攔下。
func TestComputeRejectsLongitudeTimezoneMismatch(t *testing.T) {
	opt := Default()
	opt.SolarTime = LongitudeOnly

	// 台北時間配紐約經度：明顯填錯
	b := birthAt(1990, 5, 20, 10, 30)
	ny := -74.0
	b.Longitude = &ny
	if _, err := Compute(b, opt); !errors.Is(err, ErrLongitudeMismatch) {
		t.Errorf("時區與經度嚴重不符應回 ErrLongitudeMismatch，實得 %v", err)
	}

	// 新疆喀什 76°E 配 UTC+8：合法，差距 2.9 小時
	kashgar := 76.0
	b.Longitude = &kashgar
	if _, err := Compute(b, opt); err != nil {
		t.Errorf("新疆喀什為全球合法極端情況，應可計算，實得 %v", err)
	}
}

// TestSolarTimeModes 三段修正的效果應遞增，且修正量記錄於命盤供稽核。
func TestSolarTimeModes(t *testing.T) {
	b := birthAt(1990, 5, 20, 10, 30)
	taipei := 121.5
	b.Longitude = &taipei

	opt := Default()

	opt.SolarTime = WallClock
	wall := mustCompute(t, b, opt)
	if wall.SolarTimeOffset != 0 {
		t.Errorf("鐘面時間模式的修正量應為 0，實得 %v", wall.SolarTimeOffset)
	}

	opt.SolarTime = LongitudeOnly
	lon := mustCompute(t, b, opt)
	// 台北 121.5°E 對 UTC+8 基準線 120°E，約 +6 分鐘
	wantMin := 6.0
	gotMin := lon.SolarTimeOffset.Minutes()
	if gotMin < wantMin-1 || gotMin > wantMin+1 {
		t.Errorf("經度修正應約 +%.0f 分鐘，實得 %.1f 分鐘", wantMin, gotMin)
	}

	opt.SolarTime = TrueSolar
	true_ := mustCompute(t, b, opt)
	// 均時差振幅 ±16 分，故與純經度修正應有差異
	if true_.SolarTimeOffset == lon.SolarTimeOffset {
		t.Error("真太陽時應在經度修正之上再加均時差")
	}
	// 修正後的時刻須與修正量一致
	if got := true_.EffectiveTime.Sub(b.Time); got != true_.SolarTimeOffset {
		t.Errorf("EffectiveTime 與 SolarTimeOffset 不一致：%v vs %v",
			got, true_.SolarTimeOffset)
	}
}
