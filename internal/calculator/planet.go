package calculator

import (
	"github.com/mshafiee/swephgo"
)

// DefaultSiderealFlags ค่าแฟล็กสำหรับการคำนวณตำแหน่งดาวและอายนางศะแบบนิรายนะ (Swiss Ephemeris Sidereal)
const DefaultSiderealFlags = swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal

// CalculateSiderealPlanets คำนวณตำแหน่งดาวในระบบนิรายนะ (Lahiri Sidereal) โดยใช้แฟล็กมาตรฐาน
func CalculateSiderealPlanets(jd float64) []PlanetResult {
	return CalculatePlanets(jd, DefaultSiderealFlags, GetZodiacSigns())
}

// CalculateSiderealPlanetsWithNode คำนวณตำแหน่งดาวในระบบนิรายนะโดยเลือกได้ระหว่าง Mean Node หรือ True Node
func CalculateSiderealPlanetsWithNode(jd float64, useTrueNode bool) []PlanetResult {
	return CalculatePlanetsWithNode(jd, DefaultSiderealFlags, GetZodiacSigns(), useTrueNode)
}

// CalculatePlanets ประมวลผลตำแหน่งดาวทั้งหมดและคืนค่าเป็น slice ของผลลัพธ์ (Mean Node)
func CalculatePlanets(jd float64, flags int, signs []string) []PlanetResult {
	return CalculatePlanetsWithNode(jd, flags, signs, false)
}

// CalculatePlanetsWithNode ประมวลผลตำแหน่งดาวทั้งหมดโดยเลือกรูปแบบโหนดจันทร์ (Mean Node / True Node)
func CalculatePlanetsWithNode(jd float64, flags int, signs []string, useTrueNode bool) []PlanetResult {
	sweMu.Lock()
	defer sweMu.Unlock()

	swephgo.SetSidMode(swephgo.SeSidmLahiri, 0, 0)

	planetList := getPlanetListWithNode(useTrueNode)
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

	// คำนวณ Ketu เพิ่มเติม (อยู่ตรงข้าม Rahu 180 องศาเสมอ)
	ketuNodeID := swephgo.SeMeanNode
	ketuName := "Ketu (Calculated)"
	if useTrueNode {
		ketuNodeID = swephgo.SeTrueNode
		ketuName = "Ketu (True Node)"
	}
	resNode := swephgo.CalcUt(jd, ketuNodeID, flags, xx, serr)
	if resNode >= 0 {
		ketuLon := NormalizeAngle(xx[0] + 180.0)
		results = append(results, PlanetResult{
			Name:      ketuName,
			Longitude: ketuLon,
			Speed:     xx[3],
			PosStr:    FormatSignDMS(ketuLon, signs),
			NakStr:    FormatNakshatraString(ketuLon),
			IsRetro:   false,
		})
	}

	return results
}
