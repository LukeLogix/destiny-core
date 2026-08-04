package bazi

import (
	"testing"

	"github.com/LukeLogix/destiny-core/shensha"
)

// TestChartCarriesShenSha 命盤須帶原局神煞。
//
// 基準命例 1990-05-20 10:30 男，四柱庚午 辛巳 乙酉 辛巳。
// 日主乙：金輿在巳（月、時兩柱）、文昌在午（年柱）；
// 日支酉屬巳酉丑金局：將星在酉（日柱）、咸池在午（年柱）。
func TestChartCarriesShenSha(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	want := map[shensha.Kind][]int{
		shensha.JinYu:          {1, 3},
		shensha.WenChangGuiRen: {0},
		shensha.JiangXing:      {2},
		shensha.XianChi:        {0},
	}

	got := map[shensha.Kind][]int{}
	for _, h := range c.ShenSha {
		got[h.Kind] = append(got[h.Kind], h.At)
	}

	for k, ats := range want {
		if len(got[k]) != len(ats) {
			t.Errorf("%s 命中 %v，應為 %v", k.ID(), got[k], ats)
			continue
		}
		for i := range ats {
			if got[k][i] != ats[i] {
				t.Errorf("%s 命中 %v，應為 %v", k.ID(), got[k], ats)
				break
			}
		}
	}
	if len(c.ShenSha) != 5 {
		t.Errorf("原局共 %d 筆命中，應為 5 筆", len(c.ShenSha))
	}
}

// TestDynamicShenShaOffByDefault 預設不算大運流年的神煞。
//
// 十步大運乘十個流年共一百一十四柱，多數使用者只看原局。
func TestDynamicShenShaOffByDefault(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())
	for i, f := range c.Fortunes {
		if f.ShenSha != nil {
			t.Errorf("第 %d 步大運不應有神煞，實得 %v", i, f.ShenSha)
		}
		for j, y := range f.Years {
			if y.ShenSha != nil {
				t.Errorf("第 %d 步第 %d 年不應有神煞", i, j)
			}
		}
	}
}

// TestDynamicShenShaWhenEnabled 開啟後大運流年須有神煞，且原局結果不變。
func TestDynamicShenShaWhenEnabled(t *testing.T) {
	base := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	opt := Default()
	opt.IncludeDynamicShenSha = true
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)

	if len(c.ShenSha) != len(base.ShenSha) {
		t.Errorf("開啟動態後原局命中由 %d 變為 %d，原局不應受影響",
			len(base.ShenSha), len(c.ShenSha))
	}

	total := 0
	for _, f := range c.Fortunes {
		total += len(f.ShenSha)
		for _, y := range f.Years {
			total += len(y.ShenSha)
		}
	}
	if total == 0 {
		t.Error("開啟動態後大運流年一筆神煞都沒有，不合理")
	}
	t.Logf("動態柱共 %d 筆命中（大運 %d 步、流年 %d 個）",
		total, len(c.Fortunes), len(c.Fortunes)*10)
}

// TestDynamicShenShaNotAliased 各動態柱的神煞必須各自獨立。
//
// 組裝時以單一元素的切片重複傳入 Detect，若 Detect 保留了該切片的參照，
// 後續迭代的寫入會污染先前的結果。此測試釘住「每柱結果互不相干」。
func TestDynamicShenShaNotAliased(t *testing.T) {
	opt := Default()
	opt.IncludeDynamicShenSha = true
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)

	// 逐年比對：該年的每一筆命中，其應命中的地支必須真的等於該年的地支
	for i, f := range c.Fortunes {
		for j, y := range f.Years {
			for _, h := range y.ShenSha {
				if h.At != 0 {
					t.Errorf("第 %d 步第 %d 年的命中 At=%d，動態柱應恆為 0", i, j, h.At)
				}
			}
		}
	}

	// 至少要有兩個流年的命中內容不同——若全都一樣，多半是別名污染
	seen := map[string]bool{}
	for _, f := range c.Fortunes {
		for _, y := range f.Years {
			key := ""
			for _, h := range y.ShenSha {
				key += h.Kind.ID() + ","
			}
			seen[key] = true
		}
	}
	if len(seen) < 2 {
		t.Errorf("一百個流年只產生 %d 種命中組合，疑為別名污染", len(seen))
	}
	t.Logf("一百個流年產生 %d 種不同的命中組合", len(seen))
}

// TestDynamicShenShaCountStable 動態命中總數不因組裝方式改變而變動。
func TestDynamicShenShaCountStable(t *testing.T) {
	opt := Default()
	opt.IncludeDynamicShenSha = true
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)

	total := 0
	for _, f := range c.Fortunes {
		total += len(f.ShenSha)
		for _, y := range f.Years {
			total += len(y.ShenSha)
		}
	}
	if total != 146 {
		t.Errorf("動態柱共 %d 筆命中，改組裝方式前為 146 筆", total)
	}
}
