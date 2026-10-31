package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

// HandleCalculateDasha เป็นฟังก์ชันระดับ Package สำหรับคำนวณเฉพาะวิมโสตรีทศา
func HandleCalculateDasha(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculateDasha(w, r)
}

// normalizeReq ปรับค่าเริ่มต้น เช่น แปลงปี พ.ศ. เป็น ค.ศ. อัตโนมัติหากเกิน 2400
func normalizeReq(req *CalculateReq) {
	if req.Year > 2400 {
		req.Year -= 543
	}
	if req.HouseSystem == "" {
		req.HouseSystem = "W"
	}
}

// HandleCalculateHora ประมวลผลคำขอคำนวณดวงชะตาแบบครบวงจร (ดาว, ภพ, มาตรฐานดาว, ดาวดับ, เจ้าเรือนภพ, วิมโสตรีทศา)
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
	normalizeReq(&req)
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

	// 1. คำนวณมุมภพ (Houses & Angles) เพื่อหาลัคนา (Ascendant) ก่อนนำไปจัดภพ Whole Sign
	houseSysCode := 'W'
	if len(req.HouseSystem) > 0 {
		houseSysCode = rune(req.HouseSystem[0])
	}
	houseRes := calculator.CalculateHouses(jd, req.Latitude, req.Longitude, int(houseSysCode), ayanamsa, signs)

	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	ascSignIdx, _ := calculator.ParseLongitude(ascDeg)

	// 2. คำนวณตำแหน่งดาว (Planets) พร้อมมาตรฐานดาว (Dignity), ดาวดับ (Combust) และภพที่สถิต
	planetRes := calculator.CalculatePlanets(jd, flag, signs)

	// ค้นหาลองจิจูดของดวงอาทิตย์ (เพื่อเช็คดาวดับ) และดวงจันทร์ (เพื่อคำนวณวิมโสตรีทศา)
	var sunLon, moonLon float64
	for _, p := range planetRes {
		if p.Name == "Sun" {
			sunLon = p.Longitude
		} else if p.Name == "Moon" {
			moonLon = p.Longitude
		}
	}

	// เก็บข้อมูลว่าดาวแต่ละดวง (Canonical Name) สถิตอยู่ในภพที่เท่าไหร่ เพื่อนำไปผูกกับเจ้าเรือนภพ (Whole Sign Houses)
	planetHouseMap := make(map[string]int, len(planetRes))

	apiPlanets := make([]ApiPlanet, 0, len(planetRes))
	for _, p := range planetRes {
		nak, pada, _ := calculator.GetNakshatra(p.Longitude)
		posDetail := calculator.ParsePositionDetail(p.Longitude)
		navamsaRes := calculator.CalculateNavamsa(p.Longitude)

		nameThai, symbol, canonical := calculator.TranslateBodyInfo(p.Name)
		dignity, dignityThai := calculator.CalculateDignity(p.Name, p.Longitude)
		isCombust := calculator.IsPlanetCombust(p.Name, p.Longitude, sunLon, p.IsRetro)

		pSignIdx, _ := calculator.ParseLongitude(p.Longitude)
		houseNum := ((pSignIdx - ascSignIdx + 12) % 12) + 1
		houseNameThai := calculator.GetHouseNameThai(houseNum)

		planetHouseMap[canonical] = houseNum

		apiPlanets = append(apiPlanets, ApiPlanet{
			Name:            p.Name,
			NameThai:        nameThai,
			Symbol:          symbol,
			Longitude:       p.Longitude,
			PositionStr:     p.PosStr,
			DegreeFormatted: calculator.FormatDegreeShort(posDetail),
			Position:        posDetail,
			Nakshatra:       nak.Name,
			NakshatraThai:   nak.NameThai,
			Pada:            pada,
			Lord:            nak.Lord,
			LordThai:        nak.LordThai,
			IsRetrograde:    p.IsRetro,
			IsCombust:       isCombust,
			Dignity:         dignity,
			DignityThai:     dignityThai,
			House:           houseNum,
			HouseNameThai:   houseNameThai,
			HouseThai:       fmt.Sprintf("%d (%s)", houseNum, houseNameThai),
			Navamsa:         navamsaRes,
		})
	}

	// 3. จัดเตรียมข้อมูลมุมภพ (Ascendant, MC, DSC, IC)
	apiHouses := make([]ApiHouse, 0, len(houseRes))
	for _, hRes := range houseRes {
		nak, pada, _ := calculator.GetNakshatra(hRes.Degree)
		posDetail := calculator.ParsePositionDetail(hRes.Degree)
		navamsaRes := calculator.CalculateNavamsa(hRes.Degree)

		nameThai, symbol, _ := calculator.TranslateBodyInfo(hRes.Name)
		hSignIdx, _ := calculator.ParseLongitude(hRes.Degree)
		houseNum := ((hSignIdx - ascSignIdx + 12) % 12) + 1
		houseNameThai := calculator.GetHouseNameThai(houseNum)

		apiHouses = append(apiHouses, ApiHouse{
			Name:            hRes.Name,
			NameThai:        nameThai,
			Symbol:          symbol,
			Degree:          hRes.Degree,
			PositionStr:     hRes.PosStr,
			DegreeFormatted: calculator.FormatDegreeShort(posDetail),
			Position:        posDetail,
			Nakshatra:       nak.Name,
			NakshatraThai:   nak.NameThai,
			Pada:            pada,
			Lord:            nak.Lord,
			LordThai:        nak.LordThai,
			House:           houseNum,
			HouseNameThai:   houseNameThai,
			HouseThai:       fmt.Sprintf("%d (%s)", houseNum, houseNameThai),
			Navamsa:         navamsaRes,
		})
	}

	// 4. คำนวณ Whole Sign Houses พร้อมเจ้าเรือนภพ (House Lords) และภพที่ดาวเจ้าเรือนไปสถิต
	apiWholeSigns := make([]ApiWholeSign, 0, 12)
	for i := 1; i <= 12; i++ {
		currentSign := (ascSignIdx + (i - 1)) % 12
		houseNameThai := calculator.GetHouseNameThai(i)
		lord := calculator.GetSignLord(currentSign)
		lordThai := calculator.TranslateLordThai(lord)
		placedInHouse := planetHouseMap[lord]
		placedInHouseThai := calculator.GetHouseNameThai(placedInHouse)

		apiWholeSigns = append(apiWholeSigns, ApiWholeSign{
			House:                 i,
			HouseNameThai:         houseNameThai,
			Sign:                  englishSigns[currentSign],
			SignThai:              thaiSigns[currentSign],
			Lord:                  lord,
			LordThai:              lordThai,
			LordPlacedInHouse:     placedInHouse,
			LordPlacedInHouseThai: placedInHouseThai,
		})
	}

	// 5. คำนวณวิมโสตรีทศา (Vimshottari Dasha) จากตำแหน่งดวงจันทร์และวันเวลาเกิด
	dasha := calculator.CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, time.Now(),
	)

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
		Planets:          apiPlanets,
		HousesAndAngles:  apiHouses,
		WholeSignHouses:  apiWholeSigns,
		VimshottariDasha: dasha,
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

// HandleCalculateDasha ประมวลผลเฉพาะวิมโสตรีทศา (POST /api/v1/hora/dasha)
func (h *Handler) HandleCalculateDasha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid HTTP method attempted on dasha endpoint", zap.String("method", r.Method))
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqWrapper CalculateWrapperReq
	if err := json.NewDecoder(r.Body).Decode(&reqWrapper); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	req := reqWrapper.CalculateReq
	normalizeReq(&req)
	if err := req.Validate(); err != nil {
		h.logger.Warn("Invalid request parameters", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	jd, _ := calculator.ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)

	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	signs := calculator.GetZodiacSigns()
	planetRes := calculator.CalculatePlanets(jd, flag, signs)

	var moonLon float64
	for _, p := range planetRes {
		if p.Name == "Moon" {
			moonLon = p.Longitude
			break
		}
	}

	dasha := calculator.CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, time.Now(),
	)

	resp := DashaWrapperResp{
		Status:           "success",
		VimshottariDasha: dasha,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode dasha response JSON", zap.Error(err))
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
