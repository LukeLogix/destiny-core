// Package test 以獨立 module 存放交叉驗證，避免 tyme4go／lunar-go 這類
// 僅作「第二意見」的依賴污染下游使用者的 go.mod。
package test

import (
	"testing"
	"time"

	lcal "github.com/6tail/lunar-go/calendar"
	"github.com/6tail/tyme4go/tyme"
	"github.com/LukeLogix/destiny-core/bazi"
)

var cst = time.FixedZone("UTC+8", 8*3600)

// chartOf 排出本核心的命盤
func chartOf(t *testing.T, y, mo, d, h, mi int, opt bazi.Options) *bazi.Chart {
	t.Helper()
	c, err := bazi.Compute(bazi.Birth{
		Time:   time.Date(y, time.Month(mo), d, h, mi, 0, 0, cst),
		Gender: bazi.Male,
	}, opt)
	if err != nil {
		t.Fatalf("%d-%02d-%02d %02d:%02d Compute 失敗: %v", y, mo, d, h, mi, err)
	}
	return c
}

// pillarsOf 取本核心的四柱，以六十甲子索引表示
func pillarsOf(t *testing.T, y, mo, d, h, mi int, opt bazi.Options) [4]int {
	t.Helper()
	c := chartOf(t, y, mo, d, h, mi, opt)
	return [4]int{
		int(c.Year.Sexagenary), int(c.Month.Sexagenary),
		int(c.Day.Sexagenary), int(c.Hour.Sexagenary),
	}
}

// tymePillars 取 tyme4go 的四柱
func tymePillars(y, mo, d, h, mi int) [4]int {
	st, err := tyme.SolarTime{}.FromYmdHms(y, mo, d, h, mi, 0)
	if err != nil {
		return [4]int{-1, -1, -1, -1}
	}
	e := st.GetLunarHour().GetEightChar()
	return [4]int{
		e.GetYear().GetIndex(), e.GetMonth().GetIndex(),
		e.GetDay().GetIndex(), e.GetHour().GetIndex(),
	}
}

// lunarGoPillars 取 lunar-go 的四柱。sect=1 為早子時換日，與本核心預設一致。
func lunarGoPillars(y, mo, d, h, mi int) [4]string {
	e := lcal.NewSolar(y, mo, d, h, mi, 0).GetLunar().GetEightChar()
	e.SetSect(1)
	return [4]string{e.GetYear(), e.GetMonth(), e.GetDay(), e.GetTime()}
}

// 六十甲子索引 → 簡體干支字串，用於與 lunar-go 比對
var stems = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
var branches = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

func sexToString(idx int) string {
	return stems[idx%10] + branches[idx%12]
}

// TestAgainstTyme4go 本核心與 tyme4go 的四柱全量比對。
//
// 專挑節氣交界日與時辰邊界——這些是最容易出錯的位置。
func TestAgainstTyme4go(t *testing.T) {
	opt := bazi.Default() // 早子時換日，與 tyme4go 預設一致

	days := [][2]int{
		{2, 4}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {6, 6}, {7, 7},
		{8, 7}, {9, 8}, {10, 8}, {11, 7}, {12, 7}, {1, 6}, {6, 15},
	}
	hours := [][2]int{{0, 0}, {0, 59}, {5, 30}, {11, 0}, {12, 59}, {23, 0}, {23, 59}}

	total, hardFail, criticalDiff := 0, 0, 0
	for y := bazi.MinYear; y <= bazi.MaxYear; y++ {
		for _, md := range days {
			// 1900 年立春（2/4 13:51）之前屬 1899 干支年，超出支援範圍，跳過
			if y == bazi.MinYear && md[0] <= 2 {
				continue
			}
			for _, hm := range hours {
				total++
				c := chartOf(t, y, md[0], md[1], hm[0], hm[1], opt)
				got := [4]int{
					int(c.Year.Sexagenary), int(c.Month.Sexagenary),
					int(c.Day.Sexagenary), int(c.Hour.Sexagenary),
				}
				want := tymePillars(y, md[0], md[1], hm[0], hm[1])
				if got == want {
					continue
				}

				// 兩者的節氣時刻源自不同的 ΔT 模型，2040 年後可差達兩分鐘。
				// 出生時刻落在該不確定範圍內時，判定分歧屬預期，且本核心
				// 已將其標記為臨界——這正是 BoundaryFlags 存在的理由。
				if c.Boundary.TermCritical {
					criticalDiff++
					continue
				}
				hardFail++
				if hardFail <= 5 {
					t.Errorf("%d-%02d-%02d %02d:%02d 於非臨界情況下不符\n  本核心: %v\n  tyme4go: %v\n  距節 %.1f 秒，不確定度 %.1f 秒",
						y, md[0], md[1], hm[0], hm[1], got, want,
						c.Boundary.NearTermSeconds, c.Boundary.TermUncertaintySec)
				}
			}
		}
	}
	t.Logf("與 tyme4go 比對 %d 筆：完全一致 %d，臨界分歧 %d，非臨界不符 %d",
		total, total-criticalDiff-hardFail, criticalDiff, hardFail)
	if hardFail > 0 {
		t.Errorf("非臨界情況下有 %d 筆不符", hardFail)
	}
}

// TestAgainstLunarGo 與第二個獨立實作比對。
//
// 一套實作算錯，不會與另一套錯在同一處；兩套都對得上才有信心。
func TestAgainstLunarGo(t *testing.T) {
	opt := bazi.Default()

	total, mismatch := 0, 0
	for y := 1901; y <= 2099; y += 7 {
		for _, md := range [][2]int{{2, 4}, {5, 5}, {8, 7}, {11, 7}} {
			for _, hm := range [][2]int{{0, 30}, {12, 0}, {23, 30}} {
				total++
				got := pillarsOf(t, y, md[0], md[1], hm[0], hm[1], opt)
				want := lunarGoPillars(y, md[0], md[1], hm[0], hm[1])

				for i := range got {
					if sexToString(got[i]) != want[i] {
						mismatch++
						if mismatch <= 5 {
							t.Errorf("%d-%02d-%02d %02d:%02d 第 %d 柱不符：本核心 %s，lunar-go %s",
								y, md[0], md[1], hm[0], hm[1], i, sexToString(got[i]), want[i])
						}
						break
					}
				}
			}
		}
	}
	t.Logf("與 lunar-go 比對 %d 筆，不符 %d 筆", total, mismatch)
	if mismatch > 0 {
		t.Errorf("共 %d 筆不符", mismatch)
	}
}

// TestLateZiDayPillarAgainstLunarGo 晚子時的日柱須與 lunar-go 的 sect=2 一致。
//
// 只比對日柱，因為時柱存在已知的流派分歧——見 TestLateZiHourStemIsSelfConsistent。
func TestLateZiDayPillarAgainstLunarGo(t *testing.T) {
	opt := bazi.Default()
	opt.LateZiKeepsDay = true

	total, mismatch := 0, 0
	for y := 1950; y <= 2050; y += 11 {
		for _, md := range [][2]int{{3, 15}, {9, 20}} {
			for _, hm := range [][2]int{{23, 10}, {23, 50}} {
				total++
				got := pillarsOf(t, y, md[0], md[1], hm[0], hm[1], opt)

				e := lcal.NewSolar(y, md[0], md[1], hm[0], hm[1], 0).GetLunar().GetEightChar()
				e.SetSect(2) // 晚子時不換日
				if sexToString(got[2]) != e.GetDay() {
					mismatch++
					if mismatch <= 3 {
						t.Errorf("%d-%02d-%02d %02d:%02d 日柱不符：本核心 %s，lunar-go(sect2) %s",
							y, md[0], md[1], hm[0], hm[1], sexToString(got[2]), e.GetDay())
					}
				}
			}
		}
	}
	t.Logf("晚子時日柱比對 %d 筆，不符 %d 筆", total, mismatch)
	if mismatch > 0 {
		t.Errorf("共 %d 筆不符", mismatch)
	}
}

// TestLateZiHourStemIsSelfConsistent 晚子時的時干須由「當日」日干推出。
//
// 這是與 lunar-go 的已知分歧，且本核心刻意採不同做法：
//
// lunar-go 的 sect=2 日柱不換日，時柱卻仍用次日日干推——實測 1990-05-20 23:30
// 得日柱乙酉、時柱戊子，而戊子須由次日丙戌的日干才推得出，內部並不一致。
//
// 本核心認為，晚子時派既主張 23 時後仍屬今日，時干自然應以今日日干起，
// 故日柱與時柱同用一套口徑。此測試釘住這個選擇，避免日後被誤「修正」。
func TestLateZiHourStemIsSelfConsistent(t *testing.T) {
	opt := bazi.Default()
	opt.LateZiKeepsDay = true

	for y := 1950; y <= 2050; y += 7 {
		for _, md := range [][2]int{{3, 15}, {9, 20}} {
			p := pillarsOf(t, y, md[0], md[1], 23, 30, opt)
			dayStem := p[2] % 10
			hourStem := p[3] % 10
			hourBranch := p[3] % 12

			if hourBranch != 0 {
				t.Fatalf("%d-%02d-%02d 23:30 時支應為子，得 %d", y, md[0], md[1], hourBranch)
			}
			// 五鼠遁：子時起干為 (日干 mod 5) × 2
			if want := dayStem % 5 * 2; hourStem != want {
				t.Errorf("%d-%02d-%02d 23:30 時干為 %d，依當日日干 %d 應為 %d",
					y, md[0], md[1], hourStem, dayStem, want)
			}
		}
	}
}

// TestRoundTrip 往返一致性：由命盤反查公曆，必須包含原時刻。
//
// 這是不需任何外部基準的自我檢查——tyme4go 提供八字反查功能，
// 若本核心算出的八字反查不回原時刻，即表示某處推算有誤。
func TestRoundTrip(t *testing.T) {
	opt := bazi.Default()
	checked := 0

	for y := 1920; y <= 2080; y += 13 {
		for _, md := range [][2]int{{1, 20}, {4, 10}, {7, 15}, {10, 25}} {
			h, mi := 14, 0
			got := pillarsOf(t, y, md[0], md[1], h, mi, opt)

			ec, err := tyme.EightChar{}.New(
				sexToString(got[0]), sexToString(got[1]),
				sexToString(got[2]), sexToString(got[3]))
			if err != nil {
				t.Fatalf("%d-%02d-%02d 八字組合失敗: %v", y, md[0], md[1], err)
			}

			found := false
			for _, st := range ec.GetSolarTimes(y-1, y+1) {
				sd := st.GetSolarDay()
				if sd.GetYear() == y && sd.GetMonth() == md[0] && sd.GetDay() == md[1] {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%d-%02d-%02d %02d:%02d 的八字反查不回原日期", y, md[0], md[1], h, mi)
			}
			checked++
		}
	}
	t.Logf("往返驗證 %d 筆", checked)
}
