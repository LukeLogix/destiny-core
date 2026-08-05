# destiny-core

命理計算核心。負責「算」，不負責「解」——輸出結構化資料與可追溯的判定依據，解讀交給上層。

[![Go Reference](https://pkg.go.dev/badge/github.com/LukeLogix/destiny-core.svg)](https://pkg.go.dev/github.com/LukeLogix/destiny-core)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

```go
chart, _ := bazi.Compute(bazi.Birth{
    Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, taipei),
    Gender: bazi.Male,
}, bazi.Default())

// 庚午 辛巳 乙酉 辛巳，日主乙
```

## 為什麼再寫一個

現成的八字套件算得準，但少了幾樣東西：

| 缺口 | destiny-core |
|---|---|
| 不判身強身弱 | 兩套策略並行，判定不一致本身即為訊號 |
| 不管出生地 | 真太陽時三段修正，全球適用 |
| 中文寫死在程式裡 | 計算層零文字，i18n 由編譯器強制保證 |
| 用全域變數切換流派 | per-request Options，併發安全 |
| 不確定的地方照樣給定論 | 臨界與分歧皆以結構化資料交出 |
| 拖著整包曆法框架 | **零外部依賴**，只用 Go 標準庫 |

## 特色

### 零外部依賴

主 module 只用標準庫。節氣表以 `go:embed` 內嵌（29 KB），儒略日、干支、均時差全部自算。

部署只有一個執行檔，沒有資料檔漏帶或路徑解析失敗的問題。CI 以 `go list -deps -test` 把關（含測試檔的 import），靠檢查不靠自律。

### 神煞：寫機制，讓表自己長出來

`shensha` 收 62 種神煞，但幾乎沒有對照表。三合局系八個是同一件事在十二長生的八個取位；日干系是該干的祿位偏移或長生位；天乙貴人是「起坤佈干取合氣」的四條規則。**典籍表只出現在測試裡**，與推導互相比對——與節氣的作法同構。

這不只是省去抄表。實測中推導兩次指出《三命通會》四庫本的形近訛字（「癸以**己**」應為巳、「丙祿**己**」應為巳），也抓出四廢初版八組手抄索引錯了四組。

套件本身**與位置語意無關**——只認「基準組」與「待掃組」，不知道年月日時、大運流年為何物。呼叫方自行賦予意義：八字是柱位，六爻是爻位，紫微是宮位。

```go
hits := shensha.Detect(in, []ganzhi.SexagenaryIndex{y, m, d, h}, shensha.Default())
```

### 流派分歧是一等公民

同一個神煞各家取法不同是常態，不是例外。羊刃在祿前一辰還是帝旺？天羅地網要不要納音條件？十二歲君用神峰那組名還是洞微經那組？

`shensha` 把這些收成 11 個**主題**，各主題列出有出處的取法，預設選一派並在註解寫明理由：

```go
opt := shensha.Default().With(shensha.BladeAtProsperity, shensha.TianLuoBranchOnly)
```

主題表同時是計算來源與 API 中繼資料來源——`shensha.Topics()` 可直接餵給前端生成表單，不必手寫平行清單。

收錄門檻：**各派都須有出處，且各自在自己的體系內自洽**。同一本書自相矛盾者不列為流派，那是校勘問題，該判定並記錄理由。

### 雙軌互相守護

節氣時刻有兩條獨立路徑：

| 路徑 | 精度 | 角色 |
|---|---|---|
| 內嵌 JPL 表 | ground truth | 正式計算 |
| VSOP87 自算 | 中位 2.28 秒 | 驗證表的完整性 |

表若被竄改或損壞，自算會發現；移植若有誤，表會發現。**這個檢查每次跑測試都在做，不需要連上任何外部服務。**

### 誠實標記不確定性

不只算出答案，並標明哪裡不可靠：

```go
chart.Boundary.TermCritical   // 出生時刻臨近節氣分界，年月柱不可靠
chart.StrengthConsensus()     // ConsensusDisputed：兩套旺衰演算法判定分歧
```

臨界門檻為 **該筆節氣的不確定度 + 60 秒**（出生時間通常只記到分鐘）。ΔT（地球自轉修正）的未來值本質上不可知，故不確定度隨年代放大：2030 年約 1.4 秒 → 門檻約 61 秒；2090 年約 88 秒 → 門檻約 148 秒。

這對餵給 LLM 特別重要：AI 最容易在沒把握的地方講得斬釘截鐵，把疑慮一併交給它，它才有辦法誠實地說「這張盤有疑慮」。

### 有分歧的口徑一律可配置

各家說法不同的地方，不預設立場：

```go
opt := bazi.Default()
opt.LateZiKeepsDay = true                  // 晚子時派：23 時仍算當日
opt.HiddenStem = bazi.HiddenStemWithEarth  // 亥藏壬甲戊
opt.Terrain = bazi.TerrainSameBirth        // 陰陽同生同死
opt.ChildLimit = bazi.ChildLimitChina95    // 元亨利貞起運法
opt.SolarTime = bazi.TrueSolar             // 真太陽時（須同時提供 Birth.Longitude）
```

`SolarTime` 只要不是 `WallClock`，`Birth.Longitude` 就是必填——它是 `*float64`，因為零值 `0` 是合法經度（格林威治），無法用來表達「未提供」：

```go
lon := 121.5
birth := bazi.Birth{Time: t, Gender: bazi.Male, Longitude: &lon}
```

每張命盤都攜帶所用口徑（`chart.Options`），結果永遠可重現、可稽核。

## 安裝

```bash
go get github.com/LukeLogix/destiny-core
```

需要 Go 1.25 以上。

## 使用

```go
package main

import (
    "fmt"
    "time"

    "github.com/LukeLogix/destiny-core/bazi"
    "github.com/LukeLogix/destiny-core/lang"
)

func main() {
    // 用真實的 IANA 時區，不要用 time.FixedZone("UTC+8", 8*3600)——
    // 後者會丟失歷史夏令時間，台灣在 1945-1961、1974-1975、1979 曾行用 UTC+9。
    // 那些年份的出生資料用 FixedZone 會整整差一柱，而臨界標記不會示警。
    taipei, _ := time.LoadLocation("Asia/Taipei")

    chart, err := bazi.Compute(bazi.Birth{
        Time:   time.Date(1990, 5, 20, 10, 30, 0, 0, taipei),
        Gender: bazi.Male,
    }, bazi.Default())
    if err != nil {
        panic(err)
    }

    // 直接取索引運算，或本地化成文字
    lc := lang.Localize(chart, lang.ZhTW)
    fmt.Println(lc.Year.Sexagenary, lc.Month.Sexagenary,
        lc.Day.Sexagenary, lc.Hour.Sexagenary)
    fmt.Println("日主：", lc.DayMaster, "旺衰：", lc.Consensus)

    for _, w := range lc.Warnings {
        fmt.Println("⚠", w)
    }
}
```

更多用法見 [`example_test.go`](example_test.go)。

## 開發

```bash
./ci-local.sh    # 與 CI 同一組檢查：格式、vet、測試、零依賴、交叉驗證
```

CI 跑在自架 runner 上（`.github/workflows/ci.yml`）。

## 計算範圍

- 四柱（年柱以立春分界、月柱以十二節分界）
- 天干十神、地支藏干（本／中／餘氣）及各自十神
- 十二長生、納音
- 五行分佈與旺衰（兩套策略）
- 大運、流年（四種起運流派）
- 干支關係：天干五合、天干沖、六合、三合、半合、三會、六沖、相刑、六害
  - 相破另有 `Options.Relation.IncludeDestruction`，**預設關閉**（多數流派不採用）
  - 半合預設開啟，可由 `IncludeHalfTrinity` 關閉

**只報告關係構成，不判定合化成否**——化與不化取決於得令、引化之神、是否被沖破，各家分歧極大且無定論。

支援**干支年** 1900–2100。注意下界落在干支年而非公曆年：1900 年立春（2/4 13:51）之前出生歸 1899 干支年，會回 `ErrYearOutOfRange`。超出範圍一律回錯誤，不靜默回傳可能錯誤的結果。

## 驗證

```
節氣精度 vs NASA JPL Horizons     2412 筆，1940-2039 平均偏差 1.1 秒
四柱 vs tyme4go                  19677 筆，非臨界不符 0 筆
四柱 vs lunar-go                   348 筆，零不符
往返驗證（八字反查公曆）              52 筆，全數回到原日期
```

唯一那筆與 tyme4go 的分歧（2098-10-08）經診斷為 ΔT 模型差異，且已被臨界標記正確攔下。

驗證工具在 [`test/`](test/) 獨立 module，避免僅作第二意見的依賴污染下游的 `go.mod`。

## 架構

```
internal/calendar/   曆法：內嵌節氣表 + VSOP87 自算 + 均時差
ganzhi/              干支基礎：具名索引型別、五行生剋、沖刑合會、十二長生、納音
shensha/             神煞：位置無關的偵測，62 種神煞、11 個流派主題
bazi/                八字：四柱、十神、藏干、大運、旺衰、神煞組裝
lang/                zh-TW / zh-CN 文字，依賴計算層而非相反
```

五條硬規則：

| # | 規則 | 如何把關 |
|---|---|---|
| 1 | `lang` 依賴計算層，反向絕不成立 | `lang/arch_test.go` 以 AST 解析驗證 |
| 2 | 口徑走 per-request Options，禁止套件級全域變數 | code review；`Compute` 會把解析後的策略寫回 `chart.Options` |
| 3 | `internal/calendar` 以 Go 的 internal 機制封裝 | 編譯器 |
| 4 | 共用層不得反向依賴任一核心 | `lang/arch_test.go` |
| 5 | 主 module 不得引入外部依賴 | CI 的 `go list -deps -test` |

第 1、4 條的檢查方式是用 AST 掃描計算層的字串字面量，確認不含命理術語。人會忘記，測試不會。

## 規劃中

紫微斗數、西洋星盤、六爻——各自獨立設計，共用 `internal/calendar` 與 `ganzhi`。

## 授權

MIT。第三方來源見 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。
