package api

import (
	"encoding/json"
	"net/http"

	"lahiri-api/internal/calculator"

	"github.com/mshafiee/swephgo"
	"go.uber.org/zap"
)

// ตัวแปร Log สำหรับรองรับการใช้งานภายในแพ็กเกจ api
var Log *zap.Logger

func SetLogger(l *zap.Logger) {
	Log = l
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
	signs := calculator.GetZodiacSigns()
	englishSigns := calculator.GetEnglishSigns()

	// 1. คำนวณตำแหน่งดาว
	planetRes := calculator.CalculatePlanets(jd, flag, signs)
	var apiPlanets []ApiPlanet
	for _, p := range planetRes {
		nak, pada, _ := calculator.GetNakshatra(p.Longitude)
		posDetail := calculator.ParsePositionDetail(p.Longitude)

		apiPlanets = append(apiPlanets, ApiPlanet{
			Name:         p.Name,
			Longitude:    p.Longitude,
			PositionStr:  p.PosStr,
			Position:     posDetail,
			Nakshatra:    nak.Name,
			Pada:         pada,
			Lord:         nak.Lord,
			IsRetrograde: p.IsRetro,
		})
	}

	// 2. คำนวณมุมภพ
	houseSysCode := 'W'
	if len(req.HouseSystem) > 0 {
		houseSysCode = rune(req.HouseSystem[0])
	}

	houseRes := calculator.CalculateHouses(jd, req.Latitude, req.Longitude, int(houseSysCode), ayanamsa, signs)
	var apiHouses []ApiHouse
	for _, h := range houseRes {
		nak, pada, _ := calculator.GetNakshatra(h.Degree)
		posDetail := calculator.ParsePositionDetail(h.Degree)

		apiHouses = append(apiHouses, ApiHouse{
			Name:        h.Name,
			Degree:      h.Degree,
			PositionStr: h.PosStr,
			Position:    posDetail,
			Nakshatra:   nak.Name,
			Pada:        pada,
			Lord:        nak.Lord,
		})
	}

	// 3. คำนวณ Whole Sign Houses
	var apiWholeSigns []ApiWholeSign
	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	signIdx, _ := calculator.ParseLongitude(ascDeg)

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
				AyanamsaDegrees: calculator.FormatDMS(ayanamsa),
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