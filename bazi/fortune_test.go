package bazi

import (
	"testing"
)

// TestFortuneDirection 陽男陰女順排、陰男陽女逆排。
func TestFortuneDirection(t *testing.T) {
	// 1990 為庚午年，庚屬陽
	cases := []struct {
		gender  Gender
		forward bool
		desc    string
	}{
		{Male, true, "陽年男命順排"},
		{Female, false, "陽年女命逆排"},
	}
	for _, c := range cases {
		b := birthAt(1990, 5, 20, 10, 30)
		b.Gender = c.gender
		got := mustCompute(t, b, Default())
		if len(got.Fortunes) == 0 {
			t.Fatalf("%s: 無大運", c.desc)
		}
		if got.Fortunes[0].Forward != c.forward {
			t.Errorf("%s: 順排=%v，應為 %v", c.desc, got.Fortunes[0].Forward, c.forward)
		}
	}

	// 1989 為己巳年，己屬陰
	for _, c := range []struct {
		gender  Gender
		forward bool
		desc    string
	}{
		{Male, false, "陰年男命逆排"},
		{Female, true, "陰年女命順排"},
	} {
		b := birthAt(1989, 5, 20, 10, 30)
		b.Gender = c.gender
		got := mustCompute(t, b, Default())
		if got.Fortunes[0].Forward != c.forward {
			t.Errorf("%s: 順排=%v，應為 %v", c.desc, got.Fortunes[0].Forward, c.forward)
		}
	}
}

// TestFortuneSequence 大運干支自月柱起，順排遞增、逆排遞減。
//
// 1990-05-20 男命月柱辛巳（索引 17），順排首步為壬午（18）。
func TestFortuneSequence(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	if len(c.Fortunes) < 3 {
		t.Fatalf("大運只有 %d 步，應至少 3 步", len(c.Fortunes))
	}
	month := c.Month.Sexagenary
	for i := 0; i < 3; i++ {
		want := month.Next(i + 1)
		if got := c.Fortunes[i].Sexagenary; got != want {
			t.Errorf("第 %d 步大運為 %d，應為 %d", i+1, got, want)
		}
	}

	// 女命逆排
	b := birthAt(1990, 5, 20, 10, 30)
	b.Gender = Female
	f := mustCompute(t, b, Default())
	for i := 0; i < 3; i++ {
		want := month.Next(-(i + 1))
		if got := f.Fortunes[i].Sexagenary; got != want {
			t.Errorf("女命第 %d 步大運為 %d，應為 %d", i+1, got, want)
		}
	}
}

// TestFortuneStartAge 起運年齡由出生至界節的時距換算，三日折一年。
//
// 1990-05-20 10:30 男命順排，下一個節為芒種；經實測起運約在 5 年 7 個月後，
// 即 1995 年底，故首步大運起於 6 歲。
func TestFortuneStartAge(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	first := c.Fortunes[0]

	if first.StartAge != 6 {
		t.Errorf("首步大運起於 %d 歲，應為 6 歲", first.StartAge)
	}
	if first.EndAge != 15 {
		t.Errorf("首步大運止於 %d 歲，應為 15 歲", first.EndAge)
	}
	if first.StartYear != 1995 {
		t.Errorf("首步大運起於 %d 年，應為 1995 年", first.StartYear)
	}
	// 每步大運橫跨十年
	for i, f := range c.Fortunes {
		if f.EndAge-f.StartAge != 9 {
			t.Errorf("第 %d 步大運跨 %d 年，應跨 10 年", i+1, f.EndAge-f.StartAge+1)
		}
	}
}

// TestFortuneSectsDiffer 四種起運流派的結果須有實質差異，故必須可選。
//
// 差異落在時分層級（實測 tyme4go：19:06、10:30、前一日 10:30、18:30），
// 只比對年份會看不出來，故須以精確起運時刻驗證。
func TestFortuneSectsDiffer(t *testing.T) {
	sects := []ChildLimitSect{
		ChildLimitDefault, ChildLimitChina95,
		ChildLimitLunarSect1, ChildLimitLunarSect2,
	}

	unique := map[string]bool{}
	for _, sect := range sects {
		opt := Default()
		opt.ChildLimit = sect
		c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)

		if c.FortuneStart.IsZero() {
			t.Fatalf("流派 %d 未記錄起運時刻", sect)
		}
		key := c.FortuneStart.Format("2006-01-02 15:04")
		unique[key] = true
		t.Logf("流派 %d 起運於 %s", sect, key)
	}

	if len(unique) < 3 {
		t.Errorf("四種流派只產生 %d 種起運時刻，應至少 3 種——否則流派選項形同虛設",
			len(unique))
	}
}

// TestFortuneStartIsAfterBirth 起運時刻必在出生之後。
func TestFortuneStartIsAfterBirth(t *testing.T) {
	for _, sect := range []ChildLimitSect{
		ChildLimitDefault, ChildLimitChina95,
		ChildLimitLunarSect1, ChildLimitLunarSect2,
	} {
		opt := Default()
		opt.ChildLimit = sect
		b := birthAt(1990, 5, 20, 10, 30)
		c := mustCompute(t, b, opt)
		if !c.FortuneStart.After(b.Time) {
			t.Errorf("流派 %d 的起運時刻 %v 不在出生 %v 之後",
				sect, c.FortuneStart, b.Time)
		}
	}
}

// TestAnnualYearsAreCalendarYears 流年以干支年為準，不可與小運混淆。
//
// 這是 tyme4go 的 Fortune 型別踩過的坑：該型別實為小運，
// 1995 年小運為丁亥而流年為乙亥，誤用即算錯。
func TestAnnualYearsAreCalendarYears(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	first := c.Fortunes[0]

	if len(first.Years) != 10 {
		t.Fatalf("首步大運涵蓋 %d 個流年，應為 10 個", len(first.Years))
	}
	for i, y := range first.Years {
		wantYear := first.StartYear + i
		if y.Year != wantYear {
			t.Errorf("第 %d 個流年為 %d 年，應為 %d 年", i, y.Year, wantYear)
		}
		// 干支年由西元年推算：西元 4 年為甲子年
		wantSex := mod(wantYear-4, 60)
		if int(y.Sexagenary) != wantSex {
			t.Errorf("%d 年干支索引為 %d，應為 %d", y.Year, y.Sexagenary, wantSex)
		}
	}
	// 1995 年應為乙亥（索引 11），而非小運的丁亥
	if first.Years[0].Year == 1995 {
		if want := mod(1995-4, 60); int(first.Years[0].Sexagenary) != want {
			t.Errorf("1995 流年索引為 %d，應為 %d（乙亥）", first.Years[0].Sexagenary, want)
		}
	}
}

// TestFortuneTenGods 大運與流年天干皆對日主取十神。
func TestFortuneTenGods(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	me := c.DayMaster()

	for i, f := range c.Fortunes {
		if want := TenGodOf(me, f.Sexagenary.Stem()); f.StemTenGod != want {
			t.Errorf("第 %d 步大運十神為 %d，應為 %d", i+1, f.StemTenGod, want)
		}
		if len(f.Hidden) == 0 {
			t.Errorf("第 %d 步大運缺藏干", i+1)
		}
		for j, y := range f.Years {
			if want := TenGodOf(me, y.Sexagenary.Stem()); y.StemTenGod != want {
				t.Errorf("第 %d 步大運第 %d 個流年十神為 %d，應為 %d", i+1, j, y.StemTenGod, want)
			}
		}
	}
}

// TestFortuneCountBounded 大運步數須有上限，避免無界成長。
func TestFortuneCountBounded(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	if n := len(c.Fortunes); n < 8 || n > 12 {
		t.Errorf("大運 %d 步，應在 8-12 步之間（涵蓋約百年）", n)
	}
}
