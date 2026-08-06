package shensha

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// 雙軌互驗：本檔的表抄自典籍口訣，derive.go 的推導由既有原語算出，
// 兩者逐項比對。表若抄錯，推導會發現；推導若寫錯，表會發現。
// 與 core 對節氣的作法同構——內嵌 JPL 表 vs VSOP87 自算。
//
// 此非形式主義：實作前的查證中，「羊刃＝帝旺位」這個看似自洽的推導，
// 一對照《三命通會·論羊刃》即在陰干五個全錯。表是必要的第二意見。

// 地支索引：子0 丑1 寅2 卯3 辰4 巳5 午6 未7 申8 酉9 戌10 亥11
const (
	zi, chou, yin, mao, chen, si = 0, 1, 2, 3, 4, 5
	wu, wei, shen, you, xu, hai  = 6, 7, 8, 9, 10, 11
)

var branchName = [12]string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
var stemName = [10]string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// 三合局的代表支，順序須與 trinityStart 的分群一致：
// 申子辰(≡0)、巳酉丑(≡1)、寅午戌(≡2)、亥卯未(≡3)
var trinityRep = [4]ganzhi.BranchIndex{shen, si, yin, hai}

var trinityLabel = [4]string{"申子辰水局", "巳酉丑金局", "寅午戌火局", "亥卯未木局"}

// TestTrinityShenShaAgainstClassics 三合局八神煞對照典籍。
//
// 出處：〈論將星華葢〉「以三合中位謂之將星⋯以三合底處得庫謂之華蓋」、
// 〈論驛馬〉、〈論咸池〉「五行沐浴之地名咸池」、〈論刼煞亡神〉「刼在
// 五行絶處⋯亡在五行臨官」、〈論災煞〉「衝破將星」、〈論六厄〉「死而不生」。
func TestTrinityShenShaAgainstClassics(t *testing.T) {
	// 依序為 申子辰、巳酉丑、寅午戌、亥卯未
	cases := []struct {
		kind    Kind
		terrain ganzhi.Terrain
		want    [4]ganzhi.BranchIndex
		desc    string
	}{
		{JiangXing, ganzhi.Prosperity, [4]ganzhi.BranchIndex{zi, you, wu, mao}, "將星＝三合中位（帝旺）"},
		{HuaGai, ganzhi.Tomb, [4]ganzhi.BranchIndex{chen, chou, xu, wei}, "華蓋＝三合底處得庫（墓）"},
		{YiMa, ganzhi.Sick, [4]ganzhi.BranchIndex{yin, hai, shen, si}, "驛馬＝病位"},
		{XianChi, ganzhi.Bath, [4]ganzhi.BranchIndex{you, wu, mao, zi}, "咸池＝五行沐浴之地"},
		{JieSha, ganzhi.Void, [4]ganzhi.BranchIndex{si, yin, hai, shen}, "劫煞＝五行絕處"},
		{WangShen, ganzhi.Officer, [4]ganzhi.BranchIndex{hai, shen, si, yin}, "亡神＝五行臨官"},
		{ZaiSha, ganzhi.Fetus, [4]ganzhi.BranchIndex{wu, mao, zi, you}, "災煞＝衝破將星（胎）"},
		{LiuE, ganzhi.Death, [4]ganzhi.BranchIndex{mao, zi, you, wu}, "六厄＝死而不生"},
	}

	for _, c := range cases {
		for g, rep := range trinityRep {
			got := trinityAt(rep, c.terrain)
			if got != c.want[g] {
				t.Errorf("%s：%s 推導得 %s，典籍為 %s",
					c.desc, trinityLabel[g], branchName[got], branchName[c.want[g]])
			}
		}
	}
}

// TestTrinityAllTwelveBranches 局內三支的結果必須相同——同局同結論。
func TestTrinityAllTwelveBranches(t *testing.T) {
	for b := 0; b < 12; b++ {
		rep := trinityRep[b%4]
		for tr := 0; tr < 12; tr++ {
			a := trinityAt(ganzhi.BranchIndex(b), ganzhi.Terrain(tr))
			e := trinityAt(rep, ganzhi.Terrain(tr))
			if a != e {
				t.Errorf("%s 與同局的 %s 在長生位 %d 上結果不同（%s vs %s）",
					branchName[b], branchName[rep], tr, branchName[a], branchName[e])
			}
		}
	}
}

// TestDirectionShenShaAgainstClassics 孤辰寡宿對照典籍。
//
// 出處：〈論孤辰寡宿〉「亥子丑逐方三位，進前一辰見寅為孤，退後一辰見戌為寡」。
// 同章另給「取母絕為孤辰，夫墓妻墓為寡宿」一套，殊途同歸。
func TestDirectionShenShaAgainstClassics(t *testing.T) {
	cases := []struct {
		base   ganzhi.BranchIndex
		gu, ge ganzhi.BranchIndex
		label  string
	}{
		{hai, yin, xu, "亥子丑北方"},
		{zi, yin, xu, "亥子丑北方（自子起）"},
		{chou, yin, xu, "亥子丑北方（自丑起）"},
		{yin, si, chou, "寅卯辰東方"},
		{si, shen, chen, "巳午未南方"},
		{shen, hai, wei, "申酉戌西方"},
	}
	for _, c := range cases {
		gu := (directionLast(c.base) + 1) % ganzhi.BranchCount
		ge := (directionFirst(c.base) + ganzhi.BranchCount - 1) % ganzhi.BranchCount
		if gu != c.gu {
			t.Errorf("%s：孤辰推導得 %s，典籍為 %s", c.label, branchName[gu], branchName[c.gu])
		}
		if ge != c.ge {
			t.Errorf("%s：寡宿推導得 %s，典籍為 %s", c.label, branchName[ge], branchName[c.ge])
		}
	}
}

// TestLuShenAgainstClassics 祿神對照〈論十干祿〉。
//
// 原文：「甲祿寅、乙祿卯、庚祿申、辛祿酉、壬祿亥、癸祿子、丙祿巳、丁祿午、
// 戊寄巳、己寄午」。四庫本此處作「丙祿己」「戊寄己巳寄午」，係巳／己形近
// 之訛——推導給巳，可據以校正。
func TestLuShenAgainstClassics(t *testing.T) {
	want := [10]ganzhi.BranchIndex{yin, mao, si, wu, si, wu, shen, you, hai, zi}
	for s, w := range want {
		if got := luBranch[s]; got != w {
			t.Errorf("%s 祿推導得 %s，典籍為 %s", stemName[s], branchName[got], branchName[w])
		}
	}
}

// TestYangRenAgainstClassics 羊刃對照〈論羊刃〉。
//
// 原文：「羊刃常居祿前一辰⋯卯者甲之正位，辰者乙之正位，午者丙之正位，
// 未者丁之正位，酉者庚之正位，戌者辛之正位，子者壬之正位，丑者癸之正位」。
//
// 本測試同時釘住一項曾經的錯誤：原設計推導為「帝旺位」，陽干五個相符
// 但陰干五個全錯。故一併驗證羊刃不隨 TerrainSect 變動。
func TestYangRenAgainstClassics(t *testing.T) {
	want := [10]ganzhi.BranchIndex{mao, chen, wu, wei, wu, wei, you, xu, zi, chou}
	for s, w := range want {
		if got := offsetFromLu(ganzhi.StemIndex(s), 1); got != w {
			t.Errorf("%s 刃推導得 %s，典籍為 %s", stemName[s], branchName[got], branchName[w])
		}
	}

	// 與 TerrainSect 無關：陰干帝旺在兩個口徑下分別為寅與卯，皆非典籍的辰
	for _, sect := range []ganzhi.TerrainSect{ganzhi.TerrainYinReverse, ganzhi.TerrainSameBirth} {
		for b := 0; b < 12; b++ {
			if ganzhi.TerrainOf(1, ganzhi.BranchIndex(b), sect) == ganzhi.Prosperity {
				if ganzhi.BranchIndex(b) == chen {
					t.Errorf("口徑 %d 下乙帝旺竟在辰——若如此則「羊刃＝帝旺位」本可成立，"+
						"本測試的前提須重新檢視", sect)
				}
			}
		}
	}
}

// TestJinYuAgainstClassics 金輿對照〈論十干祿〉章節標題「祿前二辰為金轝」。
func TestJinYuAgainstClassics(t *testing.T) {
	want := [10]ganzhi.BranchIndex{chen, si, wei, shen, wei, shen, xu, hai, chou, yin}
	for s, w := range want {
		if got := offsetFromLu(ganzhi.StemIndex(s), 2); got != w {
			t.Errorf("%s 金輿推導得 %s，典籍為 %s", stemName[s], branchName[got], branchName[w])
		}
	}
}

// TestWenChangAgainstClassics 文昌貴人對照《紫微斗數》口訣（經《命理探源》轉引）。
//
// 「甲乙巳午報君知，丙戊申宮丁己雞，庚豬辛鼠壬逢虎，癸人見兔入雲梯」。
func TestWenChangAgainstClassics(t *testing.T) {
	want := [10]ganzhi.BranchIndex{si, wu, shen, you, shen, you, hai, zi, yin, mao}
	for s, w := range want {
		if got := wenChangBranch[s]; got != w {
			t.Errorf("%s 文昌推導得 %s，典籍為 %s", stemName[s], branchName[got], branchName[w])
		}
	}
}

// TestTianYiAgainstClassics 天乙貴人對照通行口訣與《命理探源》引曹震圭。
//
// 口訣：「甲戊庚牛羊，乙己鼠猴鄉，丙丁豬雞位，壬癸兔蛇藏，六辛逢馬虎」。
// 曹震圭另分晝夜：「陽貴以甲加未順行，甲得未、乙得申、丙得酉、丁得亥、
// 己得子、庚得丑、辛得寅、壬得卯、癸得巳，此晝貴也」——與推導十干逐一吻合。
func TestTianYiAgainstClassics(t *testing.T) {
	wantYang := [10]ganzhi.BranchIndex{wei, shen, you, hai, chou, zi, chou, yin, mao, si}
	wantYin := [10]ganzhi.BranchIndex{chou, zi, hai, you, wei, shen, wei, wu, si, mao}

	for s := 0; s < 10; s++ {
		if got := yangNoble[s]; got != wantYang[s] {
			t.Errorf("%s 陽貴推導得 %s，曹震圭為 %s", stemName[s], branchName[got], branchName[wantYang[s]])
		}
		if got := yinNoble[s]; got != wantYin[s] {
			t.Errorf("%s 陰貴推導得 %s，曹震圭為 %s", stemName[s], branchName[got], branchName[wantYin[s]])
		}
	}

	// 口訣是兩者的聯集，不分晝夜
	mnemonic := [10][2]ganzhi.BranchIndex{
		{chou, wei}, {zi, shen}, {hai, you}, {hai, you}, {chou, wei},
		{zi, shen}, {chou, wei}, {wu, yin}, {mao, si}, {mao, si},
	}
	for s := 0; s < 10; s++ {
		got := map[ganzhi.BranchIndex]bool{yangNoble[s]: true, yinNoble[s]: true}
		if !got[mnemonic[s][0]] || !got[mnemonic[s][1]] || len(got) != 2 {
			t.Errorf("%s 推導得 {%s,%s}，口訣為 {%s,%s}", stemName[s],
				branchName[yangNoble[s]], branchName[yinNoble[s]],
				branchName[mnemonic[s][0]], branchName[mnemonic[s][1]])
		}
	}
}

// ── 以下為第二批神煞的雙軌互驗 ──

// TestTaiJiAgainstClassics 太極貴人對照〈論太極貴〉自述的結論。
//
// derive.go 已引原文的散文段（「甲乙木先造乎子⋯後終乎午」云云），此處把
// 那段話轉成的表逐項寫死。散文轉表這一步本身會出錯，故仍值得釘住；
// 但這不是獨立的第二來源，與其他測試的「典籍 vs 推導」性質不同。
func TestTaiJiAgainstClassics(t *testing.T) {
	want := [ganzhi.StemCount][]ganzhi.BranchIndex{
		{zi, wu}, {zi, wu}, // 甲乙：子午
		{mao, you}, {mao, you}, // 丙丁：卯酉（雞兔）
		{chen, xu, chou, wei}, {chen, xu, chou, wei}, // 戊己：辰戌丑未
		{yin, hai}, {yin, hai}, // 庚辛：寅亥
		{shen, si}, {shen, si}, // 壬癸：巳申
	}
	for s, w := range want {
		got := taiJiTargets[s]
		if len(got) != len(w) {
			t.Errorf("%s 太極貴推導得 %s，典籍為 %s", stemName[s], names(got), names(w))
			continue
		}
		if !sameSet(got, w) {
			t.Errorf("%s 太極貴推導得 %s，典籍為 %s", stemName[s], names(got), names(w))
		}
	}
}

// TestAnLuIsLuHarmony 暗祿即祿位之六合。
//
// 手上只有〈暗祿格〉的「甲人辛亥暗中祿」一句可直接對照：甲祿在寅，寅亥六合，
// 故甲之暗祿在亥。餘九干在語料中查無逐條列舉，此處只釘住機制的一致性——
// 十干皆為 liuHe(祿)，且甲那一項與原文相符。
//
// 不補其餘九干的「典籍值」：憑印象寫下的口訣曾在此把癸誤作卯，
// 那實為天乙貴人的「壬癸兔蛇藏」。查不到就不寫，勝過寫個看起來有據的錯誤。
func TestAnLuIsLuHarmony(t *testing.T) {
	if got := liuHe(luBranch[0]); got != hai {
		t.Errorf("甲之暗祿推導得 %s，〈暗祿格〉作亥", branchName[got])
	}
	for s := 0; s < ganzhi.StemCount; s++ {
		lu := luBranch[s]
		an := liuHe(lu)
		// 六合為對合關係，合回去須得原位
		if liuHe(an) != lu {
			t.Errorf("%s：祿 %s 之合為 %s，再合卻得 %s，六合不對稱",
				stemName[s], branchName[lu], branchName[an], branchName[liuHe(an)])
		}
		if int(lu)+int(an) != 13 && int(lu)+int(an) != 1 {
			t.Errorf("%s：%s 與 %s 之和為 %d，六合須為子丑、寅亥⋯之配",
				stemName[s], branchName[lu], branchName[an], int(lu)+int(an))
		}
	}
}

// TestFeiRenIsBladeClash 飛刃即羊刃之對衝。此為取法本身，非另一張表。
func TestFeiRenIsBladeClash(t *testing.T) {
	want := [ganzhi.StemCount]ganzhi.BranchIndex{you, xu, zi, chou, zi, chou, mao, chen, wu, wei}
	for s, w := range want {
		if got := clash(offsetFromLu(ganzhi.StemIndex(s), 1)); got != w {
			t.Errorf("%s 飛刃推導得 %s，應為 %s", stemName[s], branchName[got], branchName[w])
		}
	}
}

// TestTianYueDeAgainstClassics 天德、月德及其合對照口訣。
//
// 天德：「正丁二坤宮，三壬四辛同，五乾六甲上，七癸八艮逢，
// 九丙十居乙，子巽丑庚中」。乾亥、坤申、艮寅、巽巳。
// 月德：「寅午戌月在丙，申子辰月在壬，亥卯未月在甲，巳酉丑月在庚」。
func TestTianYueDeAgainstClassics(t *testing.T) {
	// 月序自寅起正月，故索引為月支
	tianDe := map[ganzhi.BranchIndex]int{
		yin: 3, mao: 100 + shen, chen: 8, si: 7, // 正丁 二坤 三壬 四辛
		wu: 100 + hai, wei: 0, shen: 9, you: 100 + yin, // 五乾 六甲 七癸 八艮
		xu: 2, hai: 1, zi: 100 + si, chou: 6, // 九丙 十乙 十一巽 十二庚
	}
	for m, w := range tianDe {
		if got := tianDeByMonth[m]; got != w {
			t.Errorf("%s月天德推導得 %d，典籍為 %d", branchName[m], got, w)
		}
		// 天德合僅在天德落於天干時成立
		if w < 100 {
			if got := hePartner(ganzhi.StemIndex(w)); int(got) != (w+5)%ganzhi.StemCount {
				t.Errorf("%s月天德合推導得 %s", branchName[m], stemName[got])
			}
		}
	}

	yueDe := map[ganzhi.BranchIndex]ganzhi.StemIndex{
		yin: 2, wu: 2, xu: 2, // 寅午戌→丙
		shen: 8, zi: 8, chen: 8, // 申子辰→壬
		hai: 0, mao: 0, wei: 0, // 亥卯未→甲
		si: 6, you: 6, chou: 6, // 巳酉丑→庚
	}
	for m, w := range yueDe {
		if got := yueDeStem(m); got != w {
			t.Errorf("%s月月德推導得 %s，典籍為 %s", branchName[m], stemName[got], stemName[w])
		}
	}

	// 月德合為月德之五合：丙辛、壬丁、甲己、庚乙
	heWant := map[ganzhi.StemIndex]ganzhi.StemIndex{2: 7, 8: 3, 0: 5, 6: 1}
	for de, w := range heWant {
		if got := hePartner(de); got != w {
			t.Errorf("%s 之合推導得 %s，應為 %s", stemName[de], stemName[got], stemName[w])
		}
	}
}

// TestAnJinAgainstClassics 暗金的煞對照〈論暗金的煞〉。
//
// 「子午卯酉在巳，寅申巳亥在酉，辰戌丑未在丑」——三組分群不可混。
// 曾誤以 mod 4 分群（那是三合局的模數），午年因此得酉而非巳。
func TestAnJinAgainstClassics(t *testing.T) {
	groups := []struct {
		bs   []ganzhi.BranchIndex
		want ganzhi.BranchIndex
	}{
		{[]ganzhi.BranchIndex{zi, wu, mao, you}, si},
		{[]ganzhi.BranchIndex{yin, shen, si, hai}, you},
		{[]ganzhi.BranchIndex{chen, xu, chou, wei}, chou},
	}
	for _, g := range groups {
		for _, b := range g.bs {
			if got := anJinTarget(b); got != g.want {
				t.Errorf("%s年暗金的煞推導得 %s，典籍為 %s",
					branchName[b], branchName[got], branchName[g.want])
			}
		}
	}
}

// TestSoundAxisAgainstClassics 學堂、詞館、正印三者對照典籍。
//
// 學堂為納音長生、詞館為臨官、正印為墓庫，且該柱納音須與年命同五行，
// 故各五行恰有一組甲子。〈論學堂詞館〉舉「金命見辛巳」，
// 〈論正印〉列「金命見乙丑、木癸未、火甲戌、水土壬辰丙辰」。
//
// 納音系為水土同宮，與 ganzhi.TerrainOf 的火土同宮不同——若誤用後者，
// 土命的學堂會落在寅而非申，本測試即釘住此點。
func TestSoundAxisAgainstClassics(t *testing.T) {
	// 學堂：金巳、木亥、水土申、火寅
	xueTang := [ganzhi.ElementCount]ganzhi.BranchIndex{
		ganzhi.Metal: si, ganzhi.Wood: hai, ganzhi.Water: shen,
		ganzhi.Earth: shen, ganzhi.Fire: yin,
	}
	for e, w := range xueTang {
		if got := soundTerrainAt(ganzhi.Element(e), terrainLongLife); got != w {
			t.Errorf("%s命學堂推導得 %s，典籍為 %s",
				ganzhi.Element(e).ID(), branchName[got], branchName[w])
		}
	}

	// 正印的甲子：反查納音同五行且落在墓位者，須恰為典籍那五組
	zhengYin := map[ganzhi.Element]ganzhi.SexagenaryIndex{
		ganzhi.Metal: 1, ganzhi.Wood: 19, ganzhi.Fire: 10, // 乙丑、癸未、甲戌
		ganzhi.Water: 28, ganzhi.Earth: 52, // 壬辰、丙辰
	}
	for e, want := range zhengYin {
		var got []ganzhi.SexagenaryIndex
		for x := 0; x < 60; x++ {
			sex := ganzhi.SexagenaryIndex(x)
			if ganzhi.SoundOf(sex).Element() == e && sex.Branch() == soundTerrainAt(e, terrainTomb) {
				got = append(got, sex)
			}
		}
		if len(got) != 1 || got[0] != want {
			t.Errorf("%s命正印推導得 %v，典籍為單一的 %s%s",
				e.ID(), got, stemName[want.Stem()], branchName[want.Branch()])
		}
	}

	// 詞館同理，且金命須為壬申——原文舉學堂「辛巳」，詞館在臨官申
	for x := 0; x < 60; x++ {
		sex := ganzhi.SexagenaryIndex(x)
		if ganzhi.SoundOf(sex).Element() == ganzhi.Metal &&
			sex.Branch() == soundTerrainAt(ganzhi.Metal, terrainOfficer) && sex != 8 {
			t.Errorf("金命詞館推導另得 %s%s，應僅壬申一組",
				stemName[sex.Stem()], branchName[sex.Branch()])
		}
	}
}

// TestSiFeiAgainstClassics 四廢對照《五行精紀》〈四廢日〉。
//
// 「春庚申辛酉，夏壬子癸亥，秋甲寅乙卯，冬丙午丁巳」。推導不抄這八組，
// 只記各季無氣之五行，靠「干支同五行」把干支長出來——本測試驗證兩者等價。
func TestSiFeiAgainstClassics(t *testing.T) {
	want := map[ganzhi.BranchIndex][]ganzhi.SexagenaryIndex{
		yin:  {56, 57}, // 春（寅卯辰）：庚申、辛酉
		si:   {48, 59}, // 夏（巳午未）：壬子、癸亥
		shen: {50, 51}, // 秋（申酉戌）：甲寅、乙卯
		hai:  {42, 53}, // 冬（亥子丑）：丙午、丁巳
	}
	for month, w := range want {
		e := siFeiElement[directionGroup(month)]
		var got []ganzhi.SexagenaryIndex
		for x := 0; x < 60; x++ {
			sex := ganzhi.SexagenaryIndex(x)
			if sex.Stem().Element() == e && sex.Branch().Element() == e {
				got = append(got, sex)
			}
		}
		if len(got) != len(w) || got[0] != w[0] || got[1] != w[1] {
			t.Errorf("%s月四廢推導得 %s，典籍為 %s", branchName[month], sexNames(got), sexNames(w))
		}
	}
}

// TestShiEDaBaiAgainstClassics 十惡大敗對照〈論十惡大敗〉所列的十個日柱。
//
// 「甲辰、乙巳、壬申、丙申、丁亥、庚辰、戊戌、癸亥、辛巳、己丑」。
// 推導走「祿入旬空」，不抄這十組——原文自述其理即為「祿入空亡」。
func TestShiEDaBaiAgainstClassics(t *testing.T) {
	want := map[ganzhi.SexagenaryIndex]bool{
		40: true, 41: true, 8: true, 32: true, 23: true, // 甲辰乙巳壬申丙申丁亥
		16: true, 34: true, 59: true, 17: true, 25: true, // 庚辰戊戌癸亥辛巳己丑
	}
	for x := 0; x < 60; x++ {
		sex := ganzhi.SexagenaryIndex(x)
		lu := luBranch[int(sex.Stem())]
		a, b := xunKong(sex)
		got := lu == a || lu == b
		if got != want[sex] {
			t.Errorf("%s%s 推導為 %v，典籍為 %v",
				stemName[sex.Stem()], branchName[sex.Branch()], got, want[sex])
		}
	}
}

// TestSuiJunSets 十二歲君兩組各須為太歲起順數十二位，不重不漏。
//
// 神峰組的起例有據：《神峰通考》《命理探源》口訣逐字一致並明言順數。
// 《三命通會·總論諸神煞》的官符「取太歲前五辰」照字面為 +5，該位已為
// 死符所佔——該書自身即矛盾，故採口訣第五位即 +4。本測試釘住此決定。
//
// 洞微經組的位置係由名單順序推定，非典籍明言（見 D-55）。此處驗證兩件
// 可獨立查核的事：大耗為太歲對衝（大耗即歲破之別名），以及兩組在
// 太歲、喪門、官符、死符、弔客、病符六位上同名同位。
func TestSuiJunSets(t *testing.T) {
	shenFeng := []Kind{TaiSui, TaiYang, SangMen, TaiYin, GuanFu, SiFu,
		SuiPo, LongDe, BaiHu, FuDe, DiaoKe, BingFu}
	dongWei := []Kind{TaiSui, ShengQi, SangMen, TianYiSuiJun, GuanFu, SiFu,
		DaHao, FaDao, FuDe, DaJi, DiaoKe, BingFu}

	for set, want := range [][]Kind{shenFeng, dongWei} {
		if len(want) != ganzhi.BranchCount {
			t.Fatalf("第 %d 組共 %d 位，應為 12", set, len(want))
		}
		seen := map[Kind]bool{}
		for i, k := range want {
			got, ok := suiJunOffsetOf(k, set)
			if !ok || got != i {
				t.Errorf("第 %d 組的 %s 位移為 %d（found=%v），應為第 %d 位",
					set, k.ID(), got, ok, i)
			}
			if seen[k] {
				t.Errorf("第 %d 組的 %s 重複出現", set, k.ID())
			}
			seen[k] = true
		}
	}

	// 歲破與大耗同為太歲之對衝，故兩組的第七位皆為 +6
	for set, k := range []Kind{SuiPo, DaHao} {
		n, _ := suiJunOffsetOf(k, set)
		if n != int(clash(0)) {
			t.Errorf("%s 位移為 %d，應與對衝一致（%d）", k.ID(), n, clash(0))
		}
	}

	// 兩組同名同位者：太歲、喪門、官符、死符、弔客、病符
	for _, k := range []Kind{TaiSui, SangMen, GuanFu, SiFu, DiaoKe, BingFu} {
		a, oka := suiJunOffsetOf(k, 0)
		b, okb := suiJunOffsetOf(k, 1)
		if !oka || !okb || a != b {
			t.Errorf("%s 在兩組的位次為 %d/%d，應相同——若日後改動，"+
				"洞微組位置由名單順序推定的佐證即失效", k.ID(), a, b)
		}
	}

	// 福德兩組皆有而位次不同，正是推定撐不住之處。釘住它，
	// 免得日後有人「順手對齊」而抹掉這個已知的不確定性。
	a, _ := suiJunOffsetOf(FuDe, 0)
	b, _ := suiJunOffsetOf(FuDe, 1)
	if a == b {
		t.Errorf("福德在兩組同為第 %d 位——但洞微名單排第九、神峰組排第十，"+
			"兩者本就不同。若確有依據對齊，請一併更新 D-55 的說明", a)
	}
}

// TestTianXiSiShiAgainstClassics 四時天喜對照典籍。
//
// 《三命通會·總論諸神煞》「此外又有天喜神，春戌夏丑秋辰冬未」；
// 《五行精紀》〈天神喜〉「春戌夏丑逢天喜，秋辰冬未三三指」——兩書一致。
// 與鸞喜歌的天喜同名異物，故為不同的 Kind。
func TestTianXiSiShiAgainstClassics(t *testing.T) {
	want := map[ganzhi.BranchIndex]ganzhi.BranchIndex{
		yin: xu, si: chou, shen: chen, hai: wei, // 春戌 夏丑 秋辰 冬未
	}
	for month, w := range want {
		if got := tianXiSiShi[directionGroup(month)]; got != w {
			t.Errorf("%s月四時天喜推導得 %s，典籍為 %s",
				branchName[month], branchName[got], branchName[w])
		}
	}
	// 與鸞喜的天喜不可同位混淆：後者以年支起，前者以月令季節起
	if TianXi == TianXiSiShi {
		t.Error("兩個天喜共用同一個 Kind——同名異物須分開")
	}
}

func names(bs []ganzhi.BranchIndex) string {
	s := ""
	for _, b := range bs {
		s += branchName[b]
	}
	return s
}

func sexNames(xs []ganzhi.SexagenaryIndex) string {
	s := ""
	for _, x := range xs {
		s += stemName[x.Stem()] + branchName[x.Branch()]
	}
	return s
}

func sameSet(a, b []ganzhi.BranchIndex) bool {
	m := map[ganzhi.BranchIndex]int{}
	for _, x := range a {
		m[x]++
	}
	for _, x := range b {
		m[x]--
	}
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}

// ── 第三批：食神衍生、鸞喜、柱本身 ──

// TestSexFromRoundTrip 干支反查六十甲子——這是四廢初版手抄出錯之處。
func TestSexFromRoundTrip(t *testing.T) {
	for x := 0; x < ganzhi.SexagenaryCount; x++ {
		sex := ganzhi.SexagenaryIndex(x)
		if got := sexFrom(sex.Stem(), sex.Branch()); got != sex {
			t.Errorf("%s%s 反查得 %d，應為 %d",
				stemName[sex.Stem()], branchName[sex.Branch()], got, sex)
		}
	}
}

// TestEatGod 食神為我生者而同陰陽。
func TestEatGod(t *testing.T) {
	want := [ganzhi.StemCount]ganzhi.StemIndex{2, 3, 4, 5, 6, 7, 8, 9, 0, 1}
	for s, w := range want {
		if got := eatGod(ganzhi.StemIndex(s)); got != w {
			t.Errorf("%s 之食神推導得 %s，應為 %s", stemName[s], stemName[got], stemName[w])
		}
	}
}

// TestTianChuAgainstClassics 天廚貴人對照《五行精紀》〈天廚格〉。
//
// 原文自述取法「此以食神見祿推之」，並列結果：「甲丙愛行雙女遊（巳），
// 乙丁獅子（午）己金牛（酉），戊樂陰陽（申）庚亥地，癸來天蠍（卯）
// 壬人馬（寅），辛到寶瓶（子）福自由」——推導與所列十干逐一比對。
func TestTianChuAgainstClassics(t *testing.T) {
	want := [ganzhi.StemCount]ganzhi.BranchIndex{
		si, wu, si, wu, shen, you, hai, zi, yin, mao,
	}
	for s, w := range want {
		if got := tianChuBranch[s]; got != w {
			t.Errorf("%s 天廚推導得 %s，典籍為 %s", stemName[s], branchName[got], branchName[w])
		}
	}
}

// TestFuXingAgainstClassics 福星貴人對照《協紀辨方書》。
//
// 原文自述取法「日干生時干⋯皆本日日干之食神子孫」，並逐日列出：
// 「甲日寅時必為丙寅，乙日丑時亥時必為丁丑丁亥，丙日子時戌時必為戊子戊戌，
// 丁日酉時必為己酉，戊日申時必為庚申，己日未時必為辛未，庚日午時必為壬午，
// 辛日巳時必為癸巳，壬日辰時必為甲辰，癸日卯時必為乙卯」。
//
// 《神峰通考》歌訣的丁作亥、癸作丑與此不合，見決策日誌 D-44。
func TestFuXingAgainstClassics(t *testing.T) {
	want := [ganzhi.StemCount][]ganzhi.BranchIndex{
		{yin}, {chou, hai}, {zi, xu}, {you}, {shen},
		{wei}, {wu}, {si}, {chen}, {mao},
	}
	for s, w := range want {
		got := fuXingBranches[s]
		if len(got) != len(w) || !sameSet(got, w) {
			t.Errorf("%s 福星推導得 %s，協紀辨方書為 %s", stemName[s], names(got), names(w))
		}
	}
	// 遁時的結果必然是該干的食神——推導的內在一致性
	for s := 0; s < ganzhi.StemCount; s++ {
		stem := ganzhi.StemIndex(s)
		for _, b := range fuXingBranches[s] {
			if hourStemAt(stem, b) != eatGod(stem) {
				t.Errorf("%s 日 %s 時的遁干為 %s，非食神 %s",
					stemName[s], branchName[b], stemName[hourStemAt(stem, b)], stemName[eatGod(stem)])
			}
		}
	}
}

// TestHourStemStartAgainstClassics 五鼠遁對照《五行精紀》〈起時例〉。
//
// 「甲己還生甲，乙庚丙作初，丙辛當戊子，丁壬庚子居，戊癸壬子頭」。
func TestHourStemStartAgainstClassics(t *testing.T) {
	want := [ganzhi.StemCount]ganzhi.StemIndex{0, 2, 4, 6, 8, 0, 2, 4, 6, 8}
	for s, w := range want {
		if got := hourStemStart(ganzhi.StemIndex(s)); got != w {
			t.Errorf("%s 日起子時為 %s，口訣為 %s", stemName[s], stemName[got], stemName[w])
		}
	}
}

// TestHongLuanAgainstClassics 紅鸞天喜對照《神峰通考》〈鸞喜二德解神歌〉。
//
// 「卯起紅鸞逆數通，欲知天喜是相衝」——子年紅鸞在卯，逐年逆退一位。
func TestHongLuanAgainstClassics(t *testing.T) {
	want := [ganzhi.BranchCount]ganzhi.BranchIndex{
		zi: mao, chou: yin, yin: chou, mao: zi, chen: hai, si: xu,
		wu: you, wei: shen, shen: wei, you: wu, xu: si, hai: chen,
	}
	for y, w := range want {
		if got := hongLuan(ganzhi.BranchIndex(y)); got != w {
			t.Errorf("%s年紅鸞推導得 %s，口訣為 %s", branchName[y], branchName[got], branchName[w])
		}
		if got := clash(hongLuan(ganzhi.BranchIndex(y))); got != clash(w) {
			t.Errorf("%s年天喜推導得 %s，應為紅鸞之衝 %s", branchName[y], branchName[got], branchName[clash(w)])
		}
	}
}

// TestSelfPillarSetsAgainstClassics 魁罡、天赦、金神的干支組合。
//
// 魁罡《三命通會》「壬辰庚戌與庚辰戊戌，魁罡四座神」；
// 天赦六書一致「春戊寅、夏甲午、秋戊申、冬甲子」；
// 金神《三命通會》〈金神〉「止有三時，乃癸酉、己巳、乙丑」。
//
// 此處以干支對寫出，由 sexFrom 換算——不寫死六十甲子序位，那是四廢
// 初版八組錯四組的來源。
func TestSelfPillarSetsAgainstClassics(t *testing.T) {
	pair := func(st ganzhi.StemIndex, br ganzhi.BranchIndex) ganzhi.SexagenaryIndex {
		return sexFrom(st, br)
	}
	check := func(label string, got []ganzhi.SexagenaryIndex, want []ganzhi.SexagenaryIndex) {
		t.Helper()
		if len(got) != len(want) {
			t.Errorf("%s 推導得 %s，典籍為 %s", label, sexNames(got), sexNames(want))
			return
		}
		seen := map[ganzhi.SexagenaryIndex]bool{}
		for _, x := range got {
			seen[x] = true
		}
		for _, w := range want {
			if !seen[w] {
				t.Errorf("%s 缺 %s%s（推導得 %s）", label,
					stemName[w.Stem()], branchName[w.Branch()], sexNames(got))
			}
		}
	}

	check("魁罡", kuiGangSex, []ganzhi.SexagenaryIndex{
		pair(8, chen), pair(6, xu), pair(6, chen), pair(4, xu), // 壬辰 庚戌 庚辰 戊戌
	})
	check("金神", jinShenSex, []ganzhi.SexagenaryIndex{
		pair(9, you), pair(5, si), pair(1, chou), // 癸酉 己巳 乙丑
	})

	// 天赦按季，索引同 directionGroup
	tianShe := map[ganzhi.BranchIndex]ganzhi.SexagenaryIndex{
		yin:  pair(4, yin),  // 春（寅卯辰）：戊寅
		si:   pair(0, wu),   // 夏（巳午未）：甲午
		shen: pair(4, shen), // 秋（申酉戌）：戊申
		hai:  pair(0, zi),   // 冬（亥子丑）：甲子
	}
	for month, w := range tianShe {
		if got := tianSheSex[directionGroup(month)]; got != w {
			t.Errorf("%s月天赦推導得 %s%s，典籍為 %s%s", branchName[month],
				stemName[got.Stem()], branchName[got.Branch()],
				stemName[w.Stem()], branchName[w.Branch()])
		}
	}
}

// TestXunKongAgainstClassics 空亡對照《神峰通考》〈六甲空亡〉。
//
// 「甲子旬中無戌亥，甲戌旬中無申酉，甲申旬中無午未，甲午旬中無辰巳，
// 甲辰旬中無寅卯，甲寅旬中無子丑」。
func TestXunKongAgainstClassics(t *testing.T) {
	want := []struct {
		head ganzhi.BranchIndex // 旬首之支（旬首必為甲）
		a, b ganzhi.BranchIndex
	}{
		{zi, xu, hai}, {xu, shen, you}, {shen, wu, wei},
		{wu, chen, si}, {chen, yin, mao}, {yin, zi, chou},
	}
	for _, c := range want {
		head := sexFrom(0, c.head) // 甲＋該支
		// 旬內十位皆須得同一組空亡
		for n := 0; n < 10; n++ {
			sex := ganzhi.SexagenaryIndex((int(head) + n) % ganzhi.SexagenaryCount)
			a, b := xunKong(sex)
			if a != c.a || b != c.b {
				t.Errorf("甲%s旬的 %s%s 推導得空亡 %s%s，典籍為 %s%s",
					branchName[c.head], stemName[sex.Stem()], branchName[sex.Branch()],
					branchName[a], branchName[b], branchName[c.a], branchName[c.b])
			}
		}
	}
}

// TestKuiGangRequiresKuiGangDay 魁罡以日柱為前提，他柱重見才算重疊。
//
// 《三命通會》全章以日柱立論，末尾所舉二例的魁罡皆在日柱，且稱之為
// 「魁罡日」——本測試即以那兩例為 ground truth：
//
//	張時僉事　庚午 丁亥 戊戌 丙辰
//	劉大受少卿 丁亥 癸丑 庚戌 戊寅
//
// 初版把「疊疊相逢掌大權」讀成任一柱各自成立，於日柱非魁罡的盤上仍報出
// 他柱的魁罡干支。重疊預設了日柱本身已是魁罡，見 D-59。
func TestKuiGangRequiresKuiGangDay(t *testing.T) {
	pair := func(st ganzhi.StemIndex, br ganzhi.BranchIndex) ganzhi.SexagenaryIndex {
		return sexFrom(st, br)
	}
	kuiGangOf := func(in Input, scan []ganzhi.SexagenaryIndex) []int {
		var at []int
		for _, h := range Detect(in, scan, Options{Include: []Kind{KuiGang}, Categories: []Category{}}) {
			at = append(at, h.At)
		}
		return at
	}

	// 典籍二例：日柱為魁罡，故成立
	for _, c := range []struct {
		label string
		p     [4]ganzhi.SexagenaryIndex
		day   [2]int // 日干、日支
	}{
		{"張時僉事 庚午 丁亥 戊戌 丙辰",
			[4]ganzhi.SexagenaryIndex{pair(6, wu), pair(3, hai), pair(4, xu), pair(2, chen)}, [2]int{4, xu}},
		{"劉大受少卿 丁亥 癸丑 庚戌 戊寅",
			[4]ganzhi.SexagenaryIndex{pair(3, hai), pair(9, chou), pair(6, xu), pair(4, yin)}, [2]int{6, xu}},
	} {
		in := Input{
			DayStem: ganzhi.StemIndex(c.day[0]), DayBranch: ganzhi.BranchIndex(c.day[1]),
			YearStem: c.p[0].Stem(), YearBranch: c.p[0].Branch(), MonthBranch: c.p[1].Branch(),
		}
		at := kuiGangOf(in, c.p[:])
		if len(at) == 0 {
			t.Errorf("%s：日柱為魁罡卻無命中", c.label)
			continue
		}
		found := false
		for _, x := range at {
			if x == 2 {
				found = true
			}
		}
		if !found {
			t.Errorf("%s：命中於 %v，未含日柱", c.label, at)
		}
	}

	// 反面：日柱非魁罡，即使他柱有魁罡干支也不成立。
	// 戊辰 己未 丁亥 庚戌——時柱庚戌是魁罡四組之一，但日柱丁亥不是。
	in := Input{DayStem: 3, DayBranch: hai, YearStem: 4, YearBranch: chen, MonthBranch: wei}
	scan := []ganzhi.SexagenaryIndex{pair(4, chen), pair(5, wei), pair(3, hai), pair(6, xu)}
	if at := kuiGangOf(in, scan); len(at) != 0 {
		t.Errorf("日柱丁亥非魁罡，卻於 %v 報出命中——"+
			"「疊疊相逢」的前提是日柱已為魁罡，非任一柱各自成立", at)
	}
}
