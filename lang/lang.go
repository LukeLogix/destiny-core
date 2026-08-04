// Package lang 提供命盤的文字呈現。
//
// 依賴方向刻意是 lang → 計算層，反向絕不成立：計算層（ganzhi、bazi）
// 完全不含文字，此規則由編譯器強制，而非文件約定或 lint。
// 新增語言只需新增一份資料，不動任何計算程式碼。
package lang

import (
	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/shensha"
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
	stems     [ganzhi.StemCount]string
	branches  [ganzhi.BranchCount]string
	elements  [ganzhi.ElementCount]string
	tenGods   [bazi.TenGodCount]string
	terrains  [ganzhi.TerrainCount]string
	sounds    [ganzhi.SoundCount]string
	verdicts  [3]string
	reasons   [8]string
	polarity  [2]string
	relations [10]string
	pillars   [4]string
	solarTime [3]string

	// 口徑與分類的文字。原本這幾項不需要文字（以裸數字交出），
	// 改 Term 後每個 ID 都要有對應的顯示名。
	genders         [3]string
	hiddenTypes     [3]string
	hiddenStemSects [2]string
	terrainSects    [2]string
	childLimitSects [4]string
	consensuses     [3]string

	// 神煞
	shenSha    [shensha.KindCount]string
	basis      [10]string
	variants   [13]string
	categories [4]string
	traditions [3]string

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

// Term 一個分類值：穩定識別字 ＋ 該語系的顯示文字。
//
// ID 為計算層 Go 常數名的 snake_case，不隨 locale 變動，可作為程式分支、
// 篩選、儲存的鍵；Name 僅供顯示，隨 locale 變動，不得用於判斷。
//
// 確立的原則：凡程式可能據以分支、篩選或儲存的值，必須提供穩定 ID；
// 純顯示的文字可以只給本地化字串。判斷不了時一律給 ID。
type Term struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// term 組裝一個分類值。ID 來自計算層，Name 來自本語系的文字表。
func term(id string, arr []string, i int) Term {
	return Term{ID: id, Name: pick(arr, i)}
}

func (b *bundle) Stem(s ganzhi.StemIndex) string {
	return pick(b.stems[:], int(s))
}

func (b *bundle) Branch(x ganzhi.BranchIndex) string {
	return pick(b.branches[:], int(x))
}

func (b *bundle) Element(e ganzhi.Element) Term {
	return term(e.ID(), b.elements[:], int(e))
}

func (b *bundle) TenGod(g bazi.TenGod) Term {
	return term(g.ID(), b.tenGods[:], int(g))
}

func (b *bundle) Terrain(t ganzhi.Terrain) Term {
	return term(t.ID(), b.terrains[:], int(t))
}

func (b *bundle) Sound(s ganzhi.SoundIndex) Term {
	return term(s.ID(), b.sounds[:], int(s))
}

func (b *bundle) Verdict(v bazi.Verdict) Term {
	return term(v.ID(), b.verdicts[:], int(v))
}

// Relation 干支關係的種類名稱
func (b *bundle) Relation(k ganzhi.RelationKind) Term {
	return term(k.ID(), b.relations[:], int(k))
}

// Pillar 柱位名稱。ganzhi 只回報「傳入 slice 的第幾個」，
// 由此處賦予年月日時的語意——共用層不綁定特定術數的位置概念。
func (b *bundle) Pillar(i int) Term {
	return term(pick(pillarIDs[:], i), b.pillars[:], i)
}

// pillarIDs 柱位的識別字。ganzhi 只回報「傳入 slice 的第幾個」，位置語意
// 由本層賦予，故此表不在計算層——它是八字的概念，不是干支的。
var pillarIDs = [4]string{"year", "month", "day", "hour"}

// SolarTime 真太陽時口徑名稱
func (b *bundle) SolarTime(m bazi.SolarTimeMode) Term {
	return term(m.ID(), b.solarTime[:], int(m))
}

func (b *bundle) Reason(r bazi.ReasonCode) Term {
	return term(r.ID(), b.reasons[:], int(r))
}

func (b *bundle) ShenSha(k shensha.Kind) Term {
	return term(k.ID(), b.shenSha[:], int(k))
}

func (b *bundle) Basis(x shensha.Basis) Term {
	return term(x.ID(), b.basis[:], int(x))
}

func (b *bundle) Variant(v shensha.Variant) Term {
	return term(v.ID(), b.variants[:], int(v))
}

func (b *bundle) Category(c shensha.Category) Term {
	return term(c.ID(), b.categories[:], int(c))
}

func (b *bundle) Tradition(t shensha.Tradition) Term {
	return term(t.ID(), b.traditions[:], int(t))
}

// Sexagenary 六十甲子由干支拼成
func (b *bundle) Sexagenary(x ganzhi.SexagenaryIndex) string {
	return b.Stem(x.Stem()) + b.Branch(x.Branch())
}
