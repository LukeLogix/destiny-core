package shensha

// 列舉的穩定識別字。規則與作法見 ganzhi/id.go 與 lang/id_test.go。
//
// 神煞用拼音、十神用英文（seven_killing），表面不一致但規則一致——
// ID 一律取 Go 常數名，而各領域的常數命名依其既有慣例。

var kindIDs = [KindCount]string{
	TianYiGuiRen:   "tian_yi_gui_ren",
	LuShen:         "lu_shen",
	YangRen:        "yang_ren",
	JinYu:          "jin_yu",
	WenChangGuiRen: "wen_chang_gui_ren",
	YiMa:           "yi_ma",
	JiangXing:      "jiang_xing",
	HuaGai:         "hua_gai",
	XianChi:        "xian_chi",
	JieSha:         "jie_sha",
	WangShen:       "wang_shen",
	ZaiSha:         "zai_sha",
	LiuE:           "liu_e",
	GuChen:         "gu_chen",
	GuaSu:          "gua_su",
}

func (k Kind) ID() string { return pickID(kindIDs[:], int(k)) }

var basisIDs = [basisCount]string{
	BasisDayStem:     "day_stem",
	BasisYearStem:    "year_stem",
	BasisDayBranch:   "day_branch",
	BasisYearBranch:  "year_branch",
	BasisMonthBranch: "month_branch",
	BasisSelf:        "self",
	BasisChart:       "chart",
	BasisAnnual:      "annual",
}

func (b Basis) ID() string { return pickID(basisIDs[:], int(b)) }

var variantIDs = [variantCount]string{
	VariantNone:      "none",
	VariantYangNoble: "yang_noble",
	VariantYinNoble:  "yin_noble",
	VariantYangStem:  "yang_stem",
	VariantYinStem:   "yin_stem",
}

func (v Variant) ID() string { return pickID(variantIDs[:], int(v)) }

var categoryIDs = [categoryCount]string{
	CategoryCommon:    "common",
	CategorySecondary: "secondary",
	CategoryAnnual:    "annual",
	CategoryRare:      "rare",
}

func (c Category) ID() string { return pickID(categoryIDs[:], int(c)) }

var traditionIDs = [traditionCount]string{
	TraditionZiping: "ziping",
	TraditionZiwei:  "ziwei",
	TraditionSuiJun: "sui_jun",
}

func (t Tradition) ID() string { return pickID(traditionIDs[:], int(t)) }

var branchBaseSectIDs = [4]string{
	BranchBaseCustomary: "customary",
	BranchBaseDay:       "day",
	BranchBaseYear:      "year",
	BranchBaseBoth:      "both",
}

func (s BranchBaseSect) ID() string { return pickID(branchBaseSectIDs[:], int(s)) }

var tianYiSectIDs = [2]string{
	TianYiSanMing:   "san_ming",
	TianYiYeHuiTing: "ye_hui_ting",
}

func (s TianYiSect) ID() string { return pickID(tianYiSectIDs[:], int(s)) }

var yinStemBladeSectIDs = [2]string{
	YinStemBladeMarked: "marked",
	YinStemBladeNone:   "none",
}

func (s YinStemBladeSect) ID() string { return pickID(yinStemBladeSectIDs[:], int(s)) }

func pickID(tbl []string, i int) string {
	if i < 0 || i >= len(tbl) {
		return ""
	}
	return tbl[i]
}
