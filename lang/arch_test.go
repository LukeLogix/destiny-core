package lang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// 計算層不得出現的命理術語。
//
// 只查術語而非所有中文——錯誤訊息面向開發者，不屬於使用者可見的呈現層，
// 出現中文無妨；但天干地支十神這類術語一旦寫進計算層，i18n 即告破功。
var forbiddenTerms = []string{
	"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸",
	"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥",
	"比肩", "劫財", "劫财", "食神", "傷官", "伤官",
	"偏財", "偏财", "正財", "正财", "七殺", "七杀",
	"正官", "偏印", "正印",
	"長生", "长生", "沐浴", "冠帶", "冠带", "臨官", "临官", "帝旺",
	"身強", "身强", "身弱",
	"海中金", "爐中火", "炉中火", "路旁土", "白蠟金", "白蜡金",

	// 神煞術語。新增神煞若把中文寫進計算層，i18n 即告破功。
	"神煞", "神殺", "神杀", "貴人", "贵人", "天乙", "太極", "太极",
	"桃花", "咸池", "驛馬", "驿马", "華蓋", "华盖", "將星", "将星",
	"羊刃", "陽刃", "阳刃", "祿神", "禄神", "金輿", "金舆",
	"文昌", "劫煞", "刼煞", "亡神", "災煞", "灾煞", "六厄",
	"孤辰", "寡宿", "學堂", "学堂", "詞館", "词馆", "元辰",
	"勾絞", "勾绞", "魁罡", "空亡",
}

// computeLayers 計算層的套件目錄，相對於本測試檔
var computeLayers = []string{"../bazi", "../ganzhi", "../shensha", "../internal/calendar"}

// TestComputeLayersContainNoTerminology 計算層的字串字面量不得含命理術語。
//
// 這是「計算層與語言無關」的機械化保證。人會忘記，測試不會。
func TestComputeLayersContainNoTerminology(t *testing.T) {
	for _, dir := range computeLayers {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("讀取 %s 失敗: %v", dir, err)
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			path := filepath.Join(dir, name)
			checkFileForTerms(t, path)
		}
	}
}

func checkFileForTerms(t *testing.T, path string) {
	t.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0) // 不含註解，只查程式碼
	if err != nil {
		t.Fatalf("解析 %s 失敗: %v", path, err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return true
		}
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return true
		}
		for _, term := range forbiddenTerms {
			if strings.Contains(s, term) {
				pos := fset.Position(lit.Pos())
				t.Errorf("%s:%d 的字串字面量含命理術語 %q：%q\n"+
					"  計算層須與語言無關，文字一律由 lang 套件負責",
					path, pos.Line, term, s)
				return true
			}
		}
		return true
	})
}

// TestComputeLayersDoNotImportLang 計算層不得反向依賴 lang。
//
// 依賴方向必須是 lang → 計算層。這條規則由編譯器強制（反向會造成循環依賴），
// 此測試則讓違規在意圖階段就被攔下，而非等到編譯錯誤才發現。
func TestComputeLayersDoNotImportLang(t *testing.T) {
	const langPkg = "github.com/LukeLogix/destiny-core/lang"

	for _, dir := range computeLayers {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("讀取 %s 失敗: %v", dir, err)
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") {
				continue
			}
			path := filepath.Join(dir, name)

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				t.Fatalf("解析 %s 失敗: %v", path, err)
			}
			for _, imp := range f.Imports {
				p, _ := strconv.Unquote(imp.Path.Value)
				if p == langPkg {
					t.Errorf("%s 引用了 lang，違反依賴方向——lang 應依賴計算層而非反之", path)
				}
			}
		}
	}
}

// TestGanzhiDoesNotImportBazi 共用層不得依賴任一核心。
//
// ganzhi 若引用 bazi，紫微斗數日後加入時會被迫拖進八字邏輯。
func TestGanzhiDoesNotImportBazi(t *testing.T) {
	const baziPkg = "github.com/LukeLogix/destiny-core/bazi"

	entries, err := os.ReadDir("../ganzhi")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		path := filepath.Join("../ganzhi", e.Name())
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, imp := range f.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			if p == baziPkg {
				t.Errorf("%s 引用了 bazi，共用層不得依賴特定核心", path)
			}
		}
	}
}
