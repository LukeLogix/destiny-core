// Package shensha 神煞偵測。
//
// 本套件刻意與位置語意無關——只認「基準組」與「待掃組」，不知道年月日時、
// 大運流年為何物。呼叫方自行賦予 scan 的位置意義：八字是柱位、六爻是爻位、
// 紫微是宮位。此與 ganzhi.Relation 用 Indices 而不用年月日時的理由相同。
//
// # 表是果，不是因
//
// 神煞多數有生成規律：三合局系是局五行走一圈十二長生的不同取位，日干系是
// 該干的長生位或祿位偏移，天乙貴人是起坤佈干取合氣。本套件寫機制，讓表自己
// 長出來，典籍表只在測試中作為 ground truth 與推導互相比對——與 core 對
// 節氣的作法同構（內嵌 JPL 表 vs VSOP87 自算）。
//
// 此舉不只省去抄表。實測中推導曾指出《三命通會》四庫本〈論天乙貴人〉的
// 「癸以己」應為「癸以巳」——原文形近字出錯時，推導會當場發現。
package shensha

import "github.com/LukeLogix/destiny-core/ganzhi"

// Basis 查法基準的種類。
type Basis uint8

const (
	BasisDayStem     Basis = iota // 以日干查
	BasisYearStem                 // 以年干查
	BasisDayBranch                // 以日支查
	BasisYearBranch               // 以年支查
	BasisMonthBranch              // 以月支查
	BasisSelf                     // 柱本身的干支組合，如魁罡
	BasisChart                    // 全盤性，與 scan 無關
	BasisAnnual                   // 以太歲（流年）查，僅 CategoryAnnual 適用

	basisCount
)

// Variant 同一神煞的雙軌區分。
//
// 數值意義隨 Kind 而定，不具跨 Kind 的共同語意；0 一律表示「無雙軌之分」。
//
//	TianYiGuiRen: VariantYangNoble / VariantYinNoble
//	YangRen:      VariantYangStem  / VariantYinStem
type Variant uint8

const (
	VariantNone      Variant = iota // 無雙軌之分
	VariantYangNoble                // 天乙：陽貴（晝貴）
	VariantYinNoble                 // 天乙：陰貴（夜貴）
	VariantYangStem                 // 羊刃：陽干刃
	VariantYinStem                  // 羊刃：陰干刃
	VariantGou                      // 勾絞：勾
	VariantJiao                     // 勾絞：絞
	VariantDe                       // 德秀：德
	VariantXiu                      // 德秀：秀
	VariantTianLuo                  // 天羅地網：天羅
	VariantDiWang                   // 天羅地網：地網
	VariantTianShang                // 三奇：天上三奇乙丙丁
	VariantDiXia                    // 三奇：地下三奇甲戊庚

	variantCount
)

// Category 收錄分類。
//
// 分野依《三命通會》作者萬民英自述——〈總論諸神煞〉開篇云「命中切要者已備
// 論於前矣」，卷三前段各有專章者為切要，該章所收為次要。以作者的分界為準，
// 勝於實作者主觀認定。
type Category uint8

const (
	CategoryCommon    Category = iota // 子平常用：卷三前段有專章者
	CategorySecondary                 // 次常用
	CategoryAnnual                    // 星命／太歲類，基準為流年
	CategoryRare                      // 冷僻：〈總論諸神煞〉所收

	categoryCount
)

// Tradition 體系來源。
//
// 數個神煞名在不同體系中指涉不同的東西——文昌（紫微系／子平系）、
// 大耗（元辰別名／十二歲君第七位）、官符（亡神別名／太歲類）、
// 白虎（災煞別名／十二歲君第九位）、劍鋒（歲君第一位／〈總論諸神煞〉）。
// 故除 Kind 值不同外，另以本欄位標明體系。
type Tradition uint8

const (
	TraditionZiping Tradition = iota // 子平：三命通會、五行精紀、淵海子平
	TraditionZiwei                   // 紫微斗數
	TraditionSuiJun                  // 十二歲君（星命／太歲系）

	traditionCount
)

// Input 一次偵測所需的基準值。
//
// 不吃 *Chart——照 bazi.StrengthInput 的既有先例：吃 *Chart 則拿到的必然是
// 尚未填入結果的半成品，契約模糊且難以單獨測試。
type Input struct {
	DayStem, YearStem                  ganzhi.StemIndex
	DayBranch, YearBranch, MonthBranch ganzhi.BranchIndex

	// Annual 太歲干支，僅 CategoryAnnual 使用。
	// HasAnnual 為 false 時該分類一律不產生命中——那是正常的「原局層不算
	// 太歲類」路徑，非錯誤。
	Annual    ganzhi.SexagenaryIndex
	HasAnnual bool

	// IsMale 目前僅元辰、勾絞使用，判準為陽男陰女／陰男陽女，
	// 與大運順逆同一條件。
	IsMale bool

	// YearSound 年柱納音，學堂、詞館、天羅地網需要。
	YearSound ganzhi.SoundIndex
}

// Hit 一筆成立的神煞。
type Hit struct {
	Kind    Kind
	At      int // 命中 scan 的第幾個元素；BasisChart 類為 -1
	Basis   Basis
	Variant Variant
}

// yangMale 陽男陰女為 true，陰男陽女為 false。
//
// 與大運順逆同一判準：年干陰陽與性別同號則順、異號則逆。
func (in Input) yangMale() bool {
	return (in.YearStem.Polarity() == ganzhi.Yang) == in.IsMale
}
