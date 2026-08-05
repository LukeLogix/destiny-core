package ganzhi

// 列舉的穩定識別字。
//
// ID 不隨語系變動，供上層作為分支、篩選、儲存的鍵；顯示文字由 lang 提供。
// 規則為「Go 常數名的 snake_case」，常數名以型別名開頭者去除該前綴。
//
// Go 在執行期取不到常數名，故表只能手寫；手寫就會忘記。
// lang/id_test.go 以 AST 解析原始碼補上這道保證——人會忘記，測試不會。

var elementIDs = [ElementCount]string{
	Wood:  "wood",
	Fire:  "fire",
	Earth: "earth",
	Metal: "metal",
	Water: "water",
}

func (e Element) ID() string { return pickID(elementIDs[:], int(e)) }

var polarityIDs = [2]string{
	Yang: "yang",
	Yin:  "yin",
}

func (p Polarity) ID() string { return pickID(polarityIDs[:], int(p)) }

var terrainIDs = [TerrainCount]string{
	LongLife:   "long_life",
	Bath:       "bath",
	Cap:        "cap",
	Officer:    "officer",
	Prosperity: "prosperity",
	Decline:    "decline",
	Sick:       "sick",
	Death:      "death",
	Tomb:       "tomb",
	Void:       "void",
	Fetus:      "fetus",
	Nurture:    "nurture",
}

func (t Terrain) ID() string { return pickID(terrainIDs[:], int(t)) }

var terrainSectIDs = [2]string{
	TerrainYinReverse: "yin_reverse",
	TerrainSameBirth:  "same_birth",
}

func (s TerrainSect) ID() string { return pickID(terrainSectIDs[:], int(s)) }

var relationKindIDs = [10]string{
	StemHarmony: "stem_harmony",
	StemClash:   "stem_clash",
	SixHarmony:  "six_harmony",
	Trinity:     "trinity",
	HalfTrinity: "half_trinity",
	Directional: "directional",
	Clash:       "clash",
	Punishment:  "punishment",
	Harm:        "harm",
	Destruction: "destruction",
}

func (k RelationKind) ID() string { return pickID(relationKindIDs[:], int(k)) }

var soundIDs = [SoundCount]string{
	HaiZhongJin:   "hai_zhong_jin",
	LuZhongHuo:    "lu_zhong_huo",
	DaLinMu:       "da_lin_mu",
	LuPangTu:      "lu_pang_tu",
	JianFengJin:   "jian_feng_jin",
	ShanTouHuo:    "shan_tou_huo",
	JianXiaShui:   "jian_xia_shui",
	ChengTouTu:    "cheng_tou_tu",
	BaiLaJin:      "bai_la_jin",
	YangLiuMu:     "yang_liu_mu",
	QuanZhongShui: "quan_zhong_shui",
	WuShangTu:     "wu_shang_tu",
	PiLiHuo:       "pi_li_huo",
	SongBoMu:      "song_bo_mu",
	ChangLiuShui:  "chang_liu_shui",
	ShaZhongJin:   "sha_zhong_jin",
	ShanXiaHuo:    "shan_xia_huo",
	PingDiMu:      "ping_di_mu",
	BiShangTu:     "bi_shang_tu",
	JinBoJin:      "jin_bo_jin",
	FuDengHuo:     "fu_deng_huo",
	TianHeShui:    "tian_he_shui",
	DaYiTu:        "da_yi_tu",
	ChaiChuanJin:  "chai_chuan_jin",
	SangZheMu:     "sang_zhe_mu",
	DaXiShui:      "da_xi_shui",
	ShaZhongTu:    "sha_zhong_tu",
	TianShangHuo:  "tian_shang_huo",
	ShiLiuMu:      "shi_liu_mu",
	DaHaiShui:     "da_hai_shui",
}

func (s SoundIndex) ID() string { return pickID(soundIDs[:], int(s)) }

// pickID 越界回空字串，與 lang 的 pick 同一約定：不 panic，讓缺漏顯而易見。
func pickID(tbl []string, i int) string {
	if i < 0 || i >= len(tbl) {
		return ""
	}
	return tbl[i]
}
