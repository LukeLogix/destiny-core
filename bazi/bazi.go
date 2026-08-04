package bazi

import (
	"fmt"
	"math"
	"time"

	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/internal/calendar"
)

// dayJDNOffset 儒略日數與六十甲子的對位常數。
// 由基準日校準：2000-01-01 的 JDN 為 2451545，日柱為甲午（索引 54）。
const dayJDNOffset = 49

// 十二個「節」的節氣索引，自立春起依序至大雪，末項為次年小寒。
// 月柱只依「節」分界，「氣」不參與。
var monthJie = [12]struct{ yearDelta, idx int }{
	{0, 3}, {0, 5}, {0, 7}, {0, 9}, {0, 11}, {0, 13},
	{0, 15}, {0, 17}, {0, 19}, {0, 21}, {0, 23}, {1, 1},
}

func mod(a, n int) int { return ((a % n) + n) % n }

// Compute 排出一張八字命盤。
func Compute(b Birth, opt Options) (*Chart, error) {
	if b.Gender == GenderUnset {
		return nil, ErrGenderRequired
	}
	if y := b.Time.Year(); y < MinYear || y > MaxYear {
		return nil, fmt.Errorf("%w: %d（支援 %d-%d）", ErrYearOutOfRange, y, MinYear, MaxYear)
	}

	offset, err := solarTimeOffset(b, opt)
	if err != nil {
		return nil, err
	}
	effective := b.Time.Add(offset)
	jd := calendar.CivilJD(effective)

	c := &Chart{
		Birth:           b,
		Options:         opt,
		EffectiveTime:   effective,
		SolarTimeOffset: offset,
	}

	// --- 年柱與月柱：以「節」分界 ---
	gzYear, monthOrder, err := locateByJie(effective, jd)
	if err != nil {
		return nil, err
	}

	yearSex := ganzhi.SexagenaryIndex(mod(gzYear-4, ganzhi.SexagenaryCount))
	yearStem := yearSex.Stem()

	monthBranch := ganzhi.BranchIndex(mod(2+monthOrder, ganzhi.BranchCount))
	monthStem := ganzhi.MonthStem(yearStem, monthBranch)
	monthSex, err := ganzhi.Sexagenary(monthStem, monthBranch)
	if err != nil {
		return nil, fmt.Errorf("bazi: 月柱組合失敗: %w", err)
	}

	// --- 日柱：干支紀日為不間斷循環，以儒略日數取模 ---
	ey, em, ed := effective.Date()
	jdn := calendar.GregorianJDN(ey, int(em), ed)
	if effective.Hour() >= 23 && !opt.LateZiKeepsDay {
		jdn++ // 早子時派：23 時起算次日
	}
	daySex := ganzhi.SexagenaryIndex(mod(jdn+dayJDNOffset, ganzhi.SexagenaryCount))
	dayStem := daySex.Stem()

	// --- 時柱：五鼠遁 ---
	hourBranch := ganzhi.HourBranch(effective.Hour())
	hourStem := ganzhi.HourStem(dayStem, hourBranch)
	hourSex, err := ganzhi.Sexagenary(hourStem, hourBranch)
	if err != nil {
		return nil, fmt.Errorf("bazi: 時柱組合失敗: %w", err)
	}

	c.Year = makePillar(yearSex, dayStem, opt)
	c.Month = makePillar(monthSex, dayStem, opt)
	c.Day = makePillar(daySex, dayStem, opt)
	c.Hour = makePillar(hourSex, dayStem, opt)

	// --- 干支關係 ---
	c.Relations = append(
		ganzhi.StemRelations(c.Stems(), opt.Relation),
		ganzhi.BranchRelations(c.Branches(), opt.Relation)...,
	)

	// --- 臨界標記 ---
	c.Boundary = boundaryFlags(effective, jd, gzYear)

	// --- 大運與流年 ---
	c.Fortunes = computeFortunes(c, jd, gzYear)

	// --- 旺衰：預設同時跑兩套策略，判定不一致本身即為訊號 ---
	strategies := opt.Strengths
	if strategies == nil {
		strategies = []Strength{DefaultWeighted(), ClassicalStrength{}}
	}
	// 解析後的策略寫回命盤所記錄的 Options——WeightedStrength 是值型別且內含
	// 完整權重表，故此舉同時滿足「命盤需記錄所用權重」。若只留 nil，
	// 日後有人改動 DefaultWeights 全域變數，歸檔的命盤會重播出不同結論而無跡可循。
	c.Options.Strengths = strategies

	in := StrengthInput{Pillars: c.Pillars()}
	for _, s := range strategies {
		c.Strengths = append(c.Strengths, s.Evaluate(in))
	}

	return c, nil
}

func makePillar(sex ganzhi.SexagenaryIndex, dayMaster ganzhi.StemIndex, opt Options) Pillar {
	stem, branch := sex.Stem(), sex.Branch()
	return Pillar{
		Sexagenary: sex,
		Stem:       stem,
		Branch:     branch,
		StemTenGod: TenGodOf(dayMaster, stem),
		Hidden:     HiddenTenGods(branch, dayMaster, opt.HiddenStem),
		Terrain:    ganzhi.TerrainOf(dayMaster, branch, opt.Terrain),
		Sound:      ganzhi.SoundOf(sex),
	}
}

// solarTimeOffset 計算真太陽時相對鐘面時間的修正量。
//
// 算法為 UTC + 經度/15 小時 + 均時差，不引入「時區標準經線」這個中間量——
// 後者在夏令時間會整個偏掉：西班牙冬季 UTC+1（推得 15°E）、夏季 UTC+2（推得 30°E），
// 但國土位置並未移動。
func solarTimeOffset(b Birth, opt Options) (time.Duration, error) {
	if opt.SolarTime == WallClock {
		return 0, nil
	}
	if b.Longitude == nil {
		return 0, ErrLongitudeMissing
	}
	lon := *b.Longitude
	if lon < -180 || lon > 180 {
		return 0, fmt.Errorf("%w: %.2f", ErrLongitudeInvalid, lon)
	}

	_, tzOffsetSec := b.Time.Zone()
	// 經度換算的當地平太陽時，與鐘面時間的差
	diffHours := lon/15 - float64(tzOffsetSec)/3600

	// 正規化到 (-12, 12]。反子午線兩側的時區偏移與經度符號相反，
	// 未正規化會得到約 ±24 小時的假差值——吉里巴斯萊恩群島（-157.5°E, UTC+14）、
	// 薩摩亞、東加、查塔姆群島皆會被誤判為填錯而拒絕。
	for diffHours > 12 {
		diffHours -= 24
	}
	for diffHours <= -12 {
		diffHours += 24
	}

	// 門檻 4 小時：涵蓋全球所有合法組合（最極端為中國全境單一時區，
	// 新疆達 2.9 小時），同時攔下「台北時間配紐約經度」這類填錯（12.9 小時）。
	if math.Abs(diffHours) > 4 {
		return 0, fmt.Errorf("%w: 經度 %.2f 與時區偏移 %d 秒相差 %.1f 小時",
			ErrLongitudeMismatch, lon, tzOffsetSec, diffHours)
	}

	offset := time.Duration(diffHours * float64(time.Hour))
	if opt.SolarTime == TrueSolar {
		offset += calendar.EquationOfTime(calendar.CivilJD(b.Time))
	}
	return offset, nil
}

// locateByJie 依「節」定出干支年與月序（0 為寅月）。
func locateByJie(t time.Time, jd float64) (gzYear, monthOrder int, err error) {
	y := t.Year()

	lichun, _, err := calendar.LookupTerm(y, 3)
	if err != nil {
		return 0, 0, err
	}
	gzYear = y
	if jd < lichun {
		gzYear = y - 1
	}
	if gzYear < MinYear {
		return 0, 0, fmt.Errorf("%w: 干支年 %d 早於支援起點 %d（立春前出生歸前一干支年）",
			ErrYearOutOfRange, gzYear, MinYear)
	}

	// 自後往前找第一個已到達的節
	for i := len(monthJie) - 1; i >= 0; i-- {
		ref := monthJie[i]
		tjd, _, e := calendar.LookupTerm(gzYear+ref.yearDelta, ref.idx)
		if e != nil {
			continue // 次年小寒可能超出表範圍，退回前一個節即可
		}
		if jd >= tjd {
			return gzYear, i, nil
		}
	}
	// 理論上不可達：干支年的起點即立春，jd 必不小於它
	return gzYear, 0, nil
}

// boundaryFlags 計算臨界標記。
//
// 閾值不採固定值，而是隨該筆節氣的實測不確定度調整——
// 2030 年的盤距分界 3 秒才算臨界，2090 年的盤要差 2 分鐘。
func boundaryFlags(t time.Time, jd float64, gzYear int) BoundaryFlags {
	const birthTimePrecision = 60.0 // 出生時間記錄通常僅精確到分鐘

	f := BoundaryFlags{NearTermSeconds: math.Inf(1)}

	// 掃描本干支年前後的節，取最近者
	for _, delta := range []int{-1, 0, 1} {
		for _, ref := range monthJie {
			tjd, unc, err := calendar.LookupTerm(gzYear+delta+ref.yearDelta, ref.idx)
			if err != nil {
				continue
			}
			if d := math.Abs(jd-tjd) * 86400; d < f.NearTermSeconds {
				f.NearTermSeconds = d
				f.TermUncertaintySec = unc
			}
		}
	}
	if math.IsInf(f.NearTermSeconds, 1) {
		f.NearTermSeconds = 0
	}
	f.TermCritical = f.NearTermSeconds <= f.TermUncertaintySec+birthTimePrecision

	// 時辰分界在奇數整點（1、3、…、23 時）
	secOfDay := float64(t.Hour()*3600 + t.Minute()*60 + t.Second())
	const twoHours = 7200.0
	// 子時自 23 時起算，故先平移一小時再取模
	shifted := math.Mod(secOfDay+3600, 86400)
	within := math.Mod(shifted, twoHours)
	f.NearHourSeconds = math.Min(within, twoHours-within)
	f.HourCritical = f.NearHourSeconds <= birthTimePrecision

	return f
}
