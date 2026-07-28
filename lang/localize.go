package lang

import (
	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/ganzhi"
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

// LocalizedRelation 干支關係的文字形式
type LocalizedRelation struct {
	Kind    string   `json:"kind"`
	Pillars []string `json:"pillars"`
	// Transform 僅合會類有值，且語意為「若化則化為此五行」而非「已化」——
	// 化與不化取決於得令、引化之神、是否被沖破，各家分歧極大，核心不作此判定。
	Transform string `json:"transform,omitempty"`
}

// LocalizedBoundary 臨界數值。除了人看的警示文字，數值本身也要交出，
// 上層才能自行決定門檻。
type LocalizedBoundary struct {
	NearTermSeconds    float64 `json:"near_term_seconds"`
	NearHourSeconds    float64 `json:"near_hour_seconds"`
	TermUncertaintySec float64 `json:"term_uncertainty_seconds"`
	TermCritical       bool    `json:"term_critical"`
	HourCritical       bool    `json:"hour_critical"`
}

// LocalizedOptions 排盤所用口徑。命盤的可重現性靠它，
// 故不得在序列化時被丟棄。
type LocalizedOptions struct {
	LateZiKeepsDay bool   `json:"late_zi_keeps_day"`
	SolarTime      string `json:"solar_time"`
	HiddenStemSect uint8  `json:"hidden_stem_sect"`
	TerrainSect    uint8  `json:"terrain_sect"`
	ChildLimitSect uint8  `json:"child_limit_sect"`
	HalfTrinity    bool   `json:"include_half_trinity"`
	Destruction    bool   `json:"include_destruction"`
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

	Relations []LocalizedRelation `json:"relations"`

	// 稽核用欄位
	EffectiveTime   string             `json:"effective_time"`
	SolarTimeOffset string             `json:"solar_time_offset"`
	Boundary        *LocalizedBoundary `json:"boundary"`
	Options         *LocalizedOptions  `json:"options"`

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

		EffectiveTime:   c.EffectiveTime.Format(timeLayout),
		SolarTimeOffset: c.SolarTimeOffset.String(),
		Boundary: &LocalizedBoundary{
			NearTermSeconds:    c.Boundary.NearTermSeconds,
			NearHourSeconds:    c.Boundary.NearHourSeconds,
			TermUncertaintySec: c.Boundary.TermUncertaintySec,
			TermCritical:       c.Boundary.TermCritical,
			HourCritical:       c.Boundary.HourCritical,
		},
		Options: &LocalizedOptions{
			LateZiKeepsDay: c.Options.LateZiKeepsDay,
			SolarTime:      b.SolarTime(c.Options.SolarTime),
			HiddenStemSect: uint8(c.Options.HiddenStem),
			TerrainSect:    uint8(c.Options.Terrain),
			ChildLimitSect: uint8(c.Options.ChildLimit),
			HalfTrinity:    c.Options.Relation.IncludeHalfTrinity,
			Destruction:    c.Options.Relation.IncludeDestruction,
		},
	}

	for _, r := range c.Relations {
		out.Relations = append(out.Relations, localizeRelation(b, r))
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

// localizeRelation ganzhi 只回報「傳入 slice 的第幾個」，此處賦予年月日時的語意。
func localizeRelation(b *bundle, r ganzhi.Relation) LocalizedRelation {
	lr := LocalizedRelation{Kind: b.Relation(r.Kind)}
	for _, i := range r.Indices {
		lr.Pillars = append(lr.Pillars, b.Pillar(i))
	}
	if r.Transform != nil {
		lr.Transform = b.Element(*r.Transform)
	}
	return lr
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
