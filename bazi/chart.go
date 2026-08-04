package bazi

import (
	"time"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// Pillar 一柱：干支及其衍生資訊
type Pillar struct {
	Sexagenary ganzhi.SexagenaryIndex
	Stem       ganzhi.StemIndex
	Branch     ganzhi.BranchIndex
	StemTenGod TenGod            // 天干對日主的十神；日柱為比肩（日主對自身）
	Hidden     []HiddenTenGod    // 地支藏干及各自十神
	Terrain    ganzhi.Terrain    // 日主在此柱地支的十二長生狀態
	Sound      ganzhi.SoundIndex // 納音
}

// BoundaryFlags 臨界標記。
//
// 本核心的差異化設計：不只算出答案，並誠實標記不確定性。
// 出生時間差一分鐘即換月柱的盤，整體解讀完全不同；
// 而 2040 年後節氣時刻本身就有分鐘級不可知性（ΔT 預測分歧）。
type BoundaryFlags struct {
	NearTermSeconds    float64 // 距最近的「節」多少秒，影響年柱與月柱
	NearHourSeconds    float64 // 距時辰分界多少秒，影響時柱
	TermUncertaintySec float64 // 該筆節氣時刻本身的不確定度，查表取得
	TermCritical       bool    // 距離已落入不確定範圍，年柱／月柱不可靠
	HourCritical       bool    // 同理，時柱不可靠
}

// Chart 一張八字命盤
type Chart struct {
	Birth   Birth
	Options Options // 攜帶口徑，保證任何一張盤可完整重現

	EffectiveTime   time.Time     // 真太陽時修正後，實際用於推算的時刻
	SolarTimeOffset time.Duration // 相對鐘面時間的總修正量，供稽核

	Year, Month, Day, Hour Pillar

	Relations    []ganzhi.Relation // 沖刑合會，只報構成不判定化成
	Boundary     BoundaryFlags
	Fortunes     []DecadeFortune // 大運，每步含其涵蓋的流年
	FortuneStart time.Time       // 起運時刻，四種流派可差達一整天，故明列供稽核
	Strengths    []StrengthResult
}

// DayMaster 日主，即日柱天干——十神皆以此為參照。
func (c *Chart) DayMaster() ganzhi.StemIndex { return c.Day.Stem }

// Pillars 依年月日時順序回傳四柱，便於逐柱處理。
func (c *Chart) Pillars() [4]Pillar {
	return [4]Pillar{c.Year, c.Month, c.Day, c.Hour}
}

// Branches 四柱地支，供關係計算使用。
func (c *Chart) Branches() []ganzhi.BranchIndex {
	p := c.Pillars()
	return []ganzhi.BranchIndex{p[0].Branch, p[1].Branch, p[2].Branch, p[3].Branch}
}

// Stems 四柱天干
func (c *Chart) Stems() []ganzhi.StemIndex {
	p := c.Pillars()
	return []ganzhi.StemIndex{p[0].Stem, p[1].Stem, p[2].Stem, p[3].Stem}
}
