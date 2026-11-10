package calculator

import (
	"errors"
)

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
	TargetDate        string          `json:"target_date,omitempty"`
	MoonNakshatra     string          `json:"moon_nakshatra"`
	MoonNakshatraThai string          `json:"moon_nakshatra_thai"`
	BalanceAtBirth    string          `json:"balance_at_birth"`
	CurrentDasha      CurrentDasha    `json:"current_dasha"`
	MahadashaList     []MahadashaItem `json:"mahadasha_list"`
}

// CurrentDasha เก็บข้อมูลมหาทศา อันตรทศา และปรัตยันตรทศาในเวลาปัจจุบัน
type CurrentDasha struct {
	MahadashaLord           string `json:"mahadasha_lord"`
	MahadashaLordThai       string `json:"mahadasha_lord_thai"`
	MahadashaStart          string `json:"mahadasha_start"`
	MahadashaEnd            string `json:"mahadasha_end"`
	AntardashaLord          string `json:"antardasha_lord"`
	AntardashaLordThai      string `json:"antardasha_lord_thai"`
	AntardashaStart         string `json:"antardasha_start"`
	AntardashaEnd           string `json:"antardasha_end"`
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

// --- Request Parameters ---

type CalculateReq struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	Day         int     `json:"day"`
	HourThai    float64 `json:"hour_thai"`
	TimeZone    float64 `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	HouseSystem string  `json:"house_system"`
}

func (r CalculateReq) Validate() error {
	if r.Year == 0 {
		return errors.New("year is required and cannot be 0")
	}
	if r.Month < 1 || r.Month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if r.Day < 1 || r.Day > 31 {
		return errors.New("day must be between 1 and 31")
	}
	if r.HourThai < 0 || r.HourThai >= 24 {
		return errors.New("hour_thai must be between 0 and 24")
	}
	if r.TimeZone < -14 || r.TimeZone > 14 {
		return errors.New("time_zone must be between -14 and 14")
	}
	if r.Latitude < -90 || r.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

type TransitReq struct {
	Year         int     `json:"year"`
	Month        int     `json:"month"`
	Day          int     `json:"day"`
	HourThai     float64 `json:"hour_thai"`
	TimeZone     float64 `json:"time_zone"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	HouseSystem  string  `json:"house_system"`
	TransitYear  int     `json:"transit_year"`
	TransitMonth int     `json:"transit_month,omitempty"`
}

func (r TransitReq) Validate() error {
	if r.Year == 0 {
		return errors.New("year is required and cannot be 0")
	}
	if r.Month < 1 || r.Month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if r.Day < 1 || r.Day > 31 {
		return errors.New("day must be between 1 and 31")
	}
	if r.HourThai < 0 || r.HourThai >= 24 {
		return errors.New("hour_thai must be between 0 and 24")
	}
	if r.TimeZone < -14 || r.TimeZone > 14 {
		return errors.New("time_zone must be between -14 and 14")
	}
	if r.Latitude < -90 || r.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	if r.TransitYear < 1 {
		return errors.New("transit_year must be a valid positive year")
	}
	return nil
}

func (r TransitReq) ValidateWithMonth() error {
	if err := r.Validate(); err != nil {
		return err
	}
	if r.TransitMonth < 1 || r.TransitMonth > 12 {
		return errors.New("transit_month must be between 1 and 12")
	}
	return nil
}

// --- Calculation Response Models ---

type CalculateResp struct {
	Metadata         ApiMetadata      `json:"metadata"`
	Planets          []ApiPlanet      `json:"planets"`
	HousesAndAngles  []ApiHouse       `json:"houses_and_angles"`
	WholeSignHouses  []ApiWholeSign   `json:"whole_sign_houses"`
	VimshottariDasha VimshottariDasha `json:"vimshottari_dasha"`
}

type ApiMetadata struct {
	JulianDay       float64     `json:"julian_day"`
	AyanamsaName    string      `json:"ayanamsa_name"`
	AyanamsaDegrees string      `json:"ayanamsa_degrees"`
	Location        ApiLocation `json:"location"`
}

type ApiLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ApiPlanet struct {
	Name            string         `json:"name"`
	NameThai        string         `json:"name_thai"`
	Symbol          string         `json:"symbol"`
	Longitude       float64        `json:"longitude"`
	PositionStr     string         `json:"position_string"`
	DegreeFormatted string         `json:"degree_formatted"`
	Position        PositionDetail `json:"position"`
	Nakshatra       string         `json:"nakshatra"`
	NakshatraThai   string         `json:"nakshatra_thai"`
	Pada            int            `json:"pada"`
	Lord            string         `json:"lord"`
	LordThai        string         `json:"lord_thai"`
	IsRetrograde    bool           `json:"is_retrograde"`
	IsCombust       bool           `json:"is_combust"`
	Dignity         string         `json:"dignity"`
	DignityThai     string         `json:"dignity_thai"`
	House           int            `json:"house"`
	HouseNameThai   string         `json:"house_name_thai"`
	HouseThai       string         `json:"house_thai"`
	Navamsa         NavamsaDetail  `json:"navamsa"`
}

type ApiHouse struct {
	Name            string         `json:"name"`
	NameThai        string         `json:"name_thai"`
	Symbol          string         `json:"symbol"`
	Degree          float64        `json:"degree"`
	PositionStr     string         `json:"position_string"`
	DegreeFormatted string         `json:"degree_formatted"`
	Position        PositionDetail `json:"position"`
	Nakshatra       string         `json:"nakshatra"`
	NakshatraThai   string         `json:"nakshatra_thai"`
	Pada            int            `json:"pada"`
	Lord            string         `json:"lord"`
	LordThai        string         `json:"lord_thai"`
	House           int            `json:"house"`
	HouseNameThai   string         `json:"house_name_thai"`
	HouseThai       string         `json:"house_thai"`
	Navamsa         NavamsaDetail  `json:"navamsa"`
}

type ApiWholeSign struct {
	House                 int    `json:"house"`
	HouseNameThai         string `json:"house_name_thai"`
	Sign                  string `json:"sign"`
	SignThai              string `json:"sign_thai"`
	Lord                  string `json:"lord"`
	LordThai              string `json:"lord_thai"`
	LordPlacedInHouse     int    `json:"lord_placed_in_house"`
	LordPlacedInHouseThai string `json:"lord_placed_in_house_thai"`
}

// --- Transit Models ---

type MonthlyTransit struct {
	Month   int             `json:"month"`
	Date    string          `json:"date"`
	Dasha   *TransitDasha   `json:"dasha,omitempty"`
	Planets []TransitPlanet `json:"planets"`
}

type TransitDasha struct {
	MoonNakshatra           string `json:"moon_nakshatra"`
	MoonNakshatraThai       string `json:"moon_nakshatra_thai"`
	MahadashaLord           string `json:"mahadasha_lord"`
	MahadashaLordThai       string `json:"mahadasha_lord_thai"`
	MahadashaStart          string `json:"mahadasha_start"`
	MahadashaEnd            string `json:"mahadasha_end"`
	AntardashaLord          string `json:"antardasha_lord"`
	AntardashaLordThai      string `json:"antardasha_lord_thai"`
	AntardashaStart         string `json:"antardasha_start"`
	AntardashaEnd           string `json:"antardasha_end"`
	PratyantardashaLord     string `json:"pratyantardasha_lord"`
	PratyantardashaLordThai string `json:"pratyantardasha_lord_thai"`
	PratyantardashaStart    string `json:"pratyantardasha_start"`
	PratyantardashaEnd      string `json:"pratyantardasha_end"`
}

type TransitPlanet struct {
	No            int     `json:"no"`
	Name          string  `json:"name"`
	NameThai      string  `json:"name_thai"`
	Symbol        string  `json:"symbol"`
	Longitude     float64 `json:"longitude"`
	Sign          string  `json:"sign"`
	SignThai      string  `json:"sign_thai"`
	Degree        string  `json:"degree"`
	House         int     `json:"house"`
	HouseNameThai string  `json:"house_name_thai"`
	IsRetrograde  bool    `json:"is_retrograde"`
	IsCombust     bool    `json:"is_combust"`
	Dignity       string  `json:"dignity"`
	DignityThai   string  `json:"dignity_thai"`
	Nakshatra     string  `json:"nakshatra"`
	NakshatraThai string  `json:"nakshatra_thai"`
	Pada          int     `json:"pada"`
}

// --- Transit Prompt Data Models ---

type TransitPromptBody struct {
	Code        string
	House       int
	Lon         float64
	D1SignShort string
	DMS         string
	D1DigShort  string
	NakShort    string
	Pada        int
	LordShort   string
	D9SignShort string
	D9DigShort  string
	Flags       string
	Formatted   string
}

type TransitDashaPromptData struct {
	DateISO             string
	NatalAscSignShort   string
	NatalMoonSignShort  string
	NatalMoonHouse      int
	NatalMoonNakName    string
	Dasha               VimshottariDasha
	NatalHouseByShort   map[string]int
	TransitHouseByShort map[string]int
	TransitPlanets      []TransitPromptBody
	VedicBodies         []TransitPromptBody
}
