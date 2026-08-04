package destiny_test

import (
	"encoding/json"
	"fmt"
	"time"

	// 內嵌時區資料庫：scratch／alpine 之類的容器常無系統 tzdata，
	// 少了它 LoadLocation 會失敗，形成「本機跑得動、上線炸掉」。
	_ "time/tzdata"

	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/lang"
)

// taipei 取真實的 IANA 時區。
//
// 務必用 LoadLocation 而非 time.FixedZone("UTC+8", 8*3600)——後者會丟失歷史
// 夏令時間，台灣在 1945-1961、1974-1975、1979 年間行用過 UTC+9。
// 以 1975-04-05 18:05 出生為例，用 FixedZone 會算出庚辰月，
// 用 Asia/Taipei 才是正確的己卯月，整整差一柱，而臨界標記不會示警。
func taipei() *time.Location {
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	return loc
}

// 最基本的用法：給出生時間與性別，排出四柱。
func Example_basic() {
	chart, err := bazi.Compute(bazi.Birth{
		Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, taipei()),
		Gender: bazi.Male,
	}, bazi.Default())
	if err != nil {
		panic(err)
	}

	b := lang.Bundle(lang.ZhTW)
	fmt.Printf("%s %s %s %s\n",
		b.Sexagenary(chart.Year.Sexagenary),
		b.Sexagenary(chart.Month.Sexagenary),
		b.Sexagenary(chart.Day.Sexagenary),
		b.Sexagenary(chart.Hour.Sexagenary))
	fmt.Printf("日主：%s\n", b.Stem(chart.DayMaster()))

	// Output:
	// 庚午 辛巳 乙酉 辛巳
	// 日主：乙
}

// 十神與藏干：解盤最常用的資訊。
func Example_tenGods() {
	chart, _ := bazi.Compute(bazi.Birth{
		Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, taipei()),
		Gender: bazi.Male,
	}, bazi.Default())

	b := lang.Bundle(lang.ZhTW)
	fmt.Printf("年干十神：%s\n", b.TenGod(chart.Year.StemTenGod).Name)
	fmt.Print("月支藏干：")
	for _, h := range chart.Month.Hidden {
		fmt.Printf("%s(%s) ", b.Stem(h.Stem), b.TenGod(h.TenGod).Name)
	}
	fmt.Println()

	// Output:
	// 年干十神：正官
	// 月支藏干：丙(傷官) 庚(正官) 戊(正財)
}

// 口徑可配置：子時換日等分歧處由呼叫方決定，命盤會記錄所用口徑。
func Example_options() {
	born := time.Date(1990, 5, 20, 23, 30, 0, 0, taipei())
	b := lang.Bundle(lang.ZhTW)

	for _, lateZi := range []bool{false, true} {
		opt := bazi.Default()
		opt.LateZiKeepsDay = lateZi

		chart, _ := bazi.Compute(bazi.Birth{Time: born, Gender: bazi.Male}, opt)
		label := "早子時派（23 時進次日）"
		if lateZi {
			label = "晚子時派（23 時仍算當日）"
		}
		fmt.Printf("%s：日柱 %s\n", label, b.Sexagenary(chart.Day.Sexagenary))
	}

	// Output:
	// 早子時派（23 時進次日）：日柱 丙戌
	// 晚子時派（23 時仍算當日）：日柱 乙酉
}

// 真太陽時：依出生地經度修正，卡在時辰邊界時會換柱。
func Example_solarTime() {
	lon := 121.5 // 台北

	opt := bazi.Default()
	opt.SolarTime = bazi.TrueSolar

	chart, err := bazi.Compute(bazi.Birth{
		Time:      time.Date(1990, 5, 20, 10, 30, 0, 0, taipei()),
		Longitude: &lon,
		Gender:    bazi.Male,
	}, opt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("時刻前移 %d 分鐘餘\n", int(chart.SolarTimeOffset.Minutes()))

	// Output:
	// 時刻前移 9 分鐘餘
}

// 旺衰：兩套策略並行，判定不一致本身即為有用訊號。
func Example_strength() {
	chart, _ := bazi.Compute(bazi.Birth{
		Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, taipei()),
		Gender: bazi.Male,
	}, bazi.Default())

	lc := lang.Localize(chart, lang.ZhTW)
	for _, s := range lc.Strengths {
		fmt.Printf("%s：%s\n", s.Strategy.Name, s.Verdict.Name)
	}
	fmt.Printf("共識：%s\n", lc.Consensus.Name)

	// Output:
	// 加權法：身弱
	// 傳統格局法：身弱
	// 共識：身弱
}

// 不確定性標記：臨近節氣分界的命盤會帶警示，供上層或 AI 誠實表達保留。
func Example_uncertainty() {
	// 2024 立春為 16:27:08.694，此刻距分界約 29 秒
	chart, _ := bazi.Compute(bazi.Birth{
		Time:   time.Date(2024, 2, 4, 16, 26, 40, 0, taipei()),
		Gender: bazi.Male,
	}, bazi.Default())

	// 不直接斷言秒數：實際值 28.69 秒緊貼 %.0f 的捨入邊界，
	// 節氣資料只要有零點幾秒的變動就會讓這個範例無故失敗。
	fmt.Printf("距分界不到一分鐘：%v\n", chart.Boundary.NearTermSeconds < 60)
	fmt.Printf("臨界：%v\n", chart.Boundary.TermCritical)

	lc := lang.Localize(chart, lang.ZhTW)
	for _, w := range lc.Warnings {
		fmt.Println(w)
	}

	// Output:
	// 距分界不到一分鐘：true
	// 臨界：true
	// 出生時刻臨近節氣分界，年柱與月柱可能因幾分鐘之差而不同，建議確認出生時間
}

// 給 LLM 用：輸出結構化 JSON，含結論與可追溯的依據。
func Example_json() {
	chart, _ := bazi.Compute(bazi.Birth{
		Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, taipei()),
		Gender: bazi.Male,
	}, bazi.Default())

	lc := lang.Localize(chart, lang.ZhTW)
	out, _ := json.Marshal(struct {
		DayMaster string `json:"day_master"`
		Year      string `json:"year"`
		Consensus string `json:"consensus"`
	}{lc.DayMaster, lc.Year.Sexagenary, lc.Consensus.Name})

	fmt.Println(string(out))

	// Output:
	// {"day_master":"乙","year":"庚午","consensus":"身弱"}
}
