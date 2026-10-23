package main

import (
	"testing"

	"github.com/mshafiee/swephgo"
)

func TestCalculatePlanetsAndHouses(t *testing.T) {
	// 1. กำหนดค่าจำลองสำหรับทดสอบ (วันที่ 23 ตุลาคม 2026, 12:00 UTC)
	year, month, day := 2026, 10, 23
	hourUTC := 5.0 // สมมติ 12:00 Thai - 7 ชม. TimeZone
	jd := swephgo.Julday(year, month, day, hourUTC, swephgo.SeGregCal)

	// ตั้งค่า Ayanamsa เป็น Lahiri
	swephgo.SetSidMode(swephgo.SeSidmLahiri, 0, 0)
	ayanamsa := swephgo.GetAyanamsa(jd)
	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	signs := GetZodiacSigns()

	// 2. ทดสอบฟังก์ชัน CalculatePlanets
	planets := CalculatePlanets(jd, flag, signs)
	if len(planets) == 0 {
		t.Fatalf("Expected planets calculation results, got empty slice")
	}

	// ตรวจสอบว่ามีข้อมูลดาวพื้นฐานครบถ้วน (เช่น Sun ควรมีอยู่ในลิสต์)
	foundSun := false
	for _, p := range planets {
		if p.Name == "Sun" {
			foundSun = true
			if p.Longitude < 0 || p.Longitude >= 360 {
				t.Errorf("Invalid Sun longitude: %f", p.Longitude)
			}
			if p.PosStr == "" {
				t.Errorf("Expected PosStr for Sun, got empty string")
			}
		}
	}
	if !foundSun {
		t.Errorf("Expected 'Sun' in planets calculation results")
	}

	// 3. ทดสอบฟังก์ชัน CalculateHouses
	lat, lon := 13.7563, 100.5018 // กรุงเทพฯ
	houseSys := int('W')
	houses := CalculateHouses(jd, lat, lon, houseSys, ayanamsa, signs)

	if len(houses) != 4 {
		t.Fatalf("Expected 4 main angles (ASC, MC, DSC, IC), got %d", len(houses))
	}

	// ตรวจสอบชื่อมุมหลักตัวแรกต้องเป็น Ascendant (ASC)
	if houses[0].Name != "Ascendant (ASC)" {
		t.Errorf("Expected first house result to be 'Ascendant (ASC)', got '%s'", houses[0].Name)
	}
	if houses[0].Degree < 0 || houses[0].Degree >= 360 {
		t.Errorf("Invalid Ascendant degree: %f", houses[0].Degree)
	}
}