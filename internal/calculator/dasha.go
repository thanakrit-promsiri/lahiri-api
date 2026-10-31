package calculator

import (
	"fmt"
	"math"
	"time"
)

// ลำดับดาวและจำนวนปีของมหาทศาในระบบวิมโสตรีทศา (รวม 120 ปี)
var dashaSequence = []struct {
	Lord  string
	Years float64
}{
	{Lord: "Ketu", Years: 7},
	{Lord: "Venus", Years: 20},
	{Lord: "Sun", Years: 6},
	{Lord: "Moon", Years: 10},
	{Lord: "Mars", Years: 7},
	{Lord: "Rahu", Years: 18},
	{Lord: "Jupiter", Years: 16},
	{Lord: "Saturn", Years: 19},
	{Lord: "Mercury", Years: 17},
}

const daysPerYear = 365.2425

// addYearsToTime บวกจำนวนปี (ทศนิยม) เข้ากับเวลาตั้งต้น
func addYearsToTime(t time.Time, years float64) time.Time {
	days := years * daysPerYear
	seconds := int64(math.Round(days * 86400.0))
	return t.Add(time.Duration(seconds) * time.Second)
}

// formatBalanceThai แปลงจำนวนปีที่เหลือตอนเกิดเป็นข้อความภาษาไทย เช่น "Rahu 16 ปี 4 เดือน 12 วัน"
func formatBalanceThai(lord string, remainingYears float64) string {
	totalDays := int(math.Round(remainingYears * daysPerYear))
	years := totalDays / 365
	remDays := totalDays % 365
	months := remDays / 30
	days := remDays % 30
	return fmt.Sprintf("%s %d ปี %d เดือน %d วัน", lord, years, months, days)
}

// CalculateVimshottariDasha คำนวณวิมโสตรีทศา (120 ปี) จากลองจิจูดของดวงจันทร์และวันเวลาเกิด
func CalculateVimshottariDasha(moonLon float64, year, month, day int, hourLocal, timeZone float64, targetTime time.Time) VimshottariDasha {
	nakInfo, _, offsetInNak := GetNakshatra(moonLon)
	span := 360.0 / 27.0 // 13.333333333333334 องศา

	// สัดส่วนที่เหลือในนักษัตร (0.0 - 1.0)
	fractionRemaining := (span - offsetInNak) / span
	if fractionRemaining < 0 {
		fractionRemaining = 0
	}
	if fractionRemaining > 1 {
		fractionRemaining = 1
	}

	// หา index เริ่มต้นในลำดับทศาทั้ง 9
	startIdx := (nakInfo.Index - 1) % 9
	firstDasha := dashaSequence[startIdx]
	balanceYears := firstDasha.Years * fractionRemaining

	// สร้างเวลาเกิดตาม TimeZone
	tzOffsetSec := int(math.Round(timeZone * 3600))
	loc := time.FixedZone("LocalBirthTime", tzOffsetSec)
	hourInt := int(hourLocal)
	minFloat := (hourLocal - float64(hourInt)) * 60.0
	minInt := int(minFloat)
	secInt := int(math.Round((minFloat - float64(minInt)) * 60.0))
	if secInt >= 60 {
		secInt = 59
	}

	birthTime := time.Date(year, time.Month(month), day, hourInt, minInt, secInt, 0, loc)
	if targetTime.IsZero() {
		targetTime = time.Now().In(loc)
	}

	mahadashaList := make([]MahadashaItem, 0, 9)
	var currentDasha CurrentDasha
	foundCurrent := false

	// เวลาเริ่มต้นของมหาทศาแรก (นับย้อนกลับไปก่อนเกิดเพื่อใช้คำนวณอันตรทศาได้อย่างแม่นยำ)
	elapsedFirstYears := firstDasha.Years - balanceYears
	mahaFullStart := addYearsToTime(birthTime, -elapsedFirstYears)
	cursorStart := birthTime

	for i := 0; i < 9; i++ {
		seqIdx := (startIdx + i) % 9
		maha := dashaSequence[seqIdx]

		var mahaEnd time.Time
		var fullMahaStart time.Time
		if i == 0 {
			mahaEnd = addYearsToTime(birthTime, balanceYears)
			fullMahaStart = mahaFullStart
		} else {
			mahaEnd = addYearsToTime(cursorStart, maha.Years)
			fullMahaStart = cursorStart
		}

		mahadashaList = append(mahadashaList, MahadashaItem{
			Lord:      maha.Lord,
			LordThai:  TranslateLordThai(maha.Lord),
			StartDate: cursorStart.Format("2006-01-02"),
			EndDate:   mahaEnd.Format("2006-01-02"),
		})

		// ตรวจสอบว่า targetTime อยู่ในช่วงมหาทศานี้หรือไม่
		if !foundCurrent && (!targetTime.Before(cursorStart) && targetTime.Before(mahaEnd)) {
			foundCurrent = true
			currentDasha.MahadashaLord = maha.Lord
			currentDasha.MahadashaLordThai = TranslateLordThai(maha.Lord)
			currentDasha.MahadashaStart = cursorStart.Format("2006-01-02")
			currentDasha.MahadashaEnd = mahaEnd.Format("2006-01-02")

			// คำนวณอันตรทศา (Antardasha / Bhukti) ทั้ง 9 ยุคย่อยภายในมหาทศานี้
			antarCursor := fullMahaStart
			for j := 0; j < 9; j++ {
				antarIdx := (seqIdx + j) % 9
				antar := dashaSequence[antarIdx]
				// สูตรระยะเวลาอันตรทศา = (ปีมหาทศา * ปีอันตรทศา) / 120
				antarYears := (maha.Years * antar.Years) / 120.0
				antarEnd := addYearsToTime(antarCursor, antarYears)

				if !targetTime.Before(antarCursor) && targetTime.Before(antarEnd) {
					displayAntarStart := antarCursor
					if displayAntarStart.Before(birthTime) {
						displayAntarStart = birthTime
					}
					currentDasha.AntardashaLord = antar.Lord
					currentDasha.AntardashaLordThai = TranslateLordThai(antar.Lord)
					currentDasha.AntardashaStart = displayAntarStart.Format("2006-01-02")
					currentDasha.AntardashaEnd = antarEnd.Format("2006-01-02")
					break
				}
				antarCursor = antarEnd
			}
		}

		cursorStart = mahaEnd
	}

	// กรณีเวลาปัจจุบันอยู่นอกช่วง (เช่น วันเกิดในอนาคต หรือเกิน 120 ปี) ให้ใช้มหาทศาแรกเป็นค่าเริ่มต้น
	if !foundCurrent && len(mahadashaList) > 0 {
		first := mahadashaList[0]
		currentDasha.MahadashaLord = first.Lord
		currentDasha.MahadashaLordThai = first.LordThai
		currentDasha.MahadashaStart = first.StartDate
		currentDasha.MahadashaEnd = first.EndDate
		currentDasha.AntardashaLord = first.Lord
		currentDasha.AntardashaLordThai = first.LordThai
		currentDasha.AntardashaStart = first.StartDate
		currentDasha.AntardashaEnd = first.EndDate
	}

	return VimshottariDasha{
		MoonNakshatra:     nakInfo.Name,
		MoonNakshatraThai: nakInfo.NameThai,
		BalanceAtBirth:    formatBalanceThai(firstDasha.Lord, balanceYears),
		CurrentDasha:      currentDasha,
		MahadashaList:     mahadashaList,
	}
}
