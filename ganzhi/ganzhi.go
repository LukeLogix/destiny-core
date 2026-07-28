// Package ganzhi 提供干支的基礎型別與運算。
//
// 本套件只認索引，不含任何文字——文字由 lang 套件負責，依賴方向為 lang → ganzhi，
// 反向不成立。此規則由編譯器強制，確保計算層永遠與語言無關。
//
// 天干索引 0-9 依序為甲乙丙丁戊己庚辛壬癸；
// 地支索引 0-11 依序為子丑寅卯辰巳午未申酉戌亥。
package ganzhi

import (
	"errors"
	"fmt"
)

type (
	// StemIndex 天干索引，0-9
	StemIndex uint8
	// BranchIndex 地支索引，0-11
	BranchIndex uint8
	// SexagenaryIndex 六十甲子索引，0-59
	SexagenaryIndex uint8
)

const (
	StemCount       = 10
	BranchCount     = 12
	SexagenaryCount = 60
)

var (
	ErrIndexOutOfRange  = errors.New("ganzhi: 索引超出範圍")
	ErrPolarityMismatch = errors.New("ganzhi: 干支陰陽不符")
)

// mod 恆回非負餘數。Go 的 % 對負數回負值，直接用於循環索引會出錯。
func mod(a, n int) int { return ((a % n) + n) % n }

// Next 天干循環進位，n 可為負
func (s StemIndex) Next(n int) StemIndex {
	return StemIndex(mod(int(s)+n, StemCount))
}

// Next 地支循環進位，n 可為負
func (b BranchIndex) Next(n int) BranchIndex {
	return BranchIndex(mod(int(b)+n, BranchCount))
}

// Next 六十甲子循環進位，n 可為負
func (x SexagenaryIndex) Next(n int) SexagenaryIndex {
	return SexagenaryIndex(mod(int(x)+n, SexagenaryCount))
}

// Stem 取六十甲子的天干
func (x SexagenaryIndex) Stem() StemIndex { return StemIndex(int(x) % StemCount) }

// Branch 取六十甲子的地支
func (x SexagenaryIndex) Branch() BranchIndex { return BranchIndex(int(x) % BranchCount) }

// Sexagenary 由天干地支組出六十甲子索引。
//
// 陽干只配陽支、陰干只配陰支，故 10×12 的組合中僅 60 種成立；
// 「甲丑」這類陰陽不符者回 ErrPolarityMismatch。
func Sexagenary(s StemIndex, b BranchIndex) (SexagenaryIndex, error) {
	if int(s) >= StemCount {
		return 0, fmt.Errorf("%w: 天干 %d", ErrIndexOutOfRange, s)
	}
	if int(b) >= BranchCount {
		return 0, fmt.Errorf("%w: 地支 %d", ErrIndexOutOfRange, b)
	}
	if int(s)%2 != int(b)%2 {
		return 0, fmt.Errorf("%w: 天干 %d 與地支 %d", ErrPolarityMismatch, s, b)
	}
	// 解同餘式 i ≡ s (mod 10)、i ≡ b (mod 12)。
	// 兩式的模互質部分為 2，故 d 必為偶數；k 為 10 在模 12 下的步進倍數。
	d := mod(int(b)-int(s), BranchCount)
	k := (5 * (d / 2)) % 6
	return SexagenaryIndex((int(s) + 10*k) % SexagenaryCount), nil
}

// HourBranch 由時鐘小時取時支。23-01 時為子，其後每兩小時遞進一支。
func HourBranch(hour int) BranchIndex {
	h := mod(hour, 24)
	return BranchIndex(mod((h+1)/2, BranchCount))
}

// HourStem 五鼠遁：由日干與時支推時干。
//
// 甲己還加甲、乙庚丙作初、丙辛從戊起、丁壬庚子居、戊癸壬子途——
// 即子時起干為 (日干 mod 5) × 2，其後隨時支遞進。
func HourStem(dayStem StemIndex, hourBranch BranchIndex) StemIndex {
	return StemIndex(mod(int(dayStem)%5*2+int(hourBranch), StemCount))
}

// MonthStem 五虎遁：由年干與月支推月干。
//
// 甲己之年丙作首、乙庚之歲戊為頭、丙辛必定尋庚起、丁壬壬位順行流、戊癸甲寅好追求——
// 即寅月起干為 (年干 mod 5) × 2 + 2，其後隨月支自寅順行。
func MonthStem(yearStem StemIndex, monthBranch BranchIndex) StemIndex {
	// 月支自寅（索引 2）起算為正月
	offset := mod(int(monthBranch)-2, BranchCount)
	return StemIndex(mod(int(yearStem)%5*2+2+offset, StemCount))
}
