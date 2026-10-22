package main

/*
#cgo LDFLAGS: -lswe -lm
*/
import (
	"fmt"
	"github.com/mshafiee/swephgo"
)

func main() {
	printVersion()

	// 1. โหลดค่าคอนฟิกจาก config.go
	cfg := GetDefaultConfig()
	signs := GetZodiacSigns()

	// แสดงข้อมูลพิกัดสถานที่เกิด
	printLocationInfo(cfg)

	// 2. คำนวณเวลาและค่าทางดาราศาสตร์
	jd, ayanamsa := setupEphemeris(cfg)

	// 3. ประมวลผลและแสดงผลตำแหน่งดาว
	processPlanets(jd, ayanamsa, signs)

	// 4. ประมวลผลและแสดงผลลัคนาและภพ
	processHouses(jd, cfg.Latitude, cfg.Longitude, cfg.HouseSystem, ayanamsa, signs)

	cleanup()
}

// --- Helper Functions ---

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

func processPlanets(jd float64, ayanamsa float64, signs []string) {
	fmt.Println("Planetary Positions (Sidereal - Lahiri Zodiac)")
	fmt.Println("----------------------------------------------------------------------------------")

	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	xx := make([]float64, 6)
	serr := make([]byte, 256)

	for _, p := range getPlanetList() {
		calcAndPrintPlanet(jd, p, flag, xx, serr, signs)
	}

	calcAndPrintKetu(jd, flag, xx, serr, signs)
}

func calcAndPrintPlanet(jd float64, p PlanetConfig, flag int, xx []float64, serr []byte, signs []string) {
	res := swephgo.CalcUt(jd, p.ID, flag, xx, serr)
	if res < 0 {
		fmt.Printf("%-20s: Error - %s\n", p.Name, string(serr))
		return
	}

	lon := xx[0]
	speed := xx[3]
	posStr := FormatSignDMS(lon, signs)
	nakStr := FormatNakshatraString(lon)

	retrograde := ""
	if speed < 0 {
		retrograde = " (R)"
	}

	fmt.Printf("%-20s: %9.4f° | %-32s%s | Nak: %s\n",
		p.Name, lon, posStr, retrograde, nakStr)
}

func calcAndPrintKetu(jd float64, flag int, xx []float64, serr []byte, signs []string) {
	resNode := swephgo.CalcUt(jd, swephgo.SeMeanNode, flag, xx, serr)
	if resNode >= 0 {
		ketuLon := xx[0] + 180.0
		if ketuLon >= 360.0 {
			ketuLon -= 360.0
		}
		posStr := FormatSignDMS(ketuLon, signs)
		nakStr := FormatNakshatraString(ketuLon)
		fmt.Printf("%-20s: %9.4f° | %-32s | Nak: %s\n", "Ketu (Calculated)", ketuLon, posStr, nakStr)
	}
}

func processHouses(jd float64, lat float64, lon float64, houseSys int, ayanamsa float64, signs []string) {
	cusps := make([]float64, 13)
	ascmc := make([]float64, 10)

	resHouses := swephgo.Houses(jd, lat, lon, houseSys, cusps, ascmc)
	if resHouses < 0 {
		return
	}

	fmt.Println("\nAscendant & Angles (Sidereal - Lahiri Zodiac)")
	fmt.Println("----------------------------------------------------------------------------------")

	asc := NormalizeAngle(ascmc[0] - ayanamsa)
	mc := NormalizeAngle(ascmc[1] - ayanamsa)
	dsc := NormalizeAngle(asc + 180.0)
	ic := NormalizeAngle(mc + 180.0)

	printAngle("Ascendant (ASC)", asc, signs)
	printAngle("Midheaven (MC)", mc, signs)
	printAngle("Descendant (DSC)", dsc, signs)
	printAngle("Imum Coeli (IC)", ic, signs)

	printWholeSignHouses(asc, signs)
}

func printAngle(name string, val float64, signs []string) {
	posStr := FormatSignDMS(val, signs)
	nakStr := FormatNakshatraString(val)
	fmt.Printf("%-20s: %9.4f° | %-32s | Nak: %s\n", name, val, posStr, nakStr)
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