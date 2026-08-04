package shensha

import (
	"sort"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// target 一個應命中的地支及其雙軌標記。
// target 一個應命中的干或支。
//
// 神煞多數命中地支，但天德、月德、德秀等落在天干，故兩者皆須支援。
type target struct {
	branch    ganzhi.BranchIndex
	stem      ganzhi.StemIndex
	matchStem bool
	variant   Variant
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

		switch k.CustomaryBasis() {
		case BasisSelf:
			// 柱本身：不查基準，逐一檢視各柱自己的干支組合
			for i, sex := range scan {
				if v, ok := selfHit(k, sex, in, opt); ok {
					hits = append(hits, Hit{Kind: k, At: i, Basis: BasisSelf, Variant: v})
				}
			}
			continue

		case BasisChart:
			// 全盤性：與 scan 的個別元素無關，At 為 -1
			if v, ok := chartHit(k, scan); ok {
				hits = append(hits, Hit{Kind: k, At: -1, Basis: BasisChart, Variant: v})
			}
			continue

		case BasisAnnual:
			// 太歲類離開流年即無意義
			if !in.HasAnnual {
				continue
			}
		}

		nb := basesFor(k, opt, &bases)
		for _, b := range bases[:nb] {
			nt := targetsFor(k, b, in, opt, &targets)
			for _, t := range targets[:nt] {
				for i, sex := range scan {
					hit := sex.Branch() == t.branch
					if t.matchStem {
						hit = sex.Stem() == t.stem
					}
					if hit {
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
	maxBases = 2
	// maxTargets 應命中者的上限。德秀在申子辰月有八個天干，為最寬的一項。
	maxTargets = 8
)

// basesFor 該神煞此次要用哪些基準，回傳寫入的個數。
//
// 只有以地支查者受 BranchBase 影響；以日干查者不在此列。
func basesFor(k Kind, opt Options, out *[maxBases]Basis) int {
	c := k.CustomaryBasis()

	// 只有一種基準有出處者不跟著旋鈕走——切換的前提是真的有兩派
	if k.BasisIsFixed() {
		out[0] = c
		return 1
	}

	if c == BasisDayStem || c == BasisYearStem {
		switch opt.sect(TopicStemBase) {
		case StemBaseDay:
			out[0] = BasisDayStem
		case StemBaseYear:
			out[0] = BasisYearStem
		case StemBaseBoth:
			out[0], out[1] = BasisDayStem, BasisYearStem
			return 2
		default:
			out[0] = c
		}
		return 1
	}

	// 以整柱查者（空亡）與以地支查者受同一個旋鈕影響——選的是年／日這條軸，
	// 不是「干或支」。
	if c == BasisDayPillar || c == BasisYearPillar {
		switch opt.sect(TopicBranchBase) {
		case BranchBaseDay:
			out[0] = BasisDayPillar
		case BranchBaseYear:
			out[0] = BasisYearPillar
		case BranchBaseBoth:
			out[0], out[1] = BasisDayPillar, BasisYearPillar
			return 2
		default:
			out[0] = c
		}
		return 1
	}

	if c != BasisDayBranch && c != BasisYearBranch {
		out[0] = c
		return 1
	}
	switch opt.sect(TopicBranchBase) {
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
	case BasisAnnual:
		return in.Annual.Branch()
	default:
		return in.DayBranch
	}
}

// basePillar 取指定基準的整組干支。空亡須知該柱屬哪一旬，拆開的干與支
// 各自都不足以決定旬首。
func (in Input) basePillar(b Basis) ganzhi.SexagenaryIndex {
	if b == BasisYearPillar {
		return sexFrom(in.YearStem, in.YearBranch)
	}
	return sexFrom(in.DayStem, in.DayBranch)
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

	case YangRen, FeiRen:
		stem := in.baseStem(b)
		v := VariantYangStem
		if stem.Polarity() == ganzhi.Yin {
			if opt.sect(TopicYinBlade) == YinBladeNone {
				return 0
			}
			v = VariantYinStem
		}
		br := bladeAt(stem, opt.sect(TopicBladeAt))
		if k == FeiRen {
			br = clash(br) // 飛刃即羊刃之衝，故隨羊刃的口徑連動
		}
		out[0] = target{branch: br, variant: v}
		return 1

	case TaiJiGuiRen:
		n := 0
		for _, br := range taiJiTargets[int(in.baseStem(b))%ganzhi.StemCount] {
			out[n] = target{branch: br}
			n++
		}
		return n

	case TianChuGuiRen:
		return one(tianChuBranch[int(in.baseStem(b))%ganzhi.StemCount])

	case HongYan:
		return one(hongYanBranch[int(in.baseStem(b))%ganzhi.StemCount])

	case FuXingGuiRen:
		tbl := &fuXingBranches
		if opt.sect(TopicFuXing) == FuXingShenFeng {
			tbl = &fuXingShenFeng
		}
		n := 0
		for _, br := range tbl[int(in.baseStem(b))%ganzhi.StemCount] {
			out[n] = target{branch: br}
			n++
		}
		return n

	case HongLuan:
		return one(hongLuan(in.baseBranch(b)))
	case TianXi:
		return one(clash(hongLuan(in.baseBranch(b))))

	case XunKong:
		// 〈六甲空亡〉：「甲子旬中無戌亥，甲戌旬中無申酉⋯」——該旬缺的兩支
		a, c := xunKong(in.basePillar(b))
		out[0] = target{branch: a}
		out[1] = target{branch: c}
		return 2

	case AnLu:
		return one(liuHe(luBranch[int(in.baseStem(b))%ganzhi.StemCount]))
	case PanAn:
		return one((trinityAt(in.baseBranch(b), ganzhi.Sick) + 1) % ganzhi.BranchCount)

	case GeJiao:
		// 〈論孤辰寡宿〉引玉門集：「寅申巳亥為角，辰戌丑未為隔」——
		// 命前一辰為角、命後一辰為隔，取其落在四孟或四季者
		base := in.baseBranch(b)
		n := 0
		for _, d := range [2]int{1, -1} {
			t := wrap(int(base) + d)
			if isJiaoOrGe(t) {
				out[n] = target{branch: t}
				n++
			}
		}
		return n

	case YuanChen:
		// 〈論元辰〉：陽男陰女在衝前一位，陰男陽女在衝後一位
		d := 1
		if !in.yangMale() {
			d = -1
		}
		c := clash(in.baseBranch(b))
		return one(ganzhi.BranchIndex(((int(c)+d)%ganzhi.BranchCount + ganzhi.BranchCount) % ganzhi.BranchCount))

	case GouJiao:
		// 〈論勾絞〉：陽男陰女命前三辰為勾、命後三辰為絞，陰男陽女反之
		base := int(in.baseBranch(b))
		fwd := ganzhi.BranchIndex((base + 3) % ganzhi.BranchCount)
		bwd := ganzhi.BranchIndex((base + 9) % ganzhi.BranchCount)
		front, back := VariantGou, VariantJiao
		if opt.sect(TopicGouJiao) == GouJiaoFrontIsJiao {
			front, back = VariantJiao, VariantGou
		}
		out[0] = target{branch: fwd, variant: front}
		out[1] = target{branch: bwd, variant: back}
		if !in.yangMale() {
			out[0].variant, out[1].variant = out[1].variant, out[0].variant
		}
		return 2

	case AnJinDeSha:
		return one(anJinTarget(in.baseBranch(b)))

	case TianDe:
		v := tianDeByMonth[int(in.MonthBranch)%ganzhi.BranchCount]
		if v >= 100 {
			out[0] = target{branch: ganzhi.BranchIndex(v - 100)}
		} else {
			out[0] = target{stem: ganzhi.StemIndex(v), matchStem: true}
		}
		return 1

	case TianDeHe:
		// 天德之五合。天德落於地支的四個月無合可取。
		v := tianDeByMonth[int(in.MonthBranch)%ganzhi.BranchCount]
		if v >= 100 {
			return 0
		}
		out[0] = target{stem: hePartner(ganzhi.StemIndex(v)), matchStem: true}
		return 1

	case YueDe:
		out[0] = target{stem: yueDeStem(in.MonthBranch), matchStem: true}
		return 1
	case YueDeHe:
		out[0] = target{stem: hePartner(yueDeStem(in.MonthBranch)), matchStem: true}
		return 1

	case DeXiu:
		// 〈論徳秀〉依局分德、秀兩組天干，各以 Variant 標明
		g := deXiuStems[int(in.MonthBranch)%4]
		half := len(g) / 2
		for i, x := range g {
			v := VariantDe
			if i >= half {
				v = VariantXiu
			}
			out[i] = target{stem: x, matchStem: true, variant: v}
		}
		return len(g)

	case TaiSui, TaiYang, SangMen, TaiYin, GuanFu, SiFu,
		SuiPo, LongDe, BaiHu, FuDe, DiaoKe, BingFu:
		return one(wrap(int(in.Annual.Branch()) + suiJunOffset[k]))

	case TianYiGuiRen:
		stem := in.baseStem(b)
		yang, yin := yangNoble[stem], yinNoble[stem]
		if opt.sect(TopicTianYi) == TianYiYeHuiTing {
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

// wrap 地支加減後回到 0–11
func wrap(n int) ganzhi.BranchIndex {
	return ganzhi.BranchIndex((n%ganzhi.BranchCount + ganzhi.BranchCount) % ganzhi.BranchCount)
}

// isJiaoOrGe 寅申巳亥為角、辰戌丑未為隔
func isJiaoOrGe(b ganzhi.BranchIndex) bool {
	switch b {
	case 2, 5, 8, 11, 1, 4, 7, 10:
		return true
	}
	return false
}

// selfHit 柱本身的判定——不查基準，只看該柱自己的干支（部分另需年命納音）。
//
// 十惡大敗、四廢只論主體那一柱。本套件不知道哪個元素是主體，照樣報構成，
// 並以 Kind.SubjectOnly 標記，由呼叫方過濾。
func selfHit(k Kind, sex ganzhi.SexagenaryIndex, in Input, opt Options) (Variant, bool) {
	switch k {
	case XueTang, CiGuan, ZhengYin:
		// 〈論學堂詞館〉：「長生乃學堂之正位，如金命見辛巳，金長生在巳，
		// 辛巳納音又屬金是也。臨官乃詞館正位」——納音與地支的雙重條件。
		// 〈論正印〉的「五行之正庫」同一形式，取墓位。
		elem := in.YearSound.Element()
		if ganzhi.SoundOf(sex).Element() != elem {
			return VariantNone, false
		}
		n := terrainLongLife
		switch k {
		case CiGuan:
			n = terrainOfficer
		case ZhengYin:
			n = terrainTomb
		}
		return VariantNone, sex.Branch() == soundTerrainAt(elem, n)

	case TianLuoDiWang:
		// 戌亥為天羅、辰巳為地網，兩派皆同；分歧在要不要納音條件。
		// 〈論天羅地網〉與《命理探源》引《淵海子平》：「火命人有天羅，
		// 水土命人有地網，餘金木二命無之」；《五行精紀·天羅地網歌》
		// 「凡以戌亥為天羅，辰巳為地網」，通篇未及納音。
		luo := sex.Branch() == 10 || sex.Branch() == 11
		wang := sex.Branch() == 4 || sex.Branch() == 5
		if opt.sect(TopicTianLuo) == TianLuoBranchOnly {
			switch {
			case luo:
				return VariantTianLuo, true
			case wang:
				return VariantDiWang, true
			}
			return VariantNone, false
		}
		switch in.YearSound.Element() {
		case ganzhi.Fire:
			if luo {
				return VariantTianLuo, true
			}
		case ganzhi.Water, ganzhi.Earth:
			if wang {
				return VariantDiWang, true
			}
		}
		return VariantNone, false

	case ShiEDaBai:
		// 〈論十惡大敗〉：「六甲旬中十個日値祿入空亡」——
		// 與祿神、空亡共用推導，不硬編那十個日柱。
		lu := luBranch[int(sex.Stem())%ganzhi.StemCount]
		a, b := xunKong(sex)
		return VariantNone, lu == a || lu == b

	case KuiGang:
		for _, x := range kuiGangSex {
			if sex == x {
				return VariantNone, true
			}
		}
		return VariantNone, false

	case TianShe:
		return VariantNone, sex == tianSheSex[directionGroup(in.MonthBranch)]

	case JinShen:
		for _, x := range jinShenSex {
			if sex == x {
				return VariantNone, true
			}
		}
		return VariantNone, false

	case SiFei:
		// 該季無氣之五行，且干支同屬該五行——恰好是原文那八組
		e := siFeiElement[directionGroup(in.MonthBranch)]
		return VariantNone, sex.Stem().Element() == e && sex.Branch().Element() == e
	}
	return VariantNone, false
}

// chartHit 全盤性的判定。
func chartHit(k Kind, scan []ganzhi.SexagenaryIndex) (Variant, bool) {
	if k != SanQi {
		return VariantNone, false
	}
	// 〈論三奇〉：須「三干相連而無間」——天干中須連續出現，不得跳柱
	for si, set := range sanQiSets {
		for i := 0; i+2 < len(scan); i++ {
			if scan[i].Stem() == set[0] && scan[i+1].Stem() == set[1] && scan[i+2].Stem() == set[2] {
				if si == 0 {
					return VariantTianShang, true
				}
				return VariantDiXia, true
			}
		}
	}
	return VariantNone, false
}
