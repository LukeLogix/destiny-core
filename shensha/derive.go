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
		for b := 0; b < ganzhi.BranchCount; b++ {
			if ganzhi.TerrainOf(ganzhi.StemIndex(s), ganzhi.BranchIndex(b), ganzhi.TerrainYinReverse) == ganzhi.Officer {
				out[s] = ganzhi.BranchIndex(b)
				break
			}
		}
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
func wenChangBranch(stem ganzhi.StemIndex) ganzhi.BranchIndex {
	want := ganzhi.Sick
	if stem.Polarity() == ganzhi.Yin {
		want = ganzhi.LongLife
	}
	for b := 0; b < ganzhi.BranchCount; b++ {
		if ganzhi.TerrainOf(stem, ganzhi.BranchIndex(b), ganzhi.TerrainYinReverse) == want {
			return ganzhi.BranchIndex(b)
		}
	}
	return 0 // 不可達：十二長生必然涵蓋十二支
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
	for s < ganzhi.StemCount {
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
