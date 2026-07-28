package bazi

import (
	"math"
	"time"

	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/internal/calendar"
)

// DecadeFortune 大運，每步十年。
type DecadeFortune struct {
	Sexagenary ganzhi.SexagenaryIndex
	StemTenGod TenGod
	Hidden     []HiddenTenGod
	StartAge   int
	EndAge     int
	StartYear  int
	Forward    bool         // 順排或逆排
	Years      []AnnualYear // 本步涵蓋的流年
}

// AnnualYear 流年——以干支年為準。
//
// 刻意與「小運」分成不同概念：tyme4go 的 Fortune 型別實為小運，
// 1995 年小運為丁亥而流年為乙亥，混用即算錯。本核心本版不提供小運。
type AnnualYear struct {
	Year       int
	Sexagenary ganzhi.SexagenaryIndex
	StemTenGod TenGod
}

// fortuneSteps 大運步數，涵蓋約百年。
const fortuneSteps = 10

// computeFortunes 排大運與流年。
//
// 順逆依「陽男陰女順排、陰男陽女逆排」；起運時距為出生至界節——
// 順排取下一個節，逆排取上一個節。
func computeFortunes(c *Chart, jd float64, gzYear int) []DecadeFortune {
	forward := isForward(c.Year.Stem, c.Birth.Gender)

	gapSeconds := jieGapSeconds(jd, gzYear, forward)
	startYear, startMonth, startDay, extra := convertGap(gapSeconds, c.Options.ChildLimit)

	// 起運時刻 = 出生時刻加上換算所得的年月日與零頭
	startTime := c.EffectiveTime.AddDate(startYear, startMonth, startDay).Add(extra)
	c.FortuneStart = startTime

	me := c.DayMaster()
	out := make([]DecadeFortune, 0, fortuneSteps)
	for i := 1; i <= fortuneSteps; i++ {
		step := i
		if !forward {
			step = -i
		}
		sex := c.Month.Sexagenary.Next(step)

		// 首步自起運當年起算，其後每十年一步
		sy := startTime.Year() + (i-1)*10
		age := sy - c.Birth.Time.Year() + 1 // 虛歲

		f := DecadeFortune{
			Sexagenary: sex,
			StemTenGod: TenGodOf(me, sex.Stem()),
			Hidden:     HiddenTenGods(sex.Branch(), me, c.Options.HiddenStem),
			StartAge:   age,
			EndAge:     age + 9,
			StartYear:  sy,
			Forward:    forward,
			Years:      make([]AnnualYear, 0, 10),
		}
		for k := 0; k < 10; k++ {
			y := sy + k
			ys := ganzhi.SexagenaryIndex(mod(y-4, ganzhi.SexagenaryCount))
			f.Years = append(f.Years, AnnualYear{
				Year:       y,
				Sexagenary: ys,
				StemTenGod: TenGodOf(me, ys.Stem()),
			})
		}
		out = append(out, f)
	}
	return out
}

// isForward 陽男陰女順排、陰男陽女逆排
func isForward(yearStem ganzhi.StemIndex, g Gender) bool {
	yang := yearStem.Polarity() == ganzhi.Yang
	return (yang && g == Male) || (!yang && g == Female)
}

// jieGapSeconds 出生時刻與界節的時距（秒）。
// 順排取下一個節，逆排取上一個節。
func jieGapSeconds(jd float64, gzYear int, forward bool) float64 {
	best := math.Inf(1)
	for _, delta := range []int{-1, 0, 1} {
		for _, ref := range monthJie {
			tjd, _, err := calendar.LookupTerm(gzYear+delta+ref.yearDelta, ref.idx)
			if err != nil {
				continue
			}
			gap := (tjd - jd) * 86400
			if forward && gap >= 0 && gap < best {
				best = gap
			}
			if !forward && gap <= 0 && -gap < best {
				best = -gap
			}
		}
	}
	if math.IsInf(best, 1) {
		return 0
	}
	return best
}

// convertGap 依流派將出生至界節的時距換算為起運的年月日與零頭。
//
// 四套規則皆為純算術，差別在換算的精細程度與取整方式——
// 這正是四家起運時刻可差達一整天的原因，故必須可選而非寫死。
func convertGap(seconds float64, sect ChildLimitSect) (year, month, day int, extra time.Duration) {
	s := int(math.Abs(seconds))

	switch sect {
	case ChildLimitChina95:
		// 元亨利貞：以分計，4320 分折一年、360 分折一月、12 分折一日。
		// 不再細算時分，故零頭為零。
		m := s / 60
		year = m / 4320
		m %= 4320
		month = m / 360
		m %= 360
		day = m / 12

	case ChildLimitLunarSect1:
		// 按日與時辰數：一日折四月、一時辰折十日
		days := s / 86400
		hourUnits := (s % 86400) / 7200
		monthAdd := hourUnits * 10 / 30
		month = days*4 + monthAdd
		day = hourUnits*10 - monthAdd*30
		year = month / 12
		month -= year * 12

	case ChildLimitLunarSect2:
		// 同元亨利貞，但續算時辰：餘數每 1 分折 2 時
		m := s / 60
		year = m / 4320
		m %= 4320
		month = m / 360
		m %= 360
		day = m / 12
		m %= 12
		extra = time.Duration(m*2) * time.Hour

	default: // ChildLimitDefault：以秒計，三日折一年，並續算時、分
		year = s / 259200
		s %= 259200
		month = s / 21600
		s %= 21600
		day = s / 720
		s %= 720
		hour := s / 30
		s %= 30
		extra = time.Duration(hour)*time.Hour + time.Duration(s*2)*time.Minute
	}
	return year, month, day, extra
}
