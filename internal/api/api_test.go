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
		t.Errorf("Expected planets data, got empty slice")
	}

	if len(wrapperResp.CalculateResp.HousesAndAngles) == 0 {
		t.Errorf("Expected houses and angles data, got empty slice")
	}

	if len(wrapperResp.CalculateResp.WholeSignHouses) != 12 {
		t.Errorf("Expected exactly 12 whole sign houses, got %d", len(wrapperResp.CalculateResp.WholeSignHouses))
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/hora/calculate", nil)
	getRec := httptest.NewRecorder()
	HandleCalculateHora(getRec, getReq)

	if getRec.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for GET method, got %d", http.StatusMethodNotAllowed, getRec.Result().StatusCode)
	}
}