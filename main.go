package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"lahiri-api/internal/api"

	"go.uber.org/zap"
)

func main() {
	logger := InitLogger()
	defer SyncLogger()

	api.SetLogger(logger)

	logger.Info("Initializing Vedic Astrology API Service")

	handler := api.NewHandler(logger, RunHoroscopedebug)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/hora/calculate", handler.HandleCalculateHora)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	srv := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// รัน HTTP Server ใน Goroutine และรองรับ Graceful Shutdown
	go func() {
		logger.Info("Starting HTTP Server", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server exited properly")
}

func RunHoroscopedebug(resp api.CalculateResp) {
	planetLogLines := make([]string, 0, len(resp.Planets))
	for _, p := range resp.Planets {
		retroStr := ""
		if p.IsRetrograde {
			retroStr = " (R)"
		}
		line := fmt.Sprintf("%-20s: %6.2f° (%s)%s | Nakshatra: %s (Pada %d) | Lord: %s",
			p.Name, p.Longitude, p.PositionStr, retroStr, p.Nakshatra, p.Pada, p.Lord)
		planetLogLines = append(planetLogLines, line)
	}

	houseLogLines := make([]string, 0, len(resp.HousesAndAngles))
	for _, h := range resp.HousesAndAngles {
		line := fmt.Sprintf("%-20s: %6.2f° (%s) | Nakshatra: %s (Pada %d) | Lord: %s",
			h.Name, h.Degree, h.PositionStr, h.Nakshatra, h.Pada, h.Lord)
		houseLogLines = append(houseLogLines, line)
	}

	wholeSignLogLines := make([]string, 0, len(resp.WholeSignHouses))
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
