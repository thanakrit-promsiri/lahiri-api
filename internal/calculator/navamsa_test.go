package calculator

import "testing"

func TestCalculateNavamsa(t *testing.T) {
	tests := []struct {
		name     string
		lon      float64
		wantSign string
		wantThai string
		wantIdx  int
	}{
		// Fire signs (Aries=0, Leo=4, Sagittarius=8) -> start at Aries (0)
		{name: "Aries 0 deg (Fire start)", lon: 0.0, wantSign: "Aries", wantThai: "เมษ", wantIdx: 1},
		{name: "Leo 15 deg (Fire 5th navamsa)", lon: 135.0, wantSign: "Leo", wantThai: "สิงห์", wantIdx: 5},
		{name: "Sagittarius 29.99 deg (Fire 9th navamsa)", lon: 269.99, wantSign: "Sagittarius", wantThai: "ธนู", wantIdx: 9},

		// Earth signs (Taurus=1, Virgo=5, Capricorn=9) -> start at Capricorn (9)
		{name: "Taurus 0 deg (Earth start)", lon: 30.0, wantSign: "Capricorn", wantThai: "มังกร", wantIdx: 1},
		{name: "Virgo 10 deg (Earth 4th navamsa)", lon: 160.0, wantSign: "Aries", wantThai: "เมษ", wantIdx: 4},

		// Air signs (Gemini=2, Libra=6, Aquarius=10) -> start at Libra (6)
		{name: "Gemini 0 deg (Air start)", lon: 60.0, wantSign: "Libra", wantThai: "ตุลย์", wantIdx: 1},
		{name: "Aquarius 20 deg (Air 7th navamsa)", lon: 320.0, wantSign: "Aries", wantThai: "เมษ", wantIdx: 7},

		// Water signs (Cancer=3, Scorpio=7, Pisces=11) -> start at Cancer (3)
		{name: "Cancer 0 deg (Water start)", lon: 90.0, wantSign: "Cancer", wantThai: "กรกฎ", wantIdx: 1},
		{name: "Pisces 29.9999 deg (Water 9th navamsa)", lon: 359.9999, wantSign: "Pisces", wantThai: "มีน", wantIdx: 9},

		// Angle wrap-around
		{name: "360 deg wraps to Aries 0", lon: 360.0, wantSign: "Aries", wantThai: "เมษ", wantIdx: 1},
		{name: "Negative -5 deg wraps to Pisces 25 deg (8th navamsa)", lon: -5.0, wantSign: "Aquarius", wantThai: "กุมภ์", wantIdx: 8},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateNavamsa(tc.lon)
			if got.Sign != tc.wantSign || got.SignThai != tc.wantThai || got.Index != tc.wantIdx {
				t.Errorf("CalculateNavamsa(%.4f) = %+v; want {Sign:%s SignThai:%s Index:%d}",
					tc.lon, got, tc.wantSign, tc.wantThai, tc.wantIdx)
			}
		})
	}
}
