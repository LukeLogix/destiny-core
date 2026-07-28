package bazi

import (
	"errors"
	"time"

	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/internal/calendar"
)

// Gender 性別。大運順逆依「陽男陰女順排、陰男陽女逆排」而定，故為必填。
type Gender uint8

const (
	GenderUnset Gender = iota
	Male
	Female
)

// Birth 命主資料：關於「這個人」的事實。
//
// 與 Options 分開的理由——性別與經度是事實，換了就是另一個人；
// 計算口徑是觀點，同一人可用不同口徑各排一次。
type Birth struct {
	Time time.Time // 須帶 *time.Location，時區資訊由它攜帶

	// Longitude 出生地經度，東經為正。
	//
	// 用指標而非值：零值 0 是合法經度（格林威治），無法用來表達「未提供」。
	// 若用值型別，未填經度者會被靜默套用 0°E，對台北出生者造成 −6 分偏移，
	// 卡在時辰邊界即換柱且不報錯。
	Longitude *float64

	Gender Gender
}

// SolarTimeMode 真太陽時的修正程度。
//
// 真太陽時實為兩個獨立修正疊加，故為三段而非開關：
// 均時差振幅（±16 分）大於台北的經度修正（+6 分），
// 只做經度不做均時差是常見疏漏。
type SolarTimeMode uint8

const (
	WallClock     SolarTimeMode = iota // 直接用鐘面時間
	LongitudeOnly                      // 加經度時差
	TrueSolar                          // 再加均時差
)

// TerrainSect 十二長生的陰干順逆口徑
type TerrainSect uint8

const (
	// TerrainYinReverse 陰干逆行：乙長生在午（預設，傳統派）
	TerrainYinReverse TerrainSect = iota
	// TerrainSameBirth 陰陽同生同死：乙長生在亥
	TerrainSameBirth
)

// ChildLimitSect 起運流派
type ChildLimitSect uint8

const (
	ChildLimitDefault    ChildLimitSect = iota // 以秒換算
	ChildLimitChina95                          // 元亨利貞，以分換算
	ChildLimitLunarSect1                       // 按日與時辰數
	ChildLimitLunarSect2                       // 以分換算並續算時辰
)

// Options 計算口徑：關於「怎麼算」的選擇。
//
// 一律走 per-request 傳遞，不設套件級全域變數——後端服務中不同請求
// 可能採用不同口徑，全域狀態會造成互相污染。
type Options struct {
	// 時間口徑
	LateZiKeepsDay bool // true = 晚子時派，23 時仍算當日
	SolarTime      SolarTimeMode

	// 排盤口徑
	HiddenStem HiddenStemSect
	Terrain    TerrainSect
	ChildLimit ChildLimitSect

	// 干支關係
	Relation ganzhi.RelationOptions
}

// Default 主流口徑：早子時換日、鐘面時間、標準藏干、陰干逆行。
func Default() Options {
	return Options{
		LateZiKeepsDay: false,
		SolarTime:      WallClock,
		HiddenStem:     HiddenStemStandard,
		Terrain:        TerrainYinReverse,
		ChildLimit:     ChildLimitDefault,
		Relation:       ganzhi.DefaultRelationOptions(),
	}
}

// 支援範圍，與內嵌節氣表及交叉驗證的範圍一致——只承諾驗證過的區間。
var MinYear, MaxYear = calendar.YearRange()

var (
	ErrYearOutOfRange    = calendar.ErrYearOutOfRange
	ErrLongitudeMissing  = errors.New("bazi: 真太陽時修正需要經度")
	ErrLongitudeInvalid  = errors.New("bazi: 經度超出 ±180")
	ErrLongitudeMismatch = errors.New("bazi: 經度與時區不符")
	ErrGenderRequired    = errors.New("bazi: 大運計算需要性別")
)
