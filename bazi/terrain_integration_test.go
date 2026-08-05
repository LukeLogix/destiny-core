package bazi

import (
	"testing"

	"github.com/LukeLogix/destiny-core/ganzhi"
)

// TestPillarCarriesTerrainAndSound 命盤各柱須帶十二長生與納音。
func TestPillarCarriesTerrainAndSound(t *testing.T) {
	c := mustCompute(t, birthAt(1990, 5, 20, 10, 30), Default())

	// 日主乙木：年支午為長生、日支酉為絕
	if c.Year.Terrain != ganzhi.LongLife {
		t.Errorf("年柱地勢為 %d，日主乙見午應為長生", c.Year.Terrain)
	}
	if c.Day.Terrain != ganzhi.Void {
		t.Errorf("日柱地勢為 %d，日主乙見酉應為絕", c.Day.Terrain)
	}
	// 納音
	if c.Year.Sound.Element() != ganzhi.Earth {
		t.Errorf("年柱庚午納音五行為 %d，應為土（路旁土）", c.Year.Sound.Element())
	}
	if c.Day.Sound.Element() != ganzhi.Water {
		t.Errorf("日柱乙酉納音五行為 %d，應為水（泉中水）", c.Day.Sound.Element())
	}
}
