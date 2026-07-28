# 第三方來源聲明

本專案主 module 零外部依賴，但部分資料與演算法源自他人成果，於此聲明。

## internal/calendar/astro_data.go

天文係數表（VSOP87 地球黃經截斷級數 2666 項、章動 50 項、ΔT 歷史表 157 項）
轉錄自 [6tail/tyme4go](https://github.com/6tail/tyme4go) v1.5.0，
其演算法源自許劍偉的壽星萬年曆。

    MIT License
    Copyright (c) 2024 6tail

僅取太陽相關係數；月亮的 10608 項與古代曆法修正表未納入。
係數為程式自動轉錄而非人工抄寫，正確性由 JPL Horizons golden data 交叉驗證。

VSOP87 本身為公開發表的天文算法（Bretagnon & Francou, 1988）。
均時差與黃赤交角公式取自 Jean Meeus《Astronomical Algorithms》第 22、25、28 章。

## internal/calendar/solarterms.bin

節氣時刻表由 [NASA JPL Horizons](https://ssd.jpl.nasa.gov/horizons/) 的
太陽視黃經反推產生，共 2412 筆（1900-2100 年 × 12 個節）。

JPL Horizons 為美國政府公開資料，不受著作權限制。

## test/ module

交叉驗證用的 [6tail/tyme4go](https://github.com/6tail/tyme4go) 與
[6tail/lunar-go](https://github.com/6tail/lunar-go) 皆為 MIT License。
兩者僅作為獨立第二意見，不進入主 module 的依賴。
