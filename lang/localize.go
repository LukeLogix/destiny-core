package lang

import (
	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/shensha"
)

// LocalizedHiddenStem 藏干的文字形式
type LocalizedHiddenStem struct {
	Stem   string `json:"stem"`
	Type   Term   `json:"type"`
	TenGod Term   `json:"ten_god"`
}

// LocalizedPillar 一柱的文字形式
type LocalizedPillar struct {
	Sexagenary string                `json:"sexagenary"`
	Stem       string                `json:"stem"`
	Branch     string                `json:"branch"`
	StemTenGod Term                  `json:"stem_ten_god"`
	Hidden     []LocalizedHiddenStem `json:"hidden"`
	Terrain    Term                  `json:"terrain"`
	Sound      Term                  `json:"sound"`
}

// LocalizedStrength 旺衰結果的文字形式
type LocalizedStrength struct {
	Strategy   Term               `json:"strategy"`
	Verdict    Term               `json:"verdict"`
	Ratio      float64            `json:"ratio"`
	Elements   map[string]float64 `json:"elements"`
	Supporting float64            `json:"supporting"`
	Draining   float64            `json:"draining"`
	Reasons    []Term             `json:"reasons"`
}

// LocalizedAnnualYear 流年的文字形式
type LocalizedAnnualYear struct {
	Year       int    `json:"year"`
	Sexagenary string `json:"sexagenary"`
	TenGod     Term   `json:"ten_god"`

	// ShenSha 僅 bazi.Options.IncludeDynamicShenSha 為 true 時有值。
	// 以原局為基準、落在本流年柱者，故 At 恆為 0。
	ShenSha []LocalizedShenSha `json:"shen_sha,omitempty"`

	// SuiJunShenSha 十二歲君，方向相反——以本年太歲為基準，落在原局四柱，
	// At 為柱序。與 ShenSha 分欄，因兩者的 At 語意不同。
	SuiJunShenSha []LocalizedShenSha `json:"sui_jun_shen_sha,omitempty"`
}

// LocalizedFortune 大運的文字形式
type LocalizedFortune struct {
	Sexagenary string                `json:"sexagenary"`
	TenGod     Term                  `json:"ten_god"`
	StartAge   int                   `json:"start_age"`
	EndAge     int                   `json:"end_age"`
	StartYear  int                   `json:"start_year"`
	Years      []LocalizedAnnualYear `json:"years"`

	// ShenSha 僅 bazi.Options.IncludeDynamicShenSha 為 true 時有值
	ShenSha []LocalizedShenSha `json:"shen_sha,omitempty"`
}

// LocalizedRelation 干支關係的文字形式
type LocalizedRelation struct {
	Kind    Term   `json:"kind"`
	Pillars []Term `json:"pillars"`
	// Transform 僅合會類有值，且語意為「若化則化為此五行」而非「已化」——
	// 化與不化取決於得令、引化之神、是否被沖破，各家分歧極大，核心不作此判定。
	Transform *Term `json:"transform,omitempty"`
}

// LocalizedShenSha 一筆神煞的文字形式。
//
// Basis 與 Variant 一併交出：BranchBase 開「兩者皆查」時同一柱可自年支與
// 日支各命中一次，上層須能分辨來源；天乙的陽貴陰貴、羊刃的陽干陰干刃亦然。
type LocalizedShenSha struct {
	Kind  Term `json:"kind"`
	At    int  `json:"at"`
	Basis Term `json:"basis"`
	// Variant 用指標——omitempty 對 struct 無效，用值型別會讓無雙軌之分的
	// 神煞也帶一個 {"id":"none","name":""} 的雜訊欄位。
	Variant *Term `json:"variant,omitempty"`

	// Category 收錄分類，供上層分組呈現
	Category Term `json:"category"`
	// Tradition 體系來源。數個神煞名在不同體系指涉不同的東西，
	// 如文昌（紫微系／子平系）、大耗（元辰別名／十二歲君第七位）。
	Tradition Term `json:"tradition"`
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
	LateZiKeepsDay bool `json:"late_zi_keeps_day"`
	SolarTime      Term `json:"solar_time"`
	HiddenStemSect Term `json:"hidden_stem_sect"`
	TerrainSect    Term `json:"terrain_sect"`
	ChildLimitSect Term `json:"child_limit_sect"`
	HalfTrinity    bool `json:"include_half_trinity"`
	Destruction    bool `json:"include_destruction"`
}

// LocalizedChart 可直接序列化給前端或 LLM 的命盤。
type LocalizedChart struct {
	Locale    string `json:"locale"`
	BirthTime string `json:"birth_time"`
	Gender    Term   `json:"gender"`

	DayMaster string `json:"day_master"`

	Year  LocalizedPillar `json:"year"`
	Month LocalizedPillar `json:"month"`
	Day   LocalizedPillar `json:"day"`
	Hour  LocalizedPillar `json:"hour"`

	Strengths []LocalizedStrength `json:"strengths"`
	Consensus Term                `json:"consensus"`

	Relations []LocalizedRelation `json:"relations"`

	// ShenSha 原局四柱的神煞。以太歲為基準者只出現在流年層。
	ShenSha []LocalizedShenSha `json:"shen_sha,omitempty"`

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
			HiddenStemSect: term(c.Options.HiddenStem.ID(), b.hiddenStemSects[:], int(c.Options.HiddenStem)),
			TerrainSect:    term(c.Options.Terrain.ID(), b.terrainSects[:], int(c.Options.Terrain)),
			ChildLimitSect: term(c.Options.ChildLimit.ID(), b.childLimitSects[:], int(c.Options.ChildLimit)),
			HalfTrinity:    c.Options.Relation.IncludeHalfTrinity,
			Destruction:    c.Options.Relation.IncludeDestruction,
		},
	}

	for _, r := range c.Relations {
		out.Relations = append(out.Relations, localizeRelation(b, r))
	}
	out.ShenSha = localizeShenSha(b, c.ShenSha)

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
		t := b.Element(*r.Transform)
		lr.Transform = &t
	}
	return lr
}

func localizeShenSha(b *bundle, hits []shensha.Hit) []LocalizedShenSha {
	if len(hits) == 0 {
		return nil
	}
	out := make([]LocalizedShenSha, 0, len(hits))
	for _, h := range hits {
		out = append(out, LocalizedShenSha{
			Kind:      b.ShenSha(h.Kind),
			At:        h.At,
			Basis:     b.Basis(h.Basis),
			Variant:   variantTerm(b, h.Variant),
			Category:  b.Category(h.Kind.Category()),
			Tradition: b.Tradition(h.Kind.Tradition()),
		})
	}
	return out
}

// variantTerm 無雙軌之分者回 nil，讓該欄位自 JSON 中消失。
func variantTerm(b *bundle, v shensha.Variant) *Term {
	if v == shensha.VariantNone {
		return nil
	}
	t := b.Variant(v)
	return &t
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
		Strategy:   strategyName(b, s.StrategyID),
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
		ShenSha:    localizeShenSha(b, f.ShenSha),
	}
	for _, y := range f.Years {
		lf.Years = append(lf.Years, LocalizedAnnualYear{
			Year:          y.Year,
			Sexagenary:    b.Sexagenary(y.Sexagenary),
			TenGod:        b.TenGod(y.StemTenGod),
			ShenSha:       localizeShenSha(b, y.ShenSha),
			SuiJunShenSha: localizeShenSha(b, y.SuiJunShenSha),
		})
	}
	return lf
}

func hiddenTypeName(b *bundle, t bazi.HiddenStemType) Term {
	return term(t.ID(), b.hiddenTypes[:], int(t))
}

func genderName(b *bundle, g bazi.Gender) Term {
	return term(g.ID(), b.genders[:], int(g))
}

// consensusName 綜合結論。分歧本身即為訊號，故與身強身弱並列為一個分類值，
// 而非在文字裡混一句話——上層要據以分支時有穩定的 ID 可用。
func consensusName(b *bundle, c bazi.Consensus) Term {
	return term(c.ID(), b.consensuses[:], int(c))
}

func strategyName(b *bundle, id bazi.StrategyID) Term {
	return Term{ID: id.ID(), Name: strategyDisplay(b, id)}
}

// strategyDisplay 策略的顯示名。此前只給 ID（"weighted"），使用者看到的是
// 程式識別字而非人話，destiny-lab 因而得自備一份中文對照。
func strategyDisplay(b *bundle, id bazi.StrategyID) string {
	names := [2]string{"加權法", "傳統格局法"}
	if b == bundles[ZhCN] {
		names = [2]string{"加权法", "传统格局法"}
	}
	return pick(names[:], int(id))
}
