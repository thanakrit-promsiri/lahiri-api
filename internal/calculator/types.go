package calculator

// PositionDetail เก็บรายละเอียดองศา ลิปดา ฟิลิปดา และราศี
type PositionDetail struct {
	Sign   string  `json:"sign"`   
	Degree int     `json:"degree"` 
	Minute int     `json:"minute"` 
	Second float64 `json:"second"` 
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
	Lord     string  
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