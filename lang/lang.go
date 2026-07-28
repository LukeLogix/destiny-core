// Package lang 提供命盤的文字呈現。
//
// 依賴方向刻意是 lang → 計算層，反向絕不成立：計算層（ganzhi、bazi）
// 完全不含文字，此規則由編譯器強制，而非文件約定或 lint。
// 新增語言只需新增一份資料，不動任何計算程式碼。
package lang

import (
	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/ganzhi"
)

// Locale 語言標籤
type Locale string

const (
	ZhTW Locale = "zh-TW"
	ZhCN Locale = "zh-CN"

	// DefaultLocale 未知 locale 時的回退對象
	DefaultLocale = ZhTW
)

// Locales 回傳所有已支援的語言
func Locales() []Locale { return []Locale{ZhTW, ZhCN} }

// bundle 一個語言的全部文字資料。
//
// 全為索引對應的字串陣列——新增語言即新增一份資料，
// 計算層不需要知道有幾種語言，也不需要知道文字長什麼樣。
type bundle struct {
	stems    [ganzhi.StemCount]string
	branches [ganzhi.BranchCount]string
	elements [ganzhi.ElementCount]string
	tenGods  [bazi.TenGodCount]string
	terrains [bazi.TerrainCount]string
	sounds   [bazi.SoundCount]string
	verdicts [3]string
	reasons  [8]string
	polarity [2]string
	warnTerm string
	warnHour string
}

// Bundle 取得指定語言的文字資料；未知 locale 回退至預設。
func Bundle(loc Locale) *bundle {
	if b, ok := bundles[loc]; ok {
		return b
	}
	return bundles[DefaultLocale]
}

func pick(arr []string, i int) string {
	if i < 0 || i >= len(arr) {
		return ""
	}
	return arr[i]
}

func (b *bundle) Stem(s ganzhi.StemIndex) string {
	return pick(b.stems[:], int(s))
}

func (b *bundle) Branch(x ganzhi.BranchIndex) string {
	return pick(b.branches[:], int(x))
}

func (b *bundle) Element(e ganzhi.Element) string {
	return pick(b.elements[:], int(e))
}

func (b *bundle) TenGod(g bazi.TenGod) string {
	return pick(b.tenGods[:], int(g))
}

func (b *bundle) Terrain(t bazi.Terrain) string {
	return pick(b.terrains[:], int(t))
}

func (b *bundle) Sound(s bazi.SoundIndex) string {
	return pick(b.sounds[:], int(s))
}

func (b *bundle) Verdict(v bazi.Verdict) string {
	return pick(b.verdicts[:], int(v))
}

func (b *bundle) Reason(r bazi.ReasonCode) string {
	return pick(b.reasons[:], int(r))
}

// Sexagenary 六十甲子由干支拼成
func (b *bundle) Sexagenary(x ganzhi.SexagenaryIndex) string {
	return b.Stem(x.Stem()) + b.Branch(x.Branch())
}
