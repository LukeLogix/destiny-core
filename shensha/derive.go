package shensha

import "github.com/LukeLogix/destiny-core/ganzhi"

// 本檔是「表是果不是因」的落實處：所有位置皆由既有原語推導，不硬編對照表。
// 典籍表只出現在測試中，與此處的推導互相比對。

// ── 三合局 ──
//
// 三合局的四組地支恰好以 4 為模分群：申子辰皆 ≡0、巳酉丑 ≡1、
// 寅午戌 ≡2、亥卯未 ≡3。各局的五行以其陽干為代表，而陽干的長生位
// 不受陰陽口徑影響，故此處可預先算定。
var trinityStart = [4]ganzhi.BranchIndex{
	8,  // 申子辰為水局，壬長生在申
	5,  // 巳酉丑為金局，庚長生在巳
	2,  // 寅午戌為火局，丙長生在寅
	11, // 亥卯未為木局，甲長生在亥
}

// trinityAt 三合局五行的第 t 個長生位落在哪一支。
//
// 將星、華蓋、驛馬、咸池、劫煞、亡神、災煞、六厄八者，即此函式配上
// 帝旺、墓、病、沐浴、絕、臨官、胎、死八個取位——不是八張表，
// 是同一個東西的八個取位。
func trinityAt(base ganzhi.BranchIndex, t ganzhi.Terrain) ganzhi.BranchIndex {
	start := trinityStart[int(base)%4]
	return ganzhi.BranchIndex((int(start) + int(t)) % ganzhi.BranchCount)
}

// ── 三會方 ──
//
// 三會方為三個相鄰地支：亥子丑北、寅卯辰東、巳午未南、申酉戌西。
// 以 (支+1)/3 mod 4 分群。
func directionGroup(b ganzhi.BranchIndex) int {
	return ((int(b) + 1) / 3) % 4
}

// directionFirst 本方的首支
func directionFirst(b ganzhi.BranchIndex) ganzhi.BranchIndex {
	return ganzhi.BranchIndex((11 + 3*directionGroup(b)) % ganzhi.BranchCount)
}

// directionLast 本方的末支
func directionLast(b ganzhi.BranchIndex) ganzhi.BranchIndex {
	return ganzhi.BranchIndex((1 + 3*directionGroup(b)) % ganzhi.BranchCount)
}

// ── 日干：祿位 ──
//
// 祿神即臨官位。此處固定以陰干逆行取之，不隨 Options.Terrain 變動——
// 祿的十干對照（甲祿寅、乙祿卯⋯）各家一致，不在「陰陽同生同死」的
// 爭議範圍內；若跟著口徑走，乙祿會變成寅，與所有典籍不符。
var luBranch = func() [ganzhi.StemCount]ganzhi.BranchIndex {
	var out [ganzhi.StemCount]ganzhi.BranchIndex
	for s := 0; s < ganzhi.StemCount; s++ {
		out[s] = branchAtTerrain(ganzhi.StemIndex(s), ganzhi.Officer)
	}
	return out
}()

// offsetFromLu 祿位偏移 n 支。羊刃為 +1、金輿為 +2。
//
// 羊刃並非帝旺位——《三命通會·論羊刃》載「羊刃常居祿前一辰」，並列
// 乙刃在辰、丁刃在未、辛刃在戌、癸刃在丑。以帝旺位推之，陽干五個相符
// 但陰干五個全錯，且兩個 TerrainSect 口徑皆無法產生原文結果。
func offsetFromLu(stem ganzhi.StemIndex, n int) ganzhi.BranchIndex {
	b := int(luBranch[int(stem)%ganzhi.StemCount]) + n
	return ganzhi.BranchIndex(((b % ganzhi.BranchCount) + ganzhi.BranchCount) % ganzhi.BranchCount)
}

// ── 日干：文昌貴人 ──
//
// 陽干取病位、陰干取長生位（陰干逆行）。
//
// 坊間另有「食神的臨官位」一說，於丙丁破功——丙的食神為戊、丁的食神為己，
// 而土的寄宮口徑（火土同宮 vs 水土同宮）在此分歧。《命理探源》的按語即
// 於該二干改用長生位並採水土同宮以迴避。本式無此例外。
var wenChangBranch = func() [ganzhi.StemCount]ganzhi.BranchIndex {
	var out [ganzhi.StemCount]ganzhi.BranchIndex
	for s := 0; s < ganzhi.StemCount; s++ {
		stem := ganzhi.StemIndex(s)
		want := ganzhi.Sick
		if stem.Polarity() == ganzhi.Yin {
			want = ganzhi.LongLife
		}
		out[s] = branchAtTerrain(stem, want)
	}
	return out
}()

// branchAtTerrain 反查：該干的指定長生位落在哪一支。
//
// 十二長生與十二支一一對應，故必然找得到；找不到即為 TerrainOf 出錯，
// 此時 panic 勝於靜默回傳子——後者會讓整組神煞悄悄錯位。
// 本函式只在套件初始化時呼叫，panic 會在載入階段就暴露，不會留到執行期。
func branchAtTerrain(stem ganzhi.StemIndex, want ganzhi.Terrain) ganzhi.BranchIndex {
	for b := 0; b < ganzhi.BranchCount; b++ {
		if ganzhi.TerrainOf(stem, ganzhi.BranchIndex(b), ganzhi.TerrainYinReverse) == want {
			return ganzhi.BranchIndex(b)
		}
	}
	panic("shensha: terrain does not cover all branches")
}

// ── 天乙貴人 ──
//
// 《三命通會·論天乙貴人》載明機制，四條規則：
//
//  1. 陽貴從子起甲順行，陰貴從申起甲逆行
//  2. 跳過辰（天羅）、戌（地網）、起貴支之對衝（天空）、起貴支本身（不再臨）
//  3. 佈定後不取該干，取其五合對象——「不取甲德，而取合氣」
//  4. 完
//
// 原文自述其結論為「甲戊庚牛羊，六辛逢馬虎」，本式推導與之十干全中。
// 推導比口訣多給一項資訊：口訣只說甲的貴人在丑未，推導知道未是陽貴、
// 丑是陰貴，而原文云「冬至用陽，夏至用陰」，兩者可分主次。
func nobleBranches() (yang, yin [ganzhi.StemCount]ganzhi.BranchIndex) {
	return layoutNoble(0, +1), layoutNoble(8, -1)
}

func layoutNoble(start, dir int) [ganzhi.StemCount]ganzhi.BranchIndex {
	clash := (start + 6) % ganzhi.BranchCount // 起貴支之對衝，即天空

	var out [ganzhi.StemCount]ganzhi.BranchIndex
	b, s, firstPass := start, 0, true
	// 上界：每圈至多跳過辰、戌、天空、起貴支四支，兩圈可佈的位數為
	// 9 + 8 = 17，足以容納十干。設界只為防止規則日後變動時陷入無窮迴圈——
	// 本函式在套件初始化時執行，無窮迴圈會使匯入直接卡死。
	for guard := 0; s < ganzhi.StemCount; guard++ {
		if guard > 2*ganzhi.BranchCount {
			panic("shensha: noble layout could not place all stems; check the skip rules")
		}
		skip := b == 4 || b == 10 || // 辰為天羅、戌為地網，貴人不臨
			b == clash || // 天空，貴人有獨無對
			(!firstPass && b == start) // 起貴之所，貴人不再臨
		if !skip {
			// 佈定的是「德」，取其五合對象方為貴人
			out[(s+5)%ganzhi.StemCount] = ganzhi.BranchIndex(b)
			s++
		}
		b = ((b+dir)%ganzhi.BranchCount + ganzhi.BranchCount) % ganzhi.BranchCount
		if b == start {
			firstPass = false
		}
	}
	return out
}

var yangNoble, yinNoble = nobleBranches()

// ── 天干五合 ──
//
// 甲己、乙庚、丙辛、丁壬、戊癸。天德合、月德合即德神之合。
func hePartner(s ganzhi.StemIndex) ganzhi.StemIndex {
	return ganzhi.StemIndex((int(s) + 5) % ganzhi.StemCount)
}

// ── 地支六合 ──
//
// 子丑、寅亥、卯戌、辰酉、巳申、午未。暗祿即祿位之六合——
// 〈暗祿格〉「甲人辛亥暗中祿」，甲祿在寅，寅與亥合。
func liuHe(b ganzhi.BranchIndex) ganzhi.BranchIndex {
	return ganzhi.BranchIndex((13 - int(b)) % ganzhi.BranchCount)
}

// clash 地支六沖
func clash(b ganzhi.BranchIndex) ganzhi.BranchIndex {
	return ganzhi.BranchIndex((int(b) + 6) % ganzhi.BranchCount)
}

// ── 太極貴人 ──
//
// 〈論太極貴〉：「甲乙木先造乎子⋯後終乎午；丙丁火先喜出乎卯，後喜藏乎酉；
// 庚辛金得寅⋯見亥；壬癸水先得申而生，後得巳而納；戊己土⋯得辰戌丑未為正庫」
//
// 兩處異文，本表均從《五行精紀》口訣「甲乙生人子午中，丙丁卯酉定亨通，
// 戊己兩干臨四季，庚辛寅亥祿盈豐，壬巳癸申一作壬癸巳申」：
//
//	戊己  三命通會散文另有「喜生乎申」，照其餘四組的「生位＋庫位」體例
//	      應含申；口訣只作四季，本表取四季。
//	壬癸  該書自記「一作」兩讀，三命通會散文支持壬癸各得巳申，本表從之。
//
// 基準為日干或年干兩派皆有據——古法以年、子平以日，見決策日誌 D-51。
// 本套件預設日干，由 Options.StemBase 切換。
var taiJiTargets = [ganzhi.StemCount][]ganzhi.BranchIndex{
	0: {0, 6},        // 甲：子午
	1: {0, 6},        // 乙：子午
	2: {3, 9},        // 丙：卯酉
	3: {3, 9},        // 丁：卯酉
	4: {4, 10, 1, 7}, // 戊：辰戌丑未
	5: {4, 10, 1, 7}, // 己：辰戌丑未
	6: {2, 11},       // 庚：寅亥
	7: {2, 11},       // 辛：寅亥
	8: {8, 5},        // 壬：申巳
	9: {8, 5},        // 癸：申巳
}

// ── 天德 ──
//
// 《三命通會·論天月徳》記其來由時，寅申巳亥四月落在乾坤艮巽四卦而未言
// 如何轉為干支。取《五行精紀》（宋，引神白經）與《神峰通考》（明）的同一
// 口訣比對後解決——一本用卦、一本用支，互為對照表：乾＝亥、艮＝寅、
// 巽＝巳，由對稱推得坤＝申。四者與《三命通會》自述的「寅中有艮，巳中有巽，
// 申中有坤，亥中有乾」完全吻合。
//
// 值為天干時取 0–9，為地支時取 100+支。干支混雜是原本如此，非傳抄訛誤。
var tianDeByMonth = [ganzhi.BranchCount]int{
	0:  100 + 5,  // 子月→巳
	1:  6,        // 丑月→庚
	2:  3,        // 寅月→丁
	3:  100 + 8,  // 卯月→申
	4:  8,        // 辰月→壬
	5:  7,        // 巳月→辛
	6:  100 + 11, // 午月→亥
	7:  0,        // 未月→甲
	8:  9,        // 申月→癸
	9:  100 + 2,  // 酉月→寅
	10: 2,        // 戌月→丙
	11: 1,        // 亥月→乙
}

// ── 月德 ──
//
// 〈論天月徳〉：「申子辰會酉，出庚，入垣於壬；亥卯未會午，出丙，入垣於甲；
// 寅午戌會卯，出甲，入垣於丙；巳酉丑會子，出壬，入垣於庚」
func yueDeStem(month ganzhi.BranchIndex) ganzhi.StemIndex {
	switch int(month) % 4 {
	case 0: // 申子辰
		return 8 // 壬
	case 1: // 巳酉丑
		return 6 // 庚
	case 2: // 寅午戌
		return 2 // 丙
	default: // 亥卯未
		return 0 // 甲
	}
}

// ── 德秀貴人 ──
//
// 〈論徳秀〉：「寅午戌月丙丁為徳，戊癸為秀；申子辰月壬癸戊己為徳，丙辛甲己
// 為秀；巳酉丑月庚辛為徳，乙庚為秀；亥卯未月甲乙為徳，丁壬為秀」
//
// 依局分組，索引同 trinityStart：申子辰≡0、巳酉丑≡1、寅午戌≡2、亥卯未≡3。
var deXiuStems = [4][]ganzhi.StemIndex{
	0: {8, 9, 4, 5, 2, 7, 0, 5}, // 申子辰：德壬癸戊己 秀丙辛甲己
	1: {6, 7, 1, 6},             // 巳酉丑：德庚辛 秀乙庚
	2: {2, 3, 4, 9},             // 寅午戌：德丙丁 秀戊癸
	3: {0, 1, 3, 8},             // 亥卯未：德甲乙 秀丁壬
}

// ── 暗金的煞 ──
//
// 〈論暗金的煞〉：「子午卯酉在巳，寅申巳亥在酉，辰戌丑未在丑」。
// 原文明載「此一煞而有三名：一曰吟呻，二曰破碎，三曰白衣」，故合為一項。
// 四孟四仲四季以 mod 3 分群——四孟寅申巳亥為 2、四季辰戌丑未為 1、
// 四仲子午卯酉為 0。不可用 mod 4，那是三合局的分群。
func anJinTarget(b ganzhi.BranchIndex) ganzhi.BranchIndex {
	switch int(b) % 3 {
	case 0: // 子午卯酉為四仲
		return 5 // 巳
	case 2: // 寅申巳亥為四孟
		return 9 // 酉
	default: // 辰戌丑未為四季
		return 1 // 丑
	}
}

// ── 納音系的十二長生軸 ──
//
// 學堂、詞館、正印三者是同一條軸上的三個取位：長生為學堂、臨官為詞館、
// 墓為正印。故只記長生起點，另兩者以位移取得。
//
// 納音系一律水土同宮——土隨水於申起長生。此與 ganzhi.TerrainOf 的火土同宮
// 是兩個口徑，不可借用：借用會讓土命的學堂落在寅，而典籍作申。
var soundLongLife = [ganzhi.ElementCount]ganzhi.BranchIndex{
	ganzhi.Wood:  11, // 木長生亥
	ganzhi.Fire:  2,  // 火長生寅
	ganzhi.Earth: 8,  // 土長生申（從水）
	ganzhi.Metal: 5,  // 金長生巳
	ganzhi.Water: 8,  // 水長生申
}

// soundTerrainAt 納音五行在十二長生第 n 位的地支。長生 0、臨官 3、墓 8。
func soundTerrainAt(e ganzhi.Element, n int) ganzhi.BranchIndex {
	return wrap(int(soundLongLife[e]) + n)
}

// 十二長生的序位。〈論學堂詞館〉「長生乃學堂之正位⋯臨官乃詞館正位」，
// 〈論正印〉「正印者乃五行之正庫」——庫即墓。
const (
	terrainLongLife = 0
	terrainOfficer  = 3
	terrainTomb     = 8
)

// ── 四廢 ──
//
// 《五行精紀》〈四廢日〉：「春庚申、辛酉，夏壬子、癸亥，秋甲寅、乙卯，
// 冬丙午、丁巳，此五行無氣之日」。季節由月支判定：寅卯辰春、巳午未夏、
// 申酉戌秋、亥子丑冬。
//
// 不抄那八組干支——原文的八組正是「干支同五行」的全部，故只需記下各季
// 無氣的五行，干支自己會對上：金只有庚申、辛酉兩組同氣（庚陽配申陽、
// 辛陰配酉陰），水只有壬子、癸亥，餘同。土永不入列，因三會方無土之季。
var siFeiElement = [4]ganzhi.Element{
	0: ganzhi.Fire,  // 冬（亥子丑）水旺，火廢
	1: ganzhi.Metal, // 春（寅卯辰）木旺，金廢
	2: ganzhi.Water, // 夏（巳午未）火旺，水廢
	3: ganzhi.Wood,  // 秋（申酉戌）金旺，木廢
}

// ── 三奇 ──
//
// 〈論三奇〉：天上三奇乙丙丁，地下三奇甲戊庚。原文另載人中三奇辛壬癸，
// 但自評「其說無據」，故不收。
//
// 須「三干相連而無間」——四柱天干中須連續出現，不得跳柱。
var sanQiSets = [2][3]ganzhi.StemIndex{
	{1, 2, 3}, // 乙丙丁：天上三奇
	{0, 4, 6}, // 甲戊庚：地下三奇
}

// ── 十二歲君 ──
//
// 同一條「自太歲順數十二位」的軸，兩組傳承各給各的名。
//
// 神峰組（預設）：《神峰通考》與《命理探源》所載口訣逐字一致，且明言
// 「以太歲為第一位順數」——起例有據。《三命通會·總論諸神煞》的「取太歲
// 前五辰」等表述據此定為 ±N；唯官符照字面為 +5，該位已為死符所佔，該書
// 自身即矛盾，故採口訣的第五位即 +4。
//
// 洞微經組：《五行精紀》〈太歲十二殺〉只列名單「太歲、生氣、喪門、天醫、
// 官符、死符、大耗、發盜、福德、大吉、弔客、病符（洞微經）」，另附靈台經
// 的吉凶歌，**通篇無起例**。整個語料中此組僅此一見，無從交叉比對。
//
// 故本組的十二個位置係由名單順序推定，非典籍明言。推定有兩個佐證：
// 大耗排第七位即 +6，而大耗為歲破之別名、歲破即太歲對衝，獨立可驗；
// 喪門(2)、官符(4)、死符(5)、弔客(10)、病符(11) 五位與神峰組同名同位。
// 但福德在本名單排第九、神峰組排第十，同名而異位——推定在此撐不住。
//
// 因此本組列為非預設選項，見決策日誌 D-55。
var suiJunSets = [2][ganzhi.BranchCount]Kind{
	// 神峰通考／命理探源
	{TaiSui, TaiYang, SangMen, TaiYin, GuanFu, SiFu,
		SuiPo, LongDe, BaiHu, FuDe, DiaoKe, BingFu},
	// 洞微經（位置為推定）
	{TaiSui, ShengQi, SangMen, TianYiSuiJun, GuanFu, SiFu,
		DaHao, FaDao, FuDe, DaJi, DiaoKe, BingFu},
}

// suiJunOffsetOf 該歲君在指定組別中的位移。不屬於該組則第二個回傳值為 false，
// 呼叫方據以略過——切到洞微組時太陽、太陰、歲破、龍德、白虎即不成立。
func suiJunOffsetOf(k Kind, set int) (int, bool) {
	for i, x := range suiJunSets[set] {
		if x == k {
			return i, true
		}
	}
	return 0, false
}

// ── 四時天喜 ──
//
// 《三命通會·總論諸神煞》：「此外又有天喜神，春戌夏丑秋辰冬未，遇者主歡欣」；
// 《五行精紀》〈天神喜〉：「春戌夏丑逢天喜，秋辰冬未三三指」——兩書一致。
//
// 與《神峰通考》鸞喜歌的天喜（紅鸞之對衝，以年支起）同名而異物，故另立
// 一個 Kind。此為本套件第六次遇上同名異物，前五次為文昌、大耗、官符、
// 白虎、劍鋒。
//
// 索引同 directionGroup：冬 0、春 1、夏 2、秋 3。
var tianXiSiShi = [4]ganzhi.BranchIndex{
	0: 7,  // 冬（亥子丑）：未
	1: 10, // 春（寅卯辰）：戌
	2: 1,  // 夏（巳午未）：丑
	3: 4,  // 秋（申酉戌）：辰
}

// xunKong 旬空（空亡）：該甲子所屬旬中缺的兩支。
//
// 六十甲子每旬十位，十干配十二支必餘二支——「有是位而無祿曰空，
// 有支而無干曰亡」（〈論空亡〉）。十惡大敗即「祿入空亡」，故此處需要它。
func xunKong(sex ganzhi.SexagenaryIndex) (ganzhi.BranchIndex, ganzhi.BranchIndex) {
	head := int(sex) - int(sex.Stem()) // 該旬之首（旬首必為甲）
	return wrap(head + 10), wrap(head + 11)
}

// ── 第三批：食神衍生、鸞喜、柱本身 ──

// sexFrom 由干支反查六十甲子序位。
//
// 十干十二支同陰陽者必有唯一組合。手算這個序位是錯誤高發處——四廢初版
// 八組錯了四組——故一律由此函式算，不寫死。
func sexFrom(stem ganzhi.StemIndex, branch ganzhi.BranchIndex) ganzhi.SexagenaryIndex {
	x, ok := sexFromOK(stem, branch)
	if !ok {
		panic("shensha: stem and branch have opposite polarity")
	}
	return x
}

// sexFromOK 同上，但陰陽不合時回 false 而非 panic。
//
// 兩個版本的分工是刻意的：sexFrom 用於套件初始化時建表，那裡的不合法
// 組合是實作者的錯，panic 在載入階段就會炸出來；sexFromOK 用於執行期
// 處理 Input，那裡的不合法組合是呼叫方給的，不該讓函式庫掛掉。
//
// Input 的 DayStem 與 DayBranch 是兩個獨立欄位，湊出不合法的組合很容易；
// 此時空亡與魁罡不產生命中，其餘神煞不受影響。
func sexFromOK(stem ganzhi.StemIndex, branch ganzhi.BranchIndex) (ganzhi.SexagenaryIndex, bool) {
	for x := 0; x < ganzhi.SexagenaryCount; x++ {
		if x%ganzhi.StemCount == int(stem) && x%ganzhi.BranchCount == int(branch) {
			return ganzhi.SexagenaryIndex(x), true
		}
	}
	return 0, false
}

// eatGod 食神：我生者而同陰陽。甲之食神為丙、癸之食神為乙。
func eatGod(s ganzhi.StemIndex) ganzhi.StemIndex {
	e := (int(s)/2 + 1) % ganzhi.ElementCount
	return ganzhi.StemIndex(e*2 + int(s)%2)
}

// ── 天廚貴人 ──
//
// 《五行精紀》〈天廚格〉：「此以**食神見祿**推之：甲丙愛行雙女遊（巳），
// 乙丁獅子（午）己金牛（酉），戊樂陰陽（申）庚亥地，癸來天蠍（卯）
// 壬人馬（寅），辛到寶瓶（子）福自由」。原文自述取法，故直接推導。
var tianChuBranch = func() [ganzhi.StemCount]ganzhi.BranchIndex {
	var out [ganzhi.StemCount]ganzhi.BranchIndex
	for s := 0; s < ganzhi.StemCount; s++ {
		out[s] = luBranch[eatGod(ganzhi.StemIndex(s))]
	}
	return out
}()

// ── 五鼠遁：日干起子時之干 ──
//
// 甲己還加甲、乙庚丙作初、丙辛從戊起、丁壬庚子居、戊癸壬子頭。
// 即子時干為 (日干 mod 5) × 2。
func hourStemStart(day ganzhi.StemIndex) ganzhi.StemIndex {
	return ganzhi.StemIndex((int(day) % 5) * 2)
}

// hourStemAt 該日干在指定時支的遁干
func hourStemAt(day ganzhi.StemIndex, b ganzhi.BranchIndex) ganzhi.StemIndex {
	return ganzhi.StemIndex((int(hourStemStart(day)) + int(b)) % ganzhi.StemCount)
}

// ── 福星貴人 ──
//
// 《協紀辨方書》〈福星貴人〉：「日干生時干曰福星貴人⋯甲日寅時必為丙寅，
// 乙日丑時亥時必為丁丑丁亥⋯**皆本日日干之食神子孫**」，並逐日列出結果。
// 《三命通會》同章亦云「遁得本旬中真食神」，機制一致。
//
// 故取法為：該日干遁時，遁得食神的那些時支。丙、乙兩干各有兩個時支
// （十干配十二支，兩支的遁干重複），餘皆一個。
//
// 《神峰通考》另載歌訣「甲丙相邀入虎鄉⋯戊申己未丁亥遇，乙癸逢牛福祿昌」，
// 其丁作亥、癸作丑與本推導的酉、卯不合。該歌《三命通會》明言「是以年論」
// 且「若以日遁則非」，故不採，分歧記於決策日誌。
var fuXingBranches = func() [ganzhi.StemCount][]ganzhi.BranchIndex {
	var out [ganzhi.StemCount][]ganzhi.BranchIndex
	for s := 0; s < ganzhi.StemCount; s++ {
		stem := ganzhi.StemIndex(s)
		want := eatGod(stem)
		for b := 0; b < ganzhi.BranchCount; b++ {
			if hourStemAt(stem, ganzhi.BranchIndex(b)) == want {
				out[s] = append(out[s], ganzhi.BranchIndex(b))
			}
		}
	}
	return out
}()

// ── 紅艷煞 ──
//
// 《神峰通考》〈紅艷煞〉：「多情多欲少人知，六丙逢寅辛見雞。癸臨申上
// 丁見未，眉開眼笑樂嬉嬉。甲乙午申庚見戌，世間只是眾人妻。戊己怕辰
// 壬怕子，祿馬相逢作路妓」。
//
// 此項查無取法自述，只有結果。與祿、長生、沐浴諸位皆對不上，故照抄成表——
// 抄表是不得已，不是偷懶：有機制就推導，只有結果時抄表並註明無機制可依。
var hongYanBranch = [ganzhi.StemCount]ganzhi.BranchIndex{
	0: 6,  // 甲：午
	1: 8,  // 乙：申
	2: 2,  // 丙：寅
	3: 7,  // 丁：未
	4: 4,  // 戊：辰
	5: 4,  // 己：辰
	6: 10, // 庚：戌
	7: 9,  // 辛：酉
	8: 0,  // 壬：子
	9: 8,  // 癸：申
}

// ── 紅鸞、天喜 ──
//
// 《神峰通考》〈鸞喜二德解神歌〉：「卯起紅鸞逆數通，欲知天喜是相衝」——
// 子年紅鸞在卯，年支每進一位紅鸞退一位；天喜為紅鸞之對衝。
//
// 《三命通會》另有「天喜神：春戌夏丑秋辰冬未」，同名而異物，本套件不收，
// 分歧記於決策日誌。
func hongLuan(year ganzhi.BranchIndex) ganzhi.BranchIndex {
	return wrap(3 - int(year))
}

// ── 魁罡 ──
//
// 《三命通會》：「【詩曰】壬辰庚戌與庚辰戊戌，魁罡四座神，不見財官刑煞并，
// 身行旺地貴無倫」，並自按「此格俱用辰戌，獨天干少異，內庚辰二日既曰日德
// 又曰魁罡」。四日明列，不存版本之爭。
//
// 全章以日柱立論：「魁罡四日最為先」「戊戌日無財不貴⋯若魁罡重疊有情，
// 富貴兩全」「日主獨逢沖尅重」，末尾二例（庚午 丁亥 戊戌 丙辰、
// 丁亥 癸丑 庚戌 戊寅）魁罡皆在日柱，作者稱之為「魁罡日」。
//
// 故「疊疊相逢」的前提是日柱已為魁罡，非任一柱各自成立——初版讀成後者，
// 於日柱非魁罡的盤上仍報出他柱，見決策日誌 D-59。
var kuiGangSex = []ganzhi.SexagenaryIndex{
	sexFrom(8, 4),  // 壬辰
	sexFrom(6, 10), // 庚戌
	sexFrom(6, 4),  // 庚辰
	sexFrom(4, 10), // 戊戌
}

// ── 天赦 ──
//
// 《三命通會》〈天赦日〉、《五行精紀》引《三歷會同》、《淵海子平》、
// 《神峰通考》、《命理探源》引《天寶曆》、《協紀辨方書》引《曆例》
// 六書一致：「春戊寅、夏甲午、秋戊申、冬甲子」，無分歧。
//
// 《天寶曆》述其理：「天之生育甲與戊，地之成立子午寅申，故以甲戊配成天赦」。
//
// 索引同 directionGroup：冬 0、春 1、夏 2、秋 3。
var tianSheSex = [4]ganzhi.SexagenaryIndex{
	0: sexFrom(0, 0), // 冬（亥子丑）：甲子
	1: sexFrom(4, 2), // 春（寅卯辰）：戊寅
	2: sexFrom(0, 6), // 夏（巳午未）：甲午
	3: sexFrom(4, 8), // 秋（申酉戌）：戊申
}

// ── 金神 ──
//
// 《三命通會》〈金神〉：「金神者破敗之神，即的煞，止有三時，乃癸酉、己巳、
// 乙丑」；《淵海子平》〈論金神〉、《神峰通考》所載三時全同。
//
// 《三命通會》另加「此格六甲日為主」的限制，《淵海子平》與《神峰通考》
// 皆不設此限，且神峰所舉命例正是己未日。金神格成立與否是格局判斷，
// 不在神煞偵測的粒度內，故此處只報該柱的干支構成，不加日干條件。
var jinShenSex = []ganzhi.SexagenaryIndex{
	sexFrom(9, 9), // 癸酉
	sexFrom(5, 5), // 己巳
	sexFrom(1, 1), // 乙丑
}

// ── 羊刃的兩種取位 ──
//
// 《三命通會·論羊刃》：「羊刃常居祿前一辰⋯卯者甲之正位，辰者乙之正位，
// 午者丙之正位，未者丁之正位，酉者庚之正位，戌者辛之正位，子者壬之正位，
// 丑者癸之正位」——陰干四例（乙辰、丁未、辛戌、癸丑）逐一吻合祿前一辰。
//
// 帝旺一派通行於今，陽干與祿前一辰同位，陰干五個全異：乙寅、丁巳、辛申、
// 癸亥、己巳，與上引原文不合。兩派皆有人用，故列為主題而非逕自判錯。
//
// 帝旺固定以陰干逆行取，不隨 Options 之外的口徑變動——理由同祿神（D-33）。
func bladeAt(stem ganzhi.StemIndex, sect Sect) ganzhi.BranchIndex {
	if sect == BladeAtProsperity {
		return branchAtTerrain(stem, ganzhi.Prosperity)
	}
	return offsetFromLu(stem, 1)
}

// ── 福星貴人：神峰通考歌訣版 ──
//
// 「甲丙相邀入虎鄉，更逢鼠穴最高強。戊申己未丁亥遇，乙癸逢牛福祿昌。
// 庚趕馬頭辛帶巳，壬騎龍背喜非常」——甲丙同得寅子，乙癸得丑，丁得亥。
//
// 與遁時法（fuXingBranches）在丁、癸、乙、丙四干不同。《三命通會》評此歌
// 「是以年論，故有丙寅丙子；若以日遁則非」，故遁時法為預設，此為另一派。
var fuXingShenFeng = [ganzhi.StemCount][]ganzhi.BranchIndex{
	0: {2, 0}, // 甲：寅子
	1: {1},    // 乙：丑
	2: {2, 0}, // 丙：寅子
	3: {11},   // 丁：亥
	4: {8},    // 戊：申
	5: {7},    // 己：未
	6: {6},    // 庚：午
	7: {5},    // 辛：巳
	8: {4},    // 壬：辰
	9: {1},    // 癸：丑
}

// isKuiGang 該干支是否為魁罡四組之一
func isKuiGang(sex ganzhi.SexagenaryIndex) bool {
	for _, x := range kuiGangSex {
		if sex == x {
			return true
		}
	}
	return false
}
