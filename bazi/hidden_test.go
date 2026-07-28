package bazi

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// TestHiddenStemsStandard 預設藏干表，與 tyme4go 實測一致。
func TestHiddenStemsStandard(t *testing.T) {
	cases := []struct {
		branch ganzhi.BranchIndex
		want   []ganzhi.StemIndex
		desc   string
	}{
		{0, []ganzhi.StemIndex{9}, "子藏癸"},
		{1, []ganzhi.StemIndex{5, 9, 7}, "丑藏己癸辛"},
		{2, []ganzhi.StemIndex{0, 2, 4}, "寅藏甲丙戊"},
		{3, []ganzhi.StemIndex{1}, "卯藏乙"},
		{4, []ganzhi.StemIndex{4, 1, 9}, "辰藏戊乙癸"},
		{5, []ganzhi.StemIndex{2, 6, 4}, "巳藏丙庚戊"},
		{6, []ganzhi.StemIndex{3, 5}, "午藏丁己"},
		{7, []ganzhi.StemIndex{5, 3, 1}, "未藏己丁乙"},
		{8, []ganzhi.StemIndex{6, 8, 4}, "申藏庚壬戊"},
		{9, []ganzhi.StemIndex{7}, "酉藏辛"},
		{10, []ganzhi.StemIndex{4, 7, 3}, "戌藏戊辛丁"},
		{11, []ganzhi.StemIndex{8, 0}, "亥藏壬甲"},
	}
	for _, c := range cases {
		got := HiddenStems(c.branch, HiddenStemStandard)
		if len(got) != len(c.want) {
			t.Errorf("%s: 得 %d 個藏干，應為 %d 個", c.desc, len(got), len(c.want))
			continue
		}
		for i, w := range c.want {
			if got[i].Stem != w {
				t.Errorf("%s: 第 %d 個藏干為 %d，應為 %d", c.desc, i, got[i].Stem, w)
			}
		}
	}
}

// TestHiddenStemsQiOrder 藏干依本氣、中氣、餘氣排列，本氣必為第一個。
func TestHiddenStemsQiOrder(t *testing.T) {
	wantTypes := []HiddenStemType{PrimaryQi, MiddleQi, ResidualQi}
	for b := 0; b < ganzhi.BranchCount; b++ {
		hs := HiddenStems(ganzhi.BranchIndex(b), HiddenStemStandard)
		if len(hs) == 0 {
			t.Errorf("地支 %d 無藏干", b)
			continue
		}
		for i, h := range hs {
			if h.Type != wantTypes[i] {
				t.Errorf("地支 %d 第 %d 個藏干類型為 %d，應為 %d",
					b, i, h.Type, wantTypes[i])
			}
		}
	}
}

// TestHiddenStemsPrimaryMatchesBranchElement 本氣的五行必與地支本身相同。
func TestHiddenStemsPrimaryMatchesBranchElement(t *testing.T) {
	for b := 0; b < ganzhi.BranchCount; b++ {
		br := ganzhi.BranchIndex(b)
		hs := HiddenStems(br, HiddenStemStandard)
		if got, want := hs[0].Stem.Element(), br.Element(); got != want {
			t.Errorf("地支 %d 本氣五行為 %d，地支五行為 %d，應相同", b, got, want)
		}
	}
}

// TestHiddenStemsWithEarthSect 另一流派：亥藏壬甲戊，多一個餘氣戊。
// 這是各家分歧處，故做成可選口徑而非寫死。
func TestHiddenStemsWithEarthSect(t *testing.T) {
	std := HiddenStems(11, HiddenStemStandard)
	alt := HiddenStems(11, HiddenStemWithEarth)

	if len(std) != 2 {
		t.Errorf("標準口徑亥藏 %d 干，應為 2（壬甲）", len(std))
	}
	if len(alt) != 3 {
		t.Fatalf("另一口徑亥藏 %d 干，應為 3（壬甲戊）", len(alt))
	}
	if alt[2].Stem != 4 {
		t.Errorf("另一口徑亥的第三個藏干為 %d，應為戊(4)", alt[2].Stem)
	}
	if alt[2].Type != ResidualQi {
		t.Errorf("另一口徑亥的第三個藏干類型為 %d，應為餘氣", alt[2].Type)
	}
	// 其餘地支兩套口徑應相同
	for b := 0; b < 11; b++ {
		a := HiddenStems(ganzhi.BranchIndex(b), HiddenStemStandard)
		c := HiddenStems(ganzhi.BranchIndex(b), HiddenStemWithEarth)
		if len(a) != len(c) {
			t.Errorf("地支 %d 兩套口徑藏干數不同：%d vs %d", b, len(a), len(c))
		}
	}
}

// TestHiddenStemsReturnsCopy 回傳值被修改不得污染內部表。
func TestHiddenStemsReturnsCopy(t *testing.T) {
	first := HiddenStems(2, HiddenStemStandard)
	original := first[0].Stem
	first[0].Stem = 99

	second := HiddenStems(2, HiddenStemStandard)
	if second[0].Stem != original {
		t.Errorf("內部表被外部修改污染：得 %d，應為 %d", second[0].Stem, original)
	}
}

// TestHiddenTenGods 地支藏干各自對日主取十神。
// 以 1990-05-20 的月支巳為例，日主乙木：丙傷官、庚正官、戊正財。
func TestHiddenTenGods(t *testing.T) {
	me := ganzhi.StemIndex(1) // 乙
	got := HiddenTenGods(5, me, HiddenStemStandard)

	want := []TenGod{Hurting, DirectOfficer, DirectWealth}
	if len(got) != len(want) {
		t.Fatalf("巳藏干十神得 %d 個，應為 %d 個", len(got), len(want))
	}
	for i, w := range want {
		if got[i].TenGod != w {
			t.Errorf("巳第 %d 個藏干十神為 %d，應為 %d", i, got[i].TenGod, w)
		}
	}
}
