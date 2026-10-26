package calculator

// GetDefaultConfig คืนค่าคอนฟิกเริ่มต้นสำหรับคำนวณ
func GetDefaultConfig() Config {
	return Config{
		Year:        1988,
		Month:       12,
		Day:         23,
		HourThai:    10.0 + (21.0 / 60.0), // 10:21 น.
		TimeZone:    7.0,                  // เวลาประเทศไทย GMT+7
		Latitude:    17.597591140953423,   // พิกัดละติจูด
		Longitude:   101.7205414851994,    // พิกัดลองจิจูด
		HouseSystem: int('W'),             // Whole Sign
	}
}

// GetZodiacSigns คืนค่ารายชื่อราศีทั้ง 12 (ภาษาอังกฤษ)
func GetZodiacSigns() []string {
	return []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
}

// GetEnglishSigns คืนค่ารายชื่อราศีภาษาอังกฤษมาตรฐาน
func GetEnglishSigns() []string {
	return GetZodiacSigns()
}

// GetNakshatraNames คืนค่ารายชื่อนักษัตรทั้ง 27
func GetNakshatraNames() []string {
	return []string{
		"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
		"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
		"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
		"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
		"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
	}
}

// GetNakshatraLords คืนค่าลำดับดาวเจ้าเรือน (Vimshottari Dasha Lords)
func GetNakshatraLords() []string {
	return []string{
		"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	}
}
