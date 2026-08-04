package ganzhi

// SoundIndex 納音，六十甲子每兩個共用一個，共三十種。
type SoundIndex uint8

const SoundCount = 30

// soundElements 三十種納音的五行，依納音索引排列。
var soundElements = [SoundCount]Element{
	Metal, // 海中金
	Fire,  // 爐中火
	Wood,  // 大林木
	Earth, // 路旁土
	Metal, // 劍鋒金
	Fire,  // 山頭火
	Water, // 澗下水
	Earth, // 城頭土
	Metal, // 白蠟金
	Wood,  // 楊柳木
	Water, // 泉中水
	Earth, // 屋上土
	Fire,  // 霹靂火
	Wood,  // 松柏木
	Water, // 長流水
	Metal, // 沙中金
	Fire,  // 山下火
	Wood,  // 平地木
	Earth, // 壁上土
	Metal, // 金箔金
	Fire,  // 覆燈火
	Water, // 天河水
	Earth, // 大驛土
	Metal, // 釵釧金
	Wood,  // 桑柘木
	Water, // 大溪水
	Earth, // 沙中土
	Fire,  // 天上火
	Wood,  // 石榴木
	Water, // 大海水
}

// SoundOf 取六十甲子的納音
func SoundOf(sex SexagenaryIndex) SoundIndex {
	return SoundIndex(int(sex) % SexagenaryCount / 2)
}

// Element 納音的五行
func (s SoundIndex) Element() Element {
	return soundElements[int(s)%SoundCount]
}
