package calculator

import (
	"testing"
	"time"
)

func TestCalculateDignity(t *testing.T) {
	tests := []struct {
		planet      string
		lon         float64
		wantDignity string
		wantThai    string
	}{
		{"Sun", 10.0, "Exalted", "อุจจ์"},         // Aries 10°
		{"Sun", 190.0, "Debilitated", "นิจ"},      // Libra 10°
		{"Sun", 125.0, "Moolatrikona", "มูลไตรโกณ"}, // Leo 5°
		{"Sun", 145.0, "Own Sign", "เกษตร"},       // Leo 25°
		{"Sun", 247.81, "Friend", "มิตร"},         // Sagittarius (Lord Jupiter = Friend)
		{"Mars", 298.0, "Exalted", "อุจจ์"},        // Capricorn 28°
		{"Jupiter", 95.0, "Exalted", "อุจจ์"},      // Cancer 5°
		{"Venus", 357.0, "Exalted", "อุจจ์"},       // Pisces 27°
		{"Saturn", 200.0, "Exalted", "อุจจ์"},      // Libra 20°
	}

	for _, tc := range tests {
		gotEng, gotThai := CalculateDignity(tc.planet, tc.lon)
		if gotEng != tc.wantDignity || gotThai != tc.wantThai {
			t.Errorf("CalculateDignity(%s, %.2f) = (%s, %s); want (%s, %s)",
				tc.planet, tc.lon, gotEng, gotThai, tc.wantDignity, tc.wantThai)
		}
	}
}

func TestIsPlanetCombust(t *testing.T) {
	sunLon := 247.81
	// พุธห่างอาทิตย์ 5 องศา (ไม่พักร์ เกณฑ์ 14°) -> ต้องดับ (Combust = true)
	if !IsPlanetCombust("Mercury", 252.0, sunLon, false) {
		t.Errorf("Expected Mercury at 252° (Sun 247.81°) to be combust")
	}
	// พุธห่างอาทิตย์ 20 องศา -> ไม่ดับ (Combust = false)
	if IsPlanetCombust("Mercury", 268.0, sunLon, false) {
		t.Errorf("Expected Mercury at 268° (Sun 247.81°) not to be combust")
	}
	// อาทิตย์ไม่ดับตัวเอง
	if IsPlanetCombust("Sun", sunLon, sunLon, false) {
		t.Errorf("Expected Sun never to be combust")
	}
	// ราหูไม่ดับ
	if IsPlanetCombust("Rahu (Mean Node)", sunLon, sunLon, true) {
		t.Errorf("Expected Rahu never to be combust")
	}
}

func TestCalculateVimshottariDasha(t *testing.T) {
	// ดวงจันทร์ที่ 66.724615° (ราศีเมถุน 6° 43' 28.61" -> นักษัตร Ardra เจ้าฤกษ์ Rahu)
	moonLon := 66.72461481473177
	target := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)

	dasha := CalculateVimshottariDasha(moonLon, 1988, 12, 23, 10.35, 7.0, target)

	if dasha.MoonNakshatra != "Ardra" {
		t.Errorf("Expected MoonNakshatra 'Ardra', got '%s'", dasha.MoonNakshatra)
	}
	if dasha.MoonNakshatraThai != "อารทรา" {
		t.Errorf("Expected MoonNakshatraThai 'อารทรา', got '%s'", dasha.MoonNakshatraThai)
	}
	if len(dasha.MahadashaList) != 9 {
		t.Fatalf("Expected 9 Mahadasha items, got %d", len(dasha.MahadashaList))
	}
	if dasha.MahadashaList[0].Lord != "Rahu" {
		t.Errorf("Expected first Mahadasha lord 'Rahu', got '%s'", dasha.MahadashaList[0].Lord)
	}
	if dasha.CurrentDasha.MahadashaLord == "" || dasha.CurrentDasha.AntardashaLord == "" {
		t.Errorf("Expected non-empty CurrentDasha lords, got %+v", dasha.CurrentDasha)
	}
}
