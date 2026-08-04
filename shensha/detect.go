package shensha

import (
	"sort"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// BranchBaseSect 以地支查的那組，基準取自年支或日支——最大宗的分歧，
// 且一次影響六七個神煞，故抽為共通旋鈕。
type BranchBaseSect uint8

const (
	// BranchBaseCustomary 各神煞沿用其慣用基準（預設）
	BranchBaseCustomary BranchBaseSect = iota
	BranchBaseDay                      // 一律以日支為基準
	BranchBaseYear                     // 一律以年支為基準
	BranchBaseBoth                     // 年支、日支皆查，以 Hit.Basis 區分
)

// TianYiSect 天乙貴人的口訣版本
type TianYiSect uint8

const (
	// TianYiSanMing 《三命通會》本宗：甲戊庚牛羊，六辛逢馬虎。
	// 有完整推導支撐，故為預設。
	TianYiSanMing TianYiSect = iota
	// TianYiYeHuiTing 清代葉悔亭《六壬眎斯》改本：庚辛逢馬虎。
	// 係人為調整——葉氏認為前人配置不夠平均而把庚自「甲戊庚牛羊」抽走，
	// 無推導支撐。
	TianYiYeHuiTing
)

// YinStemBladeSect 陰干有無羊刃。
//
// 位置不在爭議之列——羊刃為祿位下一支，十干皆順行，與 TerrainSect 無關。
// 本項僅管「算不算」。
type YinStemBladeSect uint8

const (
	// YinStemBladeMarked 陰干亦計刃，另標 Variant 供上層過濾（預設）
	YinStemBladeMarked YinStemBladeSect = iota
	// YinStemBladeNone 陰干無刃，市面實作多數如此
	YinStemBladeNone
)

// Options 神煞的計算口徑。
type Options struct {
	// Categories 要計算哪些分類。nil 表示全部。
	Categories []Category
	Include    []Kind // 在 Categories 之外額外加入
	Exclude    []Kind // 自結果中排除

	BranchBase   BranchBaseSect
	TianYi       TianYiSect
	YinStemBlade YinStemBladeSect
}

// Default 預設口徑：全部分類、各依慣用基準、天乙採三命通會本宗、陰干計刃。
func Default() Options { return Options{} }

// enabled 判斷某神煞是否納入計算。
func (o Options) enabled(k Kind) bool {
	for _, x := range o.Exclude {
		if x == k {
			return false
		}
	}
	for _, x := range o.Include {
		if x == k {
			return true
		}
	}
	if o.Categories == nil {
		return true
	}
	for _, c := range o.Categories {
		if c == k.Category() {
			return true
		}
	}
	return false
}

// target 一個應命中的地支及其雙軌標記。
type target struct {
	branch  ganzhi.BranchIndex
	variant Variant
}

// Detect 位置無關的神煞偵測。呼叫方自行賦予 scan 的位置語意。
//
// 回傳順序為確定性：依 (Kind, At, Basis, Variant) 之字典序排序。
// 順序不確定會使 golden test 不穩，故列為契約的一部分。
func Detect(in Input, scan []ganzhi.SexagenaryIndex, opt Options) []Hit {
	hits := make([]Hit, 0, len(scan))

	var bases [maxBases]Basis
	var targets [maxTargets]target

	for k := Kind(0); k < KindCount; k++ {
		if !opt.enabled(k) {
			continue
		}
		nb := basesFor(k, opt, &bases)
		for _, b := range bases[:nb] {
			nt := targetsFor(k, b, in, opt, &targets)
			for _, t := range targets[:nt] {
				for i, sex := range scan {
					if sex.Branch() == t.branch {
						hits = append(hits, Hit{Kind: k, At: i, Basis: b, Variant: t.variant})
					}
				}
			}
		}
	}

	sort.Slice(hits, func(x, y int) bool {
		a, b := hits[x], hits[y]
		switch {
		case a.Kind != b.Kind:
			return a.Kind < b.Kind
		case a.At != b.At:
			return a.At < b.At
		case a.Basis != b.Basis:
			return a.Basis < b.Basis
		default:
			return a.Variant < b.Variant
		}
	})
	return hits
}

// 上限：基準最多兩個（BranchBaseBoth），應命中的地支最多兩個（天乙的陽貴陰貴）。
//
// 寫入呼叫端提供的緩衝而非回傳切片——Detect 在開啟大運流年時會被呼叫上百次，
// 每次配置兩個小切片即為數千次無謂的配置。
const (
	maxBases   = 2
	maxTargets = 2
)

// basesFor 該神煞此次要用哪些基準，回傳寫入的個數。
//
// 只有以地支查者受 BranchBase 影響；以日干查者不在此列。
func basesFor(k Kind, opt Options, out *[maxBases]Basis) int {
	c := k.CustomaryBasis()
	if c != BasisDayBranch && c != BasisYearBranch {
		out[0] = c
		return 1
	}
	switch opt.BranchBase {
	case BranchBaseDay:
		out[0] = BasisDayBranch
	case BranchBaseYear:
		out[0] = BasisYearBranch
	case BranchBaseBoth:
		out[0], out[1] = BasisDayBranch, BasisYearBranch
		return 2
	default:
		out[0] = c
	}
	return 1
}

// baseBranch 取指定基準的地支值
func (in Input) baseBranch(b Basis) ganzhi.BranchIndex {
	switch b {
	case BasisYearBranch:
		return in.YearBranch
	case BasisMonthBranch:
		return in.MonthBranch
	default:
		return in.DayBranch
	}
}

// baseStem 取指定基準的天干值
func (in Input) baseStem(b Basis) ganzhi.StemIndex {
	if b == BasisYearStem {
		return in.YearStem
	}
	return in.DayStem
}

// targetsFor 算出該神煞在此基準下應命中的地支，回傳寫入的個數。
func targetsFor(k Kind, b Basis, in Input, opt Options, out *[maxTargets]target) int {
	one := func(br ganzhi.BranchIndex) int {
		out[0] = target{branch: br}
		return 1
	}

	switch k {
	// ── 三合局系：同一個東西的八個取位 ──
	case JiangXing:
		return one(trinityAt(in.baseBranch(b), ganzhi.Prosperity))
	case HuaGai:
		return one(trinityAt(in.baseBranch(b), ganzhi.Tomb))
	case YiMa:
		return one(trinityAt(in.baseBranch(b), ganzhi.Sick))
	case XianChi:
		return one(trinityAt(in.baseBranch(b), ganzhi.Bath))
	case JieSha:
		return one(trinityAt(in.baseBranch(b), ganzhi.Void))
	case WangShen:
		return one(trinityAt(in.baseBranch(b), ganzhi.Officer))
	case ZaiSha:
		return one(trinityAt(in.baseBranch(b), ganzhi.Fetus))
	case LiuE:
		return one(trinityAt(in.baseBranch(b), ganzhi.Death))

	// ── 三會方系 ──
	case GuChen:
		return one((directionLast(in.baseBranch(b)) + 1) % ganzhi.BranchCount)
	case GuaSu:
		return one((directionFirst(in.baseBranch(b)) + ganzhi.BranchCount - 1) % ganzhi.BranchCount)

	// ── 日干系 ──
	case LuShen:
		return one(luBranch[int(in.baseStem(b))%ganzhi.StemCount])
	case JinYu:
		return one(offsetFromLu(in.baseStem(b), 2))
	case WenChangGuiRen:
		return one(wenChangBranch[int(in.baseStem(b))%ganzhi.StemCount])

	case YangRen:
		stem := in.baseStem(b)
		v := VariantYangStem
		if stem.Polarity() == ganzhi.Yin {
			if opt.YinStemBlade == YinStemBladeNone {
				return 0
			}
			v = VariantYinStem
		}
		out[0] = target{branch: offsetFromLu(stem, 1), variant: v}
		return 1

	case TianYiGuiRen:
		stem := in.baseStem(b)
		yang, yin := yangNoble[stem], yinNoble[stem]
		if opt.TianYi == TianYiYeHuiTing {
			// 葉悔亭改本把庚自「甲戊庚牛羊」抽走，改與辛同歸馬虎
			if stem == 6 {
				yang, yin = yangNoble[7], yinNoble[7]
			}
		}
		out[0] = target{branch: yang, variant: VariantYangNoble}
		out[1] = target{branch: yin, variant: VariantYinNoble}
		return 2
	}
	return 0
}
