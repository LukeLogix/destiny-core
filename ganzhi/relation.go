package ganzhi

// RelationKind 干支關係的種類
type RelationKind uint8

const (
	StemHarmony RelationKind = iota // 天干五合
	StemClash                       // 天干沖
	SixHarmony                      // 地支六合
	Trinity                         // 三合局
	HalfTrinity                     // 半合
	Directional                     // 三會方局
	Clash                           // 地支六沖
	Punishment                      // 相刑
	Harm                            // 六害
	Destruction                     // 相破
)

// Relation 一組成立的干支關係。
//
// Indices 指的是「傳入 slice 中的第幾個元素」，不帶年月日時等特定術數的位置語意——
// 六爻是爻位、紫微是宮位，共用層不綁定任一核心的概念，由呼叫方自行賦予意義。
type Relation struct {
	Kind      RelationKind
	Indices   []int
	Transform *Element // 若化則化為此五行；僅合會類有值，且僅表潛在而非既成
}

// RelationOptions 各流派取捨不同的關係項目
type RelationOptions struct {
	IncludeHalfTrinity bool // 半合（三合缺一且含旺支）
	IncludeDestruction bool // 相破
}

// DefaultRelationOptions 半合開啟、相破關閉——後者多數流派不採用。
func DefaultRelationOptions() RelationOptions {
	return RelationOptions{IncludeHalfTrinity: true, IncludeDestruction: false}
}

// ---- 關係資料表 ----

// 六合，鍵為排序後的地支對
var sixHarmonyTable = map[[2]BranchIndex]Element{
	{0, 1}:  Earth, // 子丑
	{2, 11}: Wood,  // 寅亥
	{3, 10}: Fire,  // 卯戌
	{4, 9}:  Metal, // 辰酉
	{5, 8}:  Water, // 巳申
	{6, 7}:  Earth, // 午未
}

// 三合局。middle 為旺支，半合須含之。
var trinityTable = []struct {
	members [3]BranchIndex
	middle  BranchIndex
	elem    Element
}{
	{[3]BranchIndex{8, 0, 4}, 0, Water}, // 申子辰
	{[3]BranchIndex{11, 3, 7}, 3, Wood}, // 亥卯未
	{[3]BranchIndex{2, 6, 10}, 6, Fire}, // 寅午戌
	{[3]BranchIndex{5, 9, 1}, 9, Metal}, // 巳酉丑
}

// 三會方局
var directionalTable = []struct {
	members [3]BranchIndex
	elem    Element
}{
	{[3]BranchIndex{2, 3, 4}, Wood},   // 寅卯辰
	{[3]BranchIndex{5, 6, 7}, Fire},   // 巳午未
	{[3]BranchIndex{8, 9, 10}, Metal}, // 申酉戌
	{[3]BranchIndex{11, 0, 1}, Water}, // 亥子丑
}

// 三刑：寅巳申（無恩）、丑戌未（恃勢）
var punishmentTriples = [][3]BranchIndex{
	{2, 5, 8},
	{1, 10, 7},
}

// 子卯相刑（無禮之刑）
var punishmentPair = [2]BranchIndex{0, 3}

// 自刑：辰午酉亥
var selfPunishments = map[BranchIndex]bool{4: true, 6: true, 9: true, 11: true}

// 六害
var harmTable = map[[2]BranchIndex]bool{
	{0, 7}: true, {1, 6}: true, {2, 5}: true,
	{3, 4}: true, {8, 11}: true, {9, 10}: true,
}

// 相破
var destructionTable = map[[2]BranchIndex]bool{
	{0, 9}: true, {3, 6}: true, {1, 4}: true,
	{7, 10}: true, {2, 11}: true, {5, 8}: true,
}

// 天干五合
var stemHarmonyTable = map[[2]StemIndex]Element{
	{0, 5}: Earth, // 甲己
	{1, 6}: Metal, // 乙庚
	{2, 7}: Water, // 丙辛
	{3, 8}: Wood,  // 丁壬
	{4, 9}: Fire,  // 戊癸
}

// 天干四沖。戊己屬土居中，不參與沖。
var stemClashTable = map[[2]StemIndex]bool{
	{0, 6}: true, {1, 7}: true, {2, 8}: true, {3, 9}: true,
}

// ---- 計算 ----

func orderedBranch(a, b BranchIndex) [2]BranchIndex {
	if a > b {
		a, b = b, a
	}
	return [2]BranchIndex{a, b}
}

func orderedStem(a, b StemIndex) [2]StemIndex {
	if a > b {
		a, b = b, a
	}
	return [2]StemIndex{a, b}
}

func elemPtr(e Element) *Element { return &e }

// StemRelations 找出天干之間成立的關係
func StemRelations(ss []StemIndex, opt RelationOptions) []Relation {
	var out []Relation
	for i := 0; i < len(ss); i++ {
		for j := i + 1; j < len(ss); j++ {
			key := orderedStem(ss[i], ss[j])
			if e, ok := stemHarmonyTable[key]; ok {
				out = append(out, Relation{StemHarmony, []int{i, j}, elemPtr(e)})
			}
			if stemClashTable[key] {
				out = append(out, Relation{StemClash, []int{i, j}, nil})
			}
		}
	}
	return out
}

// BranchRelations 找出地支之間成立的關係。
//
// 只報告關係構成，不判定合化是否成真——化與不化取決於得令、引化之神、
// 是否被沖破等條件，各家分歧極大且無定論，故留給上層或人判斷。
func BranchRelations(bs []BranchIndex, opt RelationOptions) []Relation {
	var out []Relation

	// 兩兩關係
	for i := 0; i < len(bs); i++ {
		for j := i + 1; j < len(bs); j++ {
			a, b := bs[i], bs[j]
			key := orderedBranch(a, b)

			if e, ok := sixHarmonyTable[key]; ok {
				out = append(out, Relation{SixHarmony, []int{i, j}, elemPtr(e)})
			}
			if mod(int(a)-int(b), BranchCount) == 6 {
				out = append(out, Relation{Clash, []int{i, j}, nil})
			}
			if harmTable[key] {
				out = append(out, Relation{Harm, []int{i, j}, nil})
			}
			if opt.IncludeDestruction && destructionTable[key] {
				out = append(out, Relation{Destruction, []int{i, j}, nil})
			}
			// 子卯相刑
			if key == orderedBranch(punishmentPair[0], punishmentPair[1]) {
				out = append(out, Relation{Punishment, []int{i, j}, nil})
			}
			// 自刑須同支重複出現
			if a == b && selfPunishments[a] {
				out = append(out, Relation{Punishment, []int{i, j}, nil})
			}
		}
	}

	// 三支關係：先建「地支 → 首次出現位置」的索引
	first := map[BranchIndex]int{}
	for i, b := range bs {
		if _, ok := first[b]; !ok {
			first[b] = i
		}
	}

	for _, tri := range trinityTable {
		var idx []int
		hasMiddle := false
		for _, m := range tri.members {
			if p, ok := first[m]; ok {
				idx = append(idx, p)
				if m == tri.middle {
					hasMiddle = true
				}
			}
		}
		switch {
		case len(idx) == 3:
			out = append(out, Relation{Trinity, idx, elemPtr(tri.elem)})
		case len(idx) == 2 && hasMiddle && opt.IncludeHalfTrinity:
			// 半合須含旺支；不含旺支的兩支（如申辰）本版不視為成局
			out = append(out, Relation{HalfTrinity, idx, elemPtr(tri.elem)})
		}
	}

	for _, dir := range directionalTable {
		var idx []int
		for _, m := range dir.members {
			if p, ok := first[m]; ok {
				idx = append(idx, p)
			}
		}
		if len(idx) == 3 {
			out = append(out, Relation{Directional, idx, elemPtr(dir.elem)})
		}
	}

	for _, tri := range punishmentTriples {
		var idx []int
		for _, m := range tri {
			if p, ok := first[m]; ok {
				idx = append(idx, p)
			}
		}
		if len(idx) == 3 {
			out = append(out, Relation{Punishment, idx, nil})
		}
	}

	return out
}
