package main

import (
	"encoding/json"
	"net/http"

	"github.com/mshafiee/swephgo"
	"go.uber.org/zap"
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

type ApiPositionDetail struct {
	Sign   string  `json:"sign"`   // เช่น "Libra"
	Degree int     `json:"degree"` // องศา (0-29)
	Minute int     `json:"minute"` // ลิปดา (0-59)
	Second float64 `json:"second"` // ฟิลิปดา (0-59.99)
}

type ApiPlanet struct {
	Name         string            `json:"name"`
	Longitude    float64           `json:"longitude"`
	PositionStr  string            `json:"position_string"`
	Position     ApiPositionDetail `json:"position"`
	Nakshatra    string            `json:"nakshatra"`
	Pada         int               `json:"pada"`
	Lord         string            `json:"lord"`
	IsRetrograde bool              `json:"is_retrograde"`
}

type ApiHouse struct {
	Name        string            `json:"name"`
	Degree      float64           `json:"degree"`
	PositionStr string            `json:"position_string"`
	Position    ApiPositionDetail `json:"position"`
	Nakshatra   string            `json:"nakshatra"`
	Pada        int               `json:"pada"`
	Lord        string            `json:"lord"`
}

type ApiWholeSign struct {
	House    int    `json:"house"`
	Sign     string `json:"sign"`
	SignThai string `json:"sign_thai"`
}

// --- Helper Function: แกะค่าองศาเป็นรายละเอียดปลีกย่อย ---
func ParsePositionDetail(lon float64) ApiPositionDetail {
	englishSigns := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	
	normalizedLon := lon
	for normalizedLon < 0 {
		normalizedLon += 360.0
	}
	for normalizedLon >= 360.0 {
		normalizedLon -= 360.0
	}

	signIdx := int(normalizedLon / 30.0)
	degInSign := normalizedLon - float64(signIdx*30)
	
	deg := int(degInSign)
	minFloat := (degInSign - float64(deg)) * 60.0
	min := int(minFloat)
	sec := (minFloat - float64(min)) * 60.0

	signName := ""
	if signIdx >= 0 && signIdx < len(englishSigns) {
		signName = englishSigns[signIdx]
	}

	return ApiPositionDetail{
		Sign:   signName,
		Degree: deg,
		Minute: min,
		Second: sec,
	}
}

// --- HTTP Handler ---
func HandleCalculateHora(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if Log != nil {
			Log.Warn("Invalid HTTP method attempted", zap.String("method", r.Method))
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqWrapper CalculateWrapperReq
	err := json.NewDecoder(r.Body).Decode(&reqWrapper)
	if err != nil {
		if Log != nil {
			Log.Error("Failed to decode request body", zap.Error(err))
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := reqWrapper.CalculateReq
	if Log != nil {
		Log.Info("Received API calculation request",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Int("day", req.Day),
			zap.Float64("latitude", req.Latitude),
			zap.Float64("longitude", req.Longitude),
		)
	}

	hourUTC := req.HourThai - req.TimeZone
	jd := swephgo.Julday(req.Year, req.Month, req.Day, hourUTC, swephgo.SeGregCal)

	swephgo.SetSidMode(swephgo.SeSidmLahiri, 0, 0)
	ayanamsa := swephgo.GetAyanamsa(jd)
	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	signs := GetZodiacSigns()

	// 1. คำนวณตำแหน่งดาว
	planetRes := CalculatePlanets(jd, flag, signs)
	var apiPlanets []ApiPlanet
	for _, p := range planetRes {
		nak, pada, _ := GetNakshatra(p.Longitude)
		posDetail := ParsePositionDetail(p.Longitude)

		apiPlanets = append(apiPlanets, ApiPlanet{
			Name:         p.Name,
			Longitude:    p.Longitude,
			PositionStr:  p.PosStr,
			Position:     posDetail,
			Nakshatra:    nak.Name,
			Pada:         pada,
			Lord:         nak.Lord, // ใช้ชื่อดาวจาก Nakshatra โดยตรง
			IsRetrograde: p.IsRetro,
		})
	}

	// 2. คำนวณมุมภพ
	houseSysCode := 'W'
	if len(req.HouseSystem) > 0 {
		houseSysCode = rune(req.HouseSystem[0])
	}

	houseRes := CalculateHouses(jd, req.Latitude, req.Longitude, int(houseSysCode), ayanamsa, signs)
	var apiHouses []ApiHouse
	for _, h := range houseRes {
		nak, pada, _ := GetNakshatra(h.Degree)
		posDetail := ParsePositionDetail(h.Degree)

		apiHouses = append(apiHouses, ApiHouse{
			Name:        h.Name,
			Degree:      h.Degree,
			PositionStr: h.PosStr,
			Position:    posDetail,
			Nakshatra:   nak.Name,
			Pada:        pada,
			Lord:        nak.Lord, // ใช้ชื่อดาวจาก Nakshatra โดยตรง
		})
	}

	// 3. คำนวณ Whole Sign Houses
	var apiWholeSigns []ApiWholeSign
	englishSigns := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	signIdx, _ := ParseLongitude(ascDeg)

	for i := 1; i <= 12; i++ {
		currentSign := (signIdx + (i - 1)) % 12
		apiWholeSigns = append(apiWholeSigns, ApiWholeSign{
			House:    i,
			Sign:     englishSigns[currentSign],
			SignThai: signs[currentSign],
		})
	}

	resp := CalculateWrapperResp{
		Status: "success",
		CalculateResp: CalculateResp{
			Metadata: ApiMetadata{
				JulianDay:       jd,
				AyanamsaName:    "Lahiri",
				AyanamsaDegrees: FormatDMS(ayanamsa),
				Location: ApiLocation{
					Latitude:  req.Latitude,
					Longitude: req.Longitude,
				},
			},
			Planets:         apiPlanets,
			HousesAndAngles: apiHouses,
			WholeSignHouses: apiWholeSigns,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}