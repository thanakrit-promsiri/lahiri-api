package api

import (
	"encoding/json"
	"net/http"
	"time"

	"lahiri-api/internal/calculator"

	"go.uber.org/zap"
)

// Handler จัดการ HTTP requests (Controller Layer) พร้อม Dependency Injection ของ Logger
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

// HandleCalculateTransitDasha เป็นฟังก์ชันระดับ Package สำหรับคำนวณดาวจรในเดือนพร้อมทศา/อนุทศา ณ เวลาดาวจร
func HandleCalculateTransitDasha(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculateTransitDasha(w, r)
}

// HandleCalculateTransitDashaPrompt เป็นฟังก์ชันระดับ Package สำหรับสร้างข้อความ Roman Ultra-Compact ของดาวจร+ทศาจร
func HandleCalculateTransitDashaPrompt(w http.ResponseWriter, r *http.Request) {
	NewHandler(Log, nil).HandleCalculateTransitDashaPrompt(w, r)
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

// decodeAndValidateReq อ่านและตรวจสอบความถูกต้องของ Request พื้นฐาน
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

	calcResp := calculator.ComputeHoroscope(req)

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

	dasha := calculator.ComputeDasha(req)

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

	calcResp := calculator.ComputeHoroscope(req)
	promptText := FormatRomanCompactPrompt(calcResp)

	if r.URL.Query().Get("format") == "raw" || r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(promptText))
		return
	}

	birthDate, birthTime := calculator.FormatBirthDateTime(req.Year, req.Month, req.Day, req.HourThai)

	resp := PromptWrapperResp{
		Status:    "success",
		BirthDate: birthDate,
		BirthTime: birthTime,
		TimeZone:  req.TimeZone,
		Prompt:    promptText,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode prompt response JSON", zap.Error(err))
	}
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

	monthlyTransits := calculator.ComputeAnnualTransit(req)

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

	monthlyTransit := calculator.ComputeMonthTransit(req)

	resp := TransitMonthWrapperResp{
		Status:         "success",
		MonthlyTransit: monthlyTransit,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode monthly transit response JSON", zap.Error(err))
	}
}

// HandleCalculateTransitDasha ประมวลผลคำขอคำนวณดาวจรในเดือนพร้อมทศา/อนุทศา ณ เวลาดาวจร (POST /api/v1/hora/transitdasha)
func (h *Handler) HandleCalculateTransitDasha(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeAndValidateTransitReq(w, r, true)
	if !ok {
		return
	}

	h.logger.Info("Received monthly transit with dasha calculation request",
		zap.Int("birth_year", req.Year),
		zap.Int("birth_month", req.Month),
		zap.Int("birth_day", req.Day),
		zap.Int("transit_year", req.TransitYear),
		zap.Int("transit_month", req.TransitMonth),
		zap.Float64("latitude", req.Latitude),
		zap.Float64("longitude", req.Longitude),
	)

	monthlyTransit := calculator.ComputeMonthTransitWithDasha(req)

	resp := TransitMonthWrapperResp{
		Status:         "success",
		MonthlyTransit: monthlyTransit,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode monthly transit with dasha response JSON", zap.Error(err))
	}
}

// HandleCalculateTransitDashaPrompt ประมวลผลและคืนค่าข้อความดาวจรและทศา/อนุทศา ณ เวลาดาวจร แบบ Roman Ultra-Compact (POST /api/v1/hora/transitdasha_prompt)
func (h *Handler) HandleCalculateTransitDashaPrompt(w http.ResponseWriter, r *http.Request) {
	req, ok := h.decodeAndValidateTransitReq(w, r, true)
	if !ok {
		return
	}

	h.logger.Info("Received transit dasha AI prompt generation request",
		zap.Int("birth_year", req.Year),
		zap.Int("birth_month", req.Month),
		zap.Int("birth_day", req.Day),
		zap.Int("transit_year", req.TransitYear),
		zap.Int("transit_month", req.TransitMonth),
	)

	promptText := FormatTransitDashaCompactPrompt(req)

	if r.URL.Query().Get("format") == "raw" || r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(promptText))
		return
	}

	birthDate, birthTime := calculator.FormatBirthDateTime(req.Year, req.Month, req.Day, req.HourThai)

	resp := PromptWrapperResp{
		Status:    "success",
		BirthDate: birthDate,
		BirthTime: birthTime,
		TimeZone:  req.TimeZone,
		Prompt:    promptText,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("Failed to encode transit dasha prompt response JSON", zap.Error(err))
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
