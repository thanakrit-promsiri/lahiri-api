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
	Metadata        ApiMetadata    `json:"metadata"`
	Planets         []ApiPlanet    `json:"planets"`
	HousesAndAngles []ApiHouse     `json:"houses_and_angles"`
	WholeSignHouses []ApiWholeSign `json:"whole_sign_houses"`
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
	Name         string                    `json:"name"`
	Longitude    float64                   `json:"longitude"`
	PositionStr  string                    `json:"position_string"`
	Position     calculator.PositionDetail `json:"position"`
	Nakshatra    string                    `json:"nakshatra"`
	Pada         int                       `json:"pada"`
	Lord         string                    `json:"lord"`
	IsRetrograde bool                      `json:"is_retrograde"`
	Navamsa      calculator.NavamsaDetail  `json:"navamsa"`
}

type ApiHouse struct {
	Name        string                    `json:"name"`
	Degree      float64                   `json:"degree"`
	PositionStr string                    `json:"position_string"`
	Position    calculator.PositionDetail `json:"position"`
	Nakshatra   string                    `json:"nakshatra"`
	Pada        int                       `json:"pada"`
	Lord        string                    `json:"lord"`
	Navamsa     calculator.NavamsaDetail  `json:"navamsa"`
}

type ApiWholeSign struct {
	House    int    `json:"house"`
	Sign     string `json:"sign"`
	SignThai string `json:"sign_thai"`
}