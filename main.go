package main

/*
#cgo LDFLAGS: -lswe -lm
*/
import (
	"fmt"
	"github.com/mshafiee/swephgo"
)

func main() {
	RunHoroscope()
}

// RunHoroscope ทำหน้าที่รวบรวมขั้นตอนการคำนวณและแสดงผลทั้งหมด
func RunHoroscope() {
	printVersion()

	cfg := GetDefaultConfig()
	signs := GetZodiacSigns()

	printLocationInfo(cfg)

	jd, ayanamsa := setupEphemeris(cfg)
	flag := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal

	fmt.Println("Planetary Positions (Sidereal - Lahiri Zodiac)")
	fmt.Println("----------------------------------------------------------------------------------")
	planetResults := CalculatePlanets(jd, flag, signs)
	for _, p := range planetResults {
		retrograde := ""
		if p.IsRetro {
			retrograde = " (R)"
		}
		fmt.Printf("%-20s: %9.4f° | %-32s%s | Nak: %s\n",
			p.Name, p.Longitude, p.PosStr, retrograde, p.NakStr)
	}

	fmt.Println("\nAscendant & Angles (Sidereal - Lahiri Zodiac)")
	fmt.Println("----------------------------------------------------------------------------------")
	houseResults := CalculateHouses(jd, cfg.Latitude, cfg.Longitude, cfg.HouseSystem, ayanamsa, signs)
	for _, h := range houseResults {
		fmt.Printf("%-20s: %9.4f° | %-32s | Nak: %s\n", h.Name, h.Degree, h.PosStr, h.NakStr)
	}

	if len(houseResults) > 0 {
		printWholeSignHouses(houseResults[0].Degree, signs)
	}

	cleanup()
}