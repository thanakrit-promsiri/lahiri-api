package calculator

import "fmt"

// ==========================================
// 2. Core Calculation & Math Functions
// ==========================================

// ParseLongitude แยกค่าลองจิจูดออกเป็น index ของราศีและเศษองศา
func ParseLongitude(lon float64) (int, float64) {
	signIdx := int(lon / 30)
	degree := lon - float64(signIdx)*30
	return signIdx, degree
}

// NormalizeAngle ปรับค่ามุมให้อยู่ในช่วง 0 ถึง 360 องศาเสมอ
func NormalizeAngle(val float64) float64 {
	for val < 0 {
		val += 360
	}
	for val >= 360 {
		val -= 360
	}
	return val
}

// ==========================================
// 3. Formatting Functions
// ==========================================

// FormatDMS แปลงค่าองศบทศนิยมให้เป็น องศา°, ลิปดา', ฟิลิปดา" (ละเอียดถึงฟิลิปดา)
func FormatDMS(decimalDegree float64) string {
	deg := int(decimalDegree)
	sub := (decimalDegree - float64(deg)) * 60
	min := int(sub)
	sec := (sub - float64(min)) * 60
	return fmt.Sprintf("%2d° %2d' %05.2f\"", deg, min, sec)
}

// FormatSignDMS แปลงองศาเต็มระบบราศี ให้เป็น ชื่อราศี + องศา° ลิปดา' ฟิลิปดา"
func FormatSignDMS(lon float64, signs []string) string {
	signIdx, degree := ParseLongitude(lon)
	dmsStr := FormatDMS(degree)
	return fmt.Sprintf("%s (%s)", signs[signIdx], dmsStr)
}

// ==========================================
// 4. Domain & Helper Functions
// ==========================================

// ParsePositionDetail แปลงค่าลองจิจูดเป็นรายละเอียดปลีกย่อย (ราศี, องศา, ลิปดา, ฟิลิปดา)
func ParsePositionDetail(lon float64) PositionDetail {
	englishSigns := GetEnglishSigns()
	normalizedLon := NormalizeAngle(lon)

	signIdx, degInSign := ParseLongitude(normalizedLon)
	deg := int(degInSign)
	minFloat := (degInSign - float64(deg)) * 60.0
	min := int(minFloat)
	sec := (minFloat - float64(min)) * 60.0

	signName := ""
	if signIdx >= 0 && signIdx < len(englishSigns) {
		signName = englishSigns[signIdx]
	}

	return PositionDetail{
		Sign:   signName,
		Degree: deg,
		Minute: min,
		Second: sec,
	}
}
