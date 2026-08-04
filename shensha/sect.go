package shensha

// 流派分歧的統一表述。
//
// # 為什麼以「主題」為鍵而不是以 Kind
//
// 直覺會想用 map[Kind]Sect，但分歧的作用範圍不等於單一神煞：羊刃的位置之爭
// 同時決定飛刃（飛刃即羊刃之衝），陰干計不計刃亦然；基準之爭一次影響七八個
// Kind。以 Kind 為鍵在第二個例子就破功。
//
// # 一張表出兩份東西
//
// topicSects 既是計算時的取法來源，也是 API 的中繼資料來源——「可以選什麼」
// 與「有哪些流派」是同一件事。先前後端與前端各手寫一份平行清單，那份清單
// 沒有測試守著，加一個流派要動六個地方。
//
// # 收錄門檻
//
// 一個主題要成立，各派都須有出處，且各自在自己的體系內自洽。同一本書自相
// 矛盾者（官符「取太歲前五辰」撞死符、四庫本「癸以己」形近之訛）不列為流派——
// 那是校勘問題，該判定並記錄理由，列成選項等於把校勘丟給使用者。
type Topic uint8

const (
	// TopicBranchBase 以地支查者的基準。最大宗的分歧，一次影響六七個神煞。
	TopicBranchBase Topic = iota
	// TopicStemBase 以天干查者的基準
	TopicStemBase
	// TopicTianYi 天乙貴人的口訣版本
	TopicTianYi
	// TopicBladeAt 羊刃的位置。連帶決定飛刃——飛刃即羊刃之衝。
	TopicBladeAt
	// TopicYinBlade 陰干算不算刃。與位置正交，故分成兩個主題；同樣連帶飛刃。
	TopicYinBlade
	// TopicFuXing 福星貴人的取法
	TopicFuXing
	// TopicTianLuo 天羅地網要不要納音條件
	TopicTianLuo
	// TopicGouJiao 勾與絞的方向。只影響 Variant 標記，不影響是否命中。
	TopicGouJiao

	topicCount
)

// Sect 某主題下的一種取法。
//
// 常數全域唯一而非各主題自行編號——如此 lang 與 lang/id_test.go 的既有機制
// 可直接沿用，不必為每個主題各開一張表。
//
// SectDefault 佔 0 位是必要的：若讓各主題的首選都是 0，Options 的零值會變成
// 「第一個主題的第一個選項」被套到所有主題上。
type Sect uint8

const (
	// SectDefault 未指定，採該主題的首選
	SectDefault Sect = iota

	// ── 地支基準 ──

	BranchBaseCustomary // 各依慣用（孤寡查年支、桃花驛馬查日支）
	BranchBaseDay       // 一律以日支
	BranchBaseYear      // 一律以年支
	BranchBaseBoth      // 年支日支皆查，以 Hit.Basis 區分

	// ── 天干基準 ──

	StemBaseCustomary // 各依慣用，即日干
	StemBaseDay       // 一律以日干
	StemBaseYear      // 一律以年干
	StemBaseBoth      // 年干日干皆查

	// ── 天乙貴人 ──

	// TianYiSanMing 《三命通會》本宗：甲戊庚牛羊，六辛逢馬虎。
	// 有完整推導支撐（起坤佈干取合氣），故為預設。
	TianYiSanMing
	// TianYiYeHuiTing 清代葉悔亭《六壬眎斯》改本：庚辛逢馬虎。
	// 係人為調整——葉氏認為前人配置不夠平均而把庚自「甲戊庚牛羊」抽走，
	// 無推導支撐。
	TianYiYeHuiTing

	// ── 羊刃位置 ──

	// BladeAtLuNext 祿前一辰。《三命通會·論羊刃》「羊刃常居祿前一辰」，
	// 並自列「辰者乙之正位，未者丁之正位，戌者辛之正位，丑者癸之正位」——
	// 陰干四例逐一吻合，故為預設。
	BladeAtLuNext
	// BladeAtProsperity 帝旺位。通行於今，陽干與祿前一辰同位，
	// 但陰干五個全異：乙寅、丁巳、辛申、癸亥，與上引原文不合。
	BladeAtProsperity

	// ── 陰干羊刃 ──

	// YinBladeMarked 陰干亦計刃，另標 Variant 供上層過濾
	YinBladeMarked
	// YinBladeNone 陰干無刃，市面實作多數如此
	YinBladeNone

	// ── 福星貴人 ──

	// FuXingHourStem 遁時得食神之支。《協紀辨方書》「日干生時干⋯皆本日
	// 日干之食神子孫」並逐日列出，《三命通會》「遁得本旬中真食神」同機制。
	// 有自述取法可推導，故為預設。
	FuXingHourStem
	// FuXingShenFeng 《神峰通考》歌訣「甲丙相邀入虎鄉⋯戊申己未丁亥遇，
	// 乙癸逢牛福祿昌，庚趕馬頭辛帶巳，壬騎龍背喜非常」。丁作亥、癸作丑，
	// 與遁時法的酉、卯不同。《三命通會》評該歌「是以年論⋯若以日遁則非」。
	FuXingShenFeng

	// ── 天羅地網 ──

	// TianLuoBySound 帶納音條件。《三命通會·論天羅地網》與《命理探源》
	// 引《淵海子平》：「火命人有天羅，水土命人有地網，餘金木二命無之」。
	TianLuoBySound
	// TianLuoBranchOnly 不看納音。《五行精紀·天羅地網歌》「凡以戌亥為天羅，
	// 辰巳為地網」，通篇未及納音。
	TianLuoBranchOnly

	// ── 勾絞方向 ──

	// GouJiaoFrontIsGou 陽男陰女命前三辰為勾。《三命通會·論勾絞》
	// 「陽男陰女命前三辰為勾，命後三辰為絞」，《五行精紀》自註引珞琭子、
	// 林開五命「皆以陽命人前三辰為勾」——多數，故為預設。
	GouJiaoFrontIsGou
	// GouJiaoFrontIsJiao 《五行精紀》開篇「陽命人前辰為絞，後三辰為勾」。
	// 該書自註即與此相反，係書內異說。
	GouJiaoFrontIsJiao

	sectCount
)

// topicSects 各主題的取法，第一項為預設。
var topicSects = [topicCount][]Sect{
	TopicBranchBase: {BranchBaseCustomary, BranchBaseDay, BranchBaseYear, BranchBaseBoth},
	TopicStemBase:   {StemBaseCustomary, StemBaseDay, StemBaseYear, StemBaseBoth},
	TopicTianYi:     {TianYiSanMing, TianYiYeHuiTing},
	TopicBladeAt:    {BladeAtLuNext, BladeAtProsperity},
	TopicYinBlade:   {YinBladeMarked, YinBladeNone},
	TopicFuXing:     {FuXingHourStem, FuXingShenFeng},
	TopicTianLuo:    {TianLuoBySound, TianLuoBranchOnly},
	TopicGouJiao:    {GouJiaoFrontIsGou, GouJiaoFrontIsJiao},
}

// topicKinds 各主題影響哪些神煞。空表示影響面不限於特定幾個（基準之爭）。
//
// 供上層在畫面上標示「此項另有一派」之用；計算不依賴它。
var topicKinds = [topicCount][]Kind{
	TopicTianYi:   {TianYiGuiRen},
	TopicBladeAt:  {YangRen, FeiRen},
	TopicYinBlade: {YangRen, FeiRen},
	TopicFuXing:   {FuXingGuiRen},
	TopicTianLuo:  {TianLuoDiWang},
	TopicGouJiao:  {GouJiao},
}

// Topics 全部主題，供呼叫方列舉。
func Topics() []Topic {
	out := make([]Topic, 0, topicCount)
	for t := Topic(0); t < topicCount; t++ {
		out = append(out, t)
	}
	return out
}

// Sects 該主題可選的取法，第一項為預設。
func (t Topic) Sects() []Sect {
	if t >= topicCount {
		return nil
	}
	return topicSects[t]
}

// Kinds 該主題影響哪些神煞。回傳 nil 表示影響面不限於特定幾個。
func (t Topic) Kinds() []Kind {
	if t >= topicCount {
		return nil
	}
	return topicKinds[t]
}

// Topic 該取法屬於哪個主題，第二個回傳值為是否找到。
func (s Sect) Topic() (Topic, bool) {
	for t := Topic(0); t < topicCount; t++ {
		for _, x := range topicSects[t] {
			if x == s {
				return t, true
			}
		}
	}
	return 0, false
}

// Options 神煞的計算口徑。
type Options struct {
	// Categories 要計算哪些分類。nil 表示全部。
	Categories []Category
	Include    []Kind // 在 Categories 之外額外加入
	Exclude    []Kind // 自結果中排除

	// Sects 各主題採用的取法。零值即全部採預設，不必初始化。
	//
	// 用陣列而非 map：map 的零值是 nil，讀取雖安全但寫入會 panic，
	// 呼叫方得記得 make——那是設定結構不該有的儀式。
	Sects [topicCount]Sect
}

// Default 預設口徑：全部分類，各主題採首選。
func Default() Options { return Options{} }

// With 回傳套用指定取法後的副本。取法不屬於任何主題時原樣回傳。
//
//	opt := shensha.Default().With(shensha.BranchBaseBoth).With(shensha.FuXingShenFeng)
func (o Options) With(sects ...Sect) Options {
	for _, s := range sects {
		if t, ok := s.Topic(); ok {
			o.Sects[t] = s
		}
	}
	return o
}

// sect 取該主題實際生效的取法。
//
// 未指定或指定了不屬於本主題者，一律回主題的首選——設定寫錯時給預設值，
// 勝於給另一個主題的取法。
func (o Options) sect(t Topic) Sect {
	if t >= topicCount {
		return SectDefault
	}
	choices := topicSects[t]
	for _, x := range choices {
		if x == o.Sects[t] {
			return x
		}
	}
	return choices[0]
}

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
