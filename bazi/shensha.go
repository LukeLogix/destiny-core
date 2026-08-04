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

	// 大運：十步各為一柱
	steps := make([]ganzhi.SexagenaryIndex, len(c.Fortunes))
	for i, f := range c.Fortunes {
		steps[i] = f.Sexagenary
	}
	for i := range c.Fortunes {
		c.Fortunes[i].ShenSha = pickAt(shensha.Detect(in, steps, opt.ShenSha), i)
	}

	// 流年：每步各自成組，太歲另作為該年的基準
	for i := range c.Fortunes {
		years := make([]ganzhi.SexagenaryIndex, len(c.Fortunes[i].Years))
		for j, y := range c.Fortunes[i].Years {
			years[j] = y.Sexagenary
		}
		for j := range c.Fortunes[i].Years {
			yin := in
			yin.Annual = c.Fortunes[i].Years[j].Sexagenary
			yin.HasAnnual = true
			c.Fortunes[i].Years[j].ShenSha = pickAt(shensha.Detect(yin, years, opt.ShenSha), j)
		}
	}
}

// pickAt 自一次掃描的結果中取出落在第 at 個元素者，並把 At 歸零——
// 命中已掛在該柱的結構上，位置由所在處隱含，不需再帶索引。
func pickAt(hits []shensha.Hit, at int) []shensha.Hit {
	var out []shensha.Hit
	for _, h := range hits {
		if h.At == at {
			h.At = 0
			out = append(out, h)
		}
	}
	return out
}
