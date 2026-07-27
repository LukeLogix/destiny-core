// PoC：驗證「命理層自己寫、天文層只借節氣時刻」是否可行。
// 全程 index 運算，不出現任何中文字串——這正是 i18n 版本的核心前提。
package main

import (
	"fmt"
	"math"

	"github.com/6tail/tyme4go/tyme"
)

// ---- 唯一的外部依賴面：取節氣時刻(JD) 與 曆日轉JD ----
type AstroSource interface {
	// TermJD 回傳指定年份第 idx 個節氣的儒略日(0=冬至,3=立春,奇數為「節」)
	TermJD(year, idx int) float64
	// CivilJD 回傳民用日期時間的儒略日
	CivilJD(y, mo, d, h, mi, s int) float64
}

type tymeAstro struct{ cache map[[2]int]float64 }

func (t *tymeAstro) TermJD(year, idx int) float64 {
	k := [2]int{year, idx}
	if v, ok := t.cache[k]; ok {
		return v
	}
	v := tyme.SolarTerm{}.FromIndex(year, idx).GetJulianDay().GetDay()
	t.cache[k] = v
	return v
}
func (t *tymeAstro) CivilJD(y, mo, d, h, mi, s int) float64 {
	return tyme.JulianDay{}.FromYmdHms(y, mo, d, h, mi, s).GetDay()
}

// ---- 以下全部是「我們自己的命理層」，純 index，無文字 ----

type Options struct {
	// LateZiKeepsDay: true=晚子時派(23時仍算當日)；false=早子時派(23時進次日)
	LateZiKeepsDay bool
}

type Pillars struct{ Year, Month, Day, Hour int } // 皆為 0..59 甲子索引

// sixtyOf 由天干索引(0..9)與地支索引(0..11)組出六十甲子索引
func sixtyOf(gan, zhi int) int {
	d := ((zhi-gan)%12 + 12) % 12
	k := (5 * (d / 2)) % 6
	return (gan + 10*k) % 60
}

const dayJDNOffset = 49 // 由基準日校準得出，見 calibrate()

func Compute(a AstroSource, y, mo, d, h, mi int, opt Options) Pillars {
	jd := a.CivilJD(y, mo, d, h, mi, 0)

	// --- 年柱：以立春(idx=3)為界 ---
	gzYear := y
	if jd < a.TermJD(y, 3) {
		gzYear = y - 1
	}
	yearIdx := ((gzYear-4)%60 + 60) % 60

	// --- 月柱：以十二「節」為界，寅月起 ---
	// 節序: 立春3 驚蟄5 清明7 立夏9 芒種11 小暑13 立秋15 白露17 寒露19 立冬21 大雪23 小寒(次年1)
	m := 0
	for i := 11; i >= 0; i-- {
		var tjd float64
		if i == 11 {
			tjd = a.TermJD(gzYear+1, 1)
		} else {
			tjd = a.TermJD(gzYear, 3+i*2)
		}
		if jd >= tjd {
			m = i
			break
		}
	}
	yearGan := yearIdx % 10
	monthGan := ((yearGan%5)*2 + 2 + m) % 10
	monthZhi := (2 + m) % 12
	monthIdx := sixtyOf(monthGan, monthZhi)

	// --- 日柱：儒略日循環；子時換日口徑在此生效 ---
	jdn := int(math.Floor(jd + 0.5))
	if h >= 23 && !opt.LateZiKeepsDay {
		jdn++
	}
	dayIdx := ((jdn+dayJDNOffset)%60 + 60) % 60

	// --- 時柱：五鼠遁 ---
	hourZhi := ((h + 1) / 2) % 12
	hourGan := (dayIdx%10*2 + hourZhi) % 10
	hourIdx := sixtyOf(hourGan, hourZhi)

	return Pillars{yearIdx, monthIdx, dayIdx, hourIdx}
}

// ---- 驗證 ----

func tymePillars(y, mo, d, h, mi int) Pillars {
	t, _ := tyme.SolarTime{}.FromYmdHms(y, mo, d, h, mi, 0)
	e := t.GetLunarHour().GetEightChar()
	return Pillars{e.GetYear().GetIndex(), e.GetMonth().GetIndex(),
		e.GetDay().GetIndex(), e.GetHour().GetIndex()}
}

func calibrate(a AstroSource) {
	jd := a.CivilJD(2000, 1, 1, 12, 0, 0)
	jdn := int(math.Floor(jd + 0.5))
	want := tymePillars(2000, 1, 1, 12, 0).Day
	fmt.Printf("校準: 2000-01-01 JDN=%d, tyme日柱index=%d, 需要 offset=%d\n",
		jdn, want, ((want-jdn)%60+60)%60)
}

func main() {
	a := &tymeAstro{cache: map[[2]int]float64{}}
	calibrate(a)

	opt := Options{LateZiKeepsDay: false} // 對齊 tyme4go 預設
	total, bad := 0, 0
	var firstBad []string

	for y := 1900; y <= 2100; y++ {
		for mo := 1; mo <= 12; mo++ {
			for _, d := range []int{1, 4, 5, 6, 7, 8, 15, 21, 22, 23, 28} {
				if d > 28 {
					continue
				}
				for _, hm := range [][2]int{{0, 0}, {0, 59}, {6, 30}, {12, 59}, {13, 0}, {23, 0}, {23, 59}} {
					total++
					got := Compute(a, y, mo, d, hm[0], hm[1], opt)
					want := tymePillars(y, mo, d, hm[0], hm[1])
					if got != want {
						bad++
						if len(firstBad) < 6 {
							firstBad = append(firstBad, fmt.Sprintf(
								"  %d-%02d-%02d %02d:%02d  自算=%v  tyme=%v", y, mo, d, hm[0], hm[1], got, want))
						}
					}
				}
			}
		}
	}
	fmt.Printf("\n全量對照 %d 筆，不一致 %d 筆 (%.6f%%)\n",
		total, bad, float64(bad)*100/float64(total))
	for _, s := range firstBad {
		fmt.Println(s)
	}

	// 展示：同一時刻、兩種子時口徑，由 Options 控制（非全域變數）
	fmt.Println("\n=== per-request Options 而非全域狀態 ===")
	for _, o := range []Options{{LateZiKeepsDay: false}, {LateZiKeepsDay: true}} {
		p := Compute(a, 1990, 5, 20, 23, 30, o)
		fmt.Printf("  LateZiKeepsDay=%-5v -> 日柱index=%d 時柱index=%d\n",
			o.LateZiKeepsDay, p.Day, p.Hour)
	}
}
