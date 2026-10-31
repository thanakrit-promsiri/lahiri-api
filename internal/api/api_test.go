package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
