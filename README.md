# destiny-core

命理計算核心。輸出結構化資料與可追溯的判定依據；負責「算」，不負責「解」。

## 現況

| 模組 | 內容 | 狀態 |
|---|---|---|
| `bazi/` | 八字四柱、十神、藏干、五行旺衰、大運流年、沖刑合會 | 實作中 |
| `ziwei/` | 紫微斗數 | 規劃中 |
| `astrology/` | 西洋星盤 | 規劃中 |
| `liuyao/` | 六爻 | 規劃中 |

## 結構

```
internal/calendar/   曆法、節氣、儒略日、農曆   ← 各核心共用，對外封裝
ganzhi/              干支、五行基礎型別         ← 八字/紫微/六爻共用
lang/                i18n bundle (zh-TW/zh-CN)  ← 共用
bazi/                八字
```

### 設計原則

1. `lang` 依賴計算層，反向不成立 — 計算層零文字，由編譯器強制
2. 口徑走 per-request `Options`，禁止套件級全域變數
3. 不確定性顯式標記，不假裝精確

設計文件見 [life-chat-project](../life-chat-project/docs/specs/)。
