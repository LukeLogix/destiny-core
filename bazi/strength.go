package bazi

import "github.com/LukeLogix/destiny-core/ganzhi"

// Verdict 旺衰判定
type Verdict uint8

const (
	VerdictStrong  Verdict = iota // 身強
	VerdictWeak                   // 身弱
	VerdictNeutral                // 近中和，不判強弱；僅連續計分的策略會產生
)

// StrategyID 策略識別
type StrategyID uint8

const (
	StrategyWeighted StrategyID = iota
	StrategyClassical
)

// ReasonCode 判定依據的代碼。
//
// 刻意用代碼而非字串：文字會使中文滲入計算層，i18n 立即破功。
// 附帶價值是 LLM 取得的成了「結論 + 可追溯依據」，而非黑箱結論。
type ReasonCode uint8

const (
	ReasonInSeason     ReasonCode = iota // 得令：日主五行當月令之旺或相
	ReasonNotInSeason                    // 不得令
	ReasonRooted                         // 得地：地支藏干有同五行者
	ReasonNoRoot                         // 不得地
	ReasonAllied                         // 得勢：天干比劫印星足數
	ReasonNotAllied                      // 不得勢
	ReasonSupportScore                   // 幫身總分
	ReasonDrainScore                     // 耗身總分

	reasonCodeCount
)

// Reason 一條判定依據及其權重貢獻
type Reason struct {
	Code   ReasonCode
	Weight float64
}

// StrengthResult 旺衰判定結果，過程完全攤開供稽核。
type StrengthResult struct {
	StrategyID StrategyID
	Elements   [ganzhi.ElementCount]float64 // 依 ganzhi.Element 順序：木火土金水
	Supporting float64                      // 幫身：印 + 比劫
	Draining   float64                      // 耗身：食傷 + 財 + 官殺
	Ratio      float64                      // Supporting / (Supporting + Draining)
	Verdict    Verdict
	Reasons    []Reason
}

// StrengthInput 旺衰判定所需的最小輸入。
//
// 刻意不吃 *Chart：Chart 內含 Strengths，若吃 *Chart 則拿到的必然是
// 尚未填入旺衰結果的半成品，契約模糊且難以單獨測試。
type StrengthInput struct {
	Pillars [4]Pillar // 索引 0..3 = 年月日時
}

func (in StrengthInput) DayMaster() ganzhi.StemIndex     { return in.Pillars[2].Stem }
func (in StrengthInput) MonthBranch() ganzhi.BranchIndex { return in.Pillars[1].Branch }

// Strength 旺衰策略。真正有派別差異的只有這一層，故做成可插拔。
type Strength interface {
	ID() StrategyID
	Evaluate(in StrengthInput) StrengthResult
}

// ---- Weighted：加權計分 ----

// WeightTable 各位置的權重。多數流派的差異其實只是權重不同而非演算法不同，
// 故換流派只需換這張表，不必動程式碼。
type WeightTable struct {
	Stem     [4]float64 // 年月日時 天干
	Main     [4]float64 // 地支本氣
	Middle   [4]float64 // 地支中氣
	Residual [4]float64 // 地支餘氣
}

// DefaultWeights 月支本氣 3.0 反映「月令為重」。
//
// 此組數字無權威出處，係實務慣例綜合，定位為「一個可用的預設」而非正確答案，
// 故必須可調，且命盤需記錄所用權重。
var DefaultWeights = WeightTable{
	Stem:     [4]float64{1.0, 1.0, 1.0, 1.0},
	Main:     [4]float64{1.5, 3.0, 1.5, 1.5},
	Middle:   [4]float64{0.7, 1.4, 0.7, 0.7},
	Residual: [4]float64{0.3, 0.6, 0.3, 0.3},
}

// DefaultNeutral 近中和的判定區間，同樣無權威出處。
var DefaultNeutral = [2]float64{0.45, 0.55}

// WeightedStrength 加權計分策略
type WeightedStrength struct {
	W       WeightTable
	Neutral [2]float64
}

// DefaultWeighted 以預設權重與中間帶建立策略
func DefaultWeighted() WeightedStrength {
	return WeightedStrength{W: DefaultWeights, Neutral: DefaultNeutral}
}

func (WeightedStrength) ID() StrategyID { return StrategyWeighted }

// verdictOf 依比值判定。中間帶內不判強弱——
// 權重表既無權威出處，比值貼近 0.5 時的判定就完全取決於那組沒有依據的數字。
func (w WeightedStrength) verdictOf(ratio float64) Verdict {
	switch {
	case ratio > w.Neutral[1]:
		return VerdictStrong
	case ratio < w.Neutral[0]:
		return VerdictWeak
	default:
		return VerdictNeutral
	}
}

func (w WeightedStrength) Evaluate(in StrengthInput) StrengthResult {
	var r StrengthResult
	r.StrategyID = StrategyWeighted

	for i, p := range in.Pillars {
		r.Elements[p.Stem.Element()] += w.W.Stem[i]
		for _, h := range p.Hidden {
			var weight float64
			switch h.Type {
			case PrimaryQi:
				weight = w.W.Main[i]
			case MiddleQi:
				weight = w.W.Middle[i]
			default:
				weight = w.W.Residual[i]
			}
			r.Elements[h.Stem.Element()] += weight
		}
	}

	me := in.DayMaster().Element()
	// 幫身：同我（比劫）與生我（印）
	r.Supporting = r.Elements[me] + r.Elements[me.GeneratedBy()]
	// 耗身：我生（食傷）、我剋（財）、剋我（官殺）
	r.Draining = r.Elements[me.Generates()] +
		r.Elements[me.Restrains()] + r.Elements[me.RestrainedBy()]

	if total := r.Supporting + r.Draining; total > 0 {
		r.Ratio = r.Supporting / total
	}
	r.Verdict = w.verdictOf(r.Ratio)
	r.Reasons = []Reason{
		{ReasonSupportScore, r.Supporting},
		{ReasonDrainScore, r.Draining},
	}
	return r
}

// ---- Classical：傳統三要素 ----

// ClassicalStrength 得令、得地、得勢，三中取二為身強。
type ClassicalStrength struct{}

func (ClassicalStrength) ID() StrategyID { return StrategyClassical }

func (ClassicalStrength) Evaluate(in StrengthInput) StrengthResult {
	r := StrengthResult{StrategyID: StrategyClassical}

	me := in.DayMaster()
	meElem := me.Element()
	monthElem := in.MonthBranch().Element()

	// 得令：旺相休囚死中，日主處於「旺」（同月令）或「相」（月令所生）
	inSeason := meElem == monthElem || meElem == monthElem.Generates()

	// 得地：四支藏干中有與日主同五行者，即通根
	rooted := false
	for _, p := range in.Pillars {
		for _, h := range p.Hidden {
			if h.Stem.Element() == meElem {
				rooted = true
			}
		}
	}

	// 得勢：四干中比劫印星合計兩個以上（不含日主本身）
	allies := 0
	for i, p := range in.Pillars {
		if i == 2 {
			continue // 日干即日主本身
		}
		e := p.Stem.Element()
		if e == meElem || e == meElem.GeneratedBy() {
			allies++
		}
	}
	allied := allies >= 2

	score := 0
	add := func(ok bool, yes, no ReasonCode) {
		if ok {
			score++
			r.Reasons = append(r.Reasons, Reason{yes, 1})
		} else {
			r.Reasons = append(r.Reasons, Reason{no, 0})
		}
	}
	add(inSeason, ReasonInSeason, ReasonNotInSeason)
	add(rooted, ReasonRooted, ReasonNoRoot)
	add(allied, ReasonAllied, ReasonNotAllied)

	// 三要素為整數計數，本質上無中間地帶，故不產生 VerdictNeutral
	if score >= 2 {
		r.Verdict = VerdictStrong
	} else {
		r.Verdict = VerdictWeak
	}
	r.Ratio = float64(score) / 3
	return r
}

// ---- 共識 ----

// Consensus 多策略的共同結論
type Consensus uint8

const (
	ConsensusStrong Consensus = iota
	ConsensusWeak
	ConsensusDisputed // 策略間不一致，或有策略判為近中和
)

// StrengthConsensus 綜合各策略的結論
func (c *Chart) StrengthConsensus() Consensus { return consensusOf(c.Strengths) }

// consensusOf 刻意不採多數決。
//
// 三策略 2:1 時，真正有價值的訊息是「這是邊界命例」而非「多數說了算」——
// 多數決會把分歧抹平成假的確定性，違背標記不確定性的初衷。
func consensusOf(rs []StrengthResult) Consensus {
	if len(rs) == 0 {
		return ConsensusDisputed
	}
	first := rs[0].Verdict
	for _, r := range rs {
		if r.Verdict == VerdictNeutral || r.Verdict != first {
			return ConsensusDisputed
		}
	}
	if first == VerdictStrong {
		return ConsensusStrong
	}
	return ConsensusWeak
}
