package bazi

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// 1990-05-20 10:30 男命：庚午 辛巳 乙酉 辛巳。
// 日主乙木生於巳月（火旺），火洩木氣、金重剋木，四柱無水無根——命理上為典型身弱。

// TestClassicalThreeFactors 傳統三要素：得令、得地、得勢，三中取二為身強。
func TestClassicalThreeFactors(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	in := StrengthInput{Pillars: c.Pillars()}

	r := ClassicalStrength{}.Evaluate(in)

	// 乙木生巳月（火），木生火故為「休」，不得令
	// 四支午酉巳巳的藏干無木，不得地
	// 四干庚辛乙辛，比劫印星不足兩個，不得勢
	if r.Verdict != VerdictWeak {
		t.Errorf("此命三要素皆不具，應判身弱，實得 %d", r.Verdict)
	}
	if r.StrategyID != StrategyClassical {
		t.Errorf("策略識別為 %d，應為 Classical", r.StrategyID)
	}
	if len(r.Reasons) == 0 {
		t.Error("應附判定依據")
	}
}

// TestClassicalNeverNeutral 三要素為整數計數，本質上無中間地帶。
func TestClassicalNeverNeutral(t *testing.T) {
	dates := [][5]int{
		{1990, 5, 20, 10, 30}, {1984, 2, 15, 8, 0},
		{2000, 11, 3, 22, 0}, {1975, 7, 7, 3, 30},
	}
	for _, d := range dates {
		c := mustCompute(t, birthAt(d[0], d[1], d[2], d[3], d[4]), Default())
		r := ClassicalStrength{}.Evaluate(StrengthInput{Pillars: c.Pillars()})
		if r.Verdict == VerdictNeutral {
			t.Errorf("%v: Classical 不應產生近中和判定", d)
		}
	}
}

// TestWeightedScoring 加權計分：月支本氣權重最重，反映月令為重的原則。
func TestWeightedScoring(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	in := StrengthInput{Pillars: c.Pillars()}

	r := DefaultWeighted().Evaluate(in)

	if r.StrategyID != StrategyWeighted {
		t.Errorf("策略識別為 %d，應為 Weighted", r.StrategyID)
	}
	// 此命火金極旺、木僅日主一點、水全無
	if r.Elements[ganzhi.Water] != 0 {
		t.Errorf("四柱無水，水得分應為 0，實得 %.2f", r.Elements[ganzhi.Water])
	}
	if r.Elements[ganzhi.Fire] <= r.Elements[ganzhi.Wood] {
		t.Errorf("巳月火旺，火分 %.2f 應大於木分 %.2f",
			r.Elements[ganzhi.Fire], r.Elements[ganzhi.Wood])
	}
	// 幫身遠少於耗身
	if r.Supporting >= r.Draining {
		t.Errorf("幫身 %.2f 應遠少於耗身 %.2f", r.Supporting, r.Draining)
	}
	if r.Ratio >= 0.45 {
		t.Errorf("比值 %.3f 應低於 0.45，判為身弱", r.Ratio)
	}
	if r.Verdict != VerdictWeak {
		t.Errorf("應判身弱，實得 %d", r.Verdict)
	}
}

// TestWeightedNeutralBand 中間帶內不判強弱。
//
// 權重表本身無權威出處，比值落在 0.50 附近時判定完全取決於那組沒有依據的數字；
// 與其用假精確的門檻硬分兩邊，不如承認它就是邊界。
func TestWeightedNeutralBand(t *testing.T) {
	w := DefaultWeighted()
	if w.Neutral[0] != 0.45 || w.Neutral[1] != 0.55 {
		t.Fatalf("預設中間帶為 %v，應為 [0.45 0.55]", w.Neutral)
	}

	cases := []struct {
		ratio float64
		want  Verdict
		desc  string
	}{
		{0.60, VerdictStrong, "高於上界判身強"},
		{0.30, VerdictWeak, "低於下界判身弱"},
		{0.50, VerdictNeutral, "落在中間帶不判"},
		{0.45, VerdictNeutral, "下界含在中間帶內"},
		{0.55, VerdictNeutral, "上界含在中間帶內"},
	}
	for _, c := range cases {
		if got := w.verdictOf(c.ratio); got != c.want {
			t.Errorf("%s: 比值 %.2f 得 %d，應為 %d", c.desc, c.ratio, got, c.want)
		}
	}
}

// TestWeightedTableIsAdjustable 換流派只換權重資料，不動程式碼。
func TestWeightedTableIsAdjustable(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	in := StrengthInput{Pillars: c.Pillars()}

	base := DefaultWeighted().Evaluate(in)

	// 把月令權重拉到極高，比值應隨之改變
	alt := DefaultWeighted()
	alt.W.Main = [4]float64{1.5, 30.0, 1.5, 1.5}
	altered := alt.Evaluate(in)

	if base.Ratio == altered.Ratio {
		t.Error("調整權重表後比值未改變，權重可能未生效")
	}
}

// TestStrengthConsensusAgreement 兩策略一致時給出明確結論。
func TestStrengthConsensusAgreement(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	if len(c.Strengths) != 2 {
		t.Fatalf("預設應跑兩套策略，實得 %d 套", len(c.Strengths))
	}
	if got := c.StrengthConsensus(); got != ConsensusWeak {
		t.Errorf("兩策略皆判身弱，共識應為 ConsensusWeak，實得 %d", got)
	}
}

// TestStrengthConsensusRules 共識規則須對任意策略數有定義，且不採多數決。
//
// 三策略 2:1 時，有價值的訊息是「這是邊界命例」，
// 多數決會把分歧抹平成假的確定性。
func TestStrengthConsensusRules(t *testing.T) {
	cases := []struct {
		verdicts []Verdict
		want     Consensus
		desc     string
	}{
		{[]Verdict{VerdictStrong, VerdictStrong}, ConsensusStrong, "皆判強"},
		{[]Verdict{VerdictWeak, VerdictWeak}, ConsensusWeak, "皆判弱"},
		{[]Verdict{VerdictStrong, VerdictWeak}, ConsensusDisputed, "兩者相左"},
		{[]Verdict{VerdictStrong, VerdictNeutral}, ConsensusDisputed, "其一為近中和"},
		{[]Verdict{VerdictStrong, VerdictStrong, VerdictWeak}, ConsensusDisputed, "2:1 不採多數決"},
		{[]Verdict{VerdictWeak}, ConsensusWeak, "單一策略"},
		{nil, ConsensusDisputed, "無策略"},
	}
	for _, c := range cases {
		rs := make([]StrengthResult, len(c.verdicts))
		for i, v := range c.verdicts {
			rs[i] = StrengthResult{Verdict: v}
		}
		if got := consensusOf(rs); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestStrengthReasonsAreCoded 判定依據須為代碼而非文字，否則中文會滲入計算層。
func TestStrengthReasonsAreCoded(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	for _, r := range c.Strengths {
		if len(r.Reasons) == 0 {
			t.Errorf("策略 %d 未附依據", r.StrategyID)
		}
		for _, reason := range r.Reasons {
			if reason.Code >= reasonCodeCount {
				t.Errorf("依據代碼 %d 超出定義範圍", reason.Code)
			}
		}
	}
}

// TestStrengthInputAccessors 輸入只需四柱，不需完整命盤——
// 若吃 *Chart 會拿到尚未填入旺衰結果的半成品，契約模糊且難以單獨測試。
func TestStrengthInputAccessors(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	in := StrengthInput{Pillars: c.Pillars()}

	if in.DayMaster() != c.DayMaster() {
		t.Errorf("日主取值不符：%d vs %d", in.DayMaster(), c.DayMaster())
	}
	if in.MonthBranch() != c.Month.Branch {
		t.Errorf("月令取值不符：%d vs %d", in.MonthBranch(), c.Month.Branch)
	}
}

// TestCustomStrategies 可指定自訂策略組合。
func TestCustomStrategies(t *testing.T) {
	opt := Default()
	opt.Strengths = []Strength{ClassicalStrength{}}

	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)
	if len(c.Strengths) != 1 {
		t.Fatalf("指定單一策略，實得 %d 套", len(c.Strengths))
	}
	if c.Strengths[0].StrategyID != StrategyClassical {
		t.Errorf("策略為 %d，應為 Classical", c.Strengths[0].StrategyID)
	}
}
