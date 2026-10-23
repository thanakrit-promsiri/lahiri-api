package main

import (
	"testing"
)

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()
	if cfg.Year <= 0 || cfg.Month <= 0 || cfg.Day <= 0 {
		t.Errorf("Invalid default date configuration: %+v", cfg)
	}
}

func TestFormatDMS(t *testing.T) {
	res := FormatDMS(15.5)
	expected := "15° 30' 00.00\""
	if res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}
}

func TestParseLongitude(t *testing.T) {
	signIdx, deg := ParseLongitude(45.0)
	if signIdx != 1 || deg != 15.0 {
		t.Errorf("ParseLongitude failed: got signIdx=%d, deg=%.2f", signIdx, deg)
	}
}

func TestGetNakshatra(t *testing.T) {
	nak, pada, _ := GetNakshatra(0.0)
	if nak.Name != "Ashwini" || pada != 1 || nak.Lord != "Ketu" {
		t.Errorf("Nakshatra calculation error for 0.0 degree: got %s (Pada %d)", nak.Name, pada)
	}
}

func TestCalculatePlanets(t *testing.T) {
	cfg := GetDefaultConfig()
	signs := GetZodiacSigns()
	jd, _ := setupEphemeris(cfg)
	flag := 0

	results := CalculatePlanets(jd, flag, signs)
	if len(results) == 0 {
		t.Errorf("Expected planetary results, got 0")
	}

	foundSun := false
	for _, p := range results {
		if p.Name == "Sun" {
			foundSun = true
			break
		}
	}
	if !foundSun {
		t.Errorf("Sun position not found in calculation results")
	}
}

func TestCalculateHouses(t *testing.T) {
	cfg := GetDefaultConfig()
	signs := GetZodiacSigns()
	jd, ayanamsa := setupEphemeris(cfg)

	houses := CalculateHouses(jd, cfg.Latitude, cfg.Longitude, cfg.HouseSystem, ayanamsa, signs)
	if len(houses) == 0 {
		t.Errorf("Expected house angles results, got 0")
	}

	if houses[0].Name != "Ascendant (ASC)" {
		t.Errorf("Expected first house angle to be Ascendant (ASC), got %s", houses[0].Name)
	}
}
