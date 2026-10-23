package main

import (
	"fmt"
	"github.com/mshafiee/swephgo"
)

func printVersion() {
	sweVer := make([]byte, 12)
	swephgo.Version(sweVer)
	fmt.Printf("Swiss Ephemeris Version: %s\n", string(sweVer))
}

func printLocationInfo(cfg Config) {
	fmt.Println("Birth Location & Time Configuration")
	fmt.Println("----------------------------------------------------------------------------------")
	fmt.Printf("Target Date : %02d/%02d/%d\n", cfg.Day, cfg.Month, cfg.Year)
	fmt.Printf("Target Time : %.2f Thai Time (TimeZone: GMT+%.1f)\n", cfg.HourThai, cfg.TimeZone)
	fmt.Printf("Latitude    : %.4f°\n", cfg.Latitude)
	fmt.Printf("Longitude   : %.4f°\n\n", cfg.Longitude)
}

func setupEphemeris(cfg Config) (float64, float64) {
	hourUTC := cfg.HourThai - cfg.TimeZone
	jd := swephgo.Julday(cfg.Year, cfg.Month, cfg.Day, hourUTC, swephgo.SeGregCal)
	fmt.Printf("Julian Day (UTC): %.6f\n\n", jd)

	swephgo.SetSidMode(swephgo.SeSidmLahiri, 0, 0)
	ayanamsa := swephgo.GetAyanamsa(jd)
	fmt.Printf("Ayanamsa (Lahiri): %s\n\n", FormatDMS(ayanamsa))

	return jd, ayanamsa
}

func getPlanetList() []PlanetConfig {
	return []PlanetConfig{
		{swephgo.SeSun, "Sun"},
		{swephgo.SeMoon, "Moon"},
		{swephgo.SeMercury, "Mercury"},
		{swephgo.SeVenus, "Venus"},
		{swephgo.SeMars, "Mars"},
		{swephgo.SeJupiter, "Jupiter"},
		{swephgo.SeSaturn, "Saturn"},
		{swephgo.SeUranus, "Uranus"},
		{swephgo.SeNeptune, "Neptune"},
		{swephgo.SePluto, "Pluto"},
		{swephgo.SeMeanNode, "Rahu (Mean Node)"},
	}
}

func printWholeSignHouses(asc float64, signs []string) {
	signIdx, _ := ParseLongitude(asc)

	fmt.Println("\nWhole Sign Houses (ภพละ 1 ราศีเต็ม):")
	for i := 1; i <= 12; i++ {
		currentSign := (signIdx + (i - 1)) % 12
		fmt.Printf("House %2d            : ราศี%s (0° - 30°)\n", i, signs[currentSign])
	}
}

func cleanup() {
	swephgo.Close()
	fmt.Println("\nDone!")
}