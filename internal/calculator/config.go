package calculator

var (
	zodiacSigns = []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
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
		"Midheaven (MC)":    {NameThai: "เมษูรณ์ (MC)", Symbol: "MC", Canonical: "MC"},
		"Descendant (DSC)":  {NameThai: "ขอบฟ้าตก (DSC)", Symbol: "DSC", Canonical: "DSC"},
		"Imum Coeli (IC)":   {NameThai: "จุดใต้เท้า (IC)", Symbol: "IC", Canonical: "IC"},
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
