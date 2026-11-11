package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lahiri-api/internal/calculator"
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

	if !bytes.Contains([]byte(res), []byte("Birth:1988-12-23 10:21")) {
		t.Errorf("Expected Birth date and time in prompt output")
	}
	if !bytes.Contains([]byte(res), []byte("Drishti(Special/Mutual):")) {
		t.Errorf("Expected Drishti in prompt output")
	}
	if !bytes.Contains([]byte(res), []byte("Conj:")) {
		t.Errorf("Expected Conj in prompt output")
	}
	if !bytes.Contains([]byte(res), []byte("/PD:")) {
		t.Errorf("Expected PD in prompt output")
	}

	// 2. ทดสอบแบบ JSON
	jsonPromptReq := httptest.NewRequest(http.MethodPost, "/api/v1/hora/prompt", bytes.NewBuffer(jsonBytes))
	jsonPromptReq.Header.Set("Content-Type", "application/json")
	jsonPromptRec := httptest.NewRecorder()
	HandleCalculatePrompt(jsonPromptRec, jsonPromptReq)

	var jsonResp PromptWrapperResp
	if err := json.NewDecoder(jsonPromptRec.Body).Decode(&jsonResp); err != nil {
		t.Fatalf("Failed to decode prompt JSON response: %v", err)
	}
	if jsonResp.BirthDate != "1988-12-23" {
		t.Errorf("Expected BirthDate '1988-12-23', got '%s'", jsonResp.BirthDate)
	}
	if jsonResp.BirthTime != "10:21" {
		t.Errorf("Expected BirthTime '10:21', got '%s'", jsonResp.BirthTime)
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
	if wrapperResp.BirthDate != "1988-12-23" {
		t.Errorf("Expected BirthDate '1988-12-23', got '%s'", wrapperResp.BirthDate)
	}
	if wrapperResp.BirthTime != "10:21" {
		t.Errorf("Expected BirthTime '10:21', got '%s'", wrapperResp.BirthTime)
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
		"Birth:1988-12-23 10:21",
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

func TestPromptDMSHelpers(t *testing.T) {
	if got := formatDMSInt(calculator.PositionDetail{Degree: 10, Minute: 59, Second: 59.7}); got != `11°00'00"` {
		t.Errorf("formatDMSInt carry mismatch: got %q", got)
	}
	if got := formatDMSDiff(5.9999); got != `6°00'` {
		t.Errorf("formatDMSDiff carry mismatch: got %q", got)
	}
	if got := formatDMSDiff(3.5); got != `3°30'` {
		t.Errorf("formatDMSDiff(3.5) mismatch: got %q", got)
	}
}

func TestConjunctionsAndDrishtiEdgeCases(t *testing.T) {
	t.Run("Conjunctions3PlusWithAndWithoutSun", func(t *testing.T) {
		bodies := []bodyDrishtiItem{
			{Code: "Me", House: 1, Lon: 10.0},
			{Code: "Su", House: 1, Lon: 7.0},
			{Code: "Sa", House: 1, Lon: 17.0},
			{Code: "Mo", House: 5, Lon: 125.0},
			{Code: "Ma", House: 5, Lon: 128.0},
			{Code: "Ju", House: 5, Lon: 135.0},
		}
		conj := formatConjunctions(bodies)
		if !strings.Contains(conj, "H1(Su+Me[orb 3°00',DeepCombust], Su+Sa[orb 10°00'])") {
			t.Errorf("unexpected H1 Sun-centered conjunction string: %q", conj)
		}
		if !strings.Contains(conj, "H5(Mo+Ma[orb 3°00'], Mo+Ju[orb 10°00'])") {
			t.Errorf("unexpected H5 conjunction string: %q", conj)
		}
	})

	t.Run("ConjunctionsAndDrishtiEmpty", func(t *testing.T) {
		bodies := []bodyDrishtiItem{
			{Code: "ASC", House: 1, Lon: 5.0},
			{Code: "Su", House: 1, Lon: 10.0},
			{Code: "Mo", House: 2, Lon: 40.0},
		}
		if got := formatConjunctions(bodies); got != "H1(ASC+Su[orb 5°00',DeepCombust])" {
			t.Errorf("unexpected single-house conjunction: %q", got)
		}
		if got := formatDrishti(bodies); got != "" {
			t.Errorf("expected empty drishti without special aspects or mutual opposite planets, got %q", got)
		}
	})

	t.Run("PromptHeaderVariants", func(t *testing.T) {
		respNoBirth := CalculateResp{}
		if out := FormatRomanCompactPrompt(respNoBirth); !strings.HasPrefix(out, "[Vedic Sidereal(Lahiri)/WholeSign]\n") {
			t.Errorf("unexpected header without BirthDate: %q", out)
		}

		respDateOnly := CalculateResp{
			Metadata: ApiMetadata{BirthDate: "1988-12-23", TimeZone: 0},
		}
		if out := FormatRomanCompactPrompt(respDateOnly); !strings.HasPrefix(out, "[Vedic Sidereal(Lahiri)/WholeSign | Birth:1988-12-23]\n") {
			t.Errorf("unexpected header with BirthDate only: %q", out)
		}
	})
}
