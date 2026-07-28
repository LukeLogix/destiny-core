package ganzhi

// Element 五行。
//
// 順序由本套件自行定義，不依賴任何外部庫的排列——凡以 [ElementCount]T
// 傳遞五行資料處皆以此為準，並有測試釘住順序。
//
// 順序刻意排成相生序：下一個即所生者，隔一個即所剋者。
type Element uint8

const (
	Wood Element = iota
	Fire
	Earth
	Metal
	Water

	ElementCount = 5
)

// Generates 我生者。木生火、火生土、土生金、金生水、水生木。
func (e Element) Generates() Element { return Element(mod(int(e)+1, ElementCount)) }

// GeneratedBy 生我者
func (e Element) GeneratedBy() Element { return Element(mod(int(e)-1, ElementCount)) }

// Restrains 我剋者。木剋土、火剋金、土剋水、金剋木、水剋火。
func (e Element) Restrains() Element { return Element(mod(int(e)+2, ElementCount)) }

// RestrainedBy 剋我者
func (e Element) RestrainedBy() Element { return Element(mod(int(e)-2, ElementCount)) }

// Polarity 陰陽
type Polarity uint8

const (
	Yang Polarity = iota
	Yin
)

// Element 天干五行：甲乙木、丙丁火、戊己土、庚辛金、壬癸水。
// 天干依五行順序兩兩成對，故整除 2 即得。
func (s StemIndex) Element() Element { return Element(int(s) / 2) }

// Polarity 天干陰陽：偶數為陽（甲丙戊庚壬），奇數為陰（乙丁己辛癸）。
func (s StemIndex) Polarity() Polarity { return Polarity(int(s) % 2) }

// branchElements 地支五行。四季末月（辰戌丑未）皆屬土，故無規律可循，查表為準。
var branchElements = [BranchCount]Element{
	Water, // 子
	Earth, // 丑
	Wood,  // 寅
	Wood,  // 卯
	Earth, // 辰
	Fire,  // 巳
	Fire,  // 午
	Earth, // 未
	Metal, // 申
	Metal, // 酉
	Earth, // 戌
	Water, // 亥
}

// Element 地支五行
func (b BranchIndex) Element() Element { return branchElements[int(b)%BranchCount] }

// Polarity 地支陰陽：偶數為陽（子寅辰午申戌），奇數為陰（丑卯巳未酉亥）。
func (b BranchIndex) Polarity() Polarity { return Polarity(int(b) % 2) }
