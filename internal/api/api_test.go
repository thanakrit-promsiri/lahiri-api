package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestHandleCalculateHora(t *testing.T) {
	reqBody := CalculateWrapperReq{
		CalculateReq: CalculateReq{
			Year:        2026,
			Month:       10,
			Day:         23,
			HourThai:    12.0,
			TimeZone:    7.0,
			Latitude:    13.7563,
			Longitude:   100.5018,
			HouseSystem: "W",
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/calculate", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleCalculateHora(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var wrapperResp CalculateWrapperResp
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&wrapperResp); err != nil {
		t.Fatalf("Failed to decode response JSON: %v", err)
	}

	if wrapperResp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", wrapperResp.Status)
	}

	if wrapperResp.CalculateResp.Metadata.AyanamsaName != "Lahiri" {
		t.Errorf("Expected AyanamsaName 'Lahiri', got '%s'", wrapperResp.CalculateResp.Metadata.AyanamsaName)
	}

	if len(wrapperResp.CalculateResp.Planets) == 0 {
		t.Fatalf("Expected planets data, got empty slice")
	}

	// ตรวจสอบฟิลด์ใหม่ใน Planets (คำแปลไทย, ภพ, มาตรฐานดาว)
	firstPlanet := wrapperResp.CalculateResp.Planets[0]
	if firstPlanet.NameThai == "" || firstPlanet.Symbol == "" {
		t.Errorf("Expected Thai name and symbol for planet, got '%s' / '%s'", firstPlanet.NameThai, firstPlanet.Symbol)
	}
	if firstPlanet.Dignity == "" || firstPlanet.DignityThai == "" {
		t.Errorf("Expected dignity fields for planet, got '%s' / '%s'", firstPlanet.Dignity, firstPlanet.DignityThai)
	}
	if firstPlanet.House < 1 || firstPlanet.House > 12 || firstPlanet.HouseNameThai == "" {
		t.Errorf("Expected valid house (1-12) and Thai house name, got %d / '%s'", firstPlanet.House, firstPlanet.HouseNameThai)
	}
	if firstPlanet.Position.SignThai == "" {
		t.Errorf("Expected Position.SignThai to be populated")
	}

	if len(wrapperResp.CalculateResp.HousesAndAngles) == 0 {
		t.Errorf("Expected houses and angles data, got empty slice")
	}

	if len(wrapperResp.CalculateResp.WholeSignHouses) != 12 {
		t.Fatalf("Expected exactly 12 whole sign houses, got %d", len(wrapperResp.CalculateResp.WholeSignHouses))
	}

	// ตรวจสอบฟิลด์เจ้าเรือนภพ (House Lords)
	firstHouse := wrapperResp.CalculateResp.WholeSignHouses[0]
	if firstHouse.HouseNameThai != "ตนุ" {
		t.Errorf("Expected first house Thai name 'ตนุ', got '%s'", firstHouse.HouseNameThai)
	}
	if firstHouse.Lord == "" || firstHouse.LordThai == "" {
		t.Errorf("Expected house lord and Thai lord name, got '%s' / '%s'", firstHouse.Lord, firstHouse.LordThai)
	}
	if firstHouse.LordPlacedInHouse < 1 || firstHouse.LordPlacedInHouse > 12 || firstHouse.LordPlacedInHouseThai == "" {
		t.Errorf("Expected valid lord_placed_in_house (1-12), got %d (%s)", firstHouse.LordPlacedInHouse, firstHouse.LordPlacedInHouseThai)
	}

	// ตรวจสอบฟิลด์วิมโสตรีทศา (Vimshottari Dasha)
	dasha := wrapperResp.CalculateResp.VimshottariDasha
	if dasha.MoonNakshatra == "" || dasha.MoonNakshatraThai == "" {
		t.Errorf("Expected Moon Nakshatra in VimshottariDasha, got '%s' / '%s'", dasha.MoonNakshatra, dasha.MoonNakshatraThai)
	}
	if len(dasha.MahadashaList) != 9 {
		t.Errorf("Expected 9 Mahadasha items, got %d", len(dasha.MahadashaList))
	}

	// ทดสอบ Endpoint /api/v1/hora/dasha
	dashaReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/dasha", bytes.NewBuffer(jsonBytes))
	dashaReq.Header.Set("Content-Type", "application/json")
	dashaRec := httptest.NewRecorder()
	HandleCalculateDasha(dashaRec, dashaReq)

	if dashaRec.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d from /api/v1/hora/dasha, got %d", http.StatusOK, dashaRec.Result().StatusCode)
	}

	// ทดสอบ Endpoint /api/v1/hora/prompt (JSON format)
	promptReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/prompt", bytes.NewBuffer(jsonBytes))
	promptReq.Header.Set("Content-Type", "application/json")
	promptRec := httptest.NewRecorder()
	HandleCalculatePrompt(promptRec, promptReq)

	if promptRec.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d from /api/v1/hora/prompt, got %d", http.StatusOK, promptRec.Result().StatusCode)
	}
	var promptResp PromptWrapperResp
	if err := json.NewDecoder(promptRec.Result().Body).Decode(&promptResp); err != nil {
		t.Fatalf("Failed to decode prompt response JSON: %v", err)
	}
	if promptResp.Status != "success" || promptResp.Prompt == "" {
		t.Errorf("Expected non-empty Roman prompt, got %+v", promptResp)
	}

	// ทดสอบ Endpoint /api/v1/hora/prompt?format=raw (Plain Text format)
	rawPromptReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/prompt?format=raw", bytes.NewBuffer(jsonBytes))
	rawPromptReq.Header.Set("Content-Type", "application/json")
	rawPromptRec := httptest.NewRecorder()
	HandleCalculatePrompt(rawPromptRec, rawPromptReq)

	if rawPromptRec.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d from /api/v1/hora/prompt?format=raw, got %d", http.StatusOK, rawPromptRec.Result().StatusCode)
	}
	if rawPromptRec.Body.Len() == 0 {
		t.Errorf("Expected non-empty plain text prompt")
	}

	// ทดสอบ Method Not Allowed
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/hora/calculate", nil)
	getRec := httptest.NewRecorder()
	HandleCalculateHora(getRec, getReq)

	if getRec.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for GET method, got %d", http.StatusMethodNotAllowed, getRec.Result().StatusCode)
	}

	// ทดสอบ Invalid JSON
	badJSONReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/calculate", bytes.NewBufferString("{invalid-json"))
	badJSONRec := httptest.NewRecorder()
	HandleCalculateHora(badJSONRec, badJSONReq)

	if badJSONRec.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid JSON, got %d", http.StatusBadRequest, badJSONRec.Result().StatusCode)
	}

	// ทดสอบ Input Validation (เช่น เดือนไม่ถูกต้อง)
	invalidReqBody := CalculateWrapperReq{
		CalculateReq: CalculateReq{
			Year:     2026,
			Month:    13, // เกิน 12
			Day:      23,
			HourThai: 12.0,
			TimeZone: 7.0,
		},
	}
	invalidBytes, _ := json.Marshal(invalidReqBody)
	invalidReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/calculate", bytes.NewBuffer(invalidBytes))
	invalidRec := httptest.NewRecorder()
	HandleCalculateHora(invalidRec, invalidReq)

	if invalidRec.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid input parameters, got %d", http.StatusBadRequest, invalidRec.Result().StatusCode)
	}
}

func TestHandleCalculateTransit(t *testing.T) {
	rawInput := `{
  "calculateReq": {
    "year": 1988,
    "month": 12,
    "day": 23,
    "hour_thai": 10.35,
    "time_zone": 7.0,
    "latitude": 13.7563,
    "longitude": 100.5018,
    "house_system": "W",
    "transit_year": 2026
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transit", bytes.NewBufferString(rawInput))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleCalculateTransit(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var wrapperResp TransitWrapperResp
	if err := json.NewDecoder(res.Body).Decode(&wrapperResp); err != nil {
		t.Fatalf("Failed to decode transit response JSON: %v", err)
	}

	if wrapperResp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", wrapperResp.Status)
	}

	if len(wrapperResp.MonthlyTransits) != 12 {
		t.Fatalf("Expected 12 monthly transits, got %d", len(wrapperResp.MonthlyTransits))
	}

	if wrapperResp.MonthlyTransits[0].Date != "23/01/2026" {
		t.Errorf("Expected first month date '23/01/2026', got '%s'", wrapperResp.MonthlyTransits[0].Date)
	}
	if wrapperResp.MonthlyTransits[11].Date != "23/12/2026" {
		t.Errorf("Expected last month date '23/12/2026', got '%s'", wrapperResp.MonthlyTransits[11].Date)
	}

	expectedNames := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Rahu", "Ketu"}
	for _, mt := range wrapperResp.MonthlyTransits {
		if len(mt.Planets) != 9 {
			t.Fatalf("Month %d: expected 9 Navagraha planets, got %d", mt.Month, len(mt.Planets))
		}
		for idx, p := range mt.Planets {
			if p.No != idx+1 || p.Name != expectedNames[idx] {
				t.Errorf("Month %d planet %d: expected (%d, %s), got (%d, %s)",
					mt.Month, idx, idx+1, expectedNames[idx], p.No, p.Name)
			}
			if p.House < 1 || p.House > 12 {
				t.Errorf("Month %d planet %s: invalid house %d", mt.Month, p.Name, p.House)
			}
		}
	}
}

func TestHandleCalculateTransitMonth(t *testing.T) {
	rawInput := `{
  "calculateReq": {
    "year": 1988,
    "month": 12,
    "day": 23,
    "hour_thai": 10.35,
    "time_zone": 7.0,
    "latitude": 13.7563,
    "longitude": 100.5018,
    "house_system": "W",
    "transit_year": 2026,
    "transit_month": 12
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitmonth", bytes.NewBufferString(rawInput))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleCalculateTransitMonth(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var wrapperResp TransitMonthWrapperResp
	if err := json.NewDecoder(res.Body).Decode(&wrapperResp); err != nil {
		t.Fatalf("Failed to decode transitmonth response JSON: %v", err)
	}

	if wrapperResp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", wrapperResp.Status)
	}

	mt := wrapperResp.MonthlyTransit
	if mt.Month != 12 || mt.Date != "23/12/2026" {
		t.Errorf("Expected month 12 date '23/12/2026', got month %d date '%s'", mt.Month, mt.Date)
	}

	expectedNames := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Rahu", "Ketu"}
	if len(mt.Planets) != 9 {
		t.Fatalf("Expected 9 Navagraha planets, got %d", len(mt.Planets))
	}
	for idx, p := range mt.Planets {
		if p.No != idx+1 || p.Name != expectedNames[idx] {
			t.Errorf("Planet %d: expected (%d, %s), got (%d, %s)",
				idx, idx+1, expectedNames[idx], p.No, p.Name)
		}
		if p.House < 1 || p.House > 12 {
			t.Errorf("Planet %s: invalid house %d", p.Name, p.House)
		}
	}
}

func TestHandleCalculateDashaWithTransitDate(t *testing.T) {
	rawInput := `{
  "calculateReq": {
    "year": 1988,
    "month": 12,
    "day": 23,
    "hour_thai": 10.35,
    "time_zone": 7.0,
    "latitude": 13.7563,
    "longitude": 100.5018,
    "house_system": "W",
    "transit_year": 2026,
    "transit_month": 12
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/dasha", bytes.NewBufferString(rawInput))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleCalculateDasha(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var wrapperResp DashaWrapperResp
	if err := json.NewDecoder(res.Body).Decode(&wrapperResp); err != nil {
		t.Fatalf("Failed to decode dasha response JSON: %v", err)
	}

	if wrapperResp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", wrapperResp.Status)
	}

	d := wrapperResp.VimshottariDasha
	if d.TargetDate != "23/12/2026" {
		t.Errorf("Expected target_date '23/12/2026', got '%s'", d.TargetDate)
	}
	if d.CurrentDasha.MahadashaLord == "" || d.CurrentDasha.AntardashaLord == "" || d.CurrentDasha.PratyantardashaLord == "" {
		t.Errorf("Expected Mahadasha, Antardasha, and Pratyantardasha lords to be populated, got %+v", d.CurrentDasha)
	}
}

func TestEndOfMonthTransitDates(t *testing.T) {
	// ทดสอบกรณีคนเกิดวันที่ 31 (สิ้นเดือน) ในปีปกติ (2026) และปีอธิกสุรทิน Leap Year (2028)
	rawInput2026 := `{
  "calculateReq": {
    "year": 1988,
    "month": 12,
    "day": 31,
    "hour_thai": 10.35,
    "time_zone": 7.0,
    "latitude": 13.7563,
    "longitude": 100.5018,
    "house_system": "W",
    "transit_year": 2026
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transit", bytes.NewBufferString(rawInput2026))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	HandleCalculateTransit(rec, req)

	var resp2026 TransitWrapperResp
	if err := json.NewDecoder(rec.Result().Body).Decode(&resp2026); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedDates2026 := []string{
		"31/01/2026", "28/02/2026", "31/03/2026", "30/04/2026",
		"31/05/2026", "30/06/2026", "31/07/2026", "31/08/2026",
		"30/09/2026", "31/10/2026", "30/11/2026", "31/12/2026",
	}
	for i, exp := range expectedDates2026 {
		if resp2026.MonthlyTransits[i].Date != exp {
			t.Errorf("Month %d: expected '%s', got '%s'", i+1, exp, resp2026.MonthlyTransits[i].Date)
		}
	}

	// ทดสอบปีอธิกสุรทิน (Leap Year 2028 -> ก.พ. ต้องได้ 29/02/2028)
	rawInput2028 := `{
  "calculateReq": {
    "year": 1988,
    "month": 12,
    "day": 31,
    "hour_thai": 10.35,
    "time_zone": 7.0,
    "latitude": 13.7563,
    "longitude": 100.5018,
    "house_system": "W",
    "transit_year": 2028,
    "transit_month": 2
  }
}`
	reqLeap := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitmonth", bytes.NewBufferString(rawInput2028))
	reqLeap.Header.Set("Content-Type", "application/json")
	recLeap := httptest.NewRecorder()
	HandleCalculateTransitMonth(recLeap, reqLeap)

	var respLeap TransitMonthWrapperResp
	if err := json.NewDecoder(recLeap.Result().Body).Decode(&respLeap); err != nil {
		t.Fatalf("Failed to decode leap year response: %v", err)
	}
	if respLeap.MonthlyTransit.Date != "29/02/2028" {
		t.Errorf("Expected leap year Feb date '29/02/2028', got '%s'", respLeap.MonthlyTransit.Date)
	}
}

func TestHandleCalculateTransitDasha(t *testing.T) {
	rawInput := `{
  "calculateReq": {
    "year": 1988,
    "month": 12,
    "day": 23,
    "hour_thai": 10.35,
    "time_zone": 7.0,
    "latitude": 13.7563,
    "longitude": 100.5018,
    "house_system": "W",
    "transit_year": 2026,
    "transit_month": 12
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitdasha", bytes.NewBufferString(rawInput))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleCalculateTransitDasha(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var wrapperResp TransitMonthWrapperResp
	if err := json.NewDecoder(res.Body).Decode(&wrapperResp); err != nil {
		t.Fatalf("Failed to decode transitdasha response JSON: %v", err)
	}

	if wrapperResp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", wrapperResp.Status)
	}

	mt := wrapperResp.MonthlyTransit
	if mt.Month != 12 || mt.Date != "23/12/2026" {
		t.Errorf("Expected month 12 date '23/12/2026', got month %d date '%s'", mt.Month, mt.Date)
	}

	if mt.Dasha == nil {
		t.Fatalf("Expected dasha object in monthly_transit, got nil")
	}
	if mt.Dasha.MahadashaLord != "Saturn" || mt.Dasha.AntardashaLord != "Mercury" || mt.Dasha.PratyantardashaLord != "Sun" {
		t.Errorf("Expected MD=Saturn, AD=Mercury, PD=Sun at 23/12/2026, got %+v", mt.Dasha)
	}
	if len(mt.Planets) != 9 {
		t.Fatalf("Expected 9 Navagraha planets, got %d", len(mt.Planets))
	}
}

func TestSetLoggerAndHandlerCallback(t *testing.T) {
	prevLog := Log
	t.Cleanup(func() {
		SetLogger(prevLog)
	})

	logger := zap.NewNop()
	SetLogger(logger)
	if Log != logger {
		t.Fatalf("expected global Log to be set")
	}

	nilHandler := NewHandler(nil, nil)
	if nilHandler.logger == nil {
		t.Fatalf("expected NewHandler(nil, nil) to fallback to nop logger")
	}

	var callbackCalled bool
	var gotBirthDate string
	h := NewHandler(logger, func(resp CalculateResp) {
		callbackCalled = true
		gotBirthDate = resp.Metadata.BirthDate
	})

	body := `{"calculateReq":{"year":2531,"month":12,"day":23,"hour_thai":10.35,"time_zone":7,"latitude":13.7563,"longitude":100.5018}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/calculate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.HandleCalculateHora(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Result().StatusCode)
	}
	if !callbackCalled || gotBirthDate != "1988-12-23" {
		t.Errorf("expected onResult callback with BirthDate 1988-12-23, got called=%v date=%q", callbackCalled, gotBirthDate)
	}
}

func TestNormalizeRequests(t *testing.T) {
	calcReq := CalculateReq{Year: 2531, HouseSystem: ""}
	normalizeReq(&calcReq)
	if calcReq.Year != 1988 || calcReq.HouseSystem != "W" {
		t.Errorf("normalizeReq mismatch: %+v", calcReq)
	}

	tReq := TransitReq{Year: 2531, TransitYear: 2569, TimeZone: 0, HouseSystem: ""}
	normalizeTransitReq(&tReq)
	if tReq.Year != 1988 || tReq.TransitYear != 2026 || tReq.TimeZone != 7.0 || tReq.HouseSystem != "W" {
		t.Errorf("normalizeTransitReq mismatch: %+v", tReq)
	}

	defaultYearReq := TransitReq{Year: 1988, TransitYear: 0}
	normalizeTransitReq(&defaultYearReq)
	if defaultYearReq.TransitYear != time.Now().Year() {
		t.Errorf("expected TransitYear %d, got %d", time.Now().Year(), defaultYearReq.TransitYear)
	}
}

func TestHandleCalculateDashaVariantsAndErrors(t *testing.T) {
	successPayloads := []struct {
		name string
		body string
	}{
		{
			name: "FlatJSON",
			body: `{"year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7,"latitude":13.7563,"longitude":100.5018,"transit_year":2026,"transit_month":12}`,
		},
		{
			name: "TransitReqWrapper",
			body: `{"transitReq":{"year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7,"latitude":13.7563,"longitude":100.5018,"transit_year":2026,"transit_month":12}}`,
		},
		{
			name: "CalculateReqWithOuterTransitFallback",
			body: `{"transit_year":2026,"transit_month":12,"calculateReq":{"year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7,"latitude":13.7563,"longitude":100.5018}}`,
		},
	}

	for _, tc := range successPayloads {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/dasha", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			HandleCalculateDasha(rec, req)

			if rec.Result().StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d: %s", rec.Result().StatusCode, rec.Body.String())
			}
			var resp DashaWrapperResp
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if resp.VimshottariDasha.TargetDate != "23/12/2026" {
				t.Errorf("expected TargetDate 23/12/2026, got %q", resp.VimshottariDasha.TargetDate)
			}
		})
	}

	errorCases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"MethodNotAllowed", http.MethodGet, "", http.StatusMethodNotAllowed},
		{"InvalidJSON", http.MethodPost, "{bad-json", http.StatusBadRequest},
		{"InvalidTransitMonth", http.MethodPost, `{"year":1988,"month":12,"day":23,"hour_thai":10.35,"transit_year":2026,"transit_month":13}`, http.StatusBadRequest},
		{"InvalidBirthMonth", http.MethodPost, `{"year":1988,"month":13,"day":23,"hour_thai":10.35,"transit_year":2026,"transit_month":12}`, http.StatusBadRequest},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/v1/hora/dasha", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			HandleCalculateDasha(rec, req)
			if rec.Result().StatusCode != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, rec.Result().StatusCode)
			}
		})
	}
}

func TestTransitAndPromptFormatAndErrorPaths(t *testing.T) {
	t.Run("PromptAcceptTextPlain", func(t *testing.T) {
		calcBody := `{"calculateReq":{"year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7,"latitude":13.7563,"longitude":100.5018}}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/prompt", strings.NewReader(calcBody))
		req.Header.Set("Accept", "text/plain")
		rec := httptest.NewRecorder()
		HandleCalculatePrompt(rec, req)
		if rec.Result().StatusCode != http.StatusOK || !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/plain") {
			t.Errorf("expected 200 text/plain, got %d (%s)", rec.Result().StatusCode, rec.Header().Get("Content-Type"))
		}

		transitBody := `{"year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7,"latitude":13.7563,"longitude":100.5018,"transit_year":2026,"transit_month":12}`
		tReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitdasha_prompt", strings.NewReader(transitBody))
		tReq.Header.Set("Accept", "text/plain")
		tRec := httptest.NewRecorder()
		HandleCalculateTransitDashaPrompt(tRec, tReq)
		if tRec.Result().StatusCode != http.StatusOK || !strings.HasPrefix(tRec.Header().Get("Content-Type"), "text/plain") {
			t.Errorf("expected 200 text/plain, got %d (%s)", tRec.Result().StatusCode, tRec.Header().Get("Content-Type"))
		}
	})

	t.Run("TransitWrapperAndOuterFallback", func(t *testing.T) {
		bodyWrapper := `{"transitReq":{"year":1988,"month":12,"day":23,"hour_thai":10.35,"latitude":13.7563,"longitude":100.5018,"transit_year":2026,"transit_month":5}}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitmonth", strings.NewReader(bodyWrapper))
		rec := httptest.NewRecorder()
		HandleCalculateTransitMonth(rec, req)
		if rec.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected 200 from transitReq wrapper, got %d", rec.Result().StatusCode)
		}

		bodyOuter := `{"transit_year":2026,"transit_month":5,"calculateReq":{"year":1988,"month":12,"day":23,"hour_thai":10.35,"latitude":13.7563,"longitude":100.5018}}`
		reqOuter := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitmonth", strings.NewReader(bodyOuter))
		recOuter := httptest.NewRecorder()
		HandleCalculateTransitMonth(recOuter, reqOuter)
		if recOuter.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected 200 from outer fallback, got %d", recOuter.Result().StatusCode)
		}
	})

	handlers := []struct {
		name         string
		path         string
		fn           http.HandlerFunc
		invalidParam string
	}{
		{"Prompt", "/api/v1/hora/prompt", HandleCalculatePrompt, `{"calculateReq":{"year":1988,"month":13,"day":23,"hour_thai":10.35}}`},
		{"Transit", "/api/v1/hora/transit", HandleCalculateTransit, `{"year":1988,"month":13,"day":23,"hour_thai":10.35,"transit_year":2026}`},
		{"TransitMonth", "/api/v1/hora/transitmonth", HandleCalculateTransitMonth, `{"year":1988,"month":12,"day":23,"hour_thai":10.35,"transit_year":2026,"transit_month":0}`},
		{"TransitDasha", "/api/v1/hora/transitdasha", HandleCalculateTransitDasha, `{"year":1988,"month":12,"day":23,"hour_thai":10.35,"transit_year":2026,"transit_month":13}`},
		{"TransitDashaPrompt", "/api/v1/hora/transitdasha_prompt", HandleCalculateTransitDashaPrompt, `{"year":1988,"month":12,"day":23,"hour_thai":10.35,"transit_year":2026,"transit_month":0}`},
	}

	for _, h := range handlers {
		t.Run(h.name+"_MethodNotAllowed", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, h.path, nil)
			rec := httptest.NewRecorder()
			h.fn(rec, req)
			if rec.Result().StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("%s GET: expected 405, got %d", h.name, rec.Result().StatusCode)
			}
		})
		t.Run(h.name+"_InvalidJSON", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, h.path, strings.NewReader("{bad"))
			rec := httptest.NewRecorder()
			h.fn(rec, req)
			if rec.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("%s invalid JSON: expected 400, got %d", h.name, rec.Result().StatusCode)
			}
		})
		t.Run(h.name+"_ValidationError", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, h.path, strings.NewReader(h.invalidParam))
			rec := httptest.NewRecorder()
			h.fn(rec, req)
			if rec.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("%s validation: expected 400, got %d", h.name, rec.Result().StatusCode)
			}
		})
	}
}
