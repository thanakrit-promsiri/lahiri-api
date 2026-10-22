package main

import (
	"testing"
)

// ทดสอบฟังก์ชันแปลงองศาไม่ให้เกิน 360 หรือติดลบ
func TestNormalizeAngle(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{370.0, 10.0},
		{-10.0, 350.0},
		{45.0, 45.0},
		{725.0, 5.0},
	}

	for _, test := range tests {
		result := NormalizeAngle(test.input)
		if result != test.expected {
			t.Errorf("NormalizeAngle(%.2f) = %.2f; expected %.2f", test.input, result, test.expected)
		}
	}
}

// ทดสอบฟังก์ชันแยกองศาและราศี
func TestParseLongitude(t *testing.T) {
	signIdx, degree := ParseLongitude(45.5)
	if signIdx != 1 || degree != 15.5 {
		t.Errorf("ParseLongitude(45.5) got signIdx=%d, degree=%.2f; expected 1, 15.5", signIdx, degree)
	}
}

// ทดสอบการคำนวณ Julian Day และค่า Ayanamsa (Lahiri)
func TestSetupEphemeris(t *testing.T) {
	cfg := Config{
		Year:      1988,
		Month:     12,
		Day:       23,
		HourThai:  10.25,
		TimeZone:  7.0,
		Latitude:  13.7563,
		Longitude: 100.5018,
	}

	jd, ayanamsa := setupEphemeris(cfg)

	if jd <= 0 {
		t.Errorf("Invalid Julian Day: %.2f", jd)
	}

	if ayanamsa < 23.0 || ayanamsa > 24.0 {
		t.Errorf("Ayanamsa out of expected range: %.2f", ayanamsa)
	}
}

// ทดสอบความถูกต้องของอาเรย์รายชื่อราศี
func TestGetZodiacSigns(t *testing.T) {
	signs := GetZodiacSigns()
	if len(signs) != 12 {
		t.Errorf("Expected 12 zodiac signs, got %d", len(signs))
	}
	if signs[0] != "Aries" || signs[11] != "Pisces" {
		t.Errorf("Zodiac signs order is incorrect")
	}
}

// ทดสอบการคำนวณระบบภพ Whole Sign จากฟังก์ชัน processHouses
func TestProcessHouses(t *testing.T) {
	cfg := Config{
		Year:        1988,
		Month:       12,
		Day:         23,
		HourThai:    10.25,
		TimeZone:    7.0,
		Latitude:    13.7563,
		Longitude:   100.5018,
		HouseSystem: int('W'),
	}

	jd, ayanamsa := setupEphemeris(cfg)
	signs := GetZodiacSigns()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("processHouses panicked: %v", r)
		}
	}()

	processHouses(jd, cfg.Latitude, cfg.Longitude, cfg.HouseSystem, ayanamsa, signs)
}
