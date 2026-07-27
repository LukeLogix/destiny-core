# test

驗證用的獨立 module。刻意與主 module 分離，避免 `lunar-go` 這類「只用來當第二意見」的依賴污染下游使用者的 `go.mod`。

## 工具

| 指令 | 用途 |
|---|---|
| `go run ./cmd/genfixture` | 由 NASA JPL Horizons 取太陽視黃經，反推節氣時刻，產生 golden fixture |
| `go run ./cmd/crosscheck` | `tyme4go` 與 `lunar-go` 交叉比對四柱，驗證子時換日口徑 |
| `go run ./cmd/pillarpoc` | 自實作四柱推算與 `tyme4go` 全量對照 |

## 已驗證結果

**節氣精度**（`genfixture`，對照 JPL Horizons）

| 年份 | 節氣 | tyme4go | JPL | 誤差 |
|---|---|---|---|---|
| 1950 | 春分 | 12:35:06 | 12:35:08 | −2.8 秒 |
| 1990 | 立夏 | 02:35:26 | 02:35:27 | −1.2 秒 |
| 2024 | 立春 | 16:27:07 | 16:27:08 | −1.7 秒 |
| 2050 | 冬至 | 18:52:15 | 18:52:19 | −4.4 秒 |

**雙庫交叉**（`crosscheck`）

- lunar-go `sect=2` vs tyme4go：18291 筆，差異 5226 筆，100% 集中於 23:00–23:59
- lunar-go `sect=1` vs tyme4go：18291 筆，**差異 0 筆**

**自實作 PoC**（`pillarpoc`）

- 全量對照 185724 筆，不一致 0 筆

## 注意

`genfixture` 會呼叫外部 API。fixture 一次產生後即固化，CI 不應依賴網路。
