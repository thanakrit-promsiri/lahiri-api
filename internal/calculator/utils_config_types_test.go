package calculator

import (
	"math"
	"strings"
	"testing"
)

func TestUtilsMathAndFormatting(t *testing.T) {
	t.Run("NormalizeAngle", func(t *testing.T) {
		cases := []struct {
			in, want float64
		}{
			{0, 0},
			{360, 0},
			{725.5, 5.5},
			{-30.0, 330.0},
			{-390.0, 330.0},
		}
		for _, tc := range cases {
			if got := NormalizeAngle(tc.in); math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("NormalizeAngle(%v) = %v; want %v", tc.in, got, tc.want)
			}
		}
	})

	t.Run("ParseLongitude", func(t *testing.T) {
		sign, deg := ParseLongitude(247.8175)
		if sign != 8 || math.Abs(deg-7.8175) > 1e-6 {
			t.Errorf("ParseLongitude(247.8175) = (%d, %f); want (8, ~7.8175)", sign, deg)
		}
		signNeg, degNeg := ParseLongitude(-15.0)
		if signNeg != 11 || math.Abs(degNeg-15.0) > 1e-6 {
			t.Errorf("ParseLongitude(-15) = (%d, %f); want (11, 15)", signNeg, degNeg)
		}
	})

	t.Run("ClampDayToMonth", func(t *testing.T) {
		cases := []struct {
			year, month, day, want int
		}{
			{2026, 1, 31, 31},
			{2026, 2, 31, 28},
			{2028, 2, 31, 29},
			{2026, 4, 31, 30},
			{2026, 12, 15, 15},
		}
		for _, tc := range cases {
			if got := ClampDayToMonth(tc.year, tc.month, tc.day); got != tc.want {
				t.Errorf("ClampDayToMonth(%d, %d, %d) = %d; want %d", tc.year, tc.month, tc.day, got, tc.want)
			}
		}
	})

	t.Run("AngularDistance", func(t *testing.T) {
		cases := []struct {
			a, b, want float64
		}{
			{10, 50, 40},
			{10, 350, 20},
			{358, 2, 4},
			{0, 180, 180},
		}
		for _, tc := range cases {
			if got := AngularDistance(tc.a, tc.b); math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("AngularDistance(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		}
	})

	t.Run("FormatDegreeShortAndParsePositionDetail", func(t *testing.T) {
		pos := ParsePositionDetail(247.8175) // Sagittarius 7° 49' 03.00"
		if pos.Sign != "Sagittarius" || pos.SignThai != "ธนู" {
			t.Errorf("ParsePositionDetail sign = (%s, %s); want (Sagittarius, ธนู)", pos.Sign, pos.SignThai)
		}
		if pos.Degree != 7 || pos.Minute != 49 {
			t.Errorf("ParsePositionDetail DMS = %d° %d' %.2f\"; want 7° 49'", pos.Degree, pos.Minute, pos.Second)
		}
		short := FormatDegreeShort(pos)
		if !strings.HasPrefix(short, "07° 49'") {
			t.Errorf("FormatDegreeShort = %q; expected prefix 07° 49'", short)
		}
	})

	t.Run("FormatBirthDateTime", func(t *testing.T) {
		cases := []struct {
			y, m, d  int
			hour     float64
			wantDate string
			wantTime string
		}{
			{1988, 12, 23, 10.35, "1988-12-23", "10:21"},
			{2026, 1, 5, 10.505, "2026-01-05", "10:30:18"},
			{2026, 6, 15, 23.99999, "2026-06-15", "00:00"},
		}
		for _, tc := range cases {
			gotDate, gotTime := FormatBirthDateTime(tc.y, tc.m, tc.d, tc.hour)
			if gotDate != tc.wantDate || gotTime != tc.wantTime {
				t.Errorf("FormatBirthDateTime(%d,%d,%d,%v) = (%q, %q); want (%q, %q)",
					tc.y, tc.m, tc.d, tc.hour, gotDate, gotTime, tc.wantDate, tc.wantTime)
			}
		}
	})
}

func TestConfigLookups(t *testing.T) {
	cfg := GetDefaultConfig()
	if cfg.Year != 1988 || cfg.TimeZone != 7.0 || cfg.HouseSystem != int('W') {
		t.Errorf("Unexpected GetDefaultConfig: %+v", cfg)
	}
	if len(GetEnglishSigns()) != 12 || len(GetThaiZodiacSigns()) != 12 {
		t.Errorf("Expected 12 English and 12 Thai signs")
	}

	if got := GetSignLord(0); got != "Mars" {
		t.Errorf("GetSignLord(0) = %q; want Mars", got)
	}
	if got := GetSignLord(-1); got != "" {
		t.Errorf("GetSignLord(-1) = %q; want empty", got)
	}

	if got := GetHouseNameThai(1); got != "ตนุ" {
		t.Errorf("GetHouseNameThai(1) = %q; want ตนุ", got)
	}
	if got := GetHouseNameThai(12); got != "วินาศ" {
		t.Errorf("GetHouseNameThai(12) = %q; want วินาศ", got)
	}
	if got := GetHouseNameThai(0); got != "" {
		t.Errorf("GetHouseNameThai(0) = %q; want empty", got)
	}

	if n, s, c := TranslateBodyInfo("UnknownBody"); n != "UnknownBody" || s != "UnknownBody" || c != "UnknownBody" {
		t.Errorf("TranslateBodyInfo(UnknownBody) = (%q, %q, %q)", n, s, c)
	}
	if got := TranslateLordThai("Unknown"); got != "Unknown" {
		t.Errorf("TranslateLordThai(Unknown) = %q; want Unknown", got)
	}

	if got := TranslateNakshatraThai("Ashwini"); got != "อัศวินี" {
		t.Errorf("TranslateNakshatraThai(Ashwini) = %q; want อัศวินี", got)
	}
	if got := TranslateNakshatraThai("NoSuchNak"); got != "NoSuchNak" {
		t.Errorf("TranslateNakshatraThai(NoSuchNak) = %q; want NoSuchNak", got)
	}

	if got := ShortSignCode("Sagittarius"); got != "Sg" {
		t.Errorf("ShortSignCode(Sagittarius) = %q; want Sg", got)
	}
	if got := ShortSignCode("Unknown"); got != "Unknown" {
		t.Errorf("ShortSignCode(Unknown) = %q; want Unknown", got)
	}

	if got := ShortPlanetCode("Jupiter"); got != "Ju" {
		t.Errorf("ShortPlanetCode(Jupiter) = %q; want Ju", got)
	}
	if got := ShortPlanetCode("Unknown"); got != "Unknown" {
		t.Errorf("ShortPlanetCode(Unknown) = %q; want Unknown", got)
	}

	nakShortCases := map[string]string{
		"Purva Phalguni":    "P.Phalguni",
		"Uttara Phalguni":   "U.Phalguni",
		"Purva Ashadha":     "P.Ashadha",
		"Uttara Ashadha":    "U.Ashadha",
		"Purva Bhadrapada":  "P.Bhadra",
		"Uttara Bhadrapada": "U.Bhadra",
		"Ashwini":           "Ashwini",
	}
	for in, want := range nakShortCases {
		if got := ShortNakshatraName(in); got != want {
			t.Errorf("ShortNakshatraName(%q) = %q; want %q", in, got, want)
		}
	}

	if idx := GetSignIndexByName("Aries"); idx != 0 {
		t.Errorf("GetSignIndexByName(Aries) = %d; want 0", idx)
	}
	if idx := GetSignIndexByName("Pisces"); idx != 11 {
		t.Errorf("GetSignIndexByName(Pisces) = %d; want 11", idx)
	}
	if idx := GetSignIndexByName("Invalid"); idx != -1 {
		t.Errorf("GetSignIndexByName(Invalid) = %d; want -1", idx)
	}
}

func TestRequestValidation(t *testing.T) {
	validCalc := CalculateReq{
		Year: 1988, Month: 12, Day: 23, HourThai: 10.35,
		TimeZone: 7.0, Latitude: 13.7563, Longitude: 100.5018, HouseSystem: "W",
	}
	if err := validCalc.Validate(); err != nil {
		t.Fatalf("valid CalculateReq failed Validate(): %v", err)
	}

	calcInvalidCases := []struct {
		name string
		mod  func(*CalculateReq)
	}{
		{"Year=0", func(r *CalculateReq) { r.Year = 0 }},
		{"Month=0", func(r *CalculateReq) { r.Month = 0 }},
		{"Month=13", func(r *CalculateReq) { r.Month = 13 }},
		{"Day=0", func(r *CalculateReq) { r.Day = 0 }},
		{"Day=32", func(r *CalculateReq) { r.Day = 32 }},
		{"HourThai=-0.1", func(r *CalculateReq) { r.HourThai = -0.1 }},
		{"HourThai=24", func(r *CalculateReq) { r.HourThai = 24.0 }},
		{"TimeZone=-14.1", func(r *CalculateReq) { r.TimeZone = -14.1 }},
		{"TimeZone=14.1", func(r *CalculateReq) { r.TimeZone = 14.1 }},
		{"Latitude=-90.1", func(r *CalculateReq) { r.Latitude = -90.1 }},
		{"Latitude=90.1", func(r *CalculateReq) { r.Latitude = 90.1 }},
		{"Longitude=-180.1", func(r *CalculateReq) { r.Longitude = -180.1 }},
		{"Longitude=180.1", func(r *CalculateReq) { r.Longitude = 180.1 }},
	}
	for _, tc := range calcInvalidCases {
		t.Run("CalculateReq_"+tc.name, func(t *testing.T) {
			r := validCalc
			tc.mod(&r)
			if err := r.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
		})
	}

	validTransit := TransitReq{
		Year: 1988, Month: 12, Day: 23, HourThai: 10.35,
		TimeZone: 7.0, Latitude: 13.7563, Longitude: 100.5018,
		HouseSystem: "W", TransitYear: 2026, TransitMonth: 12,
	}
	if err := validTransit.Validate(); err != nil {
		t.Fatalf("valid TransitReq failed Validate(): %v", err)
	}
	if err := validTransit.ValidateWithMonth(); err != nil {
		t.Fatalf("valid TransitReq failed ValidateWithMonth(): %v", err)
	}

	transitInvalidCases := []struct {
		name      string
		withMonth bool
		mod       func(*TransitReq)
	}{
		{"Year=0", false, func(r *TransitReq) { r.Year = 0 }},
		{"Month=13", false, func(r *TransitReq) { r.Month = 13 }},
		{"Day=32", false, func(r *TransitReq) { r.Day = 32 }},
		{"HourThai=24", false, func(r *TransitReq) { r.HourThai = 24 }},
		{"TimeZone=15", false, func(r *TransitReq) { r.TimeZone = 15 }},
		{"Latitude=91", false, func(r *TransitReq) { r.Latitude = 91 }},
		{"Longitude=181", false, func(r *TransitReq) { r.Longitude = 181 }},
		{"TransitYear=0", false, func(r *TransitReq) { r.TransitYear = 0 }},
		{"WithMonth_BaseErr", true, func(r *TransitReq) { r.Year = 0 }},
		{"WithMonth_0", true, func(r *TransitReq) { r.TransitMonth = 0 }},
		{"WithMonth_13", true, func(r *TransitReq) { r.TransitMonth = 13 }},
	}
	for _, tc := range transitInvalidCases {
		t.Run("TransitReq_"+tc.name, func(t *testing.T) {
			r := validTransit
			tc.mod(&r)
			var err error
			if tc.withMonth {
				err = r.ValidateWithMonth()
			} else {
				err = r.Validate()
			}
			if err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
		})
	}
}
