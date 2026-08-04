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
	c.ShenSha = shensha.Detect(in, natal, opt.ShenSha)

	if !opt.IncludeDynamicShenSha {
		return
	}

	// 每個動態柱各自掃描自己一柱即可——神煞是「基準對上這一柱」的關係，
	// 把整組丟進去再挑出目標那筆，會做十倍的白工。
	// 單元素掃描的另一好處是 At 恆為 0，正合 D-34：位置已由所在結構隱含。
	one := make([]ganzhi.SexagenaryIndex, 1)

	for i := range c.Fortunes {
		one[0] = c.Fortunes[i].Sexagenary
		c.Fortunes[i].ShenSha = shensha.Detect(in, one, opt.ShenSha)

		for j := range c.Fortunes[i].Years {
			// 太歲類神煞以該年干支為基準，故逐年重設
			yin := in
			yin.Annual = c.Fortunes[i].Years[j].Sexagenary
			yin.HasAnnual = true

			one[0] = c.Fortunes[i].Years[j].Sexagenary
			c.Fortunes[i].Years[j].ShenSha = shensha.Detect(yin, one, opt.ShenSha)
		}
	}
}
