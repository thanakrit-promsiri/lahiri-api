package calculator

// ตารางความสัมพันธ์ตามธรรมชาติของดาว (Naisargika Maitri)
// มิตร (Friend), ศัตรู (Enemy), สมพล/กลาง (Neutral)
var naturalFriends = map[string]map[string]bool{
	"Sun":     {"Moon": true, "Mars": true, "Jupiter": true},
	"Moon":    {"Sun": true, "Mercury": true},
	"Mars":    {"Sun": true, "Moon": true, "Jupiter": true},
	"Mercury": {"Sun": true, "Venus": true},
	"Jupiter": {"Sun": true, "Moon": true, "Mars": true},
	"Venus":   {"Mercury": true, "Saturn": true},
	"Saturn":  {"Mercury": true, "Venus": true},
}

var naturalEnemies = map[string]map[string]bool{
	"Sun":     {"Venus": true, "Saturn": true},
	"Moon":    {},
	"Mars":    {"Mercury": true},
	"Mercury": {"Moon": true},
	"Jupiter": {"Mercury": true, "Venus": true},
	"Venus":   {"Sun": true, "Moon": true},
	"Saturn":  {"Sun": true, "Moon": true, "Mars": true},
}

// CalculateDignity คำนวณมาตรฐานดาวตามหลักโหราศาสตร์พระเวท (อุจจ์, มูลไตรโกณ, เกษตร, มิตร, สมพล, ศัตรู, นิจ)
// ราศี: 0=Aries(เมษ), 1=Taurus(พฤษภ), 2=Gemini(เมถุน), 3=Cancer(กรกฎ), 4=Leo(สิงห์), 5=Virgo(กันย์),
// 6=Libra(ตุลย์), 7=Scorpio(พิจิก), 8=Sagittarius(ธนู), 9=Capricorn(มังกร), 10=Aquarius(กุมภ์), 11=Pisces(มีน)
func CalculateDignity(planetName string, lon float64) (dignityEng string, dignityThai string) {
	_, _, canonical := TranslateBodyInfo(planetName)
	signIdx, degInSign := ParseLongitude(lon)

	// สำหรับดาวรอบนอก (Uranus, Neptune, Pluto)
	switch canonical {
	case "Uranus", "Neptune", "Pluto":
		return "Neutral", "ปกติ"
	case "Rahu":
		switch signIdx {
		case 1: // Taurus
			return "Exalted", "อุจจ์"
		case 2: // Gemini
			return "Moolatrikona", "มูลไตรโกณ"
		case 10: // Aquarius
			return "Own Sign", "เกษตร"
		case 7: // Scorpio
			return "Debilitated", "นิจ"
		default:
			return "Neutral", "ปกติ"
		}
	case "Ketu":
		switch signIdx {
		case 7: // Scorpio
			return "Exalted", "อุจจ์"
		case 8: // Sagittarius
			return "Moolatrikona", "มูลไตรโกณ"
		case 1: // Taurus
			return "Debilitated", "นิจ"
		default:
			return "Neutral", "ปกติ"
		}
	}

	// 1. ตรวจสอบ มูลไตรโกณ (Moolatrikona) และ อุจจ์ (Exalted) / เกษตร (Own Sign) / นิจ (Debilitated) ของดาวทั้ง 7
	switch canonical {
	case "Sun":
		if signIdx == 0 {
			return "Exalted", "อุจจ์"
		}
		if signIdx == 6 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 4 {
			if degInSign < 20.0 {
				return "Moolatrikona", "มูลไตรโกณ"
			}
			return "Own Sign", "เกษตร"
		}
	case "Moon":
		if signIdx == 1 {
			if degInSign <= 3.0 {
				return "Exalted", "อุจจ์"
			}
			return "Moolatrikona", "มูลไตรโกณ"
		}
		if signIdx == 7 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 3 {
			return "Own Sign", "เกษตร"
		}
	case "Mars":
		if signIdx == 9 {
			return "Exalted", "อุจจ์"
		}
		if signIdx == 3 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 0 {
			if degInSign < 12.0 {
				return "Moolatrikona", "มูลไตรโกณ"
			}
			return "Own Sign", "เกษตร"
		}
		if signIdx == 7 {
			return "Own Sign", "เกษตร"
		}
	case "Mercury":
		if signIdx == 5 {
			if degInSign < 15.0 {
				return "Exalted", "อุจจ์"
			}
			if degInSign < 20.0 {
				return "Moolatrikona", "มูลไตรโกณ"
			}
			return "Own Sign", "เกษตร"
		}
		if signIdx == 11 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 2 {
			return "Own Sign", "เกษตร"
		}
	case "Jupiter":
		if signIdx == 3 {
			return "Exalted", "อุจจ์"
		}
		if signIdx == 9 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 8 {
			if degInSign < 10.0 {
				return "Moolatrikona", "มูลไตรโกณ"
			}
			return "Own Sign", "เกษตร"
		}
		if signIdx == 11 {
			return "Own Sign", "เกษตร"
		}
	case "Venus":
		if signIdx == 11 {
			return "Exalted", "อุจจ์"
		}
		if signIdx == 5 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 6 {
			if degInSign < 15.0 {
				return "Moolatrikona", "มูลไตรโกณ"
			}
			return "Own Sign", "เกษตร"
		}
		if signIdx == 1 {
			return "Own Sign", "เกษตร"
		}
	case "Saturn":
		if signIdx == 6 {
			return "Exalted", "อุจจ์"
		}
		if signIdx == 0 {
			return "Debilitated", "นิจ"
		}
		if signIdx == 10 {
			if degInSign < 20.0 {
				return "Moolatrikona", "มูลไตรโกณ"
			}
			return "Own Sign", "เกษตร"
		}
		if signIdx == 9 {
			return "Own Sign", "เกษตร"
		}
	}

	// 2. ตรวจสอบความสัมพันธ์กับดาวเจ้าเรือนเกษตร (Friend / Neutral / Enemy)
	dispositor := GetSignLord(signIdx)
	if dispositor == "" {
		return "Neutral", "ปกติ"
	}
	if friends, ok := naturalFriends[canonical]; ok && friends[dispositor] {
		return "Friend", "มิตร"
	}
	if enemies, ok := naturalEnemies[canonical]; ok && enemies[dispositor] {
		return "Enemy", "ศัตรู"
	}

	return "Neutral", "ปกติ"
}

// IsPlanetCombust ตรวจสอบภาวะดาวดับ/อัสดงคต (Combustion / Asta) เมื่อโคจรเข้าใกล้ดวงอาทิตย์ตามเกณฑ์องศาพระเวท
func IsPlanetCombust(planetName string, planetLon, sunLon float64, isRetro bool) bool {
	_, _, canonical := TranslateBodyInfo(planetName)

	// อาทิตย์ ราหู เกตุ ไม่เกิดภาวะดาวดับ (Combust)
	if canonical == "Sun" || canonical == "Rahu" || canonical == "Ketu" {
		return false
	}

	dist := AngularDistance(planetLon, sunLon)
	var orb float64

	switch canonical {
	case "Moon":
		orb = 12.0
	case "Mars":
		orb = 17.0
	case "Mercury":
		if isRetro {
			orb = 12.0
		} else {
			orb = 14.0
		}
	case "Jupiter":
		orb = 11.0
	case "Venus":
		if isRetro {
			orb = 8.0
		} else {
			orb = 10.0
		}
	case "Saturn":
		orb = 15.0
	default:
		// ดาวรอบนอกใช้เกณฑ์ 10 องศา
		orb = 10.0
	}

	return dist <= orb
}
