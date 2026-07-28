package calendar

import (
	"math"
	"time"
)

// GregorianJDN 格里曆日期對應的儒略日數（該日中午的儒略日整數部分）。
//
// 用於日柱推算——干支紀日是不間斷的循環，以儒略日取模即得。
func GregorianJDN(year, month, day int) int {
	a := (14 - month) / 12
	y := year + 4800 - a
	m := month + 12*a - 3
	return day + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// CivilJD 民用時刻的儒略日。
//
// 以絕對時刻計算，故同一瞬間的不同時區表示會得到相同結果。
func CivilJD(t time.Time) float64 {
	return float64(t.UTC().UnixNano())/1e9/secondsPerDay + 2440587.5
}

// obliquity 黃赤交角（弧度），t 為儒略世紀。
// Meeus 22.2 的平均值式，對本用途（均時差）已足夠。
func obliquity(t float64) float64 {
	const deg = math.Pi / 180
	e := 23.439291111 - 0.0130041667*t - 1.6667e-7*t*t + 5.02778e-7*t*t*t
	return e * deg
}

// meanSolarLongitude 太陽平黃經（弧度），t 為儒略世紀
func meanSolarLongitude(t float64) float64 {
	const deg = math.Pi / 180
	l := 280.4664567 + 360007.6982779*(t/10) + 0.03032028*(t/10)*(t/10)
	return math.Mod(l, 360) * deg
}

// EquationOfTime 均時差：真太陽時減平太陽時。
//
// 成因是地球公轉軌道為橢圓且自轉軸傾斜，兩者疊加使視太陽日長短不一。
// 全年振幅約 −14 至 +17 分鐘——這比多數地區的經度修正還大，
// 故真太陽時若只做經度而略過均時差，誤差反而更嚴重。
//
// 依 Meeus 28.3：E = L₀ − 0.0057183° − α + Δψ·cos ε
func EquationOfTime(jd float64) time.Duration {
	tt := (jd + DeltaTDays(jd-J2000) - J2000) / 36525 // UT → TT，再轉儒略世紀

	l0 := meanSolarLongitude(tt)
	lambda := SolarLongitude(tt, -1) // 視黃經，已含章動與光行差
	eps := obliquity(tt) + nutationObliquity(tt)

	// 由黃經求赤經
	alpha := math.Atan2(math.Cos(eps)*math.Sin(lambda), math.Cos(lambda))

	const aberrationCorr = 0.0057183 * math.Pi / 180
	e := l0 - aberrationCorr - alpha

	// 收斂到 ±π，避免跨 0/2π 時得到近一整圈的荒謬值
	for e > math.Pi {
		e -= 2 * math.Pi
	}
	for e < -math.Pi {
		e += 2 * math.Pi
	}

	// 弧度 → 時間：一整圈對應一日
	return time.Duration(e / (2 * math.Pi) * secondsPerDay * float64(time.Second))
}

// nutationObliquity 交角章動（弧度）。振幅僅約 9 角秒，
// 對均時差的影響不到 0.05 秒，取主要項即可。
func nutationObliquity(t float64) float64 {
	const deg = math.Pi / 180
	omega := (125.04452 - 1934.136261*t) * deg
	return 9.20 * math.Cos(omega) / arcsecPerRad
}
