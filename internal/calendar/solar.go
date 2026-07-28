// Package calendar 提供曆法與天文計算。
//
// 節氣時刻有兩條獨立路徑：
//   - 內嵌 JPL 表查詢（精度即 ground truth，範圍 1900-2100）
//   - 本檔的 VSOP87 計算（可算任意年份，用於驗證表的完整性）
//
// 兩者互為守護：表若損壞或被竄改，計算會發現；計算若移植有誤，表會發現。
package calendar

import "math"

// J2000 曆元的儒略日
const J2000 = 2451545.0

const (
	secondsPerDay = 86400.0
	// 1 弧度換算為角秒
	arcsecPerRad = 180 * 3600 / math.Pi
)

// earthLongitude 地球日心黃經（弧度，累積值），t 為儒略世紀（J2000 起算，TT）。
//
// n 控制取用的級數項數：n < 0 取全部（最高精度），n > 0 依比例截斷以加速。
// xl0 的佈局為 [0]=除數、[1..19]=各階級數的起訖索引、其後為 (振幅, 相位, 角速度) 三元組。
func earthLongitude(t float64, n int) float64 {
	t /= 10 // 級數以儒略千年為單位
	v, tn := 0.0, 1.0
	const pn = 1
	m0 := xl0[pn+1] - xl0[pn]

	for i := 0; i < 6; i++ {
		n1 := int(xl0[pn+i])
		n2 := int(xl0[pn+1+i])
		n0 := n2 - n1
		if n0 == 0 {
			continue
		}
		var m int
		if n < 0 {
			m = n2 // 取全部項
		} else {
			m = int(float64(3*n*n0)/m0+0.5) + n1
			if i != 0 {
				m += 3
			}
			if m > n2 {
				m = n2
			}
		}
		c := 0.0
		for j := n1; j < m; j += 3 {
			c += xl0[j] * math.Cos(xl0[j+1]+t*xl0[j+2])
		}
		v += c * tn
		tn *= t
	}
	v /= xl0[0]

	// 長期項修正
	t2 := t * t
	return v + (-0.0728-2.7702*t-1.1019*t2-0.0996*t2*t)/arcsecPerRad
}

// nutationLongitude 黃經章動（弧度），t 為儒略世紀
func nutationLongitude(t float64) float64 {
	a := -1.742 * t
	t2 := t * t
	dl := 0.0
	for i := 0; i < len(nutB); i += 5 {
		dl += (nutB[i+3] + a) * math.Sin(nutB[i]+nutB[i+1]*t+nutB[i+2]*t2)
		a = 0
	}
	return dl / 100 / arcsecPerRad
}

// aberration 太陽光行差（弧度），t 為儒略世紀
func aberration(t float64) float64 {
	t2 := t * t
	return -20.49552 *
		(1 + (0.016708634-0.000042037*t-0.0000001267*t2)*
			math.Cos(-0.043126+628.301955*t-0.000002732*t2)) / arcsecPerRad
}

// angularVelocity 太陽視黃經的角速度（弧度/儒略世紀），供迭代求解使用
func angularVelocity(t float64) float64 {
	f := 628.307585 * t
	return 628.332 + 21*math.Sin(1.527+f) + 0.44*math.Sin(1.48+f*2) +
		0.129*math.Sin(5.82+f)*t + 0.00055*math.Sin(4.21+f)*t*t
}

// SolarLongitude 太陽視黃經（弧度，累積值），t 為儒略世紀（J2000 起算，TT）。
//
// 地球日心黃經加 π 即太陽地心黃經，再疊加章動與光行差得視黃經。
func SolarLongitude(t float64, n int) float64 {
	return earthLongitude(t, n) + nutationLongitude(t) + aberration(t) + math.Pi
}

// SolarLongitudeInverse 求太陽視黃經達 w（弧度，累積值）的時刻，回傳儒略世紀（TT）。
//
// 先以平均角速度取初值，再兩次逐步提高級數精度修正——直接用全級數迭代
// 會在初值不佳時收斂緩慢。
func SolarLongitudeInverse(w float64) float64 {
	const meanVelocity = 628.3319653318
	t := (w - 1.75347 - math.Pi) / meanVelocity
	t += (w - SolarLongitude(t, 10)) / angularVelocity(t)
	return t + (w-SolarLongitude(t, -1))/angularVelocity(t)
}

// deltaTSeconds ΔT（TT − UT，秒），year 為含小數的西元年。
//
// 1900 年起有觀測表可查；表末之後改用長期外推式。
func deltaTSeconds(year float64) float64 {
	size := len(dtAt)
	y0 := dtAt[size-2] // 表中最後一段的起始年
	t0 := dtAt[size-1] // 該年的 ΔT

	if year >= y0 {
		const secular = 31.0 // 長期加速度係數
		ext := func(y float64) float64 {
			dy := (y - 1820) / 100
			return -20 + secular*dy*dy
		}
		if year > y0+100 {
			return ext(year)
		}
		// 表末 100 年內：外推式與表值的落差線性收斂，避免銜接處跳階
		return ext(year) - (ext(y0)-t0)*(y0+100-year)/100
	}

	i := 0
	for ; i < size; i += 5 {
		if year < dtAt[i+5] {
			break
		}
	}
	t1 := (year - dtAt[i]) / (dtAt[i+5] - dtAt[i]) * 10
	t2 := t1 * t1
	return dtAt[i+1] + dtAt[i+2]*t1 + dtAt[i+3]*t2 + dtAt[i+4]*t2*t1
}

// DeltaTDays ΔT（天），days 為 J2000 起算的天數
func DeltaTDays(days float64) float64 {
	return deltaTSeconds(days/365.2425+2000) / secondsPerDay
}

// SolarTermJD 求太陽視黃經達 w（弧度，累積值）的時刻，回傳 UTC 儒略日。
//
// 注意：不加時區偏移。壽星／tyme4go 的對應函式會加 1/3 天（北京時間 UTC+8），
// 本核心一律以 UTC 為準，時區轉換交由上層處理。
func SolarTermJD(w float64) float64 {
	days := SolarLongitudeInverse(w) * 36525 // 儒略世紀 → 天（TT）
	return days - DeltaTDays(days) + J2000   // TT → UT
}

// CumulativeLongitude 將「目標黃經（度）」與「近似儒略日」換算為累積黃經（弧度）。
//
// 太陽視黃經隨時間累積增長，同一個 15° 倍數每年重現一次，故需以近似時刻定位圈數。
func CumulativeLongitude(targetDeg, approxJD float64) float64 {
	// J2000 時太陽視黃經約 280.46°，其後每儒略世紀增加 36000.76°
	centuries := (approxJD - J2000) / 36525
	approxDeg := 280.46 + 36000.76*centuries
	turns := math.Round((approxDeg - targetDeg) / 360)
	return (targetDeg + 360*turns) * math.Pi / 180
}
