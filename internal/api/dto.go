package api

import "lahiri-api/internal/calculator"

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

// --- Response Structs ---
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