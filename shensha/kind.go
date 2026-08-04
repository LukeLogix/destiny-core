package shensha

// Kind 神煞種類。
//
// 一律用拼音而非意譯：神煞名為專名，硬翻英文會製造假的語意且無法反查原名
// （童子煞譯作 ChildOmen、陰陽差錯譯作 YinYangMismatch，讀者仍不知所指）。
// 與 ganzhi 的英文常數風格不一致，但該不一致是誠實的——七殺有通行英譯故
// 叫 SevenKilling，童子煞沒有。
//
// 編號不預留區段，只承諾尾端追加。分類本身有爭議（華蓋算吉算凶？），
// 區段化會製造第二個要維護的立場。
//
// iota 順序為實作細節，外界不可觀察——序列化一律走 ID()。
type Kind uint8

const (
	// ── 日干系 ──

	TianYiGuiRen   Kind = iota // 天乙貴人：起坤佈干取合氣
	LuShen                     // 祿神：臨官位
	YangRen                    // 羊刃：祿位下一支
	JinYu                      // 金輿：祿位後二支
	WenChangGuiRen             // 文昌貴人（紫微系）：陽干病位／陰干長生位

	// ── 三合局系：局五行走一圈十二長生的不同取位 ──

	YiMa      // 驛馬：病位
	JiangXing // 將星：帝旺位
	HuaGai    // 華蓋：墓位
	XianChi   // 咸池（一名桃花）：沐浴位
	JieSha    // 劫煞：絕位
	WangShen  // 亡神：臨官位
	ZaiSha    // 災煞：胎位
	LiuE      // 六厄：死位

	// ── 三會方系 ──

	GuChen // 孤辰：本方末支的下一支
	GuaSu  // 寡宿：本方首支的上一支

	KindCount
)

// category 各神煞的收錄分類，依《三命通會》作者的「命中切要」分界。
var category = [KindCount]Category{
	TianYiGuiRen:   CategoryCommon,
	LuShen:         CategoryCommon,
	YangRen:        CategoryCommon,
	JinYu:          CategoryCommon,
	WenChangGuiRen: CategoryCommon,
	YiMa:           CategoryCommon,
	JiangXing:      CategoryCommon,
	HuaGai:         CategoryCommon,
	XianChi:        CategoryCommon,
	JieSha:         CategoryCommon,
	WangShen:       CategoryCommon,
	ZaiSha:         CategoryCommon,
	LiuE:           CategoryCommon,
	GuChen:         CategoryCommon,
	GuaSu:          CategoryCommon,
}

// tradition 體系來源。
//
// 文昌貴人取自《紫微斗數》（經《命理探源》轉引），與子平典籍的「文昌貴」
// 是兩個不同的東西——《五行精紀》（宋）與《三命通會》（明）所載的文昌貴
// 與本項僅甲、戊兩干相同。現代八字排盤軟體一律用紫微系，本套件從眾，
// 但以本欄位標明來源，並保留日後收錄子平系的空間。
var tradition = [KindCount]Tradition{
	WenChangGuiRen: TraditionZiwei,
	// 其餘皆為 TraditionZiping（零值）
}

// customaryBasis 各神煞的慣用基準。
//
// 有些神煞的基準是有共識的，不是每個都在吵：孤辰寡宿傳統查年支（講六親
// 孤剋，年為祖上），桃花驛馬今法查日支。故預設各依慣用，BranchBase
// 僅在使用者明確指定時覆寫。
var customaryBasis = [KindCount]Basis{
	TianYiGuiRen:   BasisDayStem,
	LuShen:         BasisDayStem,
	YangRen:        BasisDayStem,
	JinYu:          BasisDayStem,
	WenChangGuiRen: BasisDayStem,

	YiMa:      BasisDayBranch,
	JiangXing: BasisDayBranch,
	HuaGai:    BasisDayBranch,
	XianChi:   BasisDayBranch,
	JieSha:    BasisDayBranch,
	WangShen:  BasisDayBranch,
	ZaiSha:    BasisDayBranch,
	LiuE:      BasisDayBranch,

	GuChen: BasisYearBranch,
	GuaSu:  BasisYearBranch,
}

// Category 收錄分類
func (k Kind) Category() Category {
	if k >= KindCount {
		return CategoryRare
	}
	return category[k]
}

// Tradition 體系來源
func (k Kind) Tradition() Tradition {
	if k >= KindCount {
		return TraditionZiping
	}
	return tradition[k]
}

// CustomaryBasis 慣用基準
func (k Kind) CustomaryBasis() Basis {
	if k >= KindCount {
		return BasisDayBranch
	}
	return customaryBasis[k]
}
