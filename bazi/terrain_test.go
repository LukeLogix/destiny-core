package bazi

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// TestTerrainYangStems 陽干長生位：甲亥、丙戊寅、庚巳、壬申。
func TestTerrainYangStems(t *testing.T) {
	cases := []struct {
		stem   ganzhi.StemIndex
		branch ganzhi.BranchIndex
		desc   string
	}{
		{0, 11, "甲長生在亥"},
		{2, 2, "丙長生在寅"},
		{4, 2, "戊長生在寅"},
		{6, 5, "庚長生在巳"},
		{8, 8, "壬長生在申"},
	}
	for _, c := range cases {
		if got := TerrainOf(c.stem, c.branch, TerrainYinReverse); got != LongLife {
			t.Errorf("%s: 得 %d，應為長生 %d", c.desc, got, LongLife)
		}
	}
}

// TestTerrainYinReverse 陰干逆行（傳統派）：乙長生在午、丁己酉、辛子、癸卯。
func TestTerrainYinReverse(t *testing.T) {
	cases := []struct {
		stem   ganzhi.StemIndex
		branch ganzhi.BranchIndex
		desc   string
	}{
		{1, 6, "乙長生在午"},
		{3, 9, "丁長生在酉"},
		{5, 9, "己長生在酉"},
		{7, 0, "辛長生在子"},
		{9, 3, "癸長生在卯"},
	}
	for _, c := range cases {
		if got := TerrainOf(c.stem, c.branch, TerrainYinReverse); got != LongLife {
			t.Errorf("%s: 得 %d，應為長生", c.desc, got)
		}
	}
	// 乙逆行：午長生 → 巳沐浴 → 辰冠帶 → 卯臨官 → 寅帝旺 → … → 亥死
	seq := []struct {
		branch ganzhi.BranchIndex
		want   Terrain
		desc   string
	}{
		{6, LongLife, "乙午長生"},
		{5, Bath, "乙巳沐浴"},
		{4, Cap, "乙辰冠帶"},
		{3, Officer, "乙卯臨官"},
		{2, Prosperity, "乙寅帝旺"},
		{11, Death, "乙亥死"},
		{9, Void, "乙酉絕"},
	}
	for _, c := range seq {
		if got := TerrainOf(1, c.branch, TerrainYinReverse); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestTerrainSameBirthSect 陰陽同生同死派：陰干與同五行的陽干共用長生位且順行。
func TestTerrainSameBirthSect(t *testing.T) {
	// 乙在此派長生於亥（與甲同），而非午
	if got := TerrainOf(1, 11, TerrainSameBirth); got != LongLife {
		t.Errorf("陰陽同生派乙亥應為長生，得 %d", got)
	}
	if got := TerrainOf(1, 6, TerrainSameBirth); got == LongLife {
		t.Error("陰陽同生派乙午不應為長生")
	}
	// 陽干在兩派結果相同
	for _, s := range []ganzhi.StemIndex{0, 2, 4, 6, 8} {
		for b := 0; b < ganzhi.BranchCount; b++ {
			a := TerrainOf(s, ganzhi.BranchIndex(b), TerrainYinReverse)
			c := TerrainOf(s, ganzhi.BranchIndex(b), TerrainSameBirth)
			if a != c {
				t.Errorf("陽干 %d 在地支 %d 的兩派結果不同：%d vs %d", s, b, a, c)
			}
		}
	}
}

// TestTerrainCoversTwelve 任一天干走遍十二支，恰好經歷十二種狀態。
func TestTerrainCoversTwelve(t *testing.T) {
	for _, sect := range []TerrainSect{TerrainYinReverse, TerrainSameBirth} {
		for s := 0; s < ganzhi.StemCount; s++ {
			seen := map[Terrain]bool{}
			for b := 0; b < ganzhi.BranchCount; b++ {
				seen[TerrainOf(ganzhi.StemIndex(s), ganzhi.BranchIndex(b), sect)] = true
			}
			if len(seen) != TerrainCount {
				t.Errorf("口徑 %d 天干 %d 只經歷 %d 種狀態，應為 %d 種",
					sect, s, len(seen), TerrainCount)
			}
		}
	}
}

// TestSoundOfSexagenary 納音每兩個甲子一組，六十甲子共三十種。
func TestSoundOfSexagenary(t *testing.T) {
	cases := []struct {
		sex  ganzhi.SexagenaryIndex
		want SoundIndex
		desc string
	}{
		{0, 0, "甲子海中金"},
		{1, 0, "乙丑亦海中金"},
		{2, 1, "丙寅爐中火"},
		{6, 3, "庚午路旁土"},
		{17, 8, "辛巳白蠟金"},
		{21, 10, "乙酉泉中水"},
		{59, 29, "癸亥大海水"},
	}
	for _, c := range cases {
		if got := SoundOf(c.sex); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestSoundElement 納音各有五行，1990-05-20 命例：
// 庚午路旁土、辛巳白蠟金、乙酉泉中水、辛巳白蠟金。
func TestSoundElement(t *testing.T) {
	cases := []struct {
		sex  ganzhi.SexagenaryIndex
		want ganzhi.Element
		desc string
	}{
		{6, ganzhi.Earth, "庚午路旁土"},
		{17, ganzhi.Metal, "辛巳白蠟金"},
		{21, ganzhi.Water, "乙酉泉中水"},
		{0, ganzhi.Metal, "甲子海中金"},
		{59, ganzhi.Water, "癸亥大海水"},
	}
	for _, c := range cases {
		if got := SoundOf(c.sex).Element(); got != c.want {
			t.Errorf("%s: 五行得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestSoundCoversThirty 六十甲子恰好對應三十種納音，每種各兩個。
func TestSoundCoversThirty(t *testing.T) {
	count := map[SoundIndex]int{}
	for i := 0; i < ganzhi.SexagenaryCount; i++ {
		count[SoundOf(ganzhi.SexagenaryIndex(i))]++
	}
	if len(count) != SoundCount {
		t.Errorf("納音共 %d 種，應為 %d 種", len(count), SoundCount)
	}
	for s, n := range count {
		if n != 2 {
			t.Errorf("納音 %d 對應 %d 個甲子，應為 2 個", s, n)
		}
	}
}

// TestPillarCarriesTerrainAndSound 命盤各柱須帶十二長生與納音。
func TestPillarCarriesTerrainAndSound(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	// 日主乙木：年支午為長生、日支酉為絕
	if c.Year.Terrain != LongLife {
		t.Errorf("年柱地勢為 %d，日主乙見午應為長生", c.Year.Terrain)
	}
	if c.Day.Terrain != Void {
		t.Errorf("日柱地勢為 %d，日主乙見酉應為絕", c.Day.Terrain)
	}
	// 納音
	if c.Year.Sound.Element() != ganzhi.Earth {
		t.Errorf("年柱庚午納音五行為 %d，應為土（路旁土）", c.Year.Sound.Element())
	}
	if c.Day.Sound.Element() != ganzhi.Water {
		t.Errorf("日柱乙酉納音五行為 %d，應為水（泉中水）", c.Day.Sound.Element())
	}
}
