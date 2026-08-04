package bazi

import (
	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/shensha"
)

// 神煞的組裝層。
//
// shensha 套件本身位置無關——只認「基準組」與「待掃組」，不知道年月日時、
// 大運流年為何物。此處負責分次呼叫並賦予位置語意，如此六壬的爻位、紫微的
// 宮位未來可直接重用同一套偵測。

// shenShaInput 由命盤取出偵測所需的基準值。
//
// 基準一律取自原局：神煞的語意是「命主帶什麼」，大運流年只是被掃描的對象。
// 若允許大運當基準，基準數乘掃描數使輸出量平方成長，且「大運的桃花」與
// 「命主的桃花」語意不同，混在一起容易誤讀。
func (c *Chart) shenShaInput() shensha.Input {
	return shensha.Input{
		DayStem:     c.Day.Stem,
		YearStem:    c.Year.Stem,
		DayBranch:   c.Day.Branch,
		YearBranch:  c.Year.Branch,
		MonthBranch: c.Month.Branch,
		IsMale:      c.Birth.Gender == Male,
		YearSound:   c.Year.Sound,
	}
}

// detectShenSha 填入原局四柱的神煞；口徑開啟時另填大運與流年。
func (c *Chart) detectShenSha(opt Options) {
	in := c.shenShaInput()

	natal := make([]ganzhi.SexagenaryIndex, 4)
	for i, p := range c.Pillars() {
		natal[i] = p.Sexagenary
	}
	c.ShenSha = subjectAtDay(shensha.Detect(in, natal, opt.ShenSha))

	if !opt.IncludeDynamicShenSha {
		return
	}

	// 每個動態柱各自掃描自己一柱即可——神煞是「基準對上這一柱」的關係，
	// 把整組丟進去再挑出目標那筆，會做十倍的白工。
	// 單元素掃描的另一好處是 At 恆為 0，正合 D-34：位置已由所在結構隱含。
	one := make([]ganzhi.SexagenaryIndex, 1)

	// 十二歲君的方向相反：基準是流年，掃的是原局四柱。若照其餘神煞的作法
	// 只掃流年那一柱，太歲永遠命中自己、其餘十一位永遠落空——看似有算，
	// 實則什麼也沒說。故另行偵測、另欄存放。
	// 兩趟都用 Exclude 縮限而不改寫 Categories／Include——那三者是使用者的
	// 口徑設定，覆寫掉會讓「我明明排除了」失效。
	suiJun := opt.ShenSha
	suiJun.Exclude = withExcluded(opt.ShenSha.Exclude, func(k shensha.Kind) bool {
		return k.Category() != shensha.CategoryAnnual
	})

	// 動態柱那趟不含太歲類——歲君已由上一趟以原局為掃描對象算過。留著只會
	// 讓「太歲」每年命中流年支自己，是套套邏輯而非資訊。
	dynamic := opt.ShenSha
	dynamic.Exclude = withExcluded(opt.ShenSha.Exclude, func(k shensha.Kind) bool {
		return k.SubjectOnly() || k.Category() == shensha.CategoryAnnual
	})

	for i := range c.Fortunes {
		one[0] = c.Fortunes[i].Sexagenary
		c.Fortunes[i].ShenSha = shensha.Detect(in, one, dynamic)

		for j := range c.Fortunes[i].Years {
			// 太歲類神煞以該年干支為基準，故逐年重設
			yin := in
			yin.Annual = c.Fortunes[i].Years[j].Sexagenary
			yin.HasAnnual = true

			one[0] = c.Fortunes[i].Years[j].Sexagenary
			c.Fortunes[i].Years[j].ShenSha = shensha.Detect(yin, one, dynamic)
			c.Fortunes[i].Years[j].SuiJunShenSha = shensha.Detect(yin, natal, suiJun)
		}
	}
}

// subjectAtDay 濾掉落在非日柱的「僅論主體」神煞。
//
// shensha 位置無關，只報構成；八字的主體是日柱，故十惡大敗、四廢出現在
// 年月時柱不算數——那只是碰巧同一組干支。
func subjectAtDay(hits []shensha.Hit) []shensha.Hit {
	out := hits[:0]
	for _, h := range hits {
		if h.Kind.SubjectOnly() && h.At != dayPillarIndex {
			continue
		}
		out = append(out, h)
	}
	return out
}

// dayPillarIndex 日柱在 Chart.Pillars 中的序位
const dayPillarIndex = 2

// withExcluded 在既有排除清單上再排除符合條件者，不動原切片。
func withExcluded(base []shensha.Kind, drop func(shensha.Kind) bool) []shensha.Kind {
	out := append([]shensha.Kind(nil), base...)
	for k := shensha.Kind(0); k < shensha.KindCount; k++ {
		if drop(k) {
			out = append(out, k)
		}
	}
	return out
}
