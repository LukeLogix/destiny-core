package ganzhi

import "testing"

// TestSoundOfSexagenary 納音每兩個甲子一組，六十甲子共三十種。
func TestSoundOfSexagenary(t *testing.T) {
	cases := []struct {
		sex  SexagenaryIndex
		want SoundIndex
		desc string
	}{
		{0, 0, "甲子海中金"},
		{1, 0, "乙丑亦海中金"},
		{2, 1, "丙寅爐中火"},
		{6, 3, "庚午路旁土"},
		{17, 8, "辛巳白蠟金"},
		{21, 10, "乙酉泉中水"},
		{59, 29, "癸亥大海水"},
	}
	for _, c := range cases {
		if got := SoundOf(c.sex); got != c.want {
			t.Errorf("%s: 得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestSoundElement 納音各有五行，1990-05-20 命例：
// 庚午路旁土、辛巳白蠟金、乙酉泉中水、辛巳白蠟金。
func TestSoundElement(t *testing.T) {
	cases := []struct {
		sex  SexagenaryIndex
		want Element
		desc string
	}{
		{6, Earth, "庚午路旁土"},
		{17, Metal, "辛巳白蠟金"},
		{21, Water, "乙酉泉中水"},
		{0, Metal, "甲子海中金"},
		{59, Water, "癸亥大海水"},
	}
	for _, c := range cases {
		if got := SoundOf(c.sex).Element(); got != c.want {
			t.Errorf("%s: 五行得 %d，應為 %d", c.desc, got, c.want)
		}
	}
}

// TestSoundCoversThirty 六十甲子恰好對應三十種納音，每種各兩個。
func TestSoundCoversThirty(t *testing.T) {
	count := map[SoundIndex]int{}
	for i := 0; i < SexagenaryCount; i++ {
		count[SoundOf(SexagenaryIndex(i))]++
	}
	if len(count) != SoundCount {
		t.Errorf("納音共 %d 種，應為 %d 種", len(count), SoundCount)
	}
	for s, n := range count {
		if n != 2 {
			t.Errorf("納音 %d 對應 %d 個甲子，應為 2 個", s, n)
		}
	}
}

// TestSoundNamedConstants 具名常數須與納音索引一致。
//
// 常數順序一旦寫錯即全盤位移，且不會有任何症狀（納音仍算得出來，只是名字錯）。
// 故以命例釘住若干點，並驗證五行相符——名字與五行對不上即為錯位。
func TestSoundNamedConstants(t *testing.T) {
	cases := []struct {
		sex  SexagenaryIndex
		want SoundIndex
		elem Element
		desc string
	}{
		{0, HaiZhongJin, Metal, "甲子海中金"},
		{6, LuPangTu, Earth, "庚午路旁土"},
		{17, BaiLaJin, Metal, "辛巳白蠟金"},
		{21, QuanZhongShui, Water, "乙酉泉中水"},
		{59, DaHaiShui, Water, "癸亥大海水"},
	}
	for _, c := range cases {
		got := SoundOf(c.sex)
		if got != c.want {
			t.Errorf("%s: 納音索引得 %d，應為 %d", c.desc, got, c.want)
		}
		if e := got.Element(); e != c.elem {
			t.Errorf("%s: 五行得 %d，應為 %d——名字與五行對不上代表常數錯位", c.desc, e, c.elem)
		}
	}
	if int(DaHaiShui)+1 != SoundCount {
		t.Errorf("末項 DaHaiShui = %d，加一應等於 SoundCount %d", DaHaiShui, SoundCount)
	}
}
