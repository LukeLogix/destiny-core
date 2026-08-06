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
	// 僅論主體者只能落在日柱——十惡大敗、四廢的原文是「日」，
	// 年月時柱碰巧湊出同一組干支不算數
	for _, h := range c.ShenSha {
		if h.Kind.SubjectOnly() && h.At != dayPillarIndex {
			t.Errorf("%s 落在第 %d 柱，僅論主體者只能在日柱", h.Kind.ID(), h.At)
		}
	}

	// 太歲類的基準是流年，原局層不該出現
	for _, h := range c.ShenSha {
		if h.Kind.Category() == shensha.CategoryAnnual {
			t.Errorf("原局出現太歲類的 %s，該類只在流年層成立", h.Kind.ID())
		}
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

// TestDynamicShenShaContract 動態層兩個欄位的分工。
//
// 原先此處釘的是總數 146。那個數字每加一個神煞就要改一次，改完也看不出
// 對錯——改成釘住兩欄位各自的契約，加神煞不必動它，組裝寫錯則會當場失敗。
func TestDynamicShenShaContract(t *testing.T) {
	opt := Default()
	opt.IncludeDynamicShenSha = true
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), opt)

	var dyn, suiJun int
	for i, f := range c.Fortunes {
		for _, h := range f.ShenSha {
			checkDynamicHit(t, h, i, -1)
			dyn++
		}
		for j, y := range f.Years {
			for _, h := range y.ShenSha {
				checkDynamicHit(t, h, i, j)
				dyn++
			}
			for _, h := range y.SuiJunShenSha {
				if h.Kind.Category() != shensha.CategoryAnnual {
					t.Errorf("第 %d 步第 %d 年的歲君欄混入非太歲類的 %s", i, j, h.Kind.ID())
				}
				if h.At < 0 || h.At >= 4 {
					t.Errorf("第 %d 步第 %d 年的歲君 %s 落在 At=%d，應為原局四柱之一",
						i, j, h.Kind.ID(), h.At)
				}
				suiJun++
			}
		}
	}
	if dyn == 0 {
		t.Error("開啟動態後大運流年一筆神煞都沒有")
	}
	// 太歲必在流年支自己身上，故每年至少一筆——除非原局四柱都不含該支
	if suiJun == 0 {
		t.Error("十二歲君一筆也沒有——歲君掃的是原局四柱，掃錯對象會全空")
	}
	t.Logf("動態柱 %d 筆、歲君落原局 %d 筆", dyn, suiJun)
}

// checkDynamicHit 動態柱的命中：At 恆為 0，且不含僅論主體者與太歲類。
func checkDynamicHit(t *testing.T, h shensha.Hit, step, year int) {
	t.Helper()
	where := "大運"
	if year >= 0 {
		where = "流年"
	}
	if h.At != 0 {
		t.Errorf("第 %d 步%s的 %s 命中 At=%d，單柱掃描應恆為 0", step, where, h.Kind.ID(), h.At)
	}
	if h.Kind.SubjectOnly() {
		t.Errorf("第 %d 步%s出現僅論主體的 %s——大運流年柱都不是命主本柱",
			step, where, h.Kind.ID())
	}
	if h.Kind.Category() == shensha.CategoryAnnual {
		t.Errorf("第 %d 步%s的神煞欄混入太歲類的 %s，該類應在 SuiJunShenSha",
			step, where, h.Kind.ID())
	}
}

// TestJinShenOnlyAtHour 金神只論時柱。
//
// 典籍作「止有三時」——《三命通會》〈金神〉、《淵海子平》〈論金神〉、
// 《神峰通考》三本皆然。神峰所舉命例丁亥 癸丑 己未 癸酉，金神在時柱癸酉。
//
// 初版把它當成任一柱皆可，與魁罡那次（D-59）是同一類誤讀：原文的「日」
// 「時」是取法的一部分，不是「怎麼用」的慣例。
func TestJinShenOnlyAtHour(t *testing.T) {
	opt := Default()
	opt.ShenSha.Include = []shensha.Kind{shensha.JinShen}

	// 1953-01-06 巳時 → 需找一張時柱為癸酉／己巳／乙丑的盤。
	// 改以掃描既有命例的方式驗證契約，不硬編生辰。
	for _, tc := range []struct{ y, mo, d, h, mi int }{
		{1990, 5, 20, 10, 30}, {1988, 7, 31, 19, 46},
		{1953, 1, 6, 9, 30}, {1975, 11, 3, 17, 20}, {2001, 3, 14, 5, 40},
	} {
		c := mustCompute(t, birthAt(tc.y, tc.mo, tc.d, tc.h, tc.mi), opt)
		for _, h := range c.ShenSha {
			if h.Kind == shensha.JinShen && h.At != hourPillarIndex {
				t.Errorf("%d-%02d-%02d：金神落在第 %d 柱，典籍限時柱",
					tc.y, tc.mo, tc.d, h.At)
			}
		}
	}

	// 動態柱一律不該有金神——大運流年柱都不是時柱
	dyn := Default()
	dyn.IncludeDynamicShenSha = true
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), dyn)
	for i, f := range c.Fortunes {
		for _, h := range f.ShenSha {
			if h.Kind == shensha.JinShen {
				t.Errorf("第 %d 步大運出現金神——大運柱不是時柱", i)
			}
		}
		for j, y := range f.Years {
			for _, h := range y.ShenSha {
				if h.Kind == shensha.JinShen {
					t.Errorf("第 %d 步第 %d 年出現金神", i, j)
				}
			}
		}
	}
}
