package calculator

// PositionDetail เก็บรายละเอียดองศา ลิปดา ฟิลิปดา และราศี
type PositionDetail struct {
	Sign   string  `json:"sign"`   // เช่น "Aries", "Libra"
	Degree int     `json:"degree"` // องศา (0-29)
	Minute int     `json:"minute"` // ลิปดา (0-59)
	Second float64 `json:"second"` // ฟิลิปดา (0-59.99)
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
	Index    int     // ลำดับที่ 1-27
	Name     string  // ชื่อนักษัตร
	Lord     string  // ดาวเจ้าเรือน (Vimshottari Dasha Lord)
	StartDeg float64 // องศาเริ่มต้นในจักรวาล (0 - 360)
	EndDeg   float64 // องศาสิ้นสุด
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