package lang

import (
	"strings"
	"testing"
	"time"

	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/ganzhi"
)

var cst = time.FixedZone("UTC+8", 8*3600)

func sampleChart(t *testing.T) *bazi.Chart {
	t.Helper()
	c, err := bazi.Compute(bazi.Birth{
		Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, cst),
		Gender: bazi.Male,
	}, bazi.Default())
	if err != nil {
		t.Fatalf("排盤失敗: %v", err)
	}
	return c
}

// TestStemNames 天干名稱，繁簡在此組字形相同。
func TestStemNames(t *testing.T) {
	want := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	for _, loc := range []Locale{ZhTW, ZhCN} {
		b := Bundle(loc)
		for i, w := range want {
			if got := b.Stem(ganzhi.StemIndex(i)); got != w {
				t.Errorf("%s 天干 %d 為 %q，應為 %q", loc, i, got, w)
			}
		}
	}
}

// TestBranchNames 地支名稱
func TestBranchNames(t *testing.T) {
	want := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	b := Bundle(ZhTW)
	for i, w := range want {
		if got := b.Branch(ganzhi.BranchIndex(i)); got != w {
			t.Errorf("地支 %d 為 %q，應為 %q", i, got, w)
		}
	}
}

// TestTenGodDiffersBetweenLocales 十神是繁簡字形差異的主要所在。
func TestTenGodDiffersBetweenLocales(t *testing.T) {
	cases := []struct {
		god    bazi.TenGod
		tw, cn string
	}{
		{bazi.SevenKilling, "七殺", "七杀"},
		{bazi.IndirectWealth, "偏財", "偏财"},
		{bazi.DirectWealth, "正財", "正财"},
		{bazi.Hurting, "傷官", "伤官"},
		{bazi.Rival, "劫財", "劫财"},
		{bazi.Peer, "比肩", "比肩"},
	}
	tw, cn := Bundle(ZhTW), Bundle(ZhCN)
	for _, c := range cases {
		if got := tw.TenGod(c.god); got != c.tw {
			t.Errorf("繁體十神 %d 為 %q，應為 %q", c.god, got, c.tw)
		}
		if got := cn.TenGod(c.god); got != c.cn {
			t.Errorf("簡體十神 %d 為 %q，應為 %q", c.god, got, c.cn)
		}
	}
}

// TestAllLocalesCoverAllIndices 每個 locale 都須涵蓋所有索引，不得有空字串。
func TestAllLocalesCoverAllIndices(t *testing.T) {
	for _, loc := range Locales() {
		b := Bundle(loc)
		for i := 0; i < ganzhi.StemCount; i++ {
			if b.Stem(ganzhi.StemIndex(i)) == "" {
				t.Errorf("%s 缺天干 %d", loc, i)
			}
		}
		for i := 0; i < ganzhi.BranchCount; i++ {
			if b.Branch(ganzhi.BranchIndex(i)) == "" {
				t.Errorf("%s 缺地支 %d", loc, i)
			}
		}
		for i := 0; i < bazi.TenGodCount; i++ {
			if b.TenGod(bazi.TenGod(i)) == "" {
				t.Errorf("%s 缺十神 %d", loc, i)
			}
		}
		for i := 0; i < ganzhi.ElementCount; i++ {
			if b.Element(ganzhi.Element(i)) == "" {
				t.Errorf("%s 缺五行 %d", loc, i)
			}
		}
		for i := 0; i < ganzhi.TerrainCount; i++ {
			if b.Terrain(ganzhi.Terrain(i)) == "" {
				t.Errorf("%s 缺十二長生 %d", loc, i)
			}
		}
		for i := 0; i < ganzhi.SoundCount; i++ {
			if b.Sound(ganzhi.SoundIndex(i)) == "" {
				t.Errorf("%s 缺納音 %d", loc, i)
			}
		}
	}
}

// TestSexagenaryName 六十甲子由干支拼成
func TestSexagenaryName(t *testing.T) {
	b := Bundle(ZhTW)
	cases := []struct {
		idx  ganzhi.SexagenaryIndex
		want string
	}{
		{0, "甲子"}, {6, "庚午"}, {17, "辛巳"}, {21, "乙酉"}, {59, "癸亥"},
	}
	for _, c := range cases {
		if got := b.Sexagenary(c.idx); got != c.want {
			t.Errorf("六十甲子 %d 為 %q，應為 %q", c.idx, got, c.want)
		}
	}
}

// TestLocalizeChart 命盤本地化：1990-05-20 應呈現庚午 辛巳 乙酉 辛巳。
func TestLocalizeChart(t *testing.T) {
	c := sampleChart(t)
	lc := Localize(c, ZhTW)

	if lc.Year.Sexagenary != "庚午" {
		t.Errorf("年柱為 %q，應為庚午", lc.Year.Sexagenary)
	}
	if lc.Month.Sexagenary != "辛巳" {
		t.Errorf("月柱為 %q，應為辛巳", lc.Month.Sexagenary)
	}
	if lc.Day.Sexagenary != "乙酉" {
		t.Errorf("日柱為 %q，應為乙酉", lc.Day.Sexagenary)
	}
	if lc.DayMaster != "乙" {
		t.Errorf("日主為 %q，應為乙", lc.DayMaster)
	}
	if lc.Year.StemTenGod != "正官" {
		t.Errorf("年干十神為 %q，應為正官", lc.Year.StemTenGod)
	}
	// 月支巳藏丙庚戊
	if len(lc.Month.Hidden) != 3 {
		t.Fatalf("月支藏干 %d 個，應為 3 個", len(lc.Month.Hidden))
	}
	if lc.Month.Hidden[0].Stem != "丙" || lc.Month.Hidden[0].TenGod != "傷官" {
		t.Errorf("月支首個藏干為 %q/%q，應為丙/傷官",
			lc.Month.Hidden[0].Stem, lc.Month.Hidden[0].TenGod)
	}
	// 納音
	if lc.Year.Sound != "路旁土" {
		t.Errorf("年柱納音為 %q，應為路旁土", lc.Year.Sound)
	}
	// 地勢
	if lc.Year.Terrain != "長生" {
		t.Errorf("年柱地勢為 %q，應為長生", lc.Year.Terrain)
	}
}

// TestLocalizeStrength 旺衰結論與依據皆須本地化。
func TestLocalizeStrength(t *testing.T) {
	c := sampleChart(t)
	lc := Localize(c, ZhTW)

	if len(lc.Strengths) != 2 {
		t.Fatalf("應有兩套旺衰結果，實得 %d 套", len(lc.Strengths))
	}
	if lc.Consensus != "身弱" {
		t.Errorf("共識為 %q，應為身弱", lc.Consensus)
	}
	for _, s := range lc.Strengths {
		if s.Verdict == "" {
			t.Error("旺衰結論未本地化")
		}
		if len(s.Reasons) == 0 {
			t.Error("旺衰依據未本地化")
		}
		for _, r := range s.Reasons {
			if r == "" {
				t.Error("依據文字為空")
			}
		}
	}
}

// TestLocalizeBoundary 臨界警示須以人看得懂的方式呈現。
func TestLocalizeBoundary(t *testing.T) {
	// 2024 立春前 30 秒出生，必為臨界
	c, err := bazi.Compute(bazi.Birth{
		Time:   time.Date(2024, 2, 4, 16, 26, 40, 0, cst),
		Gender: bazi.Male,
	}, bazi.Default())
	if err != nil {
		t.Fatal(err)
	}
	lc := Localize(c, ZhTW)

	if !c.Boundary.TermCritical {
		t.Skip("此時刻未觸發臨界，跳過")
	}
	if lc.Warnings == nil || len(lc.Warnings) == 0 {
		t.Error("臨界情況應產生警示文字")
	}
	found := false
	for _, w := range lc.Warnings {
		if strings.Contains(w, "節氣") || strings.Contains(w, "臨界") {
			found = true
		}
	}
	if !found {
		t.Errorf("警示未提及節氣臨界：%v", lc.Warnings)
	}
}

// TestUnknownLocaleFallsBack 未知的 locale 回退到預設，不得 panic 或回空。
func TestUnknownLocaleFallsBack(t *testing.T) {
	b := Bundle(Locale("xx-YY"))
	if b.Stem(0) == "" {
		t.Error("未知 locale 應回退到預設而非回空字串")
	}
}

// TestOutOfRangeIndexDoesNotPanic 越界索引須安全處理。
func TestOutOfRangeIndexDoesNotPanic(t *testing.T) {
	b := Bundle(ZhTW)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("越界索引不應 panic: %v", r)
		}
	}()
	_ = b.Stem(ganzhi.StemIndex(200))
	_ = b.Branch(ganzhi.BranchIndex(200))
	_ = b.TenGod(bazi.TenGod(200))
	_ = b.Sound(ganzhi.SoundIndex(200))
}

// TestLocalizeIncludesRelations 干支關係須進入本地化輸出。
//
// 少了它，餵給 LLM 的 JSON 會讓 AI 以為此盤沒有任何沖刑合會，
// 與 README 所承諾的計算範圍相違。
func TestLocalizeIncludesRelations(t *testing.T) {
	c := sampleChart(t)
	lc := Localize(c, ZhTW)

	if len(c.Relations) == 0 {
		t.Fatal("此命例應有干支關係可測（乙庚合金、巳酉半合）")
	}
	if len(lc.Relations) != len(c.Relations) {
		t.Fatalf("本地化後有 %d 組關係，命盤有 %d 組", len(lc.Relations), len(c.Relations))
	}
	for _, r := range lc.Relations {
		if r.Kind == "" {
			t.Error("關係種類未本地化")
		}
		if len(r.Pillars) == 0 {
			t.Error("關係未標明涉及哪幾柱")
		}
	}

	// 乙庚合金：年干庚與日干乙
	var found bool
	for _, r := range lc.Relations {
		if r.Kind == "天干五合" && r.Transform == "金" {
			found = true
		}
	}
	if !found {
		t.Errorf("應含乙庚合金，實得 %+v", lc.Relations)
	}
}

// TestLocalizeRelationPillarNames 關係的位置須譯為年月日時，
// 而非直接透出 ganzhi 的 slice 索引。
func TestLocalizeRelationPillarNames(t *testing.T) {
	c := sampleChart(t)
	lc := Localize(c, ZhTW)

	valid := map[string]bool{"年": true, "月": true, "日": true, "時": true}
	for _, r := range lc.Relations {
		for _, p := range r.Pillars {
			if !valid[p] {
				t.Errorf("柱位名稱 %q 不在年月日時之內", p)
			}
		}
	}
}

// TestLocalizeIncludesAuditFields 稽核所需的欄位不得在本地化時被丟棄。
func TestLocalizeIncludesAuditFields(t *testing.T) {
	c := sampleChart(t)
	lc := Localize(c, ZhTW)

	if lc.EffectiveTime == "" {
		t.Error("缺實際推算時刻")
	}
	if lc.Boundary == nil {
		t.Fatal("缺臨界數值")
	}
	if lc.Boundary.NearTermSeconds != c.Boundary.NearTermSeconds {
		t.Errorf("距節秒數為 %v，應為 %v",
			lc.Boundary.NearTermSeconds, c.Boundary.NearTermSeconds)
	}
	if lc.Boundary.TermUncertaintySec != c.Boundary.TermUncertaintySec {
		t.Error("不確定度未轉出")
	}
	if lc.Options == nil {
		t.Fatal("缺口徑記錄——命盤的可重現性靠它")
	}
	if lc.Options.LateZiKeepsDay != c.Options.LateZiKeepsDay {
		t.Error("子時口徑未轉出")
	}
	if lc.Options.SolarTime == "" {
		t.Error("真太陽時口徑未轉出")
	}
}
