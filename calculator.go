package main

import (
	"github.com/mshafiee/swephgo"
)

// PlanetResult เก็บผลลัพธ์การคำนวณดาวแต่ละดวง
type PlanetResult struct {
	Name      string
	Longitude float64
	Speed     float64
	PosStr    string
	NakStr    string
	IsRetro   bool
}

// HouseResult เก็บผลลัพธ์มุมลัคนาและภพ
type HouseResult struct {
	Name   string
	Degree float64
	PosStr string
	NakStr string
}

// CalculatePlanets ประมวลผลตำแหน่งดาวทั้งหมดและคืนค่าเป็น slice ของผลลัพธ์
func CalculatePlanets(jd float64, flags int, signs []string) []PlanetResult {
	var results []PlanetResult
	xx := make([]float64, 6)
	serr := make([]byte, 256)

	for _, p := range getPlanetList() {
		res := swephgo.CalcUt(jd, p.ID, flags, xx, serr)
		if res < 0 {
			continue
		}

		lon := xx[0]
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
