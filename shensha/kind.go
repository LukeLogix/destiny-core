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
//
// # 收錄門檻
//
// 每一項都須有典籍出處，取法出自原文而非坊間轉述。查不到出處者一律不收，
// 不憑印象實作——寧可少一個，不可錯一個。尚待查證者見 destiny-project
// 的附錄 A。
type Kind uint8

const (
	// ── 日干系 ──

	TianYiGuiRen   Kind = iota // 天乙貴人：起坤佈干取合氣
	TaiJiGuiRen                // 太極貴人：造化始終相保
	LuShen                     // 祿神：臨官位
	AnLu                       // 暗祿：祿位之六合
	YangRen                    // 羊刃：祿位下一支
	FeiRen                     // 飛刃：羊刃之對衝
	JinYu                      // 金輿：祿位後二支
	WenChangGuiRen             // 文昌貴人（紫微系）：陽干病位／陰干長生位

	// ── 三合局系：局五行走一圈十二長生的不同取位 ──

	YiMa      // 驛馬：病位
	PanAn     // 攀鞍：驛馬後一辰
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
	GeJiao // 隔角：寅申巳亥為角、辰戌丑未為隔

	// ── 月支系 ──

	TianDe   // 天德
	YueDe    // 月德
	TianDeHe // 天德合：天德之五合
	YueDeHe  // 月德合：月德之五合
	DeXiu    // 德秀貴人

	// ── 年支系 ──

	YuanChen   // 元辰：年支之衝再前後一位
	GouJiao    // 勾絞：年支 ±3
	AnJinDeSha // 暗金的煞（一名吟呻、破碎、白衣）

	// ── 納音系 ──

	XueTang       // 學堂：年命納音之長生位，且該柱納音同五行
	CiGuan        // 詞館：年命納音之臨官位，且該柱納音同五行
	ZhengYin      // 正印：五行之正庫
	TianLuoDiWang // 天羅地網：火命看戌亥、水土命看辰巳

	// ── 日干系（食神衍生）──

	TianChuGuiRen // 天廚貴人：食神之祿位
	FuXingGuiRen  // 福星貴人：食神所在的遁時支
	HongYan       // 紅艷煞

	// ── 年支系（鸞喜）──

	HongLuan // 紅鸞：子年起卯逆數
	TianXi   // 天喜：紅鸞之對衝

	// ── 柱本身 ──

	ShiEDaBai // 十惡大敗：日柱祿神落入該旬空亡
	SiFei     // 四廢：月令當旺之五行所剋而無氣的日柱
	KuiGang   // 魁罡：壬辰、庚戌、庚辰、戊戌
	TianShe   // 天赦：春戊寅、夏甲午、秋戊申、冬甲子
	JinShen   // 金神：癸酉、己巳、乙丑

	// ── 日柱系 ──

	XunKong // 空亡（旬空）：日柱所屬旬中缺的兩支

	// ── 全盤 ──

	SanQi // 三奇：四柱天干連續見乙丙丁或甲戊庚

	// ── 四時天喜（與鸞喜的天喜同名異物）──

	TianXiSiShi // 四時天喜：春戌、夏丑、秋辰、冬未

	// ── 十二歲君（星命系，以太歲為基準）──
	//
	// 兩組並存，同一條十二位的軸上各給各的名。神峰通考／命理探源那組有
	// 「以太歲為第一位順數」的起例；洞微經那組只有名單，位置係推定。
	// 由 TopicSuiJun 切換，見決策日誌 D-55。

	TaiSui  // 一：太歲
	TaiYang // 二：太陽
	SangMen // 三：喪門
	TaiYin  // 四：太陰
	GuanFu  // 五：官符
	SiFu    // 六：死符
	SuiPo   // 七：歲破
	LongDe  // 八：龍德
	BaiHu   // 九：白虎
	FuDe    // 十：福德
	DiaoKe  // 十一：弔客
	BingFu  // 十二：病符

	// 洞微經組獨有的五位。太歲、喪門、官符、死符、弔客、病符六位兩組同名
	// 同位，福德兩組皆有而位次不同，故只需補這五個。

	ShengQi      // 洞微經第二位：生氣
	TianYiSuiJun // 洞微經第四位：天醫。非天乙——兩者拼音同形，故加體系後綴
	DaHao        // 洞微經第七位：大耗。與神峰組的歲破同位，係對衝之異名
	FaDao        // 洞微經第八位：發盜
	DaJi         // 洞微經第十位：大吉

	KindCount
)

// meta 各神煞的分類、體系與慣用基準。
//
// 分類依《三命通會》作者自述——〈總論諸神煞〉開篇云「命中切要者已備論於
// 前矣」，卷三前段各有專章者為切要，該章所收為次要。以作者的分界為準，
// 勝於實作者主觀認定。
//
// 慣用基準：有些神煞的基準是有共識的，不是每個都在吵。孤辰寡宿傳統查年支
// （講六親孤剋，年為祖上），桃花驛馬今法查日支。故預設各依慣用，
// StemBase 與 BranchBase 僅在使用者明確指定時覆寫。
var meta = [KindCount]struct {
	category  Category
	tradition Tradition
	basis     Basis
}{
	TianYiGuiRen:   {CategoryCommon, TraditionZiping, BasisDayStem},
	TaiJiGuiRen:    {CategorySecondary, TraditionZiping, BasisDayStem},
	LuShen:         {CategoryCommon, TraditionZiping, BasisDayStem},
	AnLu:           {CategoryRare, TraditionZiping, BasisDayStem},
	YangRen:        {CategoryCommon, TraditionZiping, BasisDayStem},
	FeiRen:         {CategorySecondary, TraditionZiping, BasisDayStem},
	JinYu:          {CategoryCommon, TraditionZiping, BasisDayStem},
	WenChangGuiRen: {CategoryCommon, TraditionZiwei, BasisDayStem},

	YiMa:      {CategoryCommon, TraditionZiping, BasisDayBranch},
	PanAn:     {CategorySecondary, TraditionZiping, BasisDayBranch},
	JiangXing: {CategoryCommon, TraditionZiping, BasisDayBranch},
	HuaGai:    {CategoryCommon, TraditionZiping, BasisDayBranch},
	XianChi:   {CategoryCommon, TraditionZiping, BasisDayBranch},
	JieSha:    {CategoryCommon, TraditionZiping, BasisDayBranch},
	WangShen:  {CategoryCommon, TraditionZiping, BasisDayBranch},
	ZaiSha:    {CategoryCommon, TraditionZiping, BasisDayBranch},
	LiuE:      {CategoryCommon, TraditionZiping, BasisDayBranch},

	GuChen: {CategoryCommon, TraditionZiping, BasisYearBranch},
	GuaSu:  {CategoryCommon, TraditionZiping, BasisYearBranch},
	GeJiao: {CategorySecondary, TraditionZiping, BasisYearBranch},

	TianDe:   {CategoryCommon, TraditionZiping, BasisMonthBranch},
	YueDe:    {CategoryCommon, TraditionZiping, BasisMonthBranch},
	TianDeHe: {CategorySecondary, TraditionZiping, BasisMonthBranch},
	YueDeHe:  {CategorySecondary, TraditionZiping, BasisMonthBranch},
	DeXiu:    {CategorySecondary, TraditionZiping, BasisMonthBranch},

	YuanChen:   {CategoryCommon, TraditionZiping, BasisYearBranch},
	GouJiao:    {CategoryCommon, TraditionZiping, BasisYearBranch},
	AnJinDeSha: {CategorySecondary, TraditionZiping, BasisYearBranch},

	XueTang:       {CategorySecondary, TraditionZiping, BasisSelf},
	CiGuan:        {CategorySecondary, TraditionZiping, BasisSelf},
	ZhengYin:      {CategoryRare, TraditionZiping, BasisSelf},
	TianLuoDiWang: {CategorySecondary, TraditionZiping, BasisSelf},

	ShiEDaBai: {CategorySecondary, TraditionZiping, BasisSelf},
	SiFei:     {CategoryRare, TraditionZiping, BasisSelf},
	KuiGang:   {CategoryCommon, TraditionZiping, BasisSelf},
	TianShe:   {CategorySecondary, TraditionZiping, BasisSelf},
	JinShen:   {CategoryRare, TraditionZiping, BasisSelf},

	TianChuGuiRen: {CategorySecondary, TraditionZiping, BasisDayStem},
	FuXingGuiRen:  {CategorySecondary, TraditionZiping, BasisDayStem},
	HongYan:       {CategorySecondary, TraditionZiping, BasisDayStem},

	HongLuan: {CategorySecondary, TraditionZiping, BasisYearBranch},
	TianXi:   {CategorySecondary, TraditionZiping, BasisYearBranch},

	XunKong: {CategoryCommon, TraditionZiping, BasisDayPillar},

	SanQi: {CategorySecondary, TraditionZiping, BasisChart},

	TaiSui:  {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	TaiYang: {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	SangMen: {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	TaiYin:  {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	GuanFu:  {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	SiFu:    {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	SuiPo:   {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	LongDe:  {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	BaiHu:   {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	FuDe:    {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	DiaoKe:  {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	BingFu:  {CategoryAnnual, TraditionSuiJun, BasisAnnual},

	ShengQi:      {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	TianYiSuiJun: {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	DaHao:        {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	FaDao:        {CategoryAnnual, TraditionSuiJun, BasisAnnual},
	DaJi:         {CategoryAnnual, TraditionSuiJun, BasisAnnual},

	TianXiSiShi: {CategoryRare, TraditionZiping, BasisMonthBranch},
}

// Category 收錄分類
func (k Kind) Category() Category {
	if k >= KindCount {
		return CategoryRare
	}
	return meta[k].category
}

// Tradition 體系來源。
//
// 數個神煞名在不同體系指涉不同的東西——文昌（紫微系／子平系）、大耗
// （元辰別名／十二歲君第七位）、官符（亡神別名／太歲類）、白虎（災煞
// 別名／歲君第九位）。故除 Kind 值不同外，另以本欄位標明體系。
func (k Kind) Tradition() Tradition {
	if k >= KindCount {
		return TraditionZiping
	}
	return meta[k].tradition
}

// SubjectOnly 此神煞是否只就「主體本身」的干支成立。
//
// 十惡大敗、四廢的典籍原文都是「日」——十惡大敗日、四廢日，說的是命主
// 那一柱，不是任一柱碰巧湊出同樣的干支組合。本套件位置無關，無從知道
// 第幾個元素是主體，故只標記，由呼叫方依自己的位置語意過濾：八字取日柱，
// 六爻取世爻。
func (k Kind) SubjectOnly() bool {
	switch k {
	case ShiEDaBai, SiFei, TianShe:
		return true
	}
	return false
}

// BasisIsFixed 此神煞的基準是否只有一種有出處，不受 Options 的基準旋鈕影響。
//
// BranchBase／StemBase 是把「年或日」的選擇一次套到所有神煞，但那個選擇
// 只在確實有兩派的項目上成立。三合局系（華蓋、驛馬、將星⋯）與孤辰寡宿
// 古今各有年支、日支兩種用法，該切；以下四項則典籍只給年支：
//
//	元辰   〈論元辰〉「假如甲子生男與甲午對衝」——甲子是生年
//	勾絞   〈論勾絞〉「假令甲子陽命人」、《五行精紀》「各隨本命求之」
//	紅鸞   《神峰通考》同章起例一律作「◯生人見◯字」（華蓋、將星、喪門皆然）
//	天喜   同上，且其取法就是紅鸞之對衝
//
// 套上日支會生出典籍沒有的命中。實測基準命例即多出兩筆——紅鸞年柱與
// 天喜時柱，皆以日支查得；後者更是該盤唯一的天喜，以年支根本不成立。
//
// 註：本檔的引號一律保留給典籍原文，敘述不用——tools/verify_citations.py
// 會把引號內的文字當引文丟回語料核對。
func (k Kind) BasisIsFixed() bool {
	switch k {
	case YuanChen, GouJiao, HongLuan, TianXi:
		return true
	}
	return false
}

// CustomaryBasis 慣用基準
func (k Kind) CustomaryBasis() Basis {
	if k >= KindCount {
		return BasisDayBranch
	}
	return meta[k].basis
}
