package calculator

import (
	"fmt"
	"sync"

	"github.com/mshafiee/swephgo"
)

// sweMu ป้องกันการเรียกใช้งาน C-library (Swiss Ephemeris) พร้อมกันจากหลาย Goroutine
var sweMu sync.Mutex

func printVersion() {
	sweMu.Lock()
	defer sweMu.Unlock()

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
	jd, ayanamsa := ComputeJulianDayAndAyanamsa(cfg.Year, cfg.Month, cfg.Day, cfg.HourThai, cfg.TimeZone)
	fmt.Printf("Julian Day (UTC): %.6f\n\n", jd)
	fmt.Printf("Ayanamsa (Lahiri): %s\n\n", FormatDMS(ayanamsa))

	return jd, ayanamsa
}

// ComputeJulianDayAndAyanamsa คำนวณค่า Julian Day (UTC) และค่า Ayanamsa (Lahiri) อย่างปลอดภัยต่อ Concurrency
func ComputeJulianDayAndAyanamsa(year, month, day int, hourLocal, timeZone float64) (float64, float64) {
	sweMu.Lock()
	defer sweMu.Unlock()

	hourUTC := hourLocal - timeZone
	jd := swephgo.Julday(year, month, day, hourUTC, swephgo.SeGregCal)
	swephgo.SetSidMode(swephgo.SeSidmLahiri, 0, 0)
	ayanamsa := swephgo.GetAyanamsa(jd)

	return jd, ayanamsa
}

var defaultPlanets = []PlanetConfig{
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

func getPlanetList() []PlanetConfig {
	return defaultPlanets
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
	sweMu.Lock()
	defer sweMu.Unlock()

	swephgo.Close()
	fmt.Println("\nDone!")
}
