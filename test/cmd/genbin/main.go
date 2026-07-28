// genbin 將 JPL golden JSON 轉為內嵌用的緊湊二進位表。
//
// 481 KB 的人類可讀 JSON 留作 golden data 與稽核用；
// 執行期用的是這裡產出的 29 KB 二進位檔，以 go:embed 打包進執行檔。
//
// 用法：go run ./cmd/genbin -in fixtures/solarterms_jpl.json -out ../internal/calendar/solarterms.bin
package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
)

// 12 個「節」的索引，順序即檔案中的排列順序
var jieIndices = []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}

type Term struct {
	Year         int     `json:"year"`
	Index        int     `json:"index"`
	Key          string  `json:"key"`
	JD           float64 `json:"jd"`
	TymeDeltaSec float64 `json:"tyme_delta_sec"`
}

type Fixture struct {
	YearFrom int    `json:"year_from"`
	YearTo   int    `json:"year_to"`
	Count    int    `json:"count"`
	Terms    []Term `json:"terms"`
}

func main() {
	in := flag.String("in", "fixtures/solarterms_jpl.json", "輸入 JSON")
	out := flag.String("out", "../internal/calendar/solarterms.bin", "輸出二進位")
	flag.Parse()

	b, err := os.ReadFile(*in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "讀取失敗: %v\n", err)
		os.Exit(1)
	}
	var f Fixture
	if err := json.Unmarshal(b, &f); err != nil {
		fmt.Fprintf(os.Stderr, "解析失敗: %v\n", err)
		os.Exit(1)
	}

	years := f.YearTo - f.YearFrom + 1
	want := years * len(jieIndices)
	if len(f.Terms) != want {
		fmt.Fprintf(os.Stderr, "筆數不符：JSON %d 筆，%d 年 × %d 節應為 %d 筆\n",
			len(f.Terms), years, len(jieIndices), want)
		os.Exit(1)
	}

	// 建索引：(year, idx) → Term，確保每格都有值，不容許缺漏
	type key struct{ y, i int }
	m := make(map[key]Term, len(f.Terms))
	for _, t := range f.Terms {
		m[key{t.Year, t.Index}] = t
	}

	jieOrder := make(map[int]int, len(jieIndices))
	order := append([]int(nil), jieIndices...)
	sort.Ints(order)
	for n, idx := range order {
		jieOrder[idx] = n
	}

	// header: magic(4) version(1) yearFrom(int16) yearTo(int16) count(uint32)
	buf := make([]byte, 0, 13+want*12)
	buf = append(buf, 'D', 'C', 'S', 'T', 1)
	buf = binary.LittleEndian.AppendUint16(buf, uint16(int16(f.YearFrom)))
	buf = binary.LittleEndian.AppendUint16(buf, uint16(int16(f.YearTo)))
	buf = binary.LittleEndian.AppendUint32(buf, uint32(want))

	for y := f.YearFrom; y <= f.YearTo; y++ {
		for _, idx := range order {
			t, ok := m[key{y, idx}]
			if !ok {
				fmt.Fprintf(os.Stderr, "缺少 %d 年索引 %d 的節氣\n", y, idx)
				os.Exit(1)
			}
			// 不確定度取絕對值：方向不重要，量級才是臨界判定要用的
			unc := t.TymeDeltaSec
			if unc < 0 {
				unc = -unc
			}
			buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(t.JD))
			buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(float32(unc)))
		}
	}

	if err := os.WriteFile(*out, buf, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "寫檔失敗: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("已寫出 %s\n  %d 年 × %d 節 = %d 筆，%d bytes\n",
		*out, years, len(jieIndices), want, len(buf))
}
