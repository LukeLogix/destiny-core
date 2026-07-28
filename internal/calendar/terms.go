package calendar

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// solarterms.bin 由 test/cmd/genbin 從 JPL golden JSON 產生。
// 以 go:embed 內嵌而非讀外部檔：部署只需單一執行檔，不會有資料檔漏帶或路徑解析失敗。
//
//go:embed solarterms.bin
var termData []byte

var (
	ErrYearOutOfRange = errors.New("calendar: 年份超出支援範圍")
	ErrNotJie         = errors.New("calendar: 節氣索引非「節」")
	ErrInvalidIndex   = errors.New("calendar: 節氣索引超出範圍")
)

const (
	termHeaderSize = 13 // magic(4) + version(1) + yearFrom(2) + yearTo(2) + count(4)
	termRecordSize = 12 // jd float64 + uncertainty float32
	jiePerYear     = 12 // 只存「節」，「氣」不影響四柱
	termsVersion   = 1
)

var (
	tableYearFrom int
	tableYearTo   int
)

// init 解析內嵌表的 header。
// 此處 panic 是刻意的——內嵌資料損壞屬於建置錯誤，不是執行期可恢復的狀況，
// 讓它在程式啟動時就失敗，好過在某次排盤時算出錯誤的命盤。
func init() {
	if len(termData) < termHeaderSize {
		panic("calendar: 內嵌節氣表過短")
	}
	if string(termData[0:4]) != "DCST" {
		panic("calendar: 內嵌節氣表 magic 不符")
	}
	if termData[4] != termsVersion {
		panic(fmt.Sprintf("calendar: 內嵌節氣表版本 %d，預期 %d", termData[4], termsVersion))
	}
	tableYearFrom = int(int16(binary.LittleEndian.Uint16(termData[5:7])))
	tableYearTo = int(int16(binary.LittleEndian.Uint16(termData[7:9])))
	count := int(binary.LittleEndian.Uint32(termData[9:13]))

	years := tableYearTo - tableYearFrom + 1
	if want := years * jiePerYear; count != want {
		panic(fmt.Sprintf("calendar: 節氣表筆數 %d，依 %d 年應為 %d", count, years, want))
	}
	if want := termHeaderSize + count*termRecordSize; len(termData) != want {
		panic(fmt.Sprintf("calendar: 節氣表長度 %d bytes，應為 %d", len(termData), want))
	}
}

// YearRange 回傳內嵌表支援的年份範圍
func YearRange() (from, to int) { return tableYearFrom, tableYearTo }

// LookupTerm 查詢節氣時刻。
//
// idx 須為「節」（奇數：1=小寒 3=立春 5=驚蟄 … 23=大雪）；「氣」不影響四柱，表中未收錄。
// 回傳 UTC 儒略日，以及該筆時刻的不確定度（秒）——後者隨年代放大，
// 供臨界判定使用，見 BoundaryFlags。
func LookupTerm(year, idx int) (jd float64, uncertaintySec float64, err error) {
	if idx < 0 || idx > 23 {
		return 0, 0, fmt.Errorf("%w: %d", ErrInvalidIndex, idx)
	}
	if idx%2 == 0 {
		return 0, 0, fmt.Errorf("%w: %d", ErrNotJie, idx)
	}
	if year < tableYearFrom || year > tableYearTo {
		return 0, 0, fmt.Errorf("%w: %d（支援 %d-%d）",
			ErrYearOutOfRange, year, tableYearFrom, tableYearTo)
	}

	n := (year-tableYearFrom)*jiePerYear + (idx-1)/2
	off := termHeaderSize + n*termRecordSize
	jd = math.Float64frombits(binary.LittleEndian.Uint64(termData[off : off+8]))
	uncertaintySec = float64(math.Float32frombits(
		binary.LittleEndian.Uint32(termData[off+8 : off+12])))
	return jd, uncertaintySec, nil
}
