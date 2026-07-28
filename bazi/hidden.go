package bazi

import "github.com/LukeLogix/destiny-core/ganzhi"

// HiddenStemType 藏干的分級。本氣力量最強，餘氣最弱。
type HiddenStemType uint8

const (
	PrimaryQi  HiddenStemType = iota // 本氣
	MiddleQi                         // 中氣
	ResidualQi                       // 餘氣
)

// HiddenStem 地支所藏的天干及其分級
type HiddenStem struct {
	Stem ganzhi.StemIndex
	Type HiddenStemType
}

// HiddenTenGod 藏干及其對日主的十神
type HiddenTenGod struct {
	HiddenStem
	TenGod TenGod
}

// HiddenStemSect 藏干口徑。各家對餘氣的取捨不同，故為可選而非寫死。
type HiddenStemSect uint8

const (
	// HiddenStemStandard 亥藏壬甲、午藏丁己（預設，與 tyme4go 一致）
	HiddenStemStandard HiddenStemSect = iota
	// HiddenStemWithEarth 亥另藏餘氣戊
	HiddenStemWithEarth
)

// 藏干表，依本氣、中氣、餘氣排列。本氣的五行必與地支本身相同。
var hiddenStemTable = [ganzhi.BranchCount][]ganzhi.StemIndex{
	{9},       // 子：癸
	{5, 9, 7}, // 丑：己癸辛
	{0, 2, 4}, // 寅：甲丙戊
	{1},       // 卯：乙
	{4, 1, 9}, // 辰：戊乙癸
	{2, 6, 4}, // 巳：丙庚戊
	{3, 5},    // 午：丁己
	{5, 3, 1}, // 未：己丁乙
	{6, 8, 4}, // 申：庚壬戊
	{7},       // 酉：辛
	{4, 7, 3}, // 戌：戊辛丁
	{8, 0},    // 亥：壬甲
}

// HiddenStems 取地支藏干。
//
// 回傳新的 slice，呼叫方修改不會污染內部表。
func HiddenStems(b ganzhi.BranchIndex, sect HiddenStemSect) []HiddenStem {
	stems := hiddenStemTable[int(b)%ganzhi.BranchCount]

	out := make([]HiddenStem, 0, len(stems)+1)
	for i, s := range stems {
		out = append(out, HiddenStem{Stem: s, Type: HiddenStemType(i)})
	}

	// 亥在部分流派另藏餘氣戊
	if sect == HiddenStemWithEarth && b == 11 {
		out = append(out, HiddenStem{Stem: 4, Type: ResidualQi})
	}
	return out
}

// HiddenTenGods 取地支藏干及其對日主的十神
func HiddenTenGods(b ganzhi.BranchIndex, dayMaster ganzhi.StemIndex, sect HiddenStemSect) []HiddenTenGod {
	hs := HiddenStems(b, sect)
	out := make([]HiddenTenGod, len(hs))
	for i, h := range hs {
		out[i] = HiddenTenGod{HiddenStem: h, TenGod: TenGodOf(dayMaster, h.Stem)}
	}
	return out
}
