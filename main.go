package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func main() {
	InitLogger()
	defer SyncLogger()

	Log.Info("Initializing Vedic Astrology API Service")

	http.HandleFunc("/api/v1/hora/calculate", func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseRecorder{
			ResponseWriter: w,
			body:           bytes.NewBuffer(nil),
			statusCode:     http.StatusOK,
		}

		HandleCalculateHora(recorder, r)

		rawJSON := recorder.body.Bytes()
		if len(rawJSON) > 0 {
			// สร้าง Struct ครอบชั้นนอกตามที่ API พ่นออกมาจริง
			var apiWrapper struct {
				Status        string        `json:"status"`
				CalculateResp CalculateResp `json:"calculateResp"`
			}

			if err := json.Unmarshal(rawJSON, &apiWrapper); err == nil {
				// ส่ง CalculateResp ที่แกะออกมาแล้วเข้าฟังก์ชัน Debug Log
				RunHoroscopedebug(apiWrapper.CalculateResp)
			} else {
				Log.Error("JSON Unmarshal error", zap.Error(err))
			}
		}
	})

	port := ":8080"
	Log.Info("Starting HTTP Server", zap.String("port", port))
	if err := http.ListenAndServe(port, nil); err != nil {
		Log.Fatal("Server failed to start", zap.Error(err))
	}
}

func RunHoroscopedebug(resp CalculateResp) {
	var planetLogLines []string
	for _, p := range resp.Planets {
		retroStr := ""
		if p.IsRetrograde {
			retroStr = " (R)"
		}
		line := fmt.Sprintf("%-20s: %6.2f° (%s)%s | Nakshatra: %s (Pada %d) | Lord: %s",
			p.Name, p.Longitude, p.PositionStr, retroStr, p.Nakshatra, p.Pada, p.Lord)
		planetLogLines = append(planetLogLines, line)
	}

	var houseLogLines []string
	for _, h := range resp.HousesAndAngles {
		line := fmt.Sprintf("%-20s: %6.2f° (%s) | Nakshatra: %s (Pada %d) | Lord: %s",
			h.Name, h.Degree, h.PositionStr, h.Nakshatra, h.Pada, h.Lord)
		houseLogLines = append(houseLogLines, line)
	}

	var wholeSignLogLines []string
	for _, ws := range resp.WholeSignHouses {
		line := fmt.Sprintf("House %-2d            : Sign %s (0° - 30°)", ws.House, ws.Sign)
		wholeSignLogLines = append(wholeSignLogLines, line)
	}

	var sb strings.Builder
	sb.WriteString("\nSwiss Ephemeris Version: 2.10.03\n")
	sb.WriteString(fmt.Sprintf("Julian Day (UTC): %.6f\n\n", resp.Metadata.JulianDay))
	sb.WriteString(fmt.Sprintf("Ayanamsa (Lahiri): %s\n\n", resp.Metadata.AyanamsaDegrees))
	sb.WriteString("Planetary Positions (Sidereal - Lahiri Zodiac)\n")
	sb.WriteString("------------------------------------------------------------------------------------\n")
	sb.WriteString(strings.Join(planetLogLines, "\n") + "\n\n")
	sb.WriteString("Ascendant & Angles (Sidereal - Lahiri Zodiac)\n")
	sb.WriteString("------------------------------------------------------------------------------------\n")
	sb.WriteString(strings.Join(houseLogLines, "\n") + "\n\n")
	sb.WriteString("Whole Sign Houses (1 Sign per House):\n")
	sb.WriteString(strings.Join(wholeSignLogLines, "\n"))

	if Log != nil {
		Log.Info(sb.String())
	}

}
