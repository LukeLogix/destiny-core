package shensha

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// 基準命例：1990-05-20 10:30 台北男，四柱庚午 辛巳 乙酉 辛巳。
// 與 core README 及既有測試同一例，便於交叉核對。
const (
	gengWu  = ganzhi.SexagenaryIndex(6)  // 庚午
	xinSi   = ganzhi.SexagenaryIndex(17) // 辛巳
	yiYou   = ganzhi.SexagenaryIndex(21) // 乙酉
	yearIdx = 0
	monthI  = 1
	dayIdx  = 2
	hourIdx = 3
)

func sampleInput() Input {
	return Input{
		DayStem:     1, // 乙
		YearStem:    6, // 庚
		DayBranch:   9, // 酉
		YearBranch:  6, // 午
		MonthBranch: 5, // 巳
		IsMale:      true,
		YearSound:   ganzhi.SoundOf(gengWu), // 庚午路旁土
	}
}

func sampleScan() []ganzhi.SexagenaryIndex {
	return []ganzhi.SexagenaryIndex{gengWu, xinSi, yiYou, xinSi}
}

// TestDetectSampleChart 基準命例的完整命中，逐筆手算核對。
//
//	日干乙：太極在子午（年柱午）、金輿在巳（月、時兩柱）、文昌在午（年柱）；
//	        祿卯、暗祿戌、羊刃辰、飛刃戌、天乙申子皆無
//	日支酉屬巳酉丑金局：將星在酉（日柱）、咸池在午（年柱）；
//	        華蓋丑、驛馬亥、攀鞍子、劫煞寅、亡神申、災煞卯、六厄子皆無
//	年支午屬巳午未方：孤辰申、寡宿辰皆無；隔角取前後一辰之落四孟四季者，
//	        未為隔、巳為角，巳中月時兩柱
//	月支巳：天德在辛（月、時）、月德在庚（年）、月德合在乙（日）；
//	        天德合丙無。德秀巳酉丑月德庚辛、秀乙庚——庚辛乙庚共五筆
//	年支午：元辰丑無、勾絞陽男取前三為酉（日柱）、暗金的煞在巳（月、時）
//	年納音路旁土：土命地網在辰巳（月、時）；學堂申、詞館亥、正印辰皆無
//	日干乙（食神丁）：天廚在丁祿午（年柱）；福星取遁得丁之時支為丑亥，無
//	年支午：紅鸞自子起卯逆數得酉（日柱）、天喜為其衝在卯，無
//	日柱乙酉屬甲申旬：空亡在午未，午中年柱
//	柱本身：辛巳祿酉入甲戌旬之空亡，十惡大敗（月、時兩柱）；四廢夏取壬子癸亥，無；
//	        魁罡壬辰庚戌庚辰戊戌、天赦夏甲午、金神癸酉己巳乙丑皆無
//	全盤：天干庚辛乙辛，無三奇；太歲類未給流年，不計
func TestDetectSampleChart(t *testing.T) {
	got := Detect(sampleInput(), sampleScan(), Default())

	want := []Hit{
		{Kind: TaiJiGuiRen, At: yearIdx, Basis: BasisDayStem},
		{Kind: JinYu, At: monthI, Basis: BasisDayStem},
		{Kind: JinYu, At: hourIdx, Basis: BasisDayStem},
		{Kind: WenChangGuiRen, At: yearIdx, Basis: BasisDayStem},
		{Kind: JiangXing, At: dayIdx, Basis: BasisDayBranch},
		{Kind: XianChi, At: yearIdx, Basis: BasisDayBranch},
		{Kind: GeJiao, At: monthI, Basis: BasisYearBranch},
		{Kind: GeJiao, At: hourIdx, Basis: BasisYearBranch},
		{Kind: TianDe, At: monthI, Basis: BasisMonthBranch},
		{Kind: TianDe, At: hourIdx, Basis: BasisMonthBranch},
		{Kind: YueDe, At: yearIdx, Basis: BasisMonthBranch},
		{Kind: YueDeHe, At: dayIdx, Basis: BasisMonthBranch},
		{Kind: DeXiu, At: yearIdx, Basis: BasisMonthBranch, Variant: VariantDe},
		{Kind: DeXiu, At: yearIdx, Basis: BasisMonthBranch, Variant: VariantXiu},
		{Kind: DeXiu, At: monthI, Basis: BasisMonthBranch, Variant: VariantDe},
		{Kind: DeXiu, At: dayIdx, Basis: BasisMonthBranch, Variant: VariantXiu},
		{Kind: DeXiu, At: hourIdx, Basis: BasisMonthBranch, Variant: VariantDe},
		{Kind: GouJiao, At: dayIdx, Basis: BasisYearBranch, Variant: VariantGou},
		{Kind: AnJinDeSha, At: monthI, Basis: BasisYearBranch},
		{Kind: AnJinDeSha, At: hourIdx, Basis: BasisYearBranch},
		{Kind: TianLuoDiWang, At: monthI, Basis: BasisSelf, Variant: VariantDiWang},
		{Kind: TianLuoDiWang, At: hourIdx, Basis: BasisSelf, Variant: VariantDiWang},
		{Kind: ShiEDaBai, At: monthI, Basis: BasisSelf},
		{Kind: ShiEDaBai, At: hourIdx, Basis: BasisSelf},
		{Kind: TianChuGuiRen, At: yearIdx, Basis: BasisDayStem},
		{Kind: HongLuan, At: dayIdx, Basis: BasisYearBranch},
		{Kind: XunKong, At: yearIdx, Basis: BasisDayPillar},
	}

	if len(got) != len(want) {
		t.Fatalf("命中 %d 筆，應為 %d 筆\n實得：%v", len(got), len(want), got)
	}
	// Detect 保證依 (Kind, At, Basis, Variant) 排序，故 want 亦須同序
	seen := map[Hit]bool{}
	for _, h := range got {
		seen[h] = true
	}
	for _, w := range want {
		if !seen[w] {
			t.Errorf("缺少命中 %+v\n實得：%v", w, got)
		}
	}
}

// TestDetectIsDeterministic 同一輸入重複呼叫，順序必須完全一致。
//
// 順序列為契約——不確定的順序會使 golden test 不穩。
func TestDetectIsDeterministic(t *testing.T) {
	first := Detect(sampleInput(), sampleScan(), Default())
	for i := 0; i < 20; i++ {
		again := Detect(sampleInput(), sampleScan(), Default())
		if len(again) != len(first) {
			t.Fatalf("第 %d 次呼叫得 %d 筆，首次為 %d 筆", i, len(again), len(first))
		}
		for j := range first {
			if again[j] != first[j] {
				t.Fatalf("第 %d 次呼叫的第 %d 筆為 %+v，首次為 %+v", i, j, again[j], first[j])
			}
		}
	}
}

// TestDetectSortedOrder 回傳須依 (Kind, At, Basis, Variant) 遞增。
func TestDetectSortedOrder(t *testing.T) {
	in := sampleInput()
	got := Detect(in, sampleScan(), Options{BranchBase: BranchBaseBoth})
	for i := 1; i < len(got); i++ {
		a, b := got[i-1], got[i]
		less := a.Kind < b.Kind ||
			(a.Kind == b.Kind && a.At < b.At) ||
			(a.Kind == b.Kind && a.At == b.At && a.Basis < b.Basis) ||
			(a.Kind == b.Kind && a.At == b.At && a.Basis == b.Basis && a.Variant <= b.Variant)
		if !less {
			t.Errorf("第 %d 與第 %d 筆順序不對：%+v 在 %+v 之前", i-1, i, a, b)
		}
	}
}

// TestBranchBaseBothKeepsBothSources 兩者皆查時，同一柱可自年支與日支各命中一次，
// 兩筆皆須保留並以 Basis 區分——那正是開啟該選項的目的。
func TestBranchBaseBothKeepsBothSources(t *testing.T) {
	in := sampleInput()
	both := Detect(in, sampleScan(), Options{BranchBase: BranchBaseBoth})

	var day, year int
	for _, h := range both {
		switch h.Basis {
		case BasisDayBranch:
			day++
		case BasisYearBranch:
			year++
		}
	}
	if day == 0 || year == 0 {
		t.Errorf("兩者皆查應同時產生日支與年支的命中，實得 日支 %d、年支 %d", day, year)
	}

	// 慣用基準下，以地支查者只會有單一來源
	def := Detect(in, sampleScan(), Default())
	if len(both) <= len(def) {
		t.Errorf("兩者皆查得 %d 筆，應多於慣用基準的 %d 筆", len(both), len(def))
	}
}

// TestYinStemBladeOption 陰干羊刃的開關。
//
// 位置不受影響——羊刃為祿位下一支，與 TerrainSect 無關；本項僅管算不算。
func TestYinStemBladeOption(t *testing.T) {
	// 乙日主，刃在辰。造一個帶辰的盤
	in := Input{DayStem: 1, YearStem: 0, DayBranch: 9, YearBranch: 6, MonthBranch: 5}
	scan := []ganzhi.SexagenaryIndex{16} // 庚辰：16%10=6(庚) 16%12=4(辰)

	marked := Detect(in, scan, Options{Include: []Kind{YangRen}, Categories: []Category{}})
	if len(marked) != 1 || marked[0].Variant != VariantYinStem {
		t.Fatalf("陰干計刃時應得一筆且標為陰干刃，實得 %v", marked)
	}

	none := Detect(in, scan, Options{Include: []Kind{YangRen}, Categories: []Category{}, YinStemBlade: YinStemBladeNone})
	if len(none) != 0 {
		t.Errorf("陰干無刃時不應有命中，實得 %v", none)
	}

	// 陽干不受此選項影響：甲日主刃在卯
	yang := Input{DayStem: 0, YearStem: 0, DayBranch: 9, YearBranch: 6, MonthBranch: 5}
	maoScan := []ganzhi.SexagenaryIndex{3} // 丁卯：3%12=3(卯)
	for _, sect := range []YinStemBladeSect{YinStemBladeMarked, YinStemBladeNone} {
		h := Detect(yang, maoScan, Options{Include: []Kind{YangRen}, Categories: []Category{}, YinStemBlade: sect})
		if len(h) != 1 || h[0].Variant != VariantYangStem {
			t.Errorf("口徑 %d 下陽干刃應恆為一筆且標陽干刃，實得 %v", sect, h)
		}
	}
}

// TestTianYiSectOption 天乙的兩個口訣版本只在庚上分歧。
func TestTianYiSectOption(t *testing.T) {
	// 掃描全十二支，直接看庚的貴人落在哪
	all := make([]ganzhi.SexagenaryIndex, 12)
	for i := range all {
		// 取 %12 == i 的任一甲子；i 與 i+12*k 皆可，用 i+48 保證 <60
		all[i] = ganzhi.SexagenaryIndex((i + 48) % 60)
	}
	in := Input{DayStem: 6} // 庚

	branches := func(opt Options) map[ganzhi.BranchIndex]bool {
		out := map[ganzhi.BranchIndex]bool{}
		for _, h := range Detect(in, all, Options{
			Include: []Kind{TianYiGuiRen}, Categories: []Category{}, TianYi: opt.TianYi,
		}) {
			out[all[h.At].Branch()] = true
		}
		return out
	}

	sanMing := branches(Options{TianYi: TianYiSanMing})
	if !sanMing[1] || !sanMing[7] { // 丑、未
		t.Errorf("三命通會本宗：庚貴應在丑未（甲戊庚牛羊），實得 %v", sanMing)
	}
	yeHui := branches(Options{TianYi: TianYiYeHuiTing})
	if !yeHui[6] || !yeHui[2] { // 午、寅
		t.Errorf("葉悔亭改本：庚貴應在午寅（庚辛逢馬虎），實得 %v", yeHui)
	}
}

// TestCategoryFilter 分類過濾。
func TestCategoryFilter(t *testing.T) {
	in, scan := sampleInput(), sampleScan()

	full := Detect(in, scan, Default())
	if len(full) == 0 {
		t.Fatal("預設應有命中")
	}

	// 常用組是全部的真子集，且各筆的分類須相符
	common := Detect(in, scan, Options{Categories: []Category{CategoryCommon}})
	if len(common) == 0 || len(common) >= len(full) {
		t.Errorf("僅取常用組得 %d 筆，全部為 %d 筆——應為非空的真子集", len(common), len(full))
	}
	for _, h := range common {
		if h.Kind.Category() != CategoryCommon {
			t.Errorf("常用組中混入 %s（分類 %s）", h.Kind.ID(), h.Kind.Category().ID())
		}
	}

	// 空的分類清單加空的 Include 應無命中
	none := Detect(in, scan, Options{Categories: []Category{}})
	if len(none) != 0 {
		t.Errorf("不取任何分類時應無命中，實得 %d 筆", len(none))
	}

	// Exclude 優先於 Include
	ex := Detect(in, scan, Options{Include: []Kind{JinYu}, Exclude: []Kind{JinYu}, Categories: []Category{}})
	if len(ex) != 0 {
		t.Errorf("Exclude 應優先於 Include，實得 %v", ex)
	}
}

// TestKindMetadata 每個 Kind 都須有分類、體系、慣用基準與非空 ID。
func TestKindMetadata(t *testing.T) {
	ids := map[string]Kind{}
	for k := Kind(0); k < KindCount; k++ {
		if k.ID() == "" {
			t.Errorf("Kind %d 的 ID 為空", k)
			continue
		}
		if prev, dup := ids[k.ID()]; dup {
			t.Errorf("Kind %d 與 %d 的 ID 皆為 %q", prev, k, k.ID())
		}
		ids[k.ID()] = k

		if b := k.CustomaryBasis(); b >= basisCount {
			t.Errorf("%s 的慣用基準 %d 越界", k.ID(), b)
		}
		if c := k.Category(); c >= categoryCount {
			t.Errorf("%s 的分類 %d 越界", k.ID(), c)
		}
		if tr := k.Tradition(); tr >= traditionCount {
			t.Errorf("%s 的體系 %d 越界", k.ID(), tr)
		}
	}
}

// TestTraditionsAreConsistent 體系標記須與該神煞的來源一致。
//
// 文昌貴人須標明為紫微系——子平典籍另有同名而內容不同的「文昌貴」，兩者僅
// 甲、戊兩干相同，不可混用。十二歲君整組出自星命系，且該組與 CategoryAnnual
// 互為充要——若日後有非歲君的太歲類神煞，本測試會逼實作者當場表態。
func TestTraditionsAreConsistent(t *testing.T) {
	if WenChangGuiRen.Tradition() != TraditionZiwei {
		t.Errorf("文昌貴人的體系為 %s，應為紫微斗數", WenChangGuiRen.Tradition().ID())
	}
	for k := Kind(0); k < KindCount; k++ {
		want := TraditionZiping
		switch {
		case k == WenChangGuiRen:
			want = TraditionZiwei
		case k.Category() == CategoryAnnual:
			want = TraditionSuiJun
		}
		got := k.Tradition()
		if got != want {
			t.Errorf("%s 的體系為 %s，應為 %s", k.ID(), got.ID(), want.ID())
		}
		if (got == TraditionSuiJun) != (k.Category() == CategoryAnnual) {
			t.Errorf("%s：十二歲君與太歲類須互為充要，現為 %s／%s",
				k.ID(), got.ID(), k.Category().ID())
		}
	}
}
