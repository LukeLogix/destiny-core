package lang

import (
	"github.com/LukeLogix/destiny-core/bazi"
)

// LocalizedHiddenStem 藏干的文字形式
type LocalizedHiddenStem struct {
	Stem   string `json:"stem"`
	Type   string `json:"type"`
	TenGod string `json:"ten_god"`
}

// LocalizedPillar 一柱的文字形式
type LocalizedPillar struct {
	Sexagenary string                `json:"sexagenary"`
	Stem       string                `json:"stem"`
	Branch     string                `json:"branch"`
	StemTenGod string                `json:"stem_ten_god"`
	Hidden     []LocalizedHiddenStem `json:"hidden"`
	Terrain    string                `json:"terrain"`
	Sound      string                `json:"sound"`
}

// LocalizedStrength 旺衰結果的文字形式
type LocalizedStrength struct {
	Strategy   string             `json:"strategy"`
	Verdict    string             `json:"verdict"`
	Ratio      float64            `json:"ratio"`
	Elements   map[string]float64 `json:"elements"`
	Supporting float64            `json:"supporting"`
	Draining   float64            `json:"draining"`
	Reasons    []string           `json:"reasons"`
}

// LocalizedAnnualYear 流年的文字形式
type LocalizedAnnualYear struct {
	Year       int    `json:"year"`
	Sexagenary string `json:"sexagenary"`
	TenGod     string `json:"ten_god"`
}

// LocalizedFortune 大運的文字形式
type LocalizedFortune struct {
	Sexagenary string                `json:"sexagenary"`
	TenGod     string                `json:"ten_god"`
	StartAge   int                   `json:"start_age"`
	EndAge     int                   `json:"end_age"`
	StartYear  int                   `json:"start_year"`
	Years      []LocalizedAnnualYear `json:"years"`
}

// LocalizedChart 可直接序列化給前端或 LLM 的命盤。
type LocalizedChart struct {
	Locale    string `json:"locale"`
	BirthTime string `json:"birth_time"`
	Gender    string `json:"gender"`

	DayMaster string `json:"day_master"`

	Year  LocalizedPillar `json:"year"`
	Month LocalizedPillar `json:"month"`
	Day   LocalizedPillar `json:"day"`
	Hour  LocalizedPillar `json:"hour"`

	Strengths []LocalizedStrength `json:"strengths"`
	Consensus string              `json:"consensus"`

	Fortunes     []LocalizedFortune `json:"fortunes"`
	FortuneStart string             `json:"fortune_start"`

	// Warnings 不確定性提示。刻意作為結構化欄位交出——
	// LLM 最容易在沒把握處產生斬釘截鐵的敘述，把疑慮一併給它才能誠實表達保留。
	Warnings []string `json:"warnings,omitempty"`
}

const timeLayout = "2006-01-02 15:04:05 -07:00"

// Localize 將命盤轉為指定語言的文字形式。
func Localize(c *bazi.Chart, loc Locale) *LocalizedChart {
	b := Bundle(loc)

	out := &LocalizedChart{
		Locale:       string(loc),
		BirthTime:    c.Birth.Time.Format(timeLayout),
		Gender:       genderName(b, c.Birth.Gender),
		DayMaster:    b.Stem(c.DayMaster()),
		Year:         localizePillar(b, c.Year),
		Month:        localizePillar(b, c.Month),
		Day:          localizePillar(b, c.Day),
		Hour:         localizePillar(b, c.Hour),
		Consensus:    consensusName(b, c.StrengthConsensus()),
		FortuneStart: c.FortuneStart.Format(timeLayout),
	}

	for _, s := range c.Strengths {
		out.Strengths = append(out.Strengths, localizeStrength(b, s))
	}
	for _, f := range c.Fortunes {
		out.Fortunes = append(out.Fortunes, localizeFortune(b, f))
	}

	if c.Boundary.TermCritical {
		out.Warnings = append(out.Warnings, b.warnTerm)
	}
	if c.Boundary.HourCritical {
		out.Warnings = append(out.Warnings, b.warnHour)
	}
	return out
}

func localizePillar(b *bundle, p bazi.Pillar) LocalizedPillar {
	lp := LocalizedPillar{
		Sexagenary: b.Sexagenary(p.Sexagenary),
		Stem:       b.Stem(p.Stem),
		Branch:     b.Branch(p.Branch),
		StemTenGod: b.TenGod(p.StemTenGod),
		Terrain:    b.Terrain(p.Terrain),
		Sound:      b.Sound(p.Sound),
	}
	for _, h := range p.Hidden {
		lp.Hidden = append(lp.Hidden, LocalizedHiddenStem{
			Stem:   b.Stem(h.Stem),
			Type:   hiddenTypeName(b, h.Type),
			TenGod: b.TenGod(h.TenGod),
		})
	}
	return lp
}

func localizeStrength(b *bundle, s bazi.StrengthResult) LocalizedStrength {
	ls := LocalizedStrength{
		Strategy:   strategyName(s.StrategyID),
		Verdict:    b.Verdict(s.Verdict),
		Ratio:      s.Ratio,
		Supporting: s.Supporting,
		Draining:   s.Draining,
		Elements:   map[string]float64{},
	}
	for i, v := range s.Elements {
		ls.Elements[b.elements[i]] = v
	}
	for _, r := range s.Reasons {
		ls.Reasons = append(ls.Reasons, b.Reason(r.Code))
	}
	return ls
}

func localizeFortune(b *bundle, f bazi.DecadeFortune) LocalizedFortune {
	lf := LocalizedFortune{
		Sexagenary: b.Sexagenary(f.Sexagenary),
		TenGod:     b.TenGod(f.StemTenGod),
		StartAge:   f.StartAge,
		EndAge:     f.EndAge,
		StartYear:  f.StartYear,
	}
	for _, y := range f.Years {
		lf.Years = append(lf.Years, LocalizedAnnualYear{
			Year:       y.Year,
			Sexagenary: b.Sexagenary(y.Sexagenary),
			TenGod:     b.TenGod(y.StemTenGod),
		})
	}
	return lf
}

func hiddenTypeName(b *bundle, t bazi.HiddenStemType) string {
	names := [3]string{"本氣", "中氣", "餘氣"}
	if b == bundles[ZhCN] {
		names = [3]string{"本气", "中气", "余气"}
	}
	if int(t) < len(names) {
		return names[t]
	}
	return ""
}

func genderName(b *bundle, g bazi.Gender) string {
	switch g {
	case bazi.Male:
		return "男"
	case bazi.Female:
		return "女"
	default:
		return ""
	}
}

func consensusName(b *bundle, c bazi.Consensus) string {
	switch c {
	case bazi.ConsensusStrong:
		return b.Verdict(bazi.VerdictStrong)
	case bazi.ConsensusWeak:
		return b.Verdict(bazi.VerdictWeak)
	default:
		if b == bundles[ZhCN] {
			return "有分歧"
		}
		return "有分歧"
	}
}

func strategyName(id bazi.StrategyID) string {
	switch id {
	case bazi.StrategyWeighted:
		return "weighted"
	default:
		return "classical"
	}
}
