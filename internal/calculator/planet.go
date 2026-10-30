package calculator

import (
	"github.com/mshafiee/swephgo"
)

// CalculatePlanets ประมวลผลตำแหน่งดาวทั้งหมดและคืนค่าเป็น slice ของผลลัพธ์
func CalculatePlanets(jd float64, flags int, signs []string) []PlanetResult {
	sweMu.Lock()
	defer sweMu.Unlock()

	swephgo.SetSidMode(swephgo.SeSidmLahiri, 0, 0)

	planetList := getPlanetList()
	results := make([]PlanetResult, 0, len(planetList)+1)
	xx := make([]float64, 6)
	serr := make([]byte, 256)

	for _, p := range planetList {
		res := swephgo.CalcUt(jd, p.ID, flags, xx, serr)
		if res < 0 {
			continue
		}

		lon := NormalizeAngle(xx[0])
		speed := xx[3]
		results = append(results, PlanetResult{
			Name:      p.Name,
			Longitude: lon,
			Speed:     speed,
			PosStr:    FormatSignDMS(lon, signs),
			NakStr:    FormatNakshatraString(lon),
			IsRetro:   speed < 0,
		})
	}

	// คำนวณ Ketu เพิ่มเติม
	resNode := swephgo.CalcUt(jd, swephgo.SeMeanNode, flags, xx, serr)
	if resNode >= 0 {
		ketuLon := NormalizeAngle(xx[0] + 180.0)
		results = append(results, PlanetResult{
			Name:      "Ketu (Calculated)",
			Longitude: ketuLon,
			Speed:     xx[3],
			PosStr:    FormatSignDMS(ketuLon, signs),
			NakStr:    FormatNakshatraString(ketuLon),
			IsRetro:   false,
		})
	}

	return results
}
