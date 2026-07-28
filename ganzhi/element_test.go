package ganzhi

import "testing"

// TestElementOrder 五行索引順序由本套件自行定義為木火土金水，
// 不得依賴任何外部庫的排列。凡以 [ElementCount]T 傳遞五行資料者皆依此順序。
func TestElementOrder(t *testing.T) {
	if Wood != 0 || Fire != 1 || Earth != 2 || Metal != 3 || Water != 4 {
		t.Fatalf("五行順序應為木火土金水 0-4，實得 %d %d %d %d %d",
			Wood, Fire, Earth, Metal, Water)
	}
	if ElementCount != 5 {
		t.Errorf("ElementCount 為 %d，應為 5", ElementCount)
	}
}

// TestElementGenerates 相生：木生火、火生土、土生金、金生水、水生木。
func TestElementGenerates(t *testing.T) {
	cases := []struct{ from, want Element }{
		{Wood, Fire}, {Fire, Earth}, {Earth, Metal}, {Metal, Water}, {Water, Wood},
	}
	for _, c := range cases {
		if got := c.from.Generates(); got != c.want {
			t.Errorf("%d 生 %d，應生 %d", c.from, got, c.want)
		}
		if got := c.want.GeneratedBy(); got != c.from {
			t.Errorf("生 %d 者為 %d，應為 %d", c.want, got, c.from)
		}
	}
}

// TestElementRestrains 相剋：木剋土、火剋金、土剋水、金剋木、水剋火。
func TestElementRestrains(t *testing.T) {
	cases := []struct{ from, want Element }{
		{Wood, Earth}, {Fire, Metal}, {Earth, Water}, {Metal, Wood}, {Water, Fire},
	}
	for _, c := range cases {
		if got := c.from.Restrains(); got != c.want {
			t.Errorf("%d 剋 %d，應剋 %d", c.from, got, c.want)
		}
		if got := c.want.RestrainedBy(); got != c.from {
			t.Errorf("剋 %d 者為 %d，應為 %d", c.want, got, c.from)
		}
	}
}

// TestElementCycleClosure 生剋皆為五元循環，繞五步必回原點。
func TestElementCycleClosure(t *testing.T) {
	for e := Element(0); e < ElementCount; e++ {
		x := e
		for i := 0; i < 5; i++ {
			x = x.Generates()
		}
		if x != e {
			t.Errorf("%d 相生五步後為 %d，應回原點", e, x)
		}
		y := e
		for i := 0; i < 5; i++ {
			y = y.Restrains()
		}
		if y != e {
			t.Errorf("%d 相剋五步後為 %d，應回原點", e, y)
		}
	}
}

// TestStemElement 天干五行：甲乙木、丙丁火、戊己土、庚辛金、壬癸水。
func TestStemElement(t *testing.T) {
	want := []Element{Wood, Wood, Fire, Fire, Earth, Earth, Metal, Metal, Water, Water}
	for i, w := range want {
		if got := StemIndex(i).Element(); got != w {
			t.Errorf("天干 %d 五行為 %d，應為 %d", i, got, w)
		}
	}
}

// TestBranchElement 地支五行：寅卯木、巳午火、申酉金、亥子水、辰戌丑未土。
func TestBranchElement(t *testing.T) {
	want := []Element{
		Water, Earth, Wood, Wood, Earth, Fire,
		Fire, Earth, Metal, Metal, Earth, Water,
	}
	for i, w := range want {
		if got := BranchIndex(i).Element(); got != w {
			t.Errorf("地支 %d 五行為 %d，應為 %d", i, got, w)
		}
	}
}

// TestStemPolarity 天干陰陽：甲丙戊庚壬為陽，乙丁己辛癸為陰。
func TestStemPolarity(t *testing.T) {
	for i := 0; i < StemCount; i++ {
		want := Yang
		if i%2 == 1 {
			want = Yin
		}
		if got := StemIndex(i).Polarity(); got != want {
			t.Errorf("天干 %d 陰陽為 %d，應為 %d", i, got, want)
		}
	}
}

// TestBranchPolarity 地支陰陽：子寅辰午申戌為陽，丑卯巳未酉亥為陰。
func TestBranchPolarity(t *testing.T) {
	for i := 0; i < BranchCount; i++ {
		want := Yang
		if i%2 == 1 {
			want = Yin
		}
		if got := BranchIndex(i).Polarity(); got != want {
			t.Errorf("地支 %d 陰陽為 %d，應為 %d", i, got, want)
		}
	}
}

// TestSexagenaryPolarityConsistent 六十甲子中干支陰陽必然一致，
// 這正是只有 60 組而非 120 組的原因。
func TestSexagenaryPolarityConsistent(t *testing.T) {
	for i := 0; i < SexagenaryCount; i++ {
		x := SexagenaryIndex(i)
		if x.Stem().Polarity() != x.Branch().Polarity() {
			t.Errorf("索引 %d 的干支陰陽不一致", i)
		}
	}
}
