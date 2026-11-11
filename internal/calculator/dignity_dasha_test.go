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
		// Sun
		{"Sun", 10.0, "Exalted", "อุจจ์"},
		{"Sun", 190.0, "Debilitated", "นิจ"},
		{"Sun", 125.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Sun", 145.0, "Own Sign", "เกษตร"},
		{"Sun", 247.81, "Friend", "มิตร"},
		{"Sun", 275.0, "Enemy", "ศัตรู"},
		{"Sun", 65.0, "Neutral", "ปกติ"},

		// Moon
		{"Moon", 33.0, "Exalted", "อุจจ์"},
		{"Moon", 35.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Moon", 215.0, "Debilitated", "นิจ"},
		{"Moon", 95.0, "Own Sign", "เกษตร"},
		{"Moon", 125.0, "Friend", "มิตร"},
		{"Moon", 5.0, "Neutral", "ปกติ"},

		// Mars
		{"Mars", 298.0, "Exalted", "อุจจ์"},
		{"Mars", 95.0, "Debilitated", "นิจ"},
		{"Mars", 5.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Mars", 15.0, "Own Sign", "เกษตร"},
		{"Mars", 215.0, "Own Sign", "เกษตร"},
		{"Mars", 65.0, "Enemy", "ศัตรู"},
		{"Mars", 35.0, "Neutral", "ปกติ"},

		// Mercury
		{"Mercury", 155.0, "Exalted", "อุจจ์"},
		{"Mercury", 166.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Mercury", 175.0, "Own Sign", "เกษตร"},
		{"Mercury", 65.0, "Own Sign", "เกษตร"},
		{"Mercury", 335.0, "Debilitated", "นิจ"},
		{"Mercury", 95.0, "Enemy", "ศัตรู"},
		{"Mercury", 35.0, "Friend", "มิตร"},
		{"Mercury", 5.0, "Neutral", "ปกติ"},

		// Jupiter
		{"Jupiter", 95.0, "Exalted", "อุจจ์"},
		{"Jupiter", 275.0, "Debilitated", "นิจ"},
		{"Jupiter", 245.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Jupiter", 255.0, "Own Sign", "เกษตร"},
		{"Jupiter", 335.0, "Own Sign", "เกษตร"},
		{"Jupiter", 65.0, "Enemy", "ศัตรู"},
		{"Jupiter", 305.0, "Neutral", "ปกติ"},

		// Venus
		{"Venus", 357.0, "Exalted", "อุจจ์"},
		{"Venus", 155.0, "Debilitated", "นิจ"},
		{"Venus", 185.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Venus", 200.0, "Own Sign", "เกษตร"},
		{"Venus", 35.0, "Own Sign", "เกษตร"},
		{"Venus", 125.0, "Enemy", "ศัตรู"},
		{"Venus", 5.0, "Neutral", "ปกติ"},

		// Saturn
		{"Saturn", 200.0, "Exalted", "อุจจ์"},
		{"Saturn", 5.0, "Debilitated", "นิจ"},
		{"Saturn", 305.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Saturn", 325.0, "Own Sign", "เกษตร"},
		{"Saturn", 275.0, "Own Sign", "เกษตร"},
		{"Saturn", 125.0, "Enemy", "ศัตรู"},
		{"Saturn", 245.0, "Neutral", "ปกติ"},

		// Nodes & Outer planets
		{"Rahu (Mean Node)", 35.0, "Exalted", "อุจจ์"},
		{"Rahu", 65.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Rahu", 305.0, "Own Sign", "เกษตร"},
		{"Rahu", 215.0, "Debilitated", "นิจ"},
		{"Rahu", 5.0, "Neutral", "ปกติ"},
		{"Ketu (Calculated)", 215.0, "Exalted", "อุจจ์"},
		{"Ketu", 245.0, "Moolatrikona", "มูลไตรโกณ"},
		{"Ketu", 35.0, "Debilitated", "นิจ"},
		{"Ketu", 5.0, "Neutral", "ปกติ"},
		{"Uranus", 10.0, "Neutral", "ปกติ"},
		{"Neptune", 10.0, "Neutral", "ปกติ"},
		{"Pluto", 10.0, "Neutral", "ปกติ"},
		{"Unknown", 10.0, "Neutral", "ปกติ"},
	}

	for _, tc := range tests {
		gotEng, gotThai := CalculateDignity(tc.planet, tc.lon)
		if gotEng != tc.wantDignity || gotThai != tc.wantThai {
			t.Errorf("CalculateDignity(%s, %.2f) = (%s, %s); want (%s, %s)",
				tc.planet, tc.lon, gotEng, gotThai, tc.wantDignity, tc.wantThai)
		}
	}
}

func TestCalculateSignDignityAndShortCode(t *testing.T) {
	signCases := []struct {
		planet      string
		signIdx     int
		wantDignity string
	}{
		{"Sun", 0, "Exalted"},
		{"Sun", 6, "Debilitated"},
		{"Sun", 4, "Own Sign"},
		{"Sun", 8, "Friend"},
		{"Sun", 9, "Enemy"},
		{"Sun", 2, "Neutral"},
		{"Moon", 1, "Exalted"},
		{"Moon", 7, "Debilitated"},
		{"Moon", 3, "Own Sign"},
		{"Mars", 9, "Exalted"},
		{"Mars", 3, "Debilitated"},
		{"Mars", 0, "Own Sign"},
		{"Mars", 7, "Own Sign"},
		{"Mercury", 5, "Exalted"},
		{"Mercury", 11, "Debilitated"},
		{"Mercury", 2, "Own Sign"},
		{"Jupiter", 3, "Exalted"},
		{"Jupiter", 9, "Debilitated"},
		{"Jupiter", 8, "Own Sign"},
		{"Venus", 11, "Exalted"},
		{"Venus", 5, "Debilitated"},
		{"Venus", 1, "Own Sign"},
		{"Saturn", 6, "Exalted"},
		{"Saturn", 0, "Debilitated"},
		{"Saturn", 9, "Own Sign"},
		{"Rahu", 1, "Exalted"},
		{"Rahu", 2, "Moolatrikona"},
		{"Rahu", 10, "Own Sign"},
		{"Rahu", 7, "Debilitated"},
		{"Rahu", 0, "Neutral"},
		{"Ketu", 7, "Exalted"},
		{"Ketu", 8, "Moolatrikona"},
		{"Ketu", 1, "Debilitated"},
		{"Ketu", 0, "Neutral"},
		{"Uranus", 0, "Neutral"},
		{"Sun", -1, "Neutral"},
	}

	for _, tc := range signCases {
		gotEng, _ := CalculateSignDignity(tc.planet, tc.signIdx)
		if gotEng != tc.wantDignity {
			t.Errorf("CalculateSignDignity(%s, %d) = %q; want %q", tc.planet, tc.signIdx, gotEng, tc.wantDignity)
		}
	}

	codeCases := map[string]string{
		"Exalted":      "Exa",
		"Moolatrikona": "Mool",
		"Own Sign":     "Own",
		"Friend":       "Fri",
		"Enemy":        "Ene",
		"Debilitated":  "Deb",
		"Neutral":      "Neu",
		"Other":        "Neu",
	}
	for in, want := range codeCases {
		if got := ShortDignityCode(in); got != want {
			t.Errorf("ShortDignityCode(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestIsPlanetCombust(t *testing.T) {
	sunLon := 247.81
	cases := []struct {
		name      string
		planet    string
		planetLon float64
		sunLon    float64
		isRetro   bool
		want      bool
	}{
		{"Sun never combust", "Sun", sunLon, sunLon, false, false},
		{"Rahu never combust", "Rahu (Mean Node)", sunLon, sunLon, true, false},
		{"Ketu never combust", "Ketu (Calculated)", sunLon, sunLon, false, false},
		{"Moon inside 12 deg", "Moon", sunLon + 11.0, sunLon, false, true},
		{"Moon outside 12 deg", "Moon", sunLon + 13.0, sunLon, false, false},
		{"Mars inside 17 deg", "Mars", sunLon + 16.5, sunLon, false, true},
		{"Mars outside 17 deg", "Mars", sunLon + 18.0, sunLon, false, false},
		{"Mercury direct inside 14 deg", "Mercury", 252.0, sunLon, false, true},
		{"Mercury direct outside 14 deg", "Mercury", 268.0, sunLon, false, false},
		{"Mercury retro inside 12 deg", "Mercury", sunLon + 11.5, sunLon, true, true},
		{"Mercury retro outside 12 deg", "Mercury", sunLon + 13.0, sunLon, true, false},
		{"Jupiter inside 11 deg", "Jupiter", sunLon + 10.5, sunLon, false, true},
		{"Jupiter outside 11 deg", "Jupiter", sunLon + 12.0, sunLon, false, false},
		{"Venus direct inside 10 deg", "Venus", sunLon + 9.5, sunLon, false, true},
		{"Venus direct outside 10 deg", "Venus", sunLon + 10.5, sunLon, false, false},
		{"Venus retro inside 8 deg", "Venus", sunLon + 7.5, sunLon, true, true},
		{"Venus retro outside 8 deg", "Venus", sunLon + 9.0, sunLon, true, false},
		{"Saturn inside 15 deg", "Saturn", sunLon + 14.5, sunLon, false, true},
		{"Saturn outside 15 deg", "Saturn", sunLon + 16.0, sunLon, false, false},
		{"Uranus default 10 deg", "Uranus", 2.0, 358.0, false, true},
		{"Uranus outside 10 deg", "Uranus", 15.0, 358.0, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPlanetCombust(tc.planet, tc.planetLon, tc.sunLon, tc.isRetro); got != tc.want {
				t.Errorf("IsPlanetCombust(%s, %.2f, %.2f, %v) = %v; want %v",
					tc.planet, tc.planetLon, tc.sunLon, tc.isRetro, got, tc.want)
			}
		})
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

	t.Run("TargetAtBirthClampsFirstAntarAndPratyStart", func(t *testing.T) {
		loc := time.FixedZone("ICT", 7*3600)
		birthTarget := time.Date(1988, 12, 23, 10, 21, 0, 0, loc)
		atBirth := CalculateVimshottariDasha(moonLon, 1988, 12, 23, 10.35, 7.0, birthTarget)
		if atBirth.CurrentDasha.MahadashaStart != "1988-12-23" {
			t.Errorf("expected MahadashaStart 1988-12-23, got %s", atBirth.CurrentDasha.MahadashaStart)
		}
		if atBirth.CurrentDasha.AntardashaStart != "1988-12-23" {
			t.Errorf("expected AntardashaStart clamped to 1988-12-23, got %s", atBirth.CurrentDasha.AntardashaStart)
		}
		if atBirth.CurrentDasha.PratyantardashaStart != "1988-12-23" {
			t.Errorf("expected PratyantardashaStart clamped to 1988-12-23, got %s", atBirth.CurrentDasha.PratyantardashaStart)
		}
	})

	t.Run("ZeroTargetTimeDefaultsToNow", func(t *testing.T) {
		nowDasha := CalculateVimshottariDasha(moonLon, 1988, 12, 23, 10.35, 7.0, time.Time{})
		if nowDasha.TargetDate == "" || nowDasha.CurrentDasha.MahadashaLord == "" {
			t.Errorf("expected non-empty TargetDate and CurrentDasha when targetTime is zero")
		}
	})
}
