package api

import (
	"encoding/json"
	"net/http"

	"lahiri-api/internal/calculator"

	"github.com/mshafiee/swephgo"
	"go.uber.org/zap"
)

// Handler จัดการ HTTP requests พร้อม Dependency Injection ของ Logger
type Handler struct {
	logger   *zap.Logger
	onResult func(CalculateResp)
}

// NewHandler สร้าง Handler ใหม่ หากไม่ส่ง logger จะใช้ zap.NewNop() เพื่อป้องกัน nil pointer
func NewHandler(logger *zap.Logger, onResult func(CalculateResp)) *Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Handler{
		logger:   logger,
		onResult: onResult,
	}
}

// ตัวแปร Log สำหรับรองรับ Backward Compatibility
var Log *zap.Logger

func SetLogger(l *zap.Logger) {
	Log = l
}

// HandleCalculateHora เป็นฟังก์ชันระดับ Package ที่คงไว้เพื่อความเข้ากันได้กับโค้ด/เทสต์เดิม
func HandleCalculateHora(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculateHora(w, r)
}

// HandleCalculateHora ประมวลผลคำขอคำนวณดวงชะตา
func (h *Handler) HandleCalculateHora(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid HTTP method attempted", zap.String("method", r.Method))
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqWrapper CalculateWrapperReq
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqWrapper); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	req := reqWrapper.CalculateReq
	if err := req.Validate(); err != nil {
		h.logger.Warn("Invalid request parameters", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info("Received API calculation request",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("day", req.Day),
		zap.Float64("latitude", req.Latitude),
		zap.Float64("longitude", req.Longitude),
	)

	jd, ayanamsa := calculator.ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)

	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	signs := calculator.GetZodiacSigns()
	englishSigns := calculator.GetEnglishSigns()
	thaiSigns := calculator.GetThaiZodiacSigns()

	// 1. คำนวณตำแหน่งดาว
	planetRes := calculator.CalculatePlanets(jd, flag, signs)
	apiPlanets := make([]ApiPlanet, 0, len(planetRes))
	for _, p := range planetRes {
		nak, pada, _ := calculator.GetNakshatra(p.Longitude)
		posDetail := calculator.ParsePositionDetail(p.Longitude)
		navamsaRes := calculator.CalculateNavamsa(p.Longitude)

		apiPlanets = append(apiPlanets, ApiPlanet{
			Name:         p.Name,
			Longitude:    p.Longitude,
			PositionStr:  p.PosStr,
			Position:     posDetail,
			Nakshatra:    nak.Name,
			Pada:         pada,
			Lord:         nak.Lord,
			IsRetrograde: p.IsRetro,
			Navamsa:      navamsaRes,
		})
	}

	// 2. คำนวณมุมภพ
	houseSysCode := 'W'
	if len(req.HouseSystem) > 0 {
		houseSysCode = rune(req.HouseSystem[0])
	}

	houseRes := calculator.CalculateHouses(jd, req.Latitude, req.Longitude, int(houseSysCode), ayanamsa, signs)
	apiHouses := make([]ApiHouse, 0, len(houseRes))
	for _, hRes := range houseRes {
		nak, pada, _ := calculator.GetNakshatra(hRes.Degree)
		posDetail := calculator.ParsePositionDetail(hRes.Degree)
		navamsaRes := calculator.CalculateNavamsa(hRes.Degree)

		apiHouses = append(apiHouses, ApiHouse{
			Name:        hRes.Name,
			Degree:      hRes.Degree,
			PositionStr: hRes.PosStr,
			Position:    posDetail,
			Nakshatra:   nak.Name,
			Pada:        pada,
			Lord:        nak.Lord,
			Navamsa:     navamsaRes,
		})
	}

	// 3. คำนวณ Whole Sign Houses
	apiWholeSigns := make([]ApiWholeSign, 0, 12)
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
			SignThai: thaiSigns[currentSign],
		})
	}

	calcResp := CalculateResp{
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
	}

	if h.onResult != nil {
		h.onResult(calcResp)
	}

	resp := CalculateWrapperResp{
		Status:        "success",
		CalculateResp: calcResp,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode response JSON", zap.Error(err))
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResp{
		Status:  "error",
		Message: message,
	})
}
