package calendar

import (
	"math"
	"testing"
	"time"
)

// TestGregorianJDN 民用日期轉儒略日數的定位點。
func TestGregorianJDN(t *testing.T) {
	cases := []struct {
		y, m, d int
		want    int
		desc    string
	}{
		{2000, 1, 1, 2451545, "J2000 曆元"},
		{1900, 1, 1, 2415021, "1900 元旦"},
		{2100, 12, 31, 2488434, "2100 除夕（由 J2000 加日數獨立推算核對）"},
		{1582, 10, 15, 2299161, "格里曆首日"},
	}
	for _, c := range cases {
		if got := GregorianJDN(c.y, c.m, c.d); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestGregorianJDNConsecutive 相鄰日期的 JDN 必然連續，跨月跨年跨閏日皆然。
func TestGregorianJDNConsecutive(t *testing.T) {
	start := time.Date(1999, 12, 28, 12, 0, 0, 0, time.UTC)
	prev := GregorianJDN(start.Year(), int(start.Month()), start.Day())
	for i := 1; i < 800; i++ {
		d := start.AddDate(0, 0, i)
		got := GregorianJDN(d.Year(), int(d.Month()), d.Day())
		if got != prev+1 {
			t.Fatalf("%s 的 JDN 為 %d，前一日為 %d，應連續", d.Format("2006-01-02"), got, prev)
		}
		prev = got
	}
}

// TestCivilJD 帶時刻的儒略日：J2000 曆元為 2000-01-01 12:00 UTC。
func TestCivilJD(t *testing.T) {
	epoch := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	if got := CivilJD(epoch); math.Abs(got-J2000) > 1e-9 {
		t.Errorf("J2000 曆元的儒略日為 %.9f，應為 %.1f", got, J2000)
	}

	// 半日之差
	noon := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
	if got := CivilJD(noon); math.Abs(got-(J2000+0.5)) > 1e-9 {
		t.Errorf("2000-01-02 00:00 UTC 的儒略日為 %.9f，應為 %.1f", got, J2000+0.5)
	}
}

// TestCivilJDRespectsTimezone 同一絕對時刻，不同時區表示應得相同儒略日。
func TestCivilJDRespectsTimezone(t *testing.T) {
	utc := time.Date(2024, 6, 1, 4, 0, 0, 0, time.UTC)
	cst := time.Date(2024, 6, 1, 12, 0, 0, 0, time.FixedZone("UTC+8", 8*3600))

	if math.Abs(CivilJD(utc)-CivilJD(cst)) > 1e-9 {
		t.Errorf("同一時刻的兩種時區表示應得相同儒略日：%.9f vs %.9f",
			CivilJD(utc), CivilJD(cst))
	}
}

// TestEquationOfTime 均時差的年度極值。
//
// 一年有四個駐點：約 2/11 最小（約 −14.2 分）、5/14 次大（約 +3.7 分）、
// 7/26 次小（約 −6.5 分）、11/3 最大（約 +16.4 分）。
func TestEquationOfTime(t *testing.T) {
	cases := []struct {
		month, day int
		wantMin    float64
		desc       string
	}{
		{2, 11, -14.2, "二月中旬最小"},
		{5, 14, 3.7, "五月中旬次大"},
		{7, 26, -6.5, "七月下旬次小"},
		{11, 3, 16.4, "十一月初最大"},
	}
	for _, c := range cases {
		jd := CivilJD(time.Date(2024, time.Month(c.month), c.day, 12, 0, 0, 0, time.UTC))
		got := EquationOfTime(jd).Minutes()
		if math.Abs(got-c.wantMin) > 1.0 {
			t.Errorf("%s: 均時差 %.2f 分，應約 %.1f 分（容差 1 分）",
				c.desc, got, c.wantMin)
		}
	}
}

// TestEquationOfTimeAmplitude 均時差全年振幅約 −14 至 +17 分鐘。
// 這個振幅大於多數地區的經度修正，故真太陽時不可只做經度而略過均時差。
func TestEquationOfTimeAmplitude(t *testing.T) {
	var min, max float64
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	for i := 0; i < 366; i++ {
		m := EquationOfTime(CivilJD(base.AddDate(0, 0, i))).Minutes()
		if i == 0 || m < min {
			min = m
		}
		if i == 0 || m > max {
			max = m
		}
	}
	t.Logf("均時差全年範圍：%.2f ~ %.2f 分鐘", min, max)
	if min > -13 || min < -16 {
		t.Errorf("年度最小值 %.2f 分，應約 −14 分", min)
	}
	if max < 15 || max > 18 {
		t.Errorf("年度最大值 %.2f 分，應約 +16 分", max)
	}
}

// TestEquationOfTimeZeroCrossings 均時差一年四度歸零。
func TestEquationOfTimeZeroCrossings(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	crossings := 0
	prev := EquationOfTime(CivilJD(base)).Minutes()
	for i := 1; i < 366; i++ {
		cur := EquationOfTime(CivilJD(base.AddDate(0, 0, i))).Minutes()
		if (prev < 0) != (cur < 0) {
			crossings++
		}
		prev = cur
	}
	if crossings != 4 {
		t.Errorf("均時差一年歸零 %d 次，應為 4 次", crossings)
	}
}
