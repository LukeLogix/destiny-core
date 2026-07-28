package bazi

import "github.com/LukeLogix/destiny-core/ganzhi"

// Terrain 十二長生，描述天干在各地支的旺衰狀態。
type Terrain uint8

const (
	LongLife   Terrain = iota // 長生
	Bath                      // 沐浴
	Cap                       // 冠帶
	Officer                   // 臨官
	Prosperity                // 帝旺
	Decline                   // 衰
	Sick                      // 病
	Death                     // 死
	Tomb                      // 墓
	Void                      // 絕
	Fetus                     // 胎
	Nurture                   // 養

	TerrainCount = 12
)

// longLifeStart 各天干的長生位。
// 陽干依五行本性起長生；陰干在傳統派另有起點，且行進方向相反。
var longLifeStart = [ganzhi.StemCount]ganzhi.BranchIndex{
	11, // 甲：亥
	6,  // 乙：午
	2,  // 丙：寅
	9,  // 丁：酉
	2,  // 戊：寅（同丙）
	9,  // 己：酉（同丁）
	5,  // 庚：巳
	0,  // 辛：子
	8,  // 壬：申
	3,  // 癸：卯
}

// TerrainOf 取天干在指定地支的十二長生狀態。
//
// 陰干的處理是各家分歧處：傳統派主張陰干逆行、另有起點；
// 「陰陽同生同死」派則讓陰干與同五行的陽干共用起點且同樣順行。
func TerrainOf(stem ganzhi.StemIndex, branch ganzhi.BranchIndex, sect TerrainSect) Terrain {
	start := longLifeStart[int(stem)%ganzhi.StemCount]
	step := 1

	if stem.Polarity() == ganzhi.Yin {
		switch sect {
		case TerrainYinReverse:
			step = -1
		case TerrainSameBirth:
			// 與同五行的陽干共用起點（陰干索引減一即同五行的陽干）
			start = longLifeStart[int(stem)-1]
		}
	}

	return Terrain(mod((int(branch)-int(start))*step, TerrainCount))
}

// SoundIndex 納音，六十甲子每兩個共用一個，共三十種。
type SoundIndex uint8

const SoundCount = 30

// soundElements 三十種納音的五行，依納音索引排列。
var soundElements = [SoundCount]ganzhi.Element{
	ganzhi.Metal, // 海中金
	ganzhi.Fire,  // 爐中火
	ganzhi.Wood,  // 大林木
	ganzhi.Earth, // 路旁土
	ganzhi.Metal, // 劍鋒金
	ganzhi.Fire,  // 山頭火
	ganzhi.Water, // 澗下水
	ganzhi.Earth, // 城頭土
	ganzhi.Metal, // 白蠟金
	ganzhi.Wood,  // 楊柳木
	ganzhi.Water, // 泉中水
	ganzhi.Earth, // 屋上土
	ganzhi.Fire,  // 霹靂火
	ganzhi.Wood,  // 松柏木
	ganzhi.Water, // 長流水
	ganzhi.Metal, // 沙中金
	ganzhi.Fire,  // 山下火
	ganzhi.Wood,  // 平地木
	ganzhi.Earth, // 壁上土
	ganzhi.Metal, // 金箔金
	ganzhi.Fire,  // 覆燈火
	ganzhi.Water, // 天河水
	ganzhi.Earth, // 大驛土
	ganzhi.Metal, // 釵釧金
	ganzhi.Wood,  // 桑柘木
	ganzhi.Water, // 大溪水
	ganzhi.Earth, // 沙中土
	ganzhi.Fire,  // 天上火
	ganzhi.Wood,  // 石榴木
	ganzhi.Water, // 大海水
}

// SoundOf 取六十甲子的納音
func SoundOf(sex ganzhi.SexagenaryIndex) SoundIndex {
	return SoundIndex(int(sex) % ganzhi.SexagenaryCount / 2)
}

// Element 納音的五行
func (s SoundIndex) Element() ganzhi.Element {
	return soundElements[int(s)%SoundCount]
}
