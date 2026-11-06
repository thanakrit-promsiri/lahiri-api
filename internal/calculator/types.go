package calculator

// PositionDetail เก็บรายละเอียดองศา ลิปดา ฟิลิปดา และราศี (ทั้งอังกฤษและไทย)
type PositionDetail struct {
	Sign     string  `json:"sign"`
	SignThai string  `json:"sign_thai"`
	Degree   int     `json:"degree"`
	Minute   int     `json:"minute"`
	Second   float64 `json:"second"`
}

// NavamsaDetail เก็บข้อมูลผลลัพธ์ของนวางศ์ (D9)
type NavamsaDetail struct {
	Sign     string `json:"sign"`      // ราศีนวางศ์ (ภาษาอังกฤษ)
	SignThai string `json:"sign_thai"` // ราศีนวางศ์ (ภาษาไทย)
	Index    int    `json:"index"`     // นวางศ์ที่ (1-9)
}

// PlanetResult เก็บผลลัพธ์การคำนวณดาวแต่ละดวง
type PlanetResult struct {
	Name      string
	Longitude float64
	Speed     float64
	PosStr    string
	NakStr    string
	IsRetro   bool
}

// HouseResult เก็บผลลัพธ์มุมลัคนาและภพ
type HouseResult struct {
	Name   string
	Degree float64
	PosStr string
	NakStr string
}

// NakshatraInfo เก็บข้อมูลรายละเอียดของแต่ละนักษัตร
type NakshatraInfo struct {
	Index    int
	Name     string
	NameThai string
	Lord     string
	LordThai string
	StartDeg float64
	EndDeg   float64
}

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

// VimshottariDasha เก็บผลลัพธ์การคำนวณวิมโสตรีทศา (120 ปี)
type VimshottariDasha struct {
	MoonNakshatra     string          `json:"moon_nakshatra"`
	MoonNakshatraThai string          `json:"moon_nakshatra_thai"`
	BalanceAtBirth    string          `json:"balance_at_birth"`
	CurrentDasha      CurrentDasha    `json:"current_dasha"`
	MahadashaList     []MahadashaItem `json:"mahadasha_list"`
}

// CurrentDasha เก็บข้อมูลมหาทศา อันตรทศา และปรัตยันตรทศาในเวลาปัจจุบัน
type CurrentDasha struct {
	MahadashaLord         string `json:"mahadasha_lord"`
	MahadashaLordThai     string `json:"mahadasha_lord_thai"`
	MahadashaStart        string `json:"mahadasha_start"`
	MahadashaEnd          string `json:"mahadasha_end"`
	AntardashaLord        string `json:"antardasha_lord"`
	AntardashaLordThai    string `json:"antardasha_lord_thai"`
	AntardashaStart       string `json:"antardasha_start"`
	AntardashaEnd         string `json:"antardasha_end"`
	PratyantardashaLord     string `json:"pratyantardasha_lord"`
	PratyantardashaLordThai string `json:"pratyantardasha_lord_thai"`
	PratyantardashaStart    string `json:"pratyantardasha_start"`
	PratyantardashaEnd      string `json:"pratyantardasha_end"`
}

// MahadashaItem เก็บช่วงเวลามหาทศาแต่ละยุค
type MahadashaItem struct {
	Lord      string `json:"lord"`
	LordThai  string `json:"lord_thai"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
