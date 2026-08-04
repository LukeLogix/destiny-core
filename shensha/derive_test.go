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
		if got := wenChangBranch(ganzhi.StemIndex(s)); got != w {
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
