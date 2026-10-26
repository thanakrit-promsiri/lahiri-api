package calculator

import "fmt"

// NakshatraInfo เก็บข้อมูลรายละเอียดของแต่ละนักษัตร
type NakshatraInfo struct {
	Index    int     // ลำดับที่ 1-27
	Name     string  // ชื่อนักษัตร
	Lord     string  // ดาวเจ้าเรือน (Vimshottari Dasha Lord)
	StartDeg float64 // องศาเริ่มต้นในจักรวาล (0 - 360)
	EndDeg   float64 // องศา สิ้นสุด
}

// GetNakshatra รายชื่อนักษัตรทั้ง 27 แห่งโหราศาสตร์เวท
func GetNakshatra(lon float64) (NakshatraInfo, int, float64) {
	// แต่ละนักษัตรกว้าง 13° 20' หรือ 13.333333 องศา
	span := 360.0 / 27.0 // 13.333333333333334

	nakIndex := int(lon / span)
	if nakIndex >= 27 {
		nakIndex = 26
	}

	// คำนวณเศษองศาภายในนักษัตรนั้นๆ เพื่อหา Pada (บาท 1-4)
	offsetInNak := lon - (float64(nakIndex) * span)
	padaSpan := span / 4.0 // บาทละ 3° 20'
	pada := int(offsetInNak/padaSpan) + 1
	if pada > 4 {
		pada = 4
	}

	names := []string{
		"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
		"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
		"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
		"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
		"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
	}

	// ลำดับดาวเจ้าเรือน วนซ้ำ 3 รอบตามระบบ Vimshottari Dasha (Ketu, Venus, Sun, Moon, Mars, Rahu, Jupiter, Saturn, Mercury)
	lords := []string{
		"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	}
	lord := lords[nakIndex%9]

	info := NakshatraInfo{
		Index:    nakIndex + 1,
		Name:     names[nakIndex],
		Lord:     lord,
		StartDeg: float64(nakIndex) * span,
		EndDeg:   float64(nakIndex+1) * span,
	}

	return info, pada, offsetInNak
}

// FormatNakshatraString จัดรูปแบบข้อความแสดงนักษัตรและบาท
func FormatNakshatraString(lon float64) string {
	nak, pada, _ := GetNakshatra(lon)
	return fmt.Sprintf("%s (Pad %d) [Lord: %s]", nak.Name, pada, nak.Lord)
}
