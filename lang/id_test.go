package lang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/LukeLogix/destiny-core/bazi"
	"github.com/LukeLogix/destiny-core/ganzhi"
	"github.com/LukeLogix/destiny-core/shensha"
)

// 列舉的穩定識別字必須與其 Go 常數名一致。
//
// Go 在執行期取不到常數名（編譯後只剩數值），故 ID 表只能手寫；
// 手寫就會忘記，或改了常數名卻沒改表。本測試以 AST 解析原始碼補上這道
// 保證——與 arch_test.go 同屬「靠檢查不靠自律」。
//
// 規則：ID = 常數名的 snake_case，去除同一列舉所有常數共有的 CamelCase 前綴。
//
//	VerdictStrong/VerdictWeak/VerdictNeutral    → strong/weak/neutral
//	TerrainYinReverse/TerrainSameBirth          → yin_reverse/same_birth
//	Wood/Fire/Earth/Metal/Water                 → wood/fire/…（無共同前綴）
//
// 取「共同前綴」而非「型別名」，因兩者常不相同：StrategyWeighted 的型別是
// StrategyID、TerrainYinReverse 的型別是 TerrainSect。既有的
// strategy 欄位輸出 "weighted" 即依此規則，本測試把它補齊到全部列舉。

// idEnum 一個受檢的列舉。
//
// count 與 id 由呼叫端提供而非從 AST 推得——如此「表少一項」「表多一項」
// 也會被抓到，而不只是「內容對不上」。
type idEnum struct {
	dir      string // 原始碼目錄，相對於本測試檔
	typeName string // 型別名，用於定位 const 區塊與去前綴
	count    int
	id       func(int) string
}

var idEnums = []idEnum{
	{"../ganzhi", "Element", ganzhi.ElementCount, func(i int) string { return ganzhi.Element(i).ID() }},
	{"../ganzhi", "Polarity", 2, func(i int) string { return ganzhi.Polarity(i).ID() }},
	{"../ganzhi", "Terrain", ganzhi.TerrainCount, func(i int) string { return ganzhi.Terrain(i).ID() }},
	{"../ganzhi", "TerrainSect", 2, func(i int) string { return ganzhi.TerrainSect(i).ID() }},
	{"../ganzhi", "RelationKind", 10, func(i int) string { return ganzhi.RelationKind(i).ID() }},
	{"../ganzhi", "SoundIndex", ganzhi.SoundCount, func(i int) string { return ganzhi.SoundIndex(i).ID() }},

	{"../bazi", "Gender", 3, func(i int) string { return bazi.Gender(i).ID() }},
	{"../bazi", "SolarTimeMode", 3, func(i int) string { return bazi.SolarTimeMode(i).ID() }},
	{"../bazi", "ChildLimitSect", 4, func(i int) string { return bazi.ChildLimitSect(i).ID() }},
	{"../bazi", "HiddenStemSect", 2, func(i int) string { return bazi.HiddenStemSect(i).ID() }},
	{"../bazi", "HiddenStemType", 3, func(i int) string { return bazi.HiddenStemType(i).ID() }},
	{"../bazi", "TenGod", bazi.TenGodCount, func(i int) string { return bazi.TenGod(i).ID() }},
	{"../bazi", "Verdict", 3, func(i int) string { return bazi.Verdict(i).ID() }},
	{"../bazi", "StrategyID", 2, func(i int) string { return bazi.StrategyID(i).ID() }},
	{"../bazi", "ReasonCode", 8, func(i int) string { return bazi.ReasonCode(i).ID() }},
	{"../bazi", "Consensus", 3, func(i int) string { return bazi.Consensus(i).ID() }},

	{"../shensha", "Kind", int(shensha.KindCount), func(i int) string { return shensha.Kind(i).ID() }},
	{"../shensha", "Basis", 8, func(i int) string { return shensha.Basis(i).ID() }},
	{"../shensha", "Variant", 5, func(i int) string { return shensha.Variant(i).ID() }},
	{"../shensha", "Category", 4, func(i int) string { return shensha.Category(i).ID() }},
	{"../shensha", "Tradition", 3, func(i int) string { return shensha.Tradition(i).ID() }},
	{"../shensha", "BranchBaseSect", 4, func(i int) string { return shensha.BranchBaseSect(i).ID() }},
	{"../shensha", "TianYiSect", 2, func(i int) string { return shensha.TianYiSect(i).ID() }},
	{"../shensha", "YinStemBladeSect", 2, func(i int) string { return shensha.YinStemBladeSect(i).ID() }},
}

func TestEnumIDsMatchConstantNames(t *testing.T) {
	for _, e := range idEnums {
		names := constNames(t, e.dir, e.typeName)

		if len(names) != e.count {
			t.Errorf("%s.%s：原始碼有 %d 個常數，ID 表宣告 %d 個\n"+
				"  常數：%v", e.dir, e.typeName, len(names), e.count, names)
			continue
		}

		prefix := commonCamelPrefix(names)
		for i, name := range names {
			want := snakeCase(name, prefix)
			if got := e.id(i); got != want {
				t.Errorf("%s.%s[%d]：常數 %s 的 ID 為 %q，應為 %q",
					e.dir, e.typeName, i, name, got, want)
			}
		}
	}
}

// TestEnumIDsAreUnique 同一列舉內的 ID 不得重複——重複即無法作為鍵。
func TestEnumIDsAreUnique(t *testing.T) {
	for _, e := range idEnums {
		seen := map[string]int{}
		for i := 0; i < e.count; i++ {
			id := e.id(i)
			if id == "" {
				t.Errorf("%s.%s[%d]：ID 為空", e.dir, e.typeName, i)
				continue
			}
			if prev, dup := seen[id]; dup {
				t.Errorf("%s.%s：索引 %d 與 %d 的 ID 皆為 %q", e.dir, e.typeName, prev, i, id)
			}
			seen[id] = i
		}
	}
}

// TestEnumIDsOutOfRangeSafe 越界不得 panic，回空字串。
func TestEnumIDsOutOfRangeSafe(t *testing.T) {
	for _, e := range idEnums {
		if got := e.id(e.count + 100); got != "" {
			t.Errorf("%s.%s：越界索引的 ID 為 %q，應為空字串", e.dir, e.typeName, got)
		}
	}
}

// ── AST 解析 ──

var camelBoundary = regexp.MustCompile(`(?:[a-z0-9])([A-Z])`)

// commonCamelPrefix 取所有常數共有的前綴，且須止於 CamelCase 邊界。
//
// 單一常數的列舉不作去前綴——無從判斷哪一段是型別綴詞。
func commonCamelPrefix(names []string) string {
	if len(names) < 2 {
		return ""
	}
	p := names[0]
	for _, n := range names[1:] {
		for p != "" && !strings.HasPrefix(n, p) {
			p = p[:len(p)-1]
		}
		if p == "" {
			return ""
		}
	}
	// 退到邊界：前綴之後每個名字的下一字元皆須為大寫，否則切在詞中間
	for p != "" {
		ok := true
		for _, n := range names {
			if len(n) <= len(p) || !unicode.IsUpper(rune(n[len(p)])) {
				ok = false
				break
			}
		}
		if ok {
			return p
		}
		p = p[:len(p)-1]
	}
	return ""
}

// snakeCase 轉換常數名，去除指定前綴。
func snakeCase(name, prefix string) string {
	if prefix != "" && strings.HasPrefix(name, prefix) && len(name) > len(prefix) {
		name = name[len(prefix):]
	}
	out := camelBoundary.ReplaceAllStringFunc(name, func(m string) string {
		return string(m[0]) + "_" + m[1:]
	})
	return strings.ToLower(out)
}

// constNames 取出指定型別的 iota 常數名，依宣告順序。
//
// 略過哨兵常數（名稱以 Count 結尾），以及在 iota 序列後另行賦值者
// （如 `TerrainCount = 12`）。
func constNames(t *testing.T, dir, typeName string) []string {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("讀取 %s 失敗: %v", dir, err)
	}

	var names []string
	for _, entry := range entries {
		n := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		names = append(names, namesInFile(t, filepath.Join(dir, n), typeName)...)
	}
	return names
}

func namesInFile(t *testing.T, path, typeName string) []string {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("解析 %s 失敗: %v", path, err)
	}

	var out []string
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			continue
		}

		started := false
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			if !started {
				// 序列的起點：`X TypeName = iota`
				if id, ok := vs.Type.(*ast.Ident); !ok || id.Name != typeName {
					continue
				}
				if len(vs.Values) != 1 {
					continue
				}
				if id, ok := vs.Values[0].(*ast.Ident); !ok || id.Name != "iota" {
					continue
				}
				started = true
			} else if len(vs.Values) > 0 {
				// 另行賦值即脫離 iota 序列（如 TerrainCount = 12）
				break
			}

			for _, name := range vs.Names {
				if strings.HasSuffix(name.Name, "Count") || name.Name == "_" {
					continue
				}
				out = append(out, name.Name)
			}
		}
	}
	return out
}
