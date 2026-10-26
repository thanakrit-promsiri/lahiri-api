package calculator

import (
	"github.com/mshafiee/swephgo"
)

// CalculateHouses ประมวลผลลัคนา, MC, DSC, IC และคืนค่าผลลัพธ์
func CalculateHouses(jd float64, lat float64, lon float64, houseSys int, ayanamsa float64, signs []string) []HouseResult {
	cusps := make([]float64, 13)
	ascmc := make([]float64, 10)

	resHouses := swephgo.Houses(jd, lat, lon, houseSys, cusps, ascmc)
	if resHouses < 0 {
		return nil
	}

	asc := NormalizeAngle(ascmc[0] - ayanamsa)
	mc := NormalizeAngle(ascmc[1] - ayanamsa)
	dsc := NormalizeAngle(asc + 180.0)
	ic := NormalizeAngle(mc + 180.0)

	angles := []struct {
		Name string
		Val  float64
	}{
		{Name: "Ascendant (ASC)", Val: asc},
		{Name: "Midheaven (MC)", Val: mc},
		{Name: "Descendant (DSC)", Val: dsc},
		{Name: "Imum Coeli (IC)", Val: ic},
	}

	var results []HouseResult
	for _, a := range angles {
		results = append(results, HouseResult{
			Name:   a.Name,
			Degree: a.Val,
			PosStr: FormatSignDMS(a.Val, signs),
			NakStr: FormatNakshatraString(a.Val),
		})
	}

	return results
}
