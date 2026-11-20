package calculator

import (
	"fmt"
)

// Gana classifications
const (
	GanaDeva     = "Deva"
	GanaManushya = "Manushya"
	GanaRakshasa = "Rakshasa"
)

// Yoni animal archetypes (14 animals)
const (
	YoniHorse    = "Horse"
	YoniElephant = "Elephant"
	YoniSheep    = "Sheep"
	YoniSerpent  = "Serpent"
	YoniDog      = "Dog"
	YoniCat      = "Cat"
	YoniRat      = "Rat"
	YoniCow      = "Cow"
	YoniBuffalo  = "Buffalo"
	YoniTiger    = "Tiger"
	YoniDeer     = "Deer"
	YoniMonkey   = "Monkey"
	YoniLion     = "Lion"
	YoniMongoose = "Mongoose"
)

// Nadi classifications (3 Nadis)
const (
	NadiAdi    = "Adi"    // Vata
	NadiMadhya = "Madhya" // Pitta
	NadiAntya  = "Antya"  // Kapha
)

// nakshatraGanaTable: ลำดับนักษัตร 1 ถึง 27
var nakshatraGanaTable = [27]string{
	GanaDeva,     // 1: Ashwini
	GanaManushya, // 2: Bharani
	GanaRakshasa, // 3: Krittika
	GanaManushya, // 4: Rohini
	GanaDeva,     // 5: Mrigashira
	GanaManushya, // 6: Ardra
	GanaDeva,     // 7: Punarvasu
	GanaDeva,     // 8: Pushya
	GanaRakshasa, // 9: Ashlesha
	GanaRakshasa, // 10: Magha
	GanaManushya, // 11: Purva Phalguni
	GanaManushya, // 12: Uttara Phalguni
	GanaDeva,     // 13: Hasta
	GanaRakshasa, // 14: Chitra
	GanaDeva,     // 15: Swati
	GanaRakshasa, // 16: Vishakha
	GanaDeva,     // 17: Anuradha
	GanaRakshasa, // 18: Jyeshtha
	GanaRakshasa, // 19: Mula
	GanaManushya, // 20: Purva Ashadha
	GanaManushya, // 21: Uttara Ashadha
	GanaDeva,     // 22: Shravana
	GanaRakshasa, // 23: Dhanishta
	GanaRakshasa, // 24: Shatabhisha
	GanaManushya, // 25: Purva Bhadrapada
	GanaManushya, // 26: Uttara Bhadrapada
	GanaDeva,     // 27: Revati
}

// nakshatraYoniTable: ลำดับนักษัตร 1 ถึง 27
var nakshatraYoniTable = [27]string{
	YoniHorse,    // 1: Ashwini
	YoniElephant, // 2: Bharani
	YoniSheep,    // 3: Krittika
	YoniSerpent,  // 4: Rohini
	YoniSerpent,  // 5: Mrigashira
	YoniDog,      // 6: Ardra
	YoniCat,      // 7: Punarvasu
	YoniSheep,    // 8: Pushya
	YoniCat,      // 9: Ashlesha
	YoniRat,      // 10: Magha
	YoniRat,      // 11: Purva Phalguni
	YoniCow,      // 12: Uttara Phalguni
	YoniBuffalo,  // 13: Hasta
	YoniTiger,    // 14: Chitra
	YoniBuffalo,  // 15: Swati
	YoniTiger,    // 16: Vishakha
	YoniDeer,     // 17: Anuradha
	YoniDeer,     // 18: Jyeshtha
	YoniDog,      // 19: Mula
	YoniMonkey,   // 20: Purva Ashadha
	YoniMongoose, // 21: Uttara Ashadha
	YoniMonkey,   // 22: Shravana
	YoniLion,     // 23: Dhanishta
	YoniHorse,    // 24: Shatabhisha
	YoniLion,     // 25: Purva Bhadrapada
	YoniCow,      // 26: Uttara Bhadrapada
	YoniElephant, // 27: Revati
}

// nakshatraNadiTable: ลำดับนักษัตร 1 ถึง 27
var nakshatraNadiTable = [27]string{
	NadiAdi,    // 1: Ashwini
	NadiMadhya, // 2: Bharani
	NadiAntya,  // 3: Krittika
	NadiAntya,  // 4: Rohini
	NadiMadhya, // 5: Mrigashira
	NadiAdi,    // 6: Ardra
	NadiAdi,    // 7: Punarvasu
	NadiMadhya, // 8: Pushya
	NadiAntya,  // 9: Ashlesha
	NadiAntya,  // 10: Magha
	NadiMadhya, // 11: Purva Phalguni
	NadiAdi,    // 12: Uttara Phalguni
	NadiAdi,    // 13: Hasta
	NadiMadhya, // 14: Chitra
	NadiAntya,  // 15: Swati
	NadiAntya,  // 16: Vishakha
	NadiMadhya, // 17: Anuradha
	NadiAdi,    // 18: Jyeshtha
	NadiAdi,    // 19: Mula
	NadiMadhya, // 20: Purva Ashadha
	NadiAntya,  // 21: Uttara Ashadha
	NadiAntya,  // 22: Shravana
	NadiMadhya, // 23: Dhanishta
	NadiAdi,    // 24: Shatabhisha
	NadiAdi,    // 25: Purva Bhadrapada
	NadiMadhya, // 26: Uttara Bhadrapada
	NadiAntya,  // 27: Revati
}

// คู่สัตว์ที่เป็นคู่อริกันโดยตรง (Sworn Enemies: 0 คะแนน)
var hostileYoniPairs = map[string]string{
	YoniHorse:    YoniBuffalo,
	YoniBuffalo:  YoniHorse,
	YoniElephant: YoniLion,
	YoniLion:     YoniElephant,
	YoniSheep:    YoniMonkey,
	YoniMonkey:   YoniSheep,
	YoniSerpent:  YoniMongoose,
	YoniMongoose: YoniSerpent,
	YoniDog:      YoniDeer,
	YoniDeer:     YoniDog,
	YoniCat:      YoniRat,
	YoniRat:      YoniCat,
	YoniCow:      YoniTiger,
	YoniTiger:    YoniCow,
}

var taraNames = [9]string{
	"Janma (Birth)",
	"Sampat (Wealth)",
	"Vipat (Danger)",
	"Kshema (Well-being)",
	"Pratyak (Obstacles)",
	"Sadhana (Achievement)",
	"Naidhana (Affliction)",
	"Mitra (Friend)",
	"Ati-Mitra (Intimate)",
}

// GetNakshatraGana คืนค่า Gana ของนักษัตร (1-27)
func GetNakshatraGana(nakIdx int) string {
	if nakIdx >= 1 && nakIdx <= 27 {
		return nakshatraGanaTable[nakIdx-1]
	}
	return GanaDeva
}

// GetNakshatraYoni คืนค่า Yoni (สัตว์ประจำฤกษ์) ของนักษัตร (1-27)
func GetNakshatraYoni(nakIdx int) string {
	if nakIdx >= 1 && nakIdx <= 27 {
		return nakshatraYoniTable[nakIdx-1]
	}
	return YoniHorse
}

// GetNakshatraNadi คืนค่า Nadi (ธาตุลม/ไฟ/น้ำ) ของนักษัตร (1-27)
func GetNakshatraNadi(nakIdx int) string {
	if nakIdx >= 1 && nakIdx <= 27 {
		return nakshatraNadiTable[nakIdx-1]
	}
	return NadiAdi
}

// TaraBalaInfo สรุปผลการคำนวณนวตารา
type TaraBalaInfo struct {
	Number       int
	Name         string
	IsAuspicious bool
	Score        int // +1, 0, -1
}

// ComputeTaraBala คำนวณธาราพละจากนักษัตรต้นทางไปยังนักษัตรปลายทาง
func ComputeTaraBala(nakSource, nakTarget int) TaraBalaInfo {
	if nakSource < 1 || nakSource > 27 || nakTarget < 1 || nakTarget > 27 {
		return TaraBalaInfo{Number: 1, Name: taraNames[0], IsAuspicious: true, Score: 0}
	}
	dist := ((nakTarget - nakSource + 27) % 9) + 1
	name := taraNames[dist-1]
	switch dist {
	case 2, 4, 6, 8, 9:
		return TaraBalaInfo{Number: dist, Name: name, IsAuspicious: true, Score: 1}
	case 3, 5, 7:
		return TaraBalaInfo{Number: dist, Name: name, IsAuspicious: false, Score: -1}
	default:
		return TaraBalaInfo{Number: dist, Name: name, IsAuspicious: true, Score: 0}
	}
}

// ComputeGanaCompat ประเมินความเข้ากันได้ของ Gana (คืนค่าคำอธิบายและคะแนน 0-6)
func ComputeGanaCompat(ganaA, ganaB string) (string, float64) {
	if ganaA == ganaB {
		return fmt.Sprintf("Same Gana (%s)", ganaA), 6.0
	}
	if (ganaA == GanaDeva && ganaB == GanaManushya) || (ganaA == GanaManushya && ganaB == GanaDeva) {
		return "Deva-Manushya (Harmonious)", 5.0
	}
	if ganaA == GanaDeva && ganaB == GanaRakshasa {
		return "Deva-Rakshasa (Challenging)", 1.0
	}
	if ganaA == GanaRakshasa && ganaB == GanaDeva {
		return "Rakshasa-Deva (Challenging)", 1.0
	}
	return "Manushya-Rakshasa (Challenging)", 0.0
}

// ComputeYoniCompat ประเมินความเข้ากันได้ของ Yoni (คืนค่าคำอธิบายและคะแนน 0-4)
func ComputeYoniCompat(yoniA, yoniB string) (string, float64) {
	if yoniA == yoniB {
		return fmt.Sprintf("Same Yoni (%s)", yoniA), 4.0
	}
	if hostileYoniPairs[yoniA] == yoniB {
		return fmt.Sprintf("Hostile (%s vs %s)", yoniA, yoniB), 0.0
	}
	return fmt.Sprintf("Neutral (%s & %s)", yoniA, yoniB), 2.0
}

// ComputeNadiCompat ประเมินความเข้ากันได้ของ Nadi (คืนค่าคำอธิบายและคะแนน 0 หรือ 8)
func ComputeNadiCompat(nadiA, nadiB string) (string, float64) {
	if nadiA != nadiB {
		return fmt.Sprintf("Different Nadi (%s vs %s)", nadiA, nadiB), 8.0
	}
	return fmt.Sprintf("Nadi Dosha (Same %s)", nadiA), 0.0
}

// AshtakootaResult สรุปคะแนนการสมพงศ์ 8 ด้าน (เต็ม 36 คะแนน)
type AshtakootaResult struct {
	VarnaScore      float64 `json:"varna_score"`
	VashyaScore     float64 `json:"vashya_score"`
	TaraScore       float64 `json:"tara_score"`
	YoniScore       float64 `json:"yoni_score"`
	MaitriScore     float64 `json:"maitri_score"`
	GanaScore       float64 `json:"gana_score"`
	BhakootScore    float64 `json:"bhakoot_score"`
	NadiScore       float64 `json:"nadi_score"`
	TotalScore      float64 `json:"total_score"`
	MaxScore        float64 `json:"max_score"`
	Verdict         string  `json:"verdict"`
	TaraDescription string  `json:"tara_description"`
	GanaDescription string  `json:"gana_description"`
	YoniDescription string  `json:"yoni_description"`
	NadiDescription string  `json:"nadi_description"`
}

// ComputeAshtakoota คำนวณคะแนนอัษฐกูฏ 36 คะแนนระหว่างดวงจันทร์สองดวง
func ComputeAshtakoota(moonA, moonB GephiNode) AshtakootaResult {
	// 1. Varna (1 คะแนน): สถิตราศีตามวรรณะ
	varnaMap := map[int]int{3: 4, 7: 4, 11: 4, 0: 3, 4: 3, 8: 3, 1: 2, 5: 2, 9: 2, 2: 1, 6: 1, 10: 1}
	varnaA := varnaMap[moonA.Rasi-1]
	varnaB := varnaMap[moonB.Rasi-1]
	varnaScore := 1.0
	if varnaA != varnaB {
		varnaScore = 0.5
	}

	// 2. Vashya (2 คะแนน): กลุ่มราศี
	vashyaScore := 1.0
	if moonA.Rasi == moonB.Rasi {
		vashyaScore = 2.0
	}

	// 3. Tara (3 คะแนน): ธาราพละสองทิศทาง
	taraAB := ComputeTaraBala(moonA.Nakshatra, moonB.Nakshatra)
	taraBA := ComputeTaraBala(moonB.Nakshatra, moonA.Nakshatra)
	taraScore := 0.0
	if taraAB.IsAuspicious && taraBA.IsAuspicious {
		taraScore = 3.0
	} else if taraAB.IsAuspicious || taraBA.IsAuspicious {
		taraScore = 1.5
	}
	taraDesc := fmt.Sprintf("%s <-> %s", taraAB.Name, taraBA.Name)

	// 4. Yoni (4 คะแนน)
	yoniDesc, yoniScore := ComputeYoniCompat(moonA.NakshatraYoni, moonB.NakshatraYoni)

	// 5. Graha Maitri (5 คะแนน): เจ้าเรือนราศีที่ดวงจันทร์สถิต
	lordA := GetSignLord(moonA.Rasi - 1)
	lordB := GetSignLord(moonB.Rasi - 1)
	maitriScore := 2.5
	if lordA == lordB {
		maitriScore = 5.0
	} else {
		scoreAB := grahaAttitude(lordA, lordB)
		scoreBA := grahaAttitude(lordB, lordA)
		totalAtt := scoreAB + scoreBA
		switch {
		case totalAtt >= 2:
			maitriScore = 5.0
		case totalAtt == 1:
			maitriScore = 4.0
		case totalAtt == 0:
			maitriScore = 3.0
		case totalAtt == -1:
			maitriScore = 1.0
		default:
			maitriScore = 0.0
		}
	}

	// 6. Gana (6 คะแนน)
	ganaDesc, ganaScore := ComputeGanaCompat(moonA.NakshatraGana, moonB.NakshatraGana)

	// 7. Bhakoot (7 คะแนน): ระยะห่างภพเรือนราศีจันทร์
	dist := ((moonB.Rasi - moonA.Rasi + 12) % 12) + 1
	bhakootScore := 7.0
	if dist == 2 || dist == 12 || dist == 6 || dist == 8 || dist == 5 || dist == 9 {
		// ถ้าเจ้าเรือนเป็นดาวเดียวกันหรือมิตรต่อกัน Bhakoot Dosha จะได้รับการยกเว้น
		if lordA != lordB && (grahaAttitude(lordA, lordB) < 0 || grahaAttitude(lordB, lordA) < 0) {
			bhakootScore = 0.0
		}
	}

	// 8. Nadi (8 คะแนน)
	nadiDesc, nadiScore := ComputeNadiCompat(moonA.NakshatraNadi, moonB.NakshatraNadi)

	total := varnaScore + vashyaScore + taraScore + yoniScore + maitriScore + ganaScore + bhakootScore + nadiScore
	verdict := "Good (ผ่านเกณฑ์มาตรฐาน)"
	if total >= 28.0 {
		verdict = "Excellent (อุดมสมบูรณ์)"
	} else if total < 18.0 {
		verdict = "Challenging (ต่ำกว่าเกณฑ์)"
	}

	return AshtakootaResult{
		VarnaScore:      varnaScore,
		VashyaScore:     vashyaScore,
		TaraScore:       taraScore,
		YoniScore:       yoniScore,
		MaitriScore:     maitriScore,
		GanaScore:       ganaScore,
		BhakootScore:    bhakootScore,
		NadiScore:       nadiScore,
		TotalScore:      roundTo2(total),
		MaxScore:        36.0,
		Verdict:         verdict,
		TaraDescription: taraDesc,
		GanaDescription: ganaDesc,
		YoniDescription: yoniDesc,
		NadiDescription: nadiDesc,
	}
}
