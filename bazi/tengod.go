// Package bazi 提供八字命盤的計算。
//
// 本套件只認索引，不含任何文字——文字由 lang 套件負責，依賴方向為 lang → bazi。
package bazi

import "github.com/LukeLogix/destiny-core/ganzhi"

// TenGod 十神。
//
// 依日主與對象天干的五行關係與陰陽異同而定：
// 五行關係分同我、我生、我剋、剋我、生我五類，各類再依陰陽同異分為兩種。
type TenGod uint8

const (
	Peer             TenGod = iota // 比肩：同我，陰陽同
	Rival                          // 劫財：同我，陰陽異
	Output                         // 食神：我生，陰陽同
	Hurting                        // 傷官：我生，陰陽異
	IndirectWealth                 // 偏財：我剋，陰陽同
	DirectWealth                   // 正財：我剋，陰陽異
	SevenKilling                   // 七殺：剋我，陰陽同
	DirectOfficer                  // 正官：剋我，陰陽異
	IndirectResource               // 偏印：生我，陰陽同
	DirectResource                 // 正印：生我，陰陽異

	TenGodCount = 10
)

// TenGodOf 以 me 為日主，取 other 的十神。
func TenGodOf(me, other ganzhi.StemIndex) TenGod {
	meElem, otherElem := me.Element(), other.Element()

	// 五行關係決定類別，每類佔兩個連號的十神
	var base TenGod
	switch {
	case otherElem == meElem:
		base = Peer
	case otherElem == meElem.Generates():
		base = Output
	case otherElem == meElem.Restrains():
		base = IndirectWealth
	case otherElem == meElem.RestrainedBy():
		base = SevenKilling
	default: // otherElem == meElem.GeneratedBy()
		base = IndirectResource
	}

	// 陰陽同者取前者（偏），異者取後者（正）
	if me.Polarity() != other.Polarity() {
		base++
	}
	return base
}
