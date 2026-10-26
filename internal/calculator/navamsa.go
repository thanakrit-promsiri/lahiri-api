package calculator

// CalculateNavamsa คำนวณตำแหน่งนวางศ์ (D9) ตามหลักโหราศาสตร์เวท
func CalculateNavamsa(lon float64) NavamsaDetail {
	normalizedLon := NormalizeAngle(lon)
	signIdx, degInSign := ParseLongitude(normalizedLon)

	// แต่ละนวางศ์กว้าง 3° 20' (3.333333 องศา)
	navamsaSpan := 30.0 / 9.0
	navamsaIdx := int(degInSign / navamsaSpan)
	if navamsaIdx > 8 {
		navamsaIdx = 8
	}

	// หาจุดเริ่มต้นราศีนวางศ์ตามธาตุของราศีเดิม (Moveable, Fixed, Dual)
	// ไฟ (0,4,8) เริ่มที่ Aries (0)
	// ดิน (1,5,9) เริ่มที่ Capricorn (9)
	// ลม (2,6,10) เริ่มที่ Libra (6)
	// น้ำ (3,7,11) เริ่มที่ Cancer (3)
	element := signIdx % 4
	var startSignIdx int
	switch element {
	case 0:
		startSignIdx = 0
	case 1:
		startSignIdx = 9
	case 2:
		startSignIdx = 6
	case 3:
		startSignIdx = 3
	}

	navamsaSignIdx := (startSignIdx + navamsaIdx) % 12
	englishSigns := GetEnglishSigns()
	zodiacSigns := GetZodiacSigns()

	signName := ""
	signThai := ""
	if navamsaSignIdx >= 0 && navamsaSignIdx < len(englishSigns) {
		signName = englishSigns[navamsaSignIdx]
		signThai = zodiacSigns[navamsaSignIdx]
	}

	return NavamsaDetail{
		Sign:     signName,
		SignThai: signThai,
		Index:    navamsaIdx + 1,
	}
}