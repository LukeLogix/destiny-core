package ganzhi

import "testing"

// TestTerrainYangStems 陽干長生位：甲亥、丙戊寅、庚巳、壬申。
func TestTerrainYangStems(t *testing.T) {
	cases := []struct {
		stem   StemIndex
		branch BranchIndex
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
		stem   StemIndex
		branch BranchIndex
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
		branch BranchIndex
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
	for _, s := range []StemIndex{0, 2, 4, 6, 8} {
		for b := 0; b < BranchCount; b++ {
			a := TerrainOf(s, BranchIndex(b), TerrainYinReverse)
			c := TerrainOf(s, BranchIndex(b), TerrainSameBirth)
			if a != c {
				t.Errorf("陽干 %d 在地支 %d 的兩派結果不同：%d vs %d", s, b, a, c)
			}
		}
	}
}

// TestTerrainCoversTwelve 任一天干走遍十二支，恰好經歷十二種狀態。
func TestTerrainCoversTwelve(t *testing.T) {
	for _, sect := range []TerrainSect{TerrainYinReverse, TerrainSameBirth} {
		for s := 0; s < StemCount; s++ {
			seen := map[Terrain]bool{}
			for b := 0; b < BranchCount; b++ {
				seen[TerrainOf(StemIndex(s), BranchIndex(b), sect)] = true
			}
			if len(seen) != TerrainCount {
				t.Errorf("口徑 %d 天干 %d 只經歷 %d 種狀態，應為 %d 種",
					sect, s, len(seen), TerrainCount)
			}
		}
	}
}
