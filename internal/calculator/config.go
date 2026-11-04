package calculator

var (
	zodiacSigns = []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}

	shortZodiacSigns = map[string]string{
		"Aries":       "Ar",
		"Taurus":      "Ta",
		"Gemini":      "Ge",
		"Cancer":      "Cn",
		"Leo":         "Le",
		"Virgo":       "Vi",
		"Libra":       "Li",
		"Scorpio":     "Sc",
		"Sagittarius": "Sg",
		"Capricorn":   "Cp",
		"Aquarius":    "Aq",
		"Pisces":      "Pi",
	}

	shortPlanetCodes = map[string]string{
		"Sun":               "Su",
		"Moon":              "Mo",
		"Mars":              "Ma",
		"Mercury":           "Me",
		"Jupiter":           "Ju",
		"Venus":             "Ve",
		"Saturn":            "Sa",
		"Rahu (Mean Node)":  "Ra",
		"Rahu (True Node)":  "Ra",
		"Rahu":              "Ra",
		"Ketu (Calculated)": "Ke",
		"Ketu":              "Ke",
		"Uranus":            "Ur",
		"Neptune":           "Ne",
		"Pluto":             "Pl",
		"Ascendant (ASC)":   "ASC",
		"Midheaven (MC)":    "MC",
		"Descendant (DSC)":  "DSC",
		"Imum Coeli (IC)":   "IC",
	}

	thaiZodiacSigns = []string{
		"เมษ", "พฤษภ", "เมถุน", "กรกฎ", "สิงห์", "กันย์",
		"ตุลย์", "พิจิก", "ธนู", "มังกร", "กุมภ์", "มีน",
	}

	nakshatraNames = []string{
		"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
		"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
		"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
		"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
		"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
	}

	thaiNakshatraNames = []string{
		"อัศวินี", "ภรณี", "กฤตติกา", "โรหิณี", "มฤคศิระ", "อารทรา",
		"ปุนัพสุ", "ปุษยะ", "อาศเลษา", "มาฆะ", "ปุรพผลคุนี", "อุตตรผลคุนี",
		"หัสถะ", "จิตรา", "สวาติ", "วิสาขา", "อนุราธา", "เชษฐา",
		"มูละ", "ปุรพาษาฒ", "อุตตราษาฒ", "ศรวณะ", "ธนิษฐา", "ศตภิษัช",
		"ปุรพภัทรบท", "อุตตรภัทรบท", "เรวดี",
	}

	nakshatraLords = []string{
		"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	}

	signLord = []string{
		"Mars",    // 0: Aries (เมษ)
		"Venus",   // 1: Taurus (พฤษภ)
		"Mercury", // 2: Gemini (เมถุน)
		"Moon",    // 3: Cancer (กรกฎ)
		"Sun",     // 4: Leo (สิงห์)
		"Mercury", // 5: Virgo (กันย์)
		"Venus",   // 6: Libra (ตุลย์)
		"Mars",    // 7: Scorpio (พิจิก)
		"Jupiter", // 8: Sagittarius (ธนู)
		"Saturn",  // 9: Capricorn (มังกร)
		"Saturn",  // 10: Aquarius (กุมภ์)
		"Jupiter", // 11: Pisces (มีน)
	}

	houseNamesThai = []string{
		"ตนุ", "กดุมพะ", "สหัชชะ", "พันธุ", "ปุตตะ", "อริ",
		"ปัตนิ", "มรณะ", "ศุภะ", "กัมมะ", "ลาภะ", "วินาศ",
	}

	bodyThaiMap = map[string]struct {
		NameThai string
		Symbol   string
		Canonical string
	}{
		"Sun":               {NameThai: "อาทิตย์", Symbol: "๑", Canonical: "Sun"},
		"Moon":              {NameThai: "จันทร์", Symbol: "๒", Canonical: "Moon"},
		"Mars":              {NameThai: "อังคาร", Symbol: "๓", Canonical: "Mars"},
		"Mercury":           {NameThai: "พุธ", Symbol: "๔", Canonical: "Mercury"},
		"Jupiter":           {NameThai: "พฤหัสบดี", Symbol: "๕", Canonical: "Jupiter"},
		"Venus":             {NameThai: "ศุกร์", Symbol: "๖", Canonical: "Venus"},
		"Saturn":            {NameThai: "เสาร์", Symbol: "๗", Canonical: "Saturn"},
		"Rahu (Mean Node)":  {NameThai: "ราหู", Symbol: "๘", Canonical: "Rahu"},
		"Rahu (True Node)":  {NameThai: "ราหู", Symbol: "๘", Canonical: "Rahu"},
		"Rahu":              {NameThai: "ราหู", Symbol: "๘", Canonical: "Rahu"},
		"Ketu (Calculated)": {NameThai: "เกตุ", Symbol: "๙", Canonical: "Ketu"},
		"Ketu":              {NameThai: "เกตุ", Symbol: "๙", Canonical: "Ketu"},
		"Uranus":            {NameThai: "มฤตยู", Symbol: "๐", Canonical: "Uranus"},
		"Neptune":           {NameThai: "เนปจูน", Symbol: "Ne", Canonical: "Neptune"},
		"Pluto":             {NameThai: "พลูโต", Symbol: "Pl", Canonical: "Pluto"},
		"Ascendant (ASC)":   {NameThai: "ลัคนา", Symbol: "ลั", Canonical: "ASC"},
		"Midheaven (MC)":    {NameThai: "มัธยลัคนา", Symbol: "MC", Canonical: "MC"},
		"Descendant (DSC)":  {NameThai: "อัสตลัคนา", Symbol: "DSC", Canonical: "DSC"},
		"Imum Coeli (IC)":   {NameThai: "ปาตาลลัคนา", Symbol: "IC", Canonical: "IC"},
	}

	lordThaiMap = map[string]string{
		"Sun":     "อาทิตย์ (๑)",
		"Moon":    "จันทร์ (๒)",
		"Mars":    "อังคาร (๓)",
		"Mercury": "พุธ (๔)",
		"Jupiter": "พฤหัสบดี (๕)",
		"Venus":   "ศุกร์ (๖)",
		"Saturn":  "เสาร์ (๗)",
		"Rahu":    "ราหู (๘)",
		"Ketu":    "เกตุ (๙)",
	}
)

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
	return zodiacSigns
}

// GetEnglishSigns คืนค่ารายชื่อราศีภาษาอังกฤษมาตรฐาน
func GetEnglishSigns() []string {
	return zodiacSigns
}

// GetThaiZodiacSigns คืนค่ารายชื่อราศีทั้ง 12 (ภาษาไทย)
func GetThaiZodiacSigns() []string {
	return thaiZodiacSigns
}

// GetNakshatraNames คืนค่ารายชื่อนักษัตรทั้ง 27 (ภาษาอังกฤษ)
func GetNakshatraNames() []string {
	return nakshatraNames
}

// GetThaiNakshatraNames คืนค่ารายชื่อนักษัตรทั้ง 27 (ภาษาไทย)
func GetThaiNakshatraNames() []string {
	return thaiNakshatraNames
}

// GetNakshatraLords คืนค่าลำดับดาวเจ้าเรือน (Vimshottari Dasha Lords)
func GetNakshatraLords() []string {
	return nakshatraLords
}

// GetSignLord คืนค่าดาวเจ้าเรือนเกษตรประจำราศี (0=Aries ... 11=Pisces)
func GetSignLord(signIdx int) string {
	if signIdx >= 0 && signIdx < len(signLord) {
		return signLord[signIdx]
	}
	return ""
}

// GetHouseNameThai คืนค่าชื่อภพภาษาไทย (1=ตนุ ... 12=วินาศ)
func GetHouseNameThai(houseNum int) string {
	if houseNum >= 1 && houseNum <= 12 {
		return houseNamesThai[houseNum-1]
	}
	return ""
}

// TranslateBodyInfo คืนค่าชื่อภาษาไทย สัญลักษณ์ และชื่อมาตรฐานของดาว/มุมภพ
func TranslateBodyInfo(name string) (nameThai, symbol, canonical string) {
	if info, ok := bodyThaiMap[name]; ok {
		return info.NameThai, info.Symbol, info.Canonical
	}
	return name, name, name
}

// TranslateLordThai แปลชื่อดาวเจ้าฤกษ์/เจ้าเรือนเป็นภาษาไทยพร้อมเลขไทย
func TranslateLordThai(lord string) string {
	if th, ok := lordThaiMap[lord]; ok {
		return th
	}
	return lord
}

// TranslateNakshatraThai แปลชื่อนักษัตรฤกษ์เป็นภาษาไทย
func TranslateNakshatraThai(nakName string) string {
	for i, eng := range nakshatraNames {
		if eng == nakName && i < len(thaiNakshatraNames) {
			return thaiNakshatraNames[i]
		}
	}
	return nakName
}

// ShortSignCode คืนค่ารหัสย่อราศี 2 ตัวอักษร (เช่น Aries -> Ar, Sagittarius -> Sg)
func ShortSignCode(signName string) string {
	if code, ok := shortZodiacSigns[signName]; ok {
		return code
	}
	return signName
}

// ShortPlanetCode คืนค่ารหัสย่อดาว 2 ตัวอักษร (เช่น Sun -> Su, Jupiter -> Ju)
func ShortPlanetCode(name string) string {
	if code, ok := shortPlanetCodes[name]; ok {
		return code
	}
	return name
}

// ShortNakshatraName ย่อคำว่า Purva -> P. และ Uttara -> U. เพื่อประหยัด Token
func ShortNakshatraName(nakName string) string {
	switch nakName {
	case "Purva Phalguni":
		return "P.Phalguni"
	case "Uttara Phalguni":
		return "U.Phalguni"
	case "Purva Ashadha":
		return "P.Ashadha"
	case "Uttara Ashadha":
		return "U.Ashadha"
	case "Purva Bhadrapada":
		return "P.Bhadra"
	case "Uttara Bhadrapada":
		return "U.Bhadra"
	default:
		return nakName
	}
}

// GetSignIndexByName คืนค่า index ของราศี (0-11) จากชื่อราศีภาษาอังกฤษ
func GetSignIndexByName(signName string) int {
	for i, s := range zodiacSigns {
		if s == signName {
			return i
		}
	}
	return -1
}

