package shensha

// 列舉的穩定識別字。規則與作法見 ganzhi/id.go 與 lang/id_test.go。
//
// 神煞用拼音、十神用英文（seven_killing），表面不一致但規則一致——
// ID 一律取 Go 常數名，而各領域的常數命名依其既有慣例。

var kindIDs = [KindCount]string{
	TianYiGuiRen:   "tian_yi_gui_ren",
	TaiJiGuiRen:    "tai_ji_gui_ren",
	LuShen:         "lu_shen",
	AnLu:           "an_lu",
	YangRen:        "yang_ren",
	FeiRen:         "fei_ren",
	JinYu:          "jin_yu",
	WenChangGuiRen: "wen_chang_gui_ren",
	YiMa:           "yi_ma",
	PanAn:          "pan_an",
	JiangXing:      "jiang_xing",
	HuaGai:         "hua_gai",
	XianChi:        "xian_chi",
	JieSha:         "jie_sha",
	WangShen:       "wang_shen",
	ZaiSha:         "zai_sha",
	LiuE:           "liu_e",
	GuChen:         "gu_chen",
	GuaSu:          "gua_su",
	GeJiao:         "ge_jiao",
	TianDe:         "tian_de",
	YueDe:          "yue_de",
	TianDeHe:       "tian_de_he",
	YueDeHe:        "yue_de_he",
	DeXiu:          "de_xiu",
	YuanChen:       "yuan_chen",
	GouJiao:        "gou_jiao",
	AnJinDeSha:     "an_jin_de_sha",
	XueTang:        "xue_tang",
	CiGuan:         "ci_guan",
	ZhengYin:       "zheng_yin",
	TianLuoDiWang:  "tian_luo_di_wang",
	TianChuGuiRen:  "tian_chu_gui_ren",
	FuXingGuiRen:   "fu_xing_gui_ren",
	HongYan:        "hong_yan",
	HongLuan:       "hong_luan",
	TianXi:         "tian_xi",
	ShiEDaBai:      "shi_eda_bai",
	SiFei:          "si_fei",
	KuiGang:        "kui_gang",
	TianShe:        "tian_she",
	JinShen:        "jin_shen",
	XunKong:        "xun_kong",
	SanQi:          "san_qi",
	TianXiSiShi:    "tian_xi_si_shi",
	TaiSui:         "tai_sui",
	TaiYang:        "tai_yang",
	SangMen:        "sang_men",
	TaiYin:         "tai_yin",
	GuanFu:         "guan_fu",
	SiFu:           "si_fu",
	SuiPo:          "sui_po",
	LongDe:         "long_de",
	BaiHu:          "bai_hu",
	FuDe:           "fu_de",
	DiaoKe:         "diao_ke",
	BingFu:         "bing_fu",
	ShengQi:        "sheng_qi",
	TianYiSuiJun:   "tian_yi_sui_jun",
	DaHao:          "da_hao",
	FaDao:          "fa_dao",
	DaJi:           "da_ji",
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
	BasisDayPillar:   "day_pillar",
	BasisYearPillar:  "year_pillar",
}

func (b Basis) ID() string { return pickID(basisIDs[:], int(b)) }

var variantIDs = [variantCount]string{
	VariantNone:      "none",
	VariantYangNoble: "yang_noble",
	VariantYinNoble:  "yin_noble",
	VariantYangStem:  "yang_stem",
	VariantYinStem:   "yin_stem",
	VariantGou:       "gou",
	VariantJiao:      "jiao",
	VariantDe:        "de",
	VariantXiu:       "xiu",
	VariantTianLuo:   "tian_luo",
	VariantDiWang:    "di_wang",
	VariantTianShang: "tian_shang",
	VariantDiXia:     "di_xia",
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

var topicIDs = [topicCount]string{
	TopicBranchBase: "branch_base",
	TopicStemBase:   "stem_base",
	TopicTianYi:     "tian_yi",
	TopicBladeAt:    "blade_at",
	TopicYinBlade:   "yin_blade",
	TopicFuXing:     "fu_xing",
	TopicTianLuo:    "tian_luo",
	TopicGouJiao:    "gou_jiao",
	TopicSuiJun:     "sui_jun",
	TopicJinShen:    "jin_shen",
	TopicKuiGang:    "kui_gang",
}

func (x Topic) ID() string { return pickID(topicIDs[:], int(x)) }

var sectIDs = [sectCount]string{
	SectDefault:         "sect_default",
	BranchBaseCustomary: "branch_base_customary",
	BranchBaseDay:       "branch_base_day",
	BranchBaseYear:      "branch_base_year",
	BranchBaseBoth:      "branch_base_both",
	StemBaseCustomary:   "stem_base_customary",
	StemBaseDay:         "stem_base_day",
	StemBaseYear:        "stem_base_year",
	StemBaseBoth:        "stem_base_both",
	TianYiSanMing:       "tian_yi_san_ming",
	TianYiYeHuiTing:     "tian_yi_ye_hui_ting",
	BladeAtLuNext:       "blade_at_lu_next",
	BladeAtProsperity:   "blade_at_prosperity",
	YinBladeMarked:      "yin_blade_marked",
	YinBladeNone:        "yin_blade_none",
	FuXingHourStem:      "fu_xing_hour_stem",
	FuXingShenFeng:      "fu_xing_shen_feng",
	TianLuoBySound:      "tian_luo_by_sound",
	TianLuoBranchOnly:   "tian_luo_branch_only",
	GouJiaoFrontIsGou:   "gou_jiao_front_is_gou",
	GouJiaoFrontIsJiao:  "gou_jiao_front_is_jiao",
	SuiJunShenFeng:      "sui_jun_shen_feng",
	SuiJunDongWei:       "sui_jun_dong_wei",
	JinShenAnyDay:       "jin_shen_any_day",
	JinShenJiaDayOnly:   "jin_shen_jia_day_only",
	KuiGangAllPillars:   "kui_gang_all_pillars",
	KuiGangDayOnly:      "kui_gang_day_only",
}

func (x Sect) ID() string { return pickID(sectIDs[:], int(x)) }

func pickID(tbl []string, i int) string {
	if i < 0 || i >= len(tbl) {
		return ""
	}
	return tbl[i]
}
