package ganzhi

// Terrain 十二長生，描述天干在各地支的旺衰狀態。
//
// 置於本套件而非某一術數核心：其運算只用到干支與陰陽，
// 不含任何八字專屬概念，六壬、紫微等亦需之。
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

// TerrainSect 十二長生的陰干順逆口徑
type TerrainSect uint8

const (
	// TerrainYinReverse 陰干逆行：乙長生在午（預設，傳統派）
	TerrainYinReverse TerrainSect = iota
	// TerrainSameBirth 陰陽同生同死：乙長生在亥
	TerrainSameBirth
)

// longLifeStart 各天干的長生位。
// 陽干依五行本性起長生；陰干在傳統派另有起點，且行進方向相反。
var longLifeStart = [StemCount]BranchIndex{
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
func TerrainOf(stem StemIndex, branch BranchIndex, sect TerrainSect) Terrain {
	start := longLifeStart[int(stem)%StemCount]
	step := 1

	if stem.Polarity() == Yin {
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
