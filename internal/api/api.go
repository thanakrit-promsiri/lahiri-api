package api

import (
	"encoding/json"
	"fmt"
	"math"
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

// HandleCalculatePrompt เป็นฟังก์ชันระดับ Package สำหรับสร้างข้อความ Roman Ultra-Compact สำหรับป้อน AI
func HandleCalculatePrompt(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculatePrompt(w, r)
}

// HandleCalculateTransit เป็นฟังก์ชันระดับ Package สำหรับคำนวณดาวจรประจำปี (Monthly Transits 12 เดือน)
func HandleCalculateTransit(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculateTransit(w, r)
}

// HandleCalculateTransitMonth เป็นฟังก์ชันระดับ Package สำหรับคำนวณดาวจรเฉพาะเดือนที่ระบุ
func HandleCalculateTransitMonth(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculateTransitMonth(w, r)
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

// computeHoroscope คำนวณข้อมูลดวงชะตาพระเวททั้งหมดเพื่อใช้ร่วมกันระหว่างหลาย Endpoint
func computeHoroscope(req CalculateReq) CalculateResp {
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

	var sunLon, moonLon float64
	for _, p := range planetRes {
		if p.Name == "Sun" {
			sunLon = p.Longitude
		} else if p.Name == "Moon" {
			moonLon = p.Longitude
		}
	}

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

	return CalculateResp{
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
}

// decodeAndValidateReq อ่านและตรวจสอบความถูกต้องของ Request
func (h *Handler) decodeAndValidateReq(w http.ResponseWriter, r *http.Request) (CalculateReq, bool) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid HTTP method attempted", zap.String("method", r.Method))
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return CalculateReq{}, false
	}

	var reqWrapper CalculateWrapperReq
	if err := json.NewDecoder(r.Body).Decode(&reqWrapper); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return CalculateReq{}, false
	}

	req := reqWrapper.CalculateReq
	normalizeReq(&req)
	if err := req.Validate(); err != nil {
		h.logger.Warn("Invalid request parameters", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return CalculateReq{}, false
	}

	return req, true
}

// HandleCalculateHora ประมวลผลคำขอคำนวณดวงชะตาแบบครบวงจร (POST /api/v1/hora/calculate)
func (h *Handler) HandleCalculateHora(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeAndValidateReq(w, r)
	if !ok {
		return
	}

	h.logger.Info("Received API calculation request",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("day", req.Day),
		zap.Float64("latitude", req.Latitude),
		zap.Float64("longitude", req.Longitude),
	)

	calcResp := computeHoroscope(req)

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

// HandleCalculateDasha ประมวลผลเฉพาะวิมโสตรีทศา มหาทศา อันตรทศา และปรัตยันตรทศา/อนุทศา (POST /api/v1/hora/dasha)
// รองรับทั้งการส่งแบบ Flat JSON และแบบห่อใน calculateReq รวมถึงระบุ transit_year และ transit_month เพื่อหาทศาของปี/เดือนจรนั้นๆ
func (h *Handler) HandleCalculateDasha(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid HTTP method attempted", zap.String("method", r.Method))
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var wrapper TransitWrapperReq
	if err := json.NewDecoder(r.Body).Decode(&wrapper); err != nil {
		h.logger.Error("Failed to decode dasha request body", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return
	}

	req := wrapper.TransitReq
	if wrapper.CalculateReq != nil && wrapper.CalculateReq.Year != 0 {
		if req.TransitYear != 0 && wrapper.CalculateReq.TransitYear == 0 {
			wrapper.CalculateReq.TransitYear = req.TransitYear
		}
		if req.TransitMonth != 0 && wrapper.CalculateReq.TransitMonth == 0 {
			wrapper.CalculateReq.TransitMonth = req.TransitMonth
		}
		req = *wrapper.CalculateReq
	} else if wrapper.WrappedReq != nil && wrapper.WrappedReq.Year != 0 {
		req = *wrapper.WrappedReq
	}

	hasCustomTarget := req.TransitYear != 0 || req.TransitMonth != 0
	normalizeTransitReq(&req)
	if req.TransitMonth != 0 && (req.TransitMonth < 1 || req.TransitMonth > 12) {
		writeJSONError(w, http.StatusBadRequest, "transit_month must be between 1 and 12")
		return
	}
	if err := req.Validate(); err != nil {
		h.logger.Warn("Invalid dasha request parameters", zap.Error(err))
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

	targetTime := time.Now()
	if hasCustomTarget {
		tMonth := req.TransitMonth
		if tMonth == 0 {
			tMonth = req.Month
		}
		tDay := clampDayToMonth(req.TransitYear, tMonth, req.Day)
		tzOffsetSec := int(math.Round(req.TimeZone * 3600))
		loc := time.FixedZone("LocalBirthTime", tzOffsetSec)
		hourInt := int(req.HourThai)
		minFloat := (req.HourThai - float64(hourInt)) * 60.0
		minInt := int(minFloat)
		secInt := int(math.Round((minFloat - float64(minInt)) * 60.0))
		if secInt >= 60 {
			secInt = 59
		}
		targetTime = time.Date(req.TransitYear, time.Month(tMonth), tDay, hourInt, minInt, secInt, 0, loc)
	}

	dasha := calculator.CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, targetTime,
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

// HandleCalculatePrompt ประมวลผลและคืนค่าข้อความดวงชะตาแบบ Roman Ultra-Compact สำหรับป้อน AI (POST /api/v1/hora/prompt)
// รองรับ ?format=raw หรือ Header Accept: text/plain เพื่อรับเป็น Plain Text ทันที หรือรับเป็น JSON {"status":"success","prompt":"..."}
func (h *Handler) HandleCalculatePrompt(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeAndValidateReq(w, r)
	if !ok {
		return
	}

	h.logger.Info("Received AI prompt generation request",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month),
		zap.Int("day", req.Day),
	)

	calcResp := computeHoroscope(req)
	promptText := FormatRomanCompactPrompt(calcResp)

	if r.URL.Query().Get("format") == "raw" || r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(promptText))
		return
	}

	resp := PromptWrapperResp{
		Status: "success",
		Prompt: promptText,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode prompt response JSON", zap.Error(err))
	}
}

// normalizeTransitReq ปรับค่าเริ่มต้นของ TransitReq (เช่น แปลง พ.ศ. เป็น ค.ศ., ตั้งค่า TimeZone = 7.0)
func normalizeTransitReq(req *TransitReq) {
	if req.Year > 2400 {
		req.Year -= 543
	}
	if req.TransitYear == 0 {
		req.TransitYear = time.Now().Year()
	}
	if req.TransitYear > 2400 {
		req.TransitYear -= 543
	}
	if req.TimeZone == 0 {
		req.TimeZone = 7.0
	}
	if req.HouseSystem == "" {
		req.HouseSystem = "W"
	}
}

// clampDayToMonth ป้องกันกรณีวันที่เกิด (เช่น 29, 30, 31) เกินจำนวนวันสูงสุดของเดือนในปีจร
func clampDayToMonth(year, month, day int) int {
	lastDay := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		return lastDay
	}
	return day
}

// navagrahaOrder ลำดับดาวนวเคราะห์ 9 ดวงตามหลักโหราศาสตร์พระเวท (1-อาทิตย์ ถึง 9-เกตุ)
var navagrahaOrder = []string{
	"Sun",     // 1-อาทิตย์
	"Moon",    // 2-จันทร์
	"Mars",    // 3-อังคาร
	"Mercury", // 4-พุธ
	"Jupiter", // 5-พฤหัสบดี
	"Venus",   // 6-ศุกร์
	"Saturn",  // 7-เสาร์
	"Rahu",    // 8-ราหู
	"Ketu",    // 9-เกตุ
}

// decodeAndValidateTransitReq อ่านและตรวจสอบความถูกต้องของ Request สำหรับ Endpoint ดาวจร
func (h *Handler) decodeAndValidateTransitReq(w http.ResponseWriter, r *http.Request, requireMonth bool) (TransitReq, bool) {
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid HTTP method attempted", zap.String("method", r.Method))
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return TransitReq{}, false
	}

	var wrapper TransitWrapperReq
	if err := json.NewDecoder(r.Body).Decode(&wrapper); err != nil {
		h.logger.Error("Failed to decode transit request body", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload: "+err.Error())
		return TransitReq{}, false
	}

	req := wrapper.TransitReq
	if wrapper.CalculateReq != nil && wrapper.CalculateReq.Year != 0 {
		if req.TransitYear != 0 && wrapper.CalculateReq.TransitYear == 0 {
			wrapper.CalculateReq.TransitYear = req.TransitYear
		}
		if req.TransitMonth != 0 && wrapper.CalculateReq.TransitMonth == 0 {
			wrapper.CalculateReq.TransitMonth = req.TransitMonth
		}
		req = *wrapper.CalculateReq
	} else if wrapper.WrappedReq != nil && wrapper.WrappedReq.Year != 0 {
		req = *wrapper.WrappedReq
	}

	normalizeTransitReq(&req)
	var err error
	if requireMonth {
		err = req.ValidateWithMonth()
	} else {
		err = req.Validate()
	}
	if err != nil {
		h.logger.Warn("Invalid transit request parameters", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return TransitReq{}, false
	}

	return req, true
}

// computeNatalAscSignIdx คำนวณราศีลัคนาดวงกำเนิด (0-11) เพื่อกำหนดเป็นเรือนที่ 1 (Whole Sign)
func computeNatalAscSignIdx(req TransitReq, signs []string) int {
	natalJD, natalAyanamsa := calculator.ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)
	houseRes := calculator.CalculateHouses(natalJD, req.Latitude, req.Longitude, int('W'), natalAyanamsa, signs)
	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	ascSignIdx, _ := calculator.ParseLongitude(ascDeg)
	return ascSignIdx
}

// computeSingleMonthTransit คำนวณดาวจร 9 ดวงในเดือนที่ระบุ เทียบเรือนภพกับลัคนาดวงกำเนิด
func computeSingleMonthTransit(req TransitReq, m int, ascSignIdx int, flag int, signs []string) ApiMonthlyTransit {
	tDay := clampDayToMonth(req.TransitYear, m, req.Day)
	tJD, _ := calculator.ComputeJulianDayAndAyanamsa(
		req.TransitYear, m, tDay, req.HourThai, req.TimeZone,
	)

	tPlanetRes := calculator.CalculatePlanets(tJD, flag, signs)

	var sunLon float64
	planetByCanonical := make(map[string]calculator.PlanetResult, len(tPlanetRes))
	for _, p := range tPlanetRes {
		_, _, canonical := calculator.TranslateBodyInfo(p.Name)
		planetByCanonical[canonical] = p
		if canonical == "Sun" {
			sunLon = p.Longitude
		}
	}

	transitPlanets := make([]ApiTransitPlanet, 0, len(navagrahaOrder))
	for idx, canonicalName := range navagrahaOrder {
		p, exists := planetByCanonical[canonicalName]
		if !exists {
			continue
		}

		nak, pada, _ := calculator.GetNakshatra(p.Longitude)
		posDetail := calculator.ParsePositionDetail(p.Longitude)

		nameThai, symbol, canonical := calculator.TranslateBodyInfo(p.Name)
		dignity, dignityThai := calculator.CalculateDignity(p.Name, p.Longitude)
		isCombust := calculator.IsPlanetCombust(p.Name, p.Longitude, sunLon, p.IsRetro)

		// นับเรือนเทียบกับลัคนาดวงกำเนิด: เรือน = ((ราศีดาวจร - ราศีลัคนา) mod 12) + 1
		pSignIdx, _ := calculator.ParseLongitude(p.Longitude)
		houseNum := ((pSignIdx - ascSignIdx + 12) % 12) + 1
		houseNameThai := calculator.GetHouseNameThai(houseNum)

		transitPlanets = append(transitPlanets, ApiTransitPlanet{
			No:            idx + 1,
			Name:          canonical,
			NameThai:      nameThai,
			Symbol:        symbol,
			Longitude:     p.Longitude,
			Sign:          posDetail.Sign,
			SignThai:      posDetail.SignThai,
			Degree:        calculator.FormatDegreeShort(posDetail),
			House:         houseNum,
			HouseNameThai: houseNameThai,
			IsRetrograde:  p.IsRetro,
			IsCombust:     isCombust,
			Dignity:       dignity,
			DignityThai:   dignityThai,
			Nakshatra:     nak.Name,
			NakshatraThai: nak.NameThai,
			Pada:          pada,
		})
	}

	return ApiMonthlyTransit{
		Month:   m,
		Date:    fmt.Sprintf("%02d/%02d/%04d", tDay, m, req.TransitYear),
		Planets: transitPlanets,
	}
}

// computeAnnualTransit คำนวณลัคนาดวงกำเนิด (เรือนที่ 1) และคืนค่าเฉพาะดาวจรรายเดือนทั้ง 12 เดือนในปีจร
func computeAnnualTransit(req TransitReq) []ApiMonthlyTransit {
	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	signs := calculator.GetZodiacSigns()
	ascSignIdx := computeNatalAscSignIdx(req, signs)

	monthlyTransits := make([]ApiMonthlyTransit, 0, 12)
	for m := 1; m <= 12; m++ {
		monthlyTransits = append(monthlyTransits, computeSingleMonthTransit(req, m, ascSignIdx, flag, signs))
	}

	return monthlyTransits
}

// computeMonthTransit คำนวณดาวจรเฉพาะเดือนเดียวที่ระบุใน transit_month
func computeMonthTransit(req TransitReq) ApiMonthlyTransit {
	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	signs := calculator.GetZodiacSigns()
	ascSignIdx := computeNatalAscSignIdx(req, signs)
	return computeSingleMonthTransit(req, req.TransitMonth, ascSignIdx, flag, signs)
}

// HandleCalculateTransit ประมวลผลคำขอคำนวณดาวจรประจำปี 12 เดือน (POST /api/v1/hora/transit)
func (h *Handler) HandleCalculateTransit(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeAndValidateTransitReq(w, r, false)
	if !ok {
		return
	}

	h.logger.Info("Received annual transit calculation request",
		zap.Int("birth_year", req.Year),
		zap.Int("birth_month", req.Month),
		zap.Int("birth_day", req.Day),
		zap.Int("transit_year", req.TransitYear),
		zap.Float64("latitude", req.Latitude),
		zap.Float64("longitude", req.Longitude),
	)

	monthlyTransits := computeAnnualTransit(req)

	resp := TransitWrapperResp{
		Status:          "success",
		MonthlyTransits: monthlyTransits,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode transit response JSON", zap.Error(err))
	}
}

// HandleCalculateTransitMonth ประมวลผลคำขอคำนวณดาวจรเฉพาะเดือนที่ระบุ (POST /api/v1/hora/transitmonth)
func (h *Handler) HandleCalculateTransitMonth(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeAndValidateTransitReq(w, r, true)
	if !ok {
		return
	}

	h.logger.Info("Received monthly transit calculation request",
		zap.Int("birth_year", req.Year),
		zap.Int("birth_month", req.Month),
		zap.Int("birth_day", req.Day),
		zap.Int("transit_year", req.TransitYear),
		zap.Int("transit_month", req.TransitMonth),
		zap.Float64("latitude", req.Latitude),
		zap.Float64("longitude", req.Longitude),
	)

	monthlyTransit := computeMonthTransit(req)

	resp := TransitMonthWrapperResp{
		Status:         "success",
		MonthlyTransit: monthlyTransit,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode monthly transit response JSON", zap.Error(err))
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
