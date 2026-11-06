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
