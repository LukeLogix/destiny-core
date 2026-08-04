package ganzhi

// 三十種納音的具名常數。
//
// 用拼音而非意譯：納音名為專名，硬翻英文會製造假的語意且無法反查原名
// （「路旁土」譯作 RoadsideEarth 讀者仍不知所指）。文字一律由 lang 提供。
//
// 具名的另一好處是可讀的比對——chart.Year.Sound == ganzhi.LuPangTu
// 勝於與魔術數字 3 相比。
const (
	HaiZhongJin   SoundIndex = iota // 海中金
	LuZhongHuo                      // 爐中火
	DaLinMu                         // 大林木
	LuPangTu                        // 路旁土
	JianFengJin                     // 劍鋒金
	ShanTouHuo                      // 山頭火
	JianXiaShui                     // 澗下水
	ChengTouTu                      // 城頭土
	BaiLaJin                        // 白蠟金
	YangLiuMu                       // 楊柳木
	QuanZhongShui                   // 泉中水
	WuShangTu                       // 屋上土
	PiLiHuo                         // 霹靂火
	SongBoMu                        // 松柏木
	ChangLiuShui                    // 長流水
	ShaZhongJin                     // 沙中金
	ShanXiaHuo                      // 山下火
	PingDiMu                        // 平地木
	BiShangTu                       // 壁上土
	JinBoJin                        // 金箔金
	FuDengHuo                       // 覆燈火
	TianHeShui                      // 天河水
	DaYiTu                          // 大驛土
	ChaiChuanJin                    // 釵釧金
	SangZheMu                       // 桑柘木
	DaXiShui                        // 大溪水
	ShaZhongTu                      // 沙中土
	TianShangHuo                    // 天上火
	ShiLiuMu                        // 石榴木
	DaHaiShui                       // 大海水
)
