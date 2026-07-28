package ganzhi

import (
	"errors"
	"testing"
)

// TestSexagenaryRoundTrip 六十甲子索引與干支對的往返必須無損。
func TestSexagenaryRoundTrip(t *testing.T) {
	for i := 0; i < 60; i++ {
		x := SexagenaryIndex(i)
		got, err := Sexagenary(x.Stem(), x.Branch())
		if err != nil {
			t.Fatalf("index %d 分解後無法重組: %v", i, err)
		}
		if got != x {
			t.Errorf("index %d 往返後變成 %d", i, got)
		}
	}
}

// TestSexagenaryKnownValues 釘住幾個公認的定位點，避免整體偏移。
func TestSexagenaryKnownValues(t *testing.T) {
	cases := []struct {
		idx    SexagenaryIndex
		stem   StemIndex
		branch BranchIndex
		desc   string
	}{
		{0, 0, 0, "甲子"},
		{1, 1, 1, "乙丑"},
		{10, 0, 10, "甲戌"},
		{59, 9, 11, "癸亥"},
	}
	for _, c := range cases {
		if got := c.idx.Stem(); got != c.stem {
			t.Errorf("%s: 天干為 %d，應為 %d", c.desc, got, c.stem)
		}
		if got := c.idx.Branch(); got != c.branch {
			t.Errorf("%s: 地支為 %d，應為 %d", c.desc, got, c.branch)
		}
		if got, err := Sexagenary(c.stem, c.branch); err != nil || got != c.idx {
			t.Errorf("%s: 組合得 %d (err=%v)，應為 %d", c.desc, got, err, c.idx)
		}
	}
}

// TestSexagenaryRejectsMismatchedPolarity 陽干只配陽支、陰干只配陰支。
// 故 10×12 的組合中只有 60 種成立，「甲丑」這類不存在。
func TestSexagenaryRejectsMismatchedPolarity(t *testing.T) {
	invalid := []struct {
		s StemIndex
		b BranchIndex
	}{
		{0, 1},  // 甲丑：陽干配陰支
		{1, 0},  // 乙子：陰干配陽支
		{0, 11}, // 甲亥
		{9, 0},  // 癸子
	}
	for _, c := range invalid {
		if _, err := Sexagenary(c.s, c.b); !errors.Is(err, ErrPolarityMismatch) {
			t.Errorf("干 %d 支 %d 陰陽不符，應回 ErrPolarityMismatch，實得 %v", c.s, c.b, err)
		}
	}
}

// TestSexagenaryCoversExactlySixty 十干十二支恰好組出 60 種，不多不少。
func TestSexagenaryCoversExactlySixty(t *testing.T) {
	seen := map[SexagenaryIndex]bool{}
	valid := 0
	for s := 0; s < 10; s++ {
		for b := 0; b < 12; b++ {
			idx, err := Sexagenary(StemIndex(s), BranchIndex(b))
			if err != nil {
				continue
			}
			valid++
			if seen[idx] {
				t.Errorf("索引 %d 由多組干支產生，應唯一", idx)
			}
			seen[idx] = true
		}
	}
	if valid != 60 {
		t.Errorf("有效組合 %d 種，應為 60 種", valid)
	}
	if len(seen) != 60 {
		t.Errorf("涵蓋 %d 個索引，應為 60 個", len(seen))
	}
}

// TestSexagenaryRejectsOutOfRange 越界索引須報錯而非回傳錯誤結果。
func TestSexagenaryRejectsOutOfRange(t *testing.T) {
	if _, err := Sexagenary(10, 0); !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("天干 10 越界，應回 ErrIndexOutOfRange，實得 %v", err)
	}
	if _, err := Sexagenary(0, 12); !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("地支 12 越界，應回 ErrIndexOutOfRange，實得 %v", err)
	}
}

// TestStemNext 天干循環：癸的下一個回到甲。
func TestStemNext(t *testing.T) {
	if got := StemIndex(9).Next(1); got != 0 {
		t.Errorf("癸的下一個為 %d，應為 0（甲）", got)
	}
	if got := StemIndex(0).Next(-1); got != 9 {
		t.Errorf("甲的上一個為 %d，應為 9（癸）", got)
	}
	if got := StemIndex(0).Next(23); got != 3 {
		t.Errorf("甲進 23 步為 %d，應為 3", got)
	}
}

// TestBranchNext 地支循環：亥的下一個回到子。
func TestBranchNext(t *testing.T) {
	if got := BranchIndex(11).Next(1); got != 0 {
		t.Errorf("亥的下一個為 %d，應為 0（子）", got)
	}
	if got := BranchIndex(0).Next(-1); got != 11 {
		t.Errorf("子的上一個為 %d，應為 11（亥）", got)
	}
	if got := BranchIndex(2).Next(-14); got != 0 {
		t.Errorf("寅退 14 步為 %d，應為 0", got)
	}
}

// TestSexagenaryNext 六十甲子循環。
func TestSexagenaryNext(t *testing.T) {
	if got := SexagenaryIndex(59).Next(1); got != 0 {
		t.Errorf("癸亥的下一個為 %d，應為 0（甲子）", got)
	}
	if got := SexagenaryIndex(0).Next(-1); got != 59 {
		t.Errorf("甲子的上一個為 %d，應為 59（癸亥）", got)
	}
}

// TestHourBranchFromClock 時辰劃分：23-01 為子時，餘每兩小時一支。
func TestHourBranchFromClock(t *testing.T) {
	cases := []struct {
		hour int
		want BranchIndex
		desc string
	}{
		{23, 0, "23 時為子"},
		{0, 0, "0 時為子"},
		{1, 1, "1 時為丑"},
		{2, 1, "2 時為丑"},
		{11, 6, "11 時為午"},
		{12, 6, "12 時為午"},
		{13, 7, "13 時為未"},
		{22, 11, "22 時為亥"},
	}
	for _, c := range cases {
		if got := HourBranch(c.hour); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestHourStem 五鼠遁：由日干推時干。甲己日起甲子時。
func TestHourStem(t *testing.T) {
	cases := []struct {
		dayStem StemIndex
		hourBr  BranchIndex
		want    StemIndex
		desc    string
	}{
		{0, 0, 0, "甲日子時為甲"},
		{5, 0, 0, "己日子時為甲"},
		{1, 0, 2, "乙日子時為丙"},
		{6, 0, 2, "庚日子時為丙"},
		{0, 1, 1, "甲日丑時為乙"},
		{4, 0, 8, "戊日子時為壬（戊癸何方發，壬子是真途）"},
		{4, 11, 9, "戊日亥時為癸"},
	}
	for _, c := range cases {
		if got := HourStem(c.dayStem, c.hourBr); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestMonthStem 五虎遁：由年干推月干。甲己之年丙作首（寅月為丙）。
func TestMonthStem(t *testing.T) {
	cases := []struct {
		yearStem StemIndex
		monthBr  BranchIndex
		want     StemIndex
		desc     string
	}{
		{0, 2, 2, "甲年寅月為丙"},
		{5, 2, 2, "己年寅月為丙"},
		{1, 2, 4, "乙年寅月為戊"},
		{2, 2, 6, "丙年寅月為庚"},
		{3, 2, 8, "丁年寅月為壬"},
		{4, 2, 0, "戊年寅月為甲"},
		{0, 3, 3, "甲年卯月為丁"},
		{0, 10, 0, "甲年戌月為甲"},
		{0, 1, 3, "甲年丑月為丁（寅起丙，順行至丑為第十一位）"},
	}
	for _, c := range cases {
		if got := MonthStem(c.yearStem, c.monthBr); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}
