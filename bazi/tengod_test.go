package bazi

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// TestTenGodOfSelf 日主對自身為比肩。
func TestTenGodOfSelf(t *testing.T) {
	for s := 0; s < ganzhi.StemCount; s++ {
		me := ganzhi.StemIndex(s)
		if got := TenGodOf(me, me); got != Peer {
			t.Errorf("天干 %d 對自身為 %d，應為比肩 %d", s, got, Peer)
		}
	}
}

// TestTenGodOfKnownChart 以 1990-05-20 10:30 的實際命盤驗證，
// 日主乙木對年干庚、月干辛、時干辛。
func TestTenGodOfKnownChart(t *testing.T) {
	me := ganzhi.StemIndex(1) // 乙
	cases := []struct {
		other ganzhi.StemIndex
		want  TenGod
		desc  string
	}{
		{6, DirectOfficer, "乙見庚為正官（金剋木，陰陽異）"},
		{7, SevenKilling, "乙見辛為七殺（金剋木，陰陽同）"},
		{3, Output, "乙見丁為食神（木生火，陰陽同）"},
		{2, Hurting, "乙見丙為傷官（木生火，陰陽異）"},
		{5, IndirectWealth, "乙見己為偏財（木剋土，陰陽同）"},
		{4, DirectWealth, "乙見戊為正財（木剋土，陰陽異）"},
		{9, IndirectResource, "乙見癸為偏印（水生木，陰陽同）"},
		{8, DirectResource, "乙見壬為正印（水生木，陰陽異）"},
		{1, Peer, "乙見乙為比肩"},
		{0, Rival, "乙見甲為劫財"},
	}
	for _, c := range cases {
		if got := TenGodOf(me, c.other); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestTenGodCoversAllTen 任一日主對十天干，恰好產生十種十神，不重不漏。
func TestTenGodCoversAllTen(t *testing.T) {
	for s := 0; s < ganzhi.StemCount; s++ {
		me := ganzhi.StemIndex(s)
		seen := map[TenGod]bool{}
		for o := 0; o < ganzhi.StemCount; o++ {
			seen[TenGodOf(me, ganzhi.StemIndex(o))] = true
		}
		if len(seen) != TenGodCount {
			t.Errorf("日主 %d 只產生 %d 種十神，應為 %d 種", s, len(seen), TenGodCount)
		}
	}
}

// TestTenGodPolarityRule 同陰陽者為偏（比肩/食神/偏財/七殺/偏印），
// 異陰陽者為正（劫財/傷官/正財/正官/正印）。
func TestTenGodPolarityRule(t *testing.T) {
	samePolarity := map[TenGod]bool{
		Peer: true, Output: true, IndirectWealth: true,
		SevenKilling: true, IndirectResource: true,
	}
	for s := 0; s < ganzhi.StemCount; s++ {
		for o := 0; o < ganzhi.StemCount; o++ {
			me, other := ganzhi.StemIndex(s), ganzhi.StemIndex(o)
			same := me.Polarity() == other.Polarity()
			god := TenGodOf(me, other)
			if same != samePolarity[god] {
				t.Errorf("日主 %d 見 %d 得十神 %d，陰陽同=%v 與該十神的性質不符",
					s, o, god, same)
			}
		}
	}
}

// TestTenGodRelationCategory 十神依五行關係分五類，每類兩個。
func TestTenGodRelationCategory(t *testing.T) {
	me := ganzhi.StemIndex(0) // 甲木
	cases := []struct {
		other ganzhi.StemIndex
		want  [2]TenGod // 該五行關係對應的兩個十神
		desc  string
	}{
		{0, [2]TenGod{Peer, Rival}, "同我"},
		{2, [2]TenGod{Output, Hurting}, "我生"},
		{4, [2]TenGod{IndirectWealth, DirectWealth}, "我剋"},
		{6, [2]TenGod{SevenKilling, DirectOfficer}, "剋我"},
		{8, [2]TenGod{IndirectResource, DirectResource}, "生我"},
	}
	for _, c := range cases {
		got := TenGodOf(me, c.other)
		if got != c.want[0] && got != c.want[1] {
			t.Errorf("%s: 甲見天干 %d 得 %d，應為 %d 或 %d",
				c.desc, c.other, got, c.want[0], c.want[1])
		}
	}
}
