package calculator

import "fmt"

// GetNakshatra คำนวณตำแหน่งนักษัตรทั้ง 27 แห่งโหราศาสตร์เวทจากค่าลองจิจูด
func GetNakshatra(lon float64) (NakshatraInfo, int, float64) {
	// แต่ละนักษัตรกว้าง 13° 20' หรือ 13.333333 องศา
	span := 360.0 / 27.0

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

	names := GetNakshatraNames()
	lords := GetNakshatraLords()
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