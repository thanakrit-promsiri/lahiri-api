package main

import (
	"testing"
)

func TestPrintVersion(t *testing.T) {
	// ทดสอบว่าฟังก์ชันทำงานปกติและไม่มี Panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printVersion panicked: %v", r)
		}
	}()
	printVersion()
}

func TestPrintLocationInfo(t *testing.T) {
	cfg := Config{
		Year:      2026,
		Month:     10,
		Day:       23,
		HourThai:  12.0,
		TimeZone:  7.0,
		Latitude:  13.7563,
		Longitude: 100.5018,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printLocationInfo panicked: %v", r)
		}
	}()
	printLocationInfo(cfg)
}

func TestSetupEphemeris(t *testing.T) {
	cfg := Config{
		Year:      2026,
		Month:     10,
		Day:       23,
		HourThai:  12.0,
		TimeZone:  7.0,
		Latitude:  13.7563,
		Longitude: 100.5018,
	}

	jd, ayanamsa := setupEphemeris(cfg)
	if jd <= 0 {
		t.Errorf("Expected valid Julian Day, got %f", jd)
	}
	if ayanamsa == 0 {
		t.Errorf("Expected valid Ayanamsa value")
	}
}

func TestGetPlanetList(t *testing.T) {
	planets := getPlanetList()
	if len(planets) == 0 {
		t.Fatalf("Expected planet list to not be empty")
	}

	foundSun := false
	for _, p := range planets {
		if p.Name == "Sun" {
			foundSun = true
			break
		}
	}
	if !foundSun {
		t.Errorf("Expected 'Sun' in planet list")
	}
}

func TestPrintWholeSignHouses(t *testing.T) {
	signs := []string{"เมษ", "พฤษภ", "มิถุน", "กรกฎ", "สิงห์", "กันย์", "ตุลย์", "พิจิก", "ธนู", "มังกร", "กุมภ์", "มีน"}
	ascDegree := 272.07

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printWholeSignHouses panicked: %v", r)
		}
	}()
	printWholeSignHouses(ascDegree, signs)
}

func TestCleanup(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("cleanup panicked: %v", r)
		}
	}()
	cleanup()
}
