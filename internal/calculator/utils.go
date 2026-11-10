package calculator

import (
	"fmt"
	"math"
	"time"
)

// ==========================================
// 2. Core Calculation & Math Functions
// ==========================================

// ParseLongitude แยกค่าลองจิจูดออกเป็น index ของราศีและเศษองศา
func ParseLongitude(lon float64) (int, float64) {
	normalized := NormalizeAngle(lon)
	signIdx := int(normalized / 30)
	if signIdx >= 12 {
		signIdx = 11
	}
	degree := normalized - float64(signIdx)*30
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

// ClampDayToMonth ป้องกันกรณีวันที่เกิด (เช่น 29, 30, 31) เกินจำนวนวันสูงสุดของเดือนเป้าหมาย
func ClampDayToMonth(year, month, day int) int {
	lastDay := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		return lastDay
	}
	return day
}

// AngularDistance คำนวณระยะห่างเชิงมุมที่สั้นที่สุดระหว่างลองจิจูด 2 ค่า (0 - 180 องศา)
func AngularDistance(lon1, lon2 float64) float64 {
	diff := math.Abs(NormalizeAngle(lon1) - NormalizeAngle(lon2))
	if diff > 180.0 {
		diff = 360.0 - diff
	}
	return diff
}

// ==========================================
// 3. Formatting Functions
// ==========================================

// FormatDMS แปลงค่าองศาทศนิยมให้เป็น องศา°, ลิปดา', ฟิลิปดา" (ละเอียดถึงฟิลิปดา)
func FormatDMS(decimalDegree float64) string {
	deg := int(decimalDegree)
	sub := (decimalDegree - float64(deg)) * 60
	min := int(sub)
	sec := (sub - float64(min)) * 60
	return fmt.Sprintf("%2d° %2d' %05.2f\"", deg, min, sec)
}

// FormatDegreeShort แปลง PositionDetail ให้เป็นรูปแบบกระชับละเอียดถึงฟิลิปดา เช่น "07° 49' 03.39\""
func FormatDegreeShort(pos PositionDetail) string {
	return fmt.Sprintf("%02d° %02d' %05.2f\"", pos.Degree, pos.Minute, pos.Second)
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

// ParsePositionDetail แปลงค่าลองจิจูดเป็นรายละเอียดปลีกย่อย (ราศีอังกฤษ, ราศีไทย, องศา, ลิปดา, ฟิลิปดา)
func ParsePositionDetail(lon float64) PositionDetail {
	englishSigns := GetEnglishSigns()
	thaiSigns := GetThaiZodiacSigns()
	normalizedLon := NormalizeAngle(lon)

	signIdx, degInSign := ParseLongitude(normalizedLon)
	deg := int(degInSign)
	minFloat := (degInSign - float64(deg)) * 60.0
	min := int(minFloat)
	sec := (minFloat - float64(min)) * 60.0

	signName := ""
	signThai := ""
	if signIdx >= 0 && signIdx < len(englishSigns) {
		signName = englishSigns[signIdx]
		signThai = thaiSigns[signIdx]
	}

	return PositionDetail{
		Sign:     signName,
		SignThai: signThai,
		Degree:   deg,
		Minute:   min,
		Second:   sec,
	}
}
