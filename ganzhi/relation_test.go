package ganzhi

import (
	"sort"
	"testing"
)

// find 取出指定種類的關係，並以 Indices 排序後便於比對
func find(rs []Relation, kind RelationKind) []Relation {
	var out []Relation
	for _, r := range rs {
		if r.Kind == kind {
			idx := append([]int(nil), r.Indices...)
			sort.Ints(idx)
			r.Indices = idx
			out = append(out, r)
		}
	}
	return out
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestSixHarmony 六合：子丑合土、寅亥合木、卯戌合火、辰酉合金、巳申合水、午未合土。
func TestSixHarmony(t *testing.T) {
	cases := []struct {
		a, b BranchIndex
		want Element
		desc string
	}{
		{0, 1, Earth, "子丑合土"},
		{2, 11, Wood, "寅亥合木"},
		{3, 10, Fire, "卯戌合火"},
		{4, 9, Metal, "辰酉合金"},
		{5, 8, Water, "巳申合水"},
		{6, 7, Earth, "午未合土"},
	}
	for _, c := range cases {
		rs := find(BranchRelations([]BranchIndex{c.a, c.b}, DefaultRelationOptions()), SixHarmony)
		if len(rs) != 1 {
			t.Errorf("%s: 得 %d 組六合，應為 1 組", c.desc, len(rs))
			continue
		}
		if !sameInts(rs[0].Indices, []int{0, 1}) {
			t.Errorf("%s: Indices 為 %v，應為 [0 1]", c.desc, rs[0].Indices)
		}
		if rs[0].Transform == nil || *rs[0].Transform != c.want {
			t.Errorf("%s: 化神為 %v，應為 %d", c.desc, rs[0].Transform, c.want)
		}
	}
}

// TestTrinity 三合局須三支齊全：申子辰水、亥卯未木、寅午戌火、巳酉丑金。
func TestTrinity(t *testing.T) {
	cases := []struct {
		bs   []BranchIndex
		want Element
		desc string
	}{
		{[]BranchIndex{8, 0, 4}, Water, "申子辰合水"},
		{[]BranchIndex{11, 3, 7}, Wood, "亥卯未合木"},
		{[]BranchIndex{2, 6, 10}, Fire, "寅午戌合火"},
		{[]BranchIndex{5, 9, 1}, Metal, "巳酉丑合金"},
	}
	for _, c := range cases {
		rs := find(BranchRelations(c.bs, DefaultRelationOptions()), Trinity)
		if len(rs) != 1 {
			t.Errorf("%s: 得 %d 組三合，應為 1 組", c.desc, len(rs))
			continue
		}
		if !sameInts(rs[0].Indices, []int{0, 1, 2}) {
			t.Errorf("%s: Indices 為 %v，應為 [0 1 2]", c.desc, rs[0].Indices)
		}
		if rs[0].Transform == nil || *rs[0].Transform != c.want {
			t.Errorf("%s: 化神為 %v，應為 %d", c.desc, rs[0].Transform, c.want)
		}
	}
}

// TestTrinityRequiresAllThree 缺一支不成三合，只可能成半合。
func TestTrinityRequiresAllThree(t *testing.T) {
	// 申子（缺辰）
	rs := find(BranchRelations([]BranchIndex{8, 0}, DefaultRelationOptions()), Trinity)
	if len(rs) != 0 {
		t.Errorf("申子缺辰，不應構成三合，實得 %d 組", len(rs))
	}
}

// TestHalfTrinityOption 半合預設開啟，關閉後不再回報。
func TestHalfTrinityOption(t *testing.T) {
	opt := DefaultRelationOptions()
	if !opt.IncludeHalfTrinity {
		t.Fatal("半合應預設開啟")
	}
	rs := find(BranchRelations([]BranchIndex{8, 0}, opt), HalfTrinity)
	if len(rs) != 1 {
		t.Errorf("申子為半合，得 %d 組", len(rs))
	} else if rs[0].Transform == nil || *rs[0].Transform != Water {
		t.Errorf("申子半合化神為 %v，應為水", rs[0].Transform)
	}

	opt.IncludeHalfTrinity = false
	if rs := find(BranchRelations([]BranchIndex{8, 0}, opt), HalfTrinity); len(rs) != 0 {
		t.Errorf("關閉半合後仍回報 %d 組", len(rs))
	}
}

// TestDirectional 三會方局：寅卯辰東方木、巳午未南方火、申酉戌西方金、亥子丑北方水。
func TestDirectional(t *testing.T) {
	cases := []struct {
		bs   []BranchIndex
		want Element
		desc string
	}{
		{[]BranchIndex{2, 3, 4}, Wood, "寅卯辰會木"},
		{[]BranchIndex{5, 6, 7}, Fire, "巳午未會火"},
		{[]BranchIndex{8, 9, 10}, Metal, "申酉戌會金"},
		{[]BranchIndex{11, 0, 1}, Water, "亥子丑會水"},
	}
	for _, c := range cases {
		rs := find(BranchRelations(c.bs, DefaultRelationOptions()), Directional)
		if len(rs) != 1 {
			t.Errorf("%s: 得 %d 組三會，應為 1 組", c.desc, len(rs))
			continue
		}
		if rs[0].Transform == nil || *rs[0].Transform != c.want {
			t.Errorf("%s: 化神為 %v，應為 %d", c.desc, rs[0].Transform, c.want)
		}
	}
}

// TestClash 六沖：子午、丑未、寅申、卯酉、辰戌、巳亥，即相隔六位。
func TestClash(t *testing.T) {
	for a := 0; a < 6; a++ {
		b := a + 6
		rs := find(BranchRelations([]BranchIndex{BranchIndex(a), BranchIndex(b)},
			DefaultRelationOptions()), Clash)
		if len(rs) != 1 {
			t.Errorf("地支 %d 與 %d 應相沖，得 %d 組", a, b, len(rs))
			continue
		}
		if rs[0].Transform != nil {
			t.Errorf("沖不應有化神，實得 %v", rs[0].Transform)
		}
	}
}

// TestPunishment 相刑：寅巳申無恩、丑戌未恃勢、子卯無禮、辰午酉亥自刑。
func TestPunishment(t *testing.T) {
	// 三刑須三支齊全
	for _, c := range []struct {
		bs   []BranchIndex
		desc string
	}{
		{[]BranchIndex{2, 5, 8}, "寅巳申"},
		{[]BranchIndex{1, 10, 7}, "丑戌未"},
	} {
		if rs := find(BranchRelations(c.bs, DefaultRelationOptions()), Punishment); len(rs) != 1 {
			t.Errorf("%s 三刑，得 %d 組", c.desc, len(rs))
		}
	}
	// 子卯相刑（兩支）
	if rs := find(BranchRelations([]BranchIndex{0, 3}, DefaultRelationOptions()), Punishment); len(rs) != 1 {
		t.Errorf("子卯相刑，得 %d 組", len(rs))
	}
	// 自刑須同支出現兩次
	for _, b := range []BranchIndex{4, 6, 9, 11} {
		rs := find(BranchRelations([]BranchIndex{b, b}, DefaultRelationOptions()), Punishment)
		if len(rs) != 1 {
			t.Errorf("地支 %d 自刑，得 %d 組", b, len(rs))
		}
	}
	// 非自刑的地支重複出現不構成刑
	if rs := find(BranchRelations([]BranchIndex{0, 0}, DefaultRelationOptions()), Punishment); len(rs) != 0 {
		t.Errorf("子子非自刑，不應回報，得 %d 組", len(rs))
	}
}

// TestHarm 六害：子未、丑午、寅巳、卯辰、申亥、酉戌。
func TestHarm(t *testing.T) {
	cases := [][2]BranchIndex{{0, 7}, {1, 6}, {2, 5}, {3, 4}, {8, 11}, {9, 10}}
	for _, c := range cases {
		if rs := find(BranchRelations([]BranchIndex{c[0], c[1]},
			DefaultRelationOptions()), Harm); len(rs) != 1 {
			t.Errorf("地支 %d 與 %d 相害，得 %d 組", c[0], c[1], len(rs))
		}
	}
}

// TestDestructionOption 相破多數流派不採用，故預設關閉。
func TestDestructionOption(t *testing.T) {
	opt := DefaultRelationOptions()
	if opt.IncludeDestruction {
		t.Fatal("相破應預設關閉")
	}
	if rs := find(BranchRelations([]BranchIndex{0, 9}, opt), Destruction); len(rs) != 0 {
		t.Errorf("預設不應回報相破，得 %d 組", len(rs))
	}

	opt.IncludeDestruction = true
	cases := [][2]BranchIndex{{0, 9}, {3, 6}, {4, 1}, {7, 10}, {2, 11}, {5, 8}}
	for _, c := range cases {
		if rs := find(BranchRelations([]BranchIndex{c[0], c[1]}, opt), Destruction); len(rs) != 1 {
			t.Errorf("地支 %d 與 %d 相破，得 %d 組", c[0], c[1], len(rs))
		}
	}
}

// TestStemHarmony 天干五合：甲己合土、乙庚合金、丙辛合水、丁壬合木、戊癸合火。
func TestStemHarmony(t *testing.T) {
	cases := []struct {
		a, b StemIndex
		want Element
		desc string
	}{
		{0, 5, Earth, "甲己合土"},
		{1, 6, Metal, "乙庚合金"},
		{2, 7, Water, "丙辛合水"},
		{3, 8, Wood, "丁壬合木"},
		{4, 9, Fire, "戊癸合火"},
	}
	for _, c := range cases {
		rs := find(StemRelations([]StemIndex{c.a, c.b}, DefaultRelationOptions()), StemHarmony)
		if len(rs) != 1 {
			t.Errorf("%s: 得 %d 組，應為 1 組", c.desc, len(rs))
			continue
		}
		if rs[0].Transform == nil || *rs[0].Transform != c.want {
			t.Errorf("%s: 化神為 %v，應為 %d", c.desc, rs[0].Transform, c.want)
		}
	}
}

// TestStemClash 天干四沖：甲庚、乙辛、丙壬、丁癸。戊己屬土居中，不沖。
func TestStemClash(t *testing.T) {
	for _, c := range [][2]StemIndex{{0, 6}, {1, 7}, {2, 8}, {3, 9}} {
		if rs := find(StemRelations([]StemIndex{c[0], c[1]},
			DefaultRelationOptions()), StemClash); len(rs) != 1 {
			t.Errorf("天干 %d 與 %d 相沖，得 %d 組", c[0], c[1], len(rs))
		}
	}
	// 戊己不參與沖
	for _, c := range [][2]StemIndex{{4, 0}, {4, 9}, {5, 1}} {
		if rs := find(StemRelations([]StemIndex{c[0], c[1]},
			DefaultRelationOptions()), StemClash); len(rs) != 0 {
			t.Errorf("天干 %d 與 %d 不應相沖，得 %d 組", c[0], c[1], len(rs))
		}
	}
}

// TestIndicesAreCallerDefined 關係只回報「傳入 slice 中的第幾個」，
// 不帶年月日時等特定術數的位置語意——共用層不得綁定單一核心的概念。
func TestIndicesAreCallerDefined(t *testing.T) {
	// 把子丑放在第 2、3 個位置，Indices 應隨之改變
	rs := find(BranchRelations([]BranchIndex{6, 6, 0, 1}, DefaultRelationOptions()), SixHarmony)
	if len(rs) != 1 {
		t.Fatalf("應有 1 組六合，得 %d 組", len(rs))
	}
	if !sameInts(rs[0].Indices, []int{2, 3}) {
		t.Errorf("Indices 為 %v，應為 [2 3]", rs[0].Indices)
	}
}

// TestNoDuplicateReports 同一組關係只回報一次，不因遍歷順序重複。
func TestNoDuplicateReports(t *testing.T) {
	rs := BranchRelations([]BranchIndex{0, 1}, DefaultRelationOptions())
	seen := map[RelationKind]int{}
	for _, r := range rs {
		seen[r.Kind]++
	}
	for kind, n := range seen {
		if n > 1 {
			t.Errorf("種類 %d 回報 %d 次，應僅一次", kind, n)
		}
	}
}

// TestTransformIsPotentialNotActual 合會的 Transform 表示「若化則化為此五行」，
// 而非斷定已化——化與不化各家分歧極大，核心不做此判定。
func TestTransformIsPotentialNotActual(t *testing.T) {
	// 沖刑害破皆無化神
	for _, kind := range []RelationKind{Clash, Punishment, Harm} {
		var bs []BranchIndex
		switch kind {
		case Clash:
			bs = []BranchIndex{0, 6}
		case Punishment:
			bs = []BranchIndex{0, 3}
		case Harm:
			bs = []BranchIndex{0, 7}
		}
		for _, r := range find(BranchRelations(bs, DefaultRelationOptions()), kind) {
			if r.Transform != nil {
				t.Errorf("種類 %d 不應有化神，實得 %v", kind, r.Transform)
			}
		}
	}
}
