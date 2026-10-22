package main

// Config เก็บค่าพารามิเตอร์ตั้งต้นทั้งหมดของดวงชะตา
type Config struct {
	Year        int
	Month       int
	Day         int
	HourThai    float64
	TimeZone    float64
	Latitude    float64
	Longitude   float64
	HouseSystem int
}

// PlanetConfig เก็บข้อมูลรหัสและชื่อดาวเคราะห์
type PlanetConfig struct {
	ID   int
	Name string
}

// GetDefaultConfig คืนค่าคอนฟิกเริ่มต้นสำหรับคำนวณ
func GetDefaultConfig() Config {
	return Config{
		Year:        1988,
		Month:       12,
		Day:         23,
		HourThai:    10.0 + (21.0 / 60.0), // 10:21 น.
		TimeZone:    7.0,                  // เวลาประเทศไทย GMT+7
		Latitude:    17.597591140953423,   // พิกัดละติจูดใหม่
		Longitude:   101.7205414851994,    // พิกัดลองจิจูดใหม่
		HouseSystem: int('W'),             // Whole Sign
	}

}

// GetZodiacSigns คืนค่ารายชื่อราศีทั้ง 12
func GetZodiacSigns() []string {
	return []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
}
