package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPromptOutputFormat(t *testing.T) {
	reqBody := CalculateWrapperReq{
		CalculateReq: CalculateReq{
			Year:        1988,
			Month:       12,
			Day:         23,
			HourThai:    10.35,
			TimeZone:    7.0,
			Latitude:    13.7563,
			Longitude:   100.5018,
			HouseSystem: "W",
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	rawPromptReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/prompt?format=raw", bytes.NewBuffer(jsonBytes))
	rawPromptReq.Header.Set("Content-Type", "application/json")
	rawPromptRec := httptest.NewRecorder()
	HandleCalculatePrompt(rawPromptRec, rawPromptReq)

	res := rawPromptRec.Body.String()
	t.Logf("Generated Prompt:\n%s", res)

	if !bytes.Contains([]byte(res), []byte("Drishti(Special/Mutual):")) {
		t.Errorf("Expected Drishti in prompt output")
	}
	if !bytes.Contains([]byte(res), []byte("Conj:")) {
		t.Errorf("Expected Conj in prompt output")
	}
	if !bytes.Contains([]byte(res), []byte("/PD:")) {
		t.Errorf("Expected PD in prompt output")
	}
}

func TestTransitDashaPromptOutputFormat(t *testing.T) {
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

	// 1. ทดสอบแบบ JSON
	jsonReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitdasha_prompt", bytes.NewBufferString(rawInput))
	jsonReq.Header.Set("Content-Type", "application/json")
	jsonRec := httptest.NewRecorder()
	HandleCalculateTransitDashaPrompt(jsonRec, jsonReq)

	if jsonRec.Result().StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, jsonRec.Result().StatusCode)
	}
	var wrapperResp PromptWrapperResp
	if err := json.NewDecoder(jsonRec.Result().Body).Decode(&wrapperResp); err != nil {
		t.Fatalf("Failed to decode transitdasha_prompt JSON: %v", err)
	}
	if wrapperResp.Status != "success" || wrapperResp.Prompt == "" {
		t.Errorf("Expected non-empty prompt in JSON response, got %+v", wrapperResp)
	}

	// 2. ทดสอบแบบ ?format=raw (Plain Text)
	rawReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/transitdasha_prompt?format=raw", bytes.NewBufferString(rawInput))
	rawReq.Header.Set("Content-Type", "application/json")
	rawRec := httptest.NewRecorder()
	HandleCalculateTransitDashaPrompt(rawRec, rawReq)

	res := rawRec.Body.String()
	t.Logf("Generated TransitDasha Prompt:\n%s", res)

	expectedSubstrings := []string{
		"[Vedic Transit(Lahiri)/WholeSign | Date:2026-12-23]",
		"NatalRef: ASC:Aq(H1)",
		"Dasha@2026-12-23:",
		"TransitPlanets(D1:Sign,DMS,Dig|House|Nak(Pada,Lord)|D9:Sign,Dig):",
		"TransitDrishti(Special/Mutual):",
	}
	for _, sub := range expectedSubstrings {
		if !bytes.Contains([]byte(res), []byte(sub)) {
			t.Errorf("Expected '%s' in transitdasha_prompt output", sub)
		}
	}
}
