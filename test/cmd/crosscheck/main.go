package main

import (
	"fmt"

	lcal "github.com/6tail/lunar-go/calendar"
	"github.com/6tail/tyme4go/tyme"
)

func lunarGoPillars(y, mo, d, h, mi, sect int) string {
	e := lcal.NewSolar(y, mo, d, h, mi, 0).GetLunar().GetEightChar()
	e.SetSect(sect)
	return e.GetYear() + " " + e.GetMonth() + " " + e.GetDay() + " " + e.GetTime()
}

func tymePillars(y, mo, d, h, mi int) string {
	t, err := tyme.SolarTime{}.FromYmdHms(y, mo, d, h, mi, 0)
	if err != nil {
		return "ERR"
	}
	e := t.GetLunarHour().GetEightChar()
	return e.GetYear().GetName() + " " + e.GetMonth().GetName() + " " +
		e.GetDay().GetName() + " " + e.GetHour().GetName()
}

func run(sect int) (int, int, map[string]int) {
	days := [][2]int{{2, 4}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, {6, 6}, {7, 7},
		{8, 7}, {9, 8}, {10, 8}, {11, 7}, {12, 7}, {1, 6}}
	hours := [][2]int{{0, 0}, {0, 59}, {5, 30}, {11, 0}, {12, 59}, {23, 0}, {23, 59}}

	total, mismatch := 0, 0
	byHour := map[string]int{}
	for y := 1900; y <= 2100; y++ {
		for _, md := range days {
			for _, hm := range hours {
				total++
				a := lunarGoPillars(y, md[0], md[1], hm[0], hm[1], sect)
				b := tymePillars(y, md[0], md[1], hm[0], hm[1])
				if a != b {
					mismatch++
					byHour[fmt.Sprintf("%02d:%02d", hm[0], hm[1])]++
				}
			}
		}
	}
	return total, mismatch, byHour
}

func main() {
	for _, sect := range []int{2, 1} {
		total, mis, byHour := run(sect)
		label := "晚子時不換日"
		if sect == 1 {
			label = "早子時換日"
		}
		fmt.Printf("lunar-go sect=%d (%s) vs tyme4go 預設：\n", sect, label)
		fmt.Printf("  總比對 %d，差異 %d (%.4f%%)\n", total, mis, float64(mis)*100/float64(total))
		if len(byHour) > 0 {
			fmt.Printf("  差異分佈: %v\n", byHour)
		}
		fmt.Println()
	}

	// 具體驗證 23:30 這一刻，兩庫兩種口徑全排列
	fmt.Println("=== 1990-05-20 23:30 四種口徑對照 ===")
	fmt.Printf("  lunar-go sect=2 : %s\n", lunarGoPillars(1990, 5, 20, 23, 30, 2))
	fmt.Printf("  lunar-go sect=1 : %s\n", lunarGoPillars(1990, 5, 20, 23, 30, 1))
	fmt.Printf("  tyme4go 預設    : %s\n", tymePillars(1990, 5, 20, 23, 30))
	fmt.Printf("  (次日 00:30 對照): %s\n", tymePillars(1990, 5, 21, 0, 30))
}
