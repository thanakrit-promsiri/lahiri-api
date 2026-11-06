package api

import (
	"errors"
	"lahiri-api/internal/calculator"
)

// --- Request Structs ---
type CalculateWrapperReq struct {
	CalculateReq CalculateReq `json:"calculateReq"`
}

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

// Validate ตรวจสอบความถูกต้องของพารามิเตอร์ที่รับเข้ามา
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

// --- Response Structs ---
type ErrorResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CalculateWrapperResp struct {
	Status        string        `json:"status"`
	CalculateResp CalculateResp `json:"calculateResp"`
}

type CalculateResp struct {
	Metadata         ApiMetadata                 `json:"metadata"`
	Planets          []ApiPlanet                 `json:"planets"`
	HousesAndAngles  []ApiHouse                  `json:"houses_and_angles"`
	WholeSignHouses  []ApiWholeSign              `json:"whole_sign_houses"`
	VimshottariDasha calculator.VimshottariDasha `json:"vimshottari_dasha"`
}

type DashaWrapperResp struct {
	Status           string                      `json:"status"`
	VimshottariDasha calculator.VimshottariDasha `json:"vimshottari_dasha"`
}

type PromptWrapperResp struct {
	Status string `json:"status"`
	Prompt string `json:"prompt"`
}

// --- Transit Request & Response Structs ---

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

// TransitWrapperReq รองรับทั้ง JSON แบบ Flat ตามสเปกใหม่ และแบบห่อใน calculateReq / transitReq
type TransitWrapperReq struct {
	TransitReq
	CalculateReq *TransitReq `json:"calculateReq,omitempty"`
	WrappedReq   *TransitReq `json:"transitReq,omitempty"`
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

type TransitWrapperResp struct {
	Status          string              `json:"status"`
	MonthlyTransits []ApiMonthlyTransit `json:"monthly_transits"`
}

type TransitMonthWrapperResp struct {
	Status         string            `json:"status"`
	MonthlyTransit ApiMonthlyTransit `json:"monthly_transit"`
}

type ApiMonthlyTransit struct {
	Month   int                `json:"month"`
	Date    string             `json:"date"`
	Planets []ApiTransitPlanet `json:"planets"`
}

type ApiTransitPlanet struct {
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
	Name            string                    `json:"name"`
	NameThai        string                    `json:"name_thai"`
	Symbol          string                    `json:"symbol"`
	Longitude       float64                   `json:"longitude"`
	PositionStr     string                    `json:"position_string"`
	DegreeFormatted string                    `json:"degree_formatted"`
	Position        calculator.PositionDetail `json:"position"`
	Nakshatra       string                    `json:"nakshatra"`
	NakshatraThai   string                    `json:"nakshatra_thai"`
	Pada            int                       `json:"pada"`
	Lord            string                    `json:"lord"`
	LordThai        string                    `json:"lord_thai"`
	IsRetrograde    bool                      `json:"is_retrograde"`
	IsCombust       bool                      `json:"is_combust"`
	Dignity         string                    `json:"dignity"`
	DignityThai     string                    `json:"dignity_thai"`
	House           int                       `json:"house"`
	HouseNameThai   string                    `json:"house_name_thai"`
	HouseThai       string                    `json:"house_thai"`
	Navamsa         calculator.NavamsaDetail  `json:"navamsa"`
}

type ApiHouse struct {
	Name            string                    `json:"name"`
	NameThai        string                    `json:"name_thai"`
	Symbol          string                    `json:"symbol"`
	Degree          float64                   `json:"degree"`
	PositionStr     string                    `json:"position_string"`
	DegreeFormatted string                    `json:"degree_formatted"`
	Position        calculator.PositionDetail `json:"position"`
	Nakshatra       string                    `json:"nakshatra"`
	NakshatraThai   string                    `json:"nakshatra_thai"`
	Pada            int                       `json:"pada"`
	Lord            string                    `json:"lord"`
	LordThai        string                    `json:"lord_thai"`
	House           int                       `json:"house"`
	HouseNameThai   string                    `json:"house_name_thai"`
	HouseThai       string                    `json:"house_thai"`
	Navamsa         calculator.NavamsaDetail  `json:"navamsa"`
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
