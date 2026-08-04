package bazi

// 列舉的穩定識別字。規則與作法見 ganzhi/id.go 與 lang/id_test.go。

var genderIDs = [3]string{
	GenderUnset: "gender_unset",
	Male:        "male",
	Female:      "female",
}

func (g Gender) ID() string { return pickID(genderIDs[:], int(g)) }

var solarTimeModeIDs = [3]string{
	WallClock:     "wall_clock",
	LongitudeOnly: "longitude_only",
	TrueSolar:     "true_solar",
}

func (m SolarTimeMode) ID() string { return pickID(solarTimeModeIDs[:], int(m)) }

var childLimitSectIDs = [4]string{
	ChildLimitDefault:    "default",
	ChildLimitChina95:    "china95",
	ChildLimitLunarSect1: "lunar_sect1",
	ChildLimitLunarSect2: "lunar_sect2",
}

func (s ChildLimitSect) ID() string { return pickID(childLimitSectIDs[:], int(s)) }

var hiddenStemSectIDs = [2]string{
	HiddenStemStandard:  "standard",
	HiddenStemWithEarth: "with_earth",
}

func (s HiddenStemSect) ID() string { return pickID(hiddenStemSectIDs[:], int(s)) }

var hiddenStemTypeIDs = [3]string{
	PrimaryQi:  "primary_qi",
	MiddleQi:   "middle_qi",
	ResidualQi: "residual_qi",
}

func (t HiddenStemType) ID() string { return pickID(hiddenStemTypeIDs[:], int(t)) }

var tenGodIDs = [TenGodCount]string{
	Peer:             "peer",
	Rival:            "rival",
	Output:           "output",
	Hurting:          "hurting",
	IndirectWealth:   "indirect_wealth",
	DirectWealth:     "direct_wealth",
	SevenKilling:     "seven_killing",
	DirectOfficer:    "direct_officer",
	IndirectResource: "indirect_resource",
	DirectResource:   "direct_resource",
}

func (g TenGod) ID() string { return pickID(tenGodIDs[:], int(g)) }

var verdictIDs = [3]string{
	VerdictStrong:  "strong",
	VerdictWeak:    "weak",
	VerdictNeutral: "neutral",
}

func (v Verdict) ID() string { return pickID(verdictIDs[:], int(v)) }

var strategyIDIDs = [2]string{
	StrategyWeighted:  "weighted",
	StrategyClassical: "classical",
}

func (s StrategyID) ID() string { return pickID(strategyIDIDs[:], int(s)) }

var reasonCodeIDs = [reasonCodeCount]string{
	ReasonInSeason:     "in_season",
	ReasonNotInSeason:  "not_in_season",
	ReasonRooted:       "rooted",
	ReasonNoRoot:       "no_root",
	ReasonAllied:       "allied",
	ReasonNotAllied:    "not_allied",
	ReasonSupportScore: "support_score",
	ReasonDrainScore:   "drain_score",
}

func (c ReasonCode) ID() string { return pickID(reasonCodeIDs[:], int(c)) }

var consensusIDs = [3]string{
	ConsensusStrong:   "strong",
	ConsensusWeak:     "weak",
	ConsensusDisputed: "disputed",
}

func (c Consensus) ID() string { return pickID(consensusIDs[:], int(c)) }

// pickID 越界回空字串，與 lang 的 pick 同一約定。
func pickID(tbl []string, i int) string {
	if i < 0 || i >= len(tbl) {
		return ""
	}
	return tbl[i]
}
