package calculator

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ComputeHoroscope คำนวณข้อมูลดวงชะตาพระเวทแบบครบวงจร (ลัคนา, ภพ, ดาวเคราะห์, มาตรฐานดาว, ดาวดับ, เจ้าเรือนภพ, วิมโสตรีทศา)
func ComputeHoroscope(req CalculateReq) CalculateResp {
	jd, ayanamsa := ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)

	signs := GetZodiacSigns()
	englishSigns := GetEnglishSigns()
	thaiSigns := GetThaiZodiacSigns()

	// 1. คำนวณมุมภพ (Houses & Angles) เพื่อหาลัคนา (Ascendant) ก่อนนำไปจัดภพ Whole Sign
	houseSysCode := 'W'
	if len(req.HouseSystem) > 0 {
		houseSysCode = rune(req.HouseSystem[0])
	}
	houseRes := CalculateHouses(jd, req.Latitude, req.Longitude, int(houseSysCode), ayanamsa, signs)

	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	ascSignIdx, _ := ParseLongitude(ascDeg)

	// 2. คำนวณตำแหน่งดาว (Planets) พร้อมมาตรฐานดาว (Dignity), ดาวดับ (Combust) และภพที่สถิต
	planetRes := CalculateSiderealPlanets(jd)

	var sunLon, moonLon float64
	for _, p := range planetRes {
		if p.Name == "Sun" {
			sunLon = p.Longitude
		} else if p.Name == "Moon" {
			moonLon = p.Longitude
		}
	}

	planetHouseMap := make(map[string]int, len(planetRes))
	apiPlanets := make([]ApiPlanet, 0, len(planetRes))

	for _, p := range planetRes {
		nak, pada, _ := GetNakshatra(p.Longitude)
		posDetail := ParsePositionDetail(p.Longitude)
		navamsaRes := CalculateNavamsa(p.Longitude)

		nameThai, symbol, canonical := TranslateBodyInfo(p.Name)
		dignity, dignityThai := CalculateDignity(p.Name, p.Longitude)
		isCombust := IsPlanetCombust(p.Name, p.Longitude, sunLon, p.IsRetro)

		pSignIdx, _ := ParseLongitude(p.Longitude)
		houseNum := ((pSignIdx - ascSignIdx + 12) % 12) + 1
		houseNameThai := GetHouseNameThai(houseNum)

		planetHouseMap[canonical] = houseNum

		apiPlanets = append(apiPlanets, ApiPlanet{
			Name:            p.Name,
			NameThai:        nameThai,
			Symbol:          symbol,
			Longitude:       p.Longitude,
			PositionStr:     p.PosStr,
			DegreeFormatted: FormatDegreeShort(posDetail),
			Position:        posDetail,
			Nakshatra:       nak.Name,
			NakshatraThai:   nak.NameThai,
			Pada:            pada,
			Lord:            nak.Lord,
			LordThai:        nak.LordThai,
			IsRetrograde:    p.IsRetro,
			IsCombust:       isCombust,
			Dignity:         dignity,
			DignityThai:     dignityThai,
			House:           houseNum,
			HouseNameThai:   houseNameThai,
			HouseThai:       fmt.Sprintf("%d (%s)", houseNum, houseNameThai),
			Navamsa:         navamsaRes,
		})
	}

	// 3. จัดเตรียมข้อมูลมุมภพ (Ascendant, MC, DSC, IC)
	apiHouses := make([]ApiHouse, 0, len(houseRes))
	for _, hRes := range houseRes {
		nak, pada, _ := GetNakshatra(hRes.Degree)
		posDetail := ParsePositionDetail(hRes.Degree)
		navamsaRes := CalculateNavamsa(hRes.Degree)

		nameThai, symbol, _ := TranslateBodyInfo(hRes.Name)
		hSignIdx, _ := ParseLongitude(hRes.Degree)
		houseNum := ((hSignIdx - ascSignIdx + 12) % 12) + 1
		houseNameThai := GetHouseNameThai(houseNum)

		apiHouses = append(apiHouses, ApiHouse{
			Name:            hRes.Name,
			NameThai:        nameThai,
			Symbol:          symbol,
			Degree:          hRes.Degree,
			PositionStr:     hRes.PosStr,
			DegreeFormatted: FormatDegreeShort(posDetail),
			Position:        posDetail,
			Nakshatra:       nak.Name,
			NakshatraThai:   nak.NameThai,
			Pada:            pada,
			Lord:            nak.Lord,
			LordThai:        nak.LordThai,
			House:           houseNum,
			HouseNameThai:   houseNameThai,
			HouseThai:       fmt.Sprintf("%d (%s)", houseNum, houseNameThai),
			Navamsa:         navamsaRes,
		})
	}

	// 4. คำนวณ Whole Sign Houses พร้อมเจ้าเรือนภพ (House Lords) และภพที่ดาวเจ้าเรือนไปสถิต
	apiWholeSigns := make([]ApiWholeSign, 0, 12)
	for i := 1; i <= 12; i++ {
		currentSign := (ascSignIdx + (i - 1)) % 12
		houseNameThai := GetHouseNameThai(i)
		lord := GetSignLord(currentSign)
		lordThai := TranslateLordThai(lord)
		placedInHouse := planetHouseMap[lord]
		placedInHouseThai := GetHouseNameThai(placedInHouse)

		apiWholeSigns = append(apiWholeSigns, ApiWholeSign{
			House:                 i,
			HouseNameThai:         houseNameThai,
			Sign:                  englishSigns[currentSign],
			SignThai:              thaiSigns[currentSign],
			Lord:                  lord,
			LordThai:              lordThai,
			LordPlacedInHouse:     placedInHouse,
			LordPlacedInHouseThai: placedInHouseThai,
		})
	}

	// 5. คำนวณวิมโสตรีทศา (Vimshottari Dasha) จากตำแหน่งดวงจันทร์และวันเวลาเกิด
	dasha := CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, time.Now(),
	)

	return CalculateResp{
		Metadata: ApiMetadata{
			JulianDay:       jd,
			AyanamsaName:    "Lahiri",
			AyanamsaDegrees: FormatDMS(ayanamsa),
			Location: ApiLocation{
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
			},
		},
		Planets:          apiPlanets,
		HousesAndAngles:  apiHouses,
		WholeSignHouses:  apiWholeSigns,
		VimshottariDasha: dasha,
	}
}

// ComputeDasha คำนวณวิมโสตรีทศา ณ เวลาเป้าหมาย (รองรับปี/เดือนจร)
func ComputeDasha(req TransitReq) VimshottariDasha {
	jd, _ := ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)

	planetRes := CalculateSiderealPlanets(jd)

	var moonLon float64
	for _, p := range planetRes {
		if p.Name == "Moon" {
			moonLon = p.Longitude
			break
		}
	}

	hasCustomTarget := req.TransitYear != 0 || req.TransitMonth != 0
	targetTime := time.Now()
	if hasCustomTarget {
		tMonth := req.TransitMonth
		if tMonth == 0 {
			tMonth = req.Month
		}
		tDay := ClampDayToMonth(req.TransitYear, tMonth, req.Day)
		tzOffsetSec := int(math.Round(req.TimeZone * 3600))
		loc := time.FixedZone("LocalBirthTime", tzOffsetSec)
		hourInt := int(req.HourThai)
		minFloat := (req.HourThai - float64(hourInt)) * 60.0
		minInt := int(minFloat)
		secInt := int(math.Round((minFloat - float64(minInt)) * 60.0))
		if secInt >= 60 {
			secInt = 59
		}
		targetTime = time.Date(req.TransitYear, time.Month(tMonth), tDay, hourInt, minInt, secInt, 0, loc)
	}

	return CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, targetTime,
	)
}

// ComputeNatalAscSignIdx คำนวณราศีลัคนาดวงกำเนิด (0-11) เพื่อกำหนดเป็นเรือนที่ 1 (Whole Sign)
func ComputeNatalAscSignIdx(req TransitReq) int {
	natalJD, natalAyanamsa := ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)
	houseRes := CalculateHouses(natalJD, req.Latitude, req.Longitude, int('W'), natalAyanamsa, GetZodiacSigns())
	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	ascSignIdx, _ := ParseLongitude(ascDeg)
	return ascSignIdx
}

// ComputeSingleMonthTransit คำนวณดาวจร 9 ดวงในเดือนที่ระบุ เทียบเรือนภพกับลัคนาดวงกำเนิด
func ComputeSingleMonthTransit(req TransitReq, m int, ascSignIdx int) MonthlyTransit {
	tDay := ClampDayToMonth(req.TransitYear, m, req.Day)
	tJD, _ := ComputeJulianDayAndAyanamsa(
		req.TransitYear, m, tDay, req.HourThai, req.TimeZone,
	)

	tPlanetRes := CalculateSiderealPlanets(tJD)

	var sunLon float64
	planetByCanonical := make(map[string]PlanetResult, len(tPlanetRes))
	for _, p := range tPlanetRes {
		_, _, canonical := TranslateBodyInfo(p.Name)
		planetByCanonical[canonical] = p
		if canonical == "Sun" {
			sunLon = p.Longitude
		}
	}

	transitPlanets := make([]TransitPlanet, 0, len(NavagrahaOrder))
	for idx, canonicalName := range NavagrahaOrder {
		p, exists := planetByCanonical[canonicalName]
		if !exists {
			continue
		}

		nak, pada, _ := GetNakshatra(p.Longitude)
		posDetail := ParsePositionDetail(p.Longitude)

		nameThai, symbol, canonical := TranslateBodyInfo(p.Name)
		dignity, dignityThai := CalculateDignity(p.Name, p.Longitude)
		isCombust := IsPlanetCombust(p.Name, p.Longitude, sunLon, p.IsRetro)

		// นับเรือนเทียบกับลัคนาดวงกำเนิด: เรือน = ((ราศีดาวจร - ราศีลัคนา) mod 12) + 1
		pSignIdx, _ := ParseLongitude(p.Longitude)
		houseNum := ((pSignIdx - ascSignIdx + 12) % 12) + 1
		houseNameThai := GetHouseNameThai(houseNum)

		transitPlanets = append(transitPlanets, TransitPlanet{
			No:            idx + 1,
			Name:          canonical,
			NameThai:      nameThai,
			Symbol:        symbol,
			Longitude:     p.Longitude,
			Sign:          posDetail.Sign,
			SignThai:      posDetail.SignThai,
			Degree:        FormatDegreeShort(posDetail),
			House:         houseNum,
			HouseNameThai: houseNameThai,
			IsRetrograde:  p.IsRetro,
			IsCombust:     isCombust,
			Dignity:       dignity,
			DignityThai:   dignityThai,
			Nakshatra:     nak.Name,
			NakshatraThai: nak.NameThai,
			Pada:          pada,
		})
	}

	return MonthlyTransit{
		Month:   m,
		Date:    fmt.Sprintf("%02d/%02d/%04d", tDay, m, req.TransitYear),
		Planets: transitPlanets,
	}
}

// ComputeAnnualTransit คำนวณลัคนาดวงกำเนิด (เรือนที่ 1) และคืนค่าเฉพาะดาวจรรายเดือนทั้ง 12 เดือนในปีจร
func ComputeAnnualTransit(req TransitReq) []MonthlyTransit {
	ascSignIdx := ComputeNatalAscSignIdx(req)

	monthlyTransits := make([]MonthlyTransit, 0, 12)
	for m := 1; m <= 12; m++ {
		monthlyTransits = append(monthlyTransits, ComputeSingleMonthTransit(req, m, ascSignIdx))
	}

	return monthlyTransits
}

// ComputeMonthTransit คำนวณดาวจรเฉพาะเดือนเดียวที่ระบุใน transit_month
func ComputeMonthTransit(req TransitReq) MonthlyTransit {
	ascSignIdx := ComputeNatalAscSignIdx(req)
	return ComputeSingleMonthTransit(req, req.TransitMonth, ascSignIdx)
}

// ComputeMonthTransitWithDasha คำนวณดาวจรเฉพาะเดือนที่ระบุ พร้อมทศา/อนุทศาของเจ้าชะตา ณ เวลาดาวจรนั้น
func ComputeMonthTransitWithDasha(req TransitReq) MonthlyTransit {
	// 1. คำนวณลัคนาดวงกำเนิด (เรือนที่ 1) และองศาดวงจันทร์กำเนิดสำหรับวิมโสตรีทศา
	natalJD, natalAyanamsa := ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)
	houseRes := CalculateHouses(natalJD, req.Latitude, req.Longitude, int('W'), natalAyanamsa, GetZodiacSigns())
	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	ascSignIdx, _ := ParseLongitude(ascDeg)

	natalPlanets := CalculateSiderealPlanets(natalJD)
	var moonLon float64
	for _, p := range natalPlanets {
		if p.Name == "Moon" {
			moonLon = p.Longitude
			break
		}
	}

	// 2. คำนวณดาวจรของเดือนที่ระบุ
	mt := ComputeSingleMonthTransit(req, req.TransitMonth, ascSignIdx)

	// 3. คำนวณทศา (มหาทศา, อันตรทศา/อนุทศา, ปรัตยันตรทศา) ณ วันเวลาดาวจรของเดือนนั้น
	tDay := ClampDayToMonth(req.TransitYear, req.TransitMonth, req.Day)
	tzOffsetSec := int(math.Round(req.TimeZone * 3600))
	loc := time.FixedZone("LocalBirthTime", tzOffsetSec)
	hourInt := int(req.HourThai)
	minFloat := (req.HourThai - float64(hourInt)) * 60.0
	minInt := int(minFloat)
	secInt := int(math.Round((minFloat - float64(minInt)) * 60.0))
	if secInt >= 60 {
		secInt = 59
	}
	targetTime := time.Date(req.TransitYear, time.Month(req.TransitMonth), tDay, hourInt, minInt, secInt, 0, loc)

	fullDasha := CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, targetTime,
	)

	mt.Dasha = &TransitDasha{
		MoonNakshatra:           fullDasha.MoonNakshatra,
		MoonNakshatraThai:       fullDasha.MoonNakshatraThai,
		MahadashaLord:           fullDasha.CurrentDasha.MahadashaLord,
		MahadashaLordThai:       fullDasha.CurrentDasha.MahadashaLordThai,
		MahadashaStart:          fullDasha.CurrentDasha.MahadashaStart,
		MahadashaEnd:            fullDasha.CurrentDasha.MahadashaEnd,
		AntardashaLord:          fullDasha.CurrentDasha.AntardashaLord,
		AntardashaLordThai:      fullDasha.CurrentDasha.AntardashaLordThai,
		AntardashaStart:         fullDasha.CurrentDasha.AntardashaStart,
		AntardashaEnd:           fullDasha.CurrentDasha.AntardashaEnd,
		PratyantardashaLord:     fullDasha.CurrentDasha.PratyantardashaLord,
		PratyantardashaLordThai: fullDasha.CurrentDasha.PratyantardashaLordThai,
		PratyantardashaStart:    fullDasha.CurrentDasha.PratyantardashaStart,
		PratyantardashaEnd:      fullDasha.CurrentDasha.PratyantardashaEnd,
	}

	return mt
}

// ComputeTransitDashaPromptData คำนวณข้อมูลทั้งหมดที่จำเป็นสำหรับการสร้างข้อความ Prompt ดาวจร+ทศาจร (Roman Ultra-Compact)
func ComputeTransitDashaPromptData(req TransitReq) TransitDashaPromptData {
	englishSigns := GetEnglishSigns()

	// 1. คำนวณพื้นดวงกำเนิด (ลัคนากำเนิด, ตำแหน่งดาวกำเนิด, ดวงจันทร์กำเนิด)
	natalJD, natalAyanamsa := ComputeJulianDayAndAyanamsa(
		req.Year, req.Month, req.Day, req.HourThai, req.TimeZone,
	)
	houseRes := CalculateHouses(natalJD, req.Latitude, req.Longitude, int('W'), natalAyanamsa, GetZodiacSigns())
	var ascDeg float64
	if len(houseRes) > 0 {
		ascDeg = houseRes[0].Degree
	}
	ascSignIdx, _ := ParseLongitude(ascDeg)
	ascSignShort := ShortSignCode(englishSigns[ascSignIdx])

	natalPlanets := CalculateSiderealPlanets(natalJD)
	natalHouseByShort := make(map[string]int, len(natalPlanets))
	var moonLon float64
	var moonSignShort string
	var moonHouse int
	for _, p := range natalPlanets {
		pShort := ShortPlanetCode(p.Name)
		pSignIdx, _ := ParseLongitude(p.Longitude)
		hNum := ((pSignIdx - ascSignIdx + 12) % 12) + 1
		natalHouseByShort[pShort] = hNum
		if p.Name == "Moon" {
			moonLon = p.Longitude
			moonSignShort = ShortSignCode(englishSigns[pSignIdx])
			moonHouse = hNum
		}
	}

	// 2. กำหนดวันเวลาเป้าหมายของเดือนจร และคำนวณทศา ณ เวลาดาวจร
	tDay := ClampDayToMonth(req.TransitYear, req.TransitMonth, req.Day)
	dateISO := fmt.Sprintf("%04d-%02d-%02d", req.TransitYear, req.TransitMonth, tDay)

	tzOffsetSec := int(math.Round(req.TimeZone * 3600))
	loc := time.FixedZone("LocalBirthTime", tzOffsetSec)
	hourInt := int(req.HourThai)
	minFloat := (req.HourThai - float64(hourInt)) * 60.0
	minInt := int(minFloat)
	secInt := int(math.Round((minFloat - float64(minInt)) * 60.0))
	if secInt >= 60 {
		secInt = 59
	}
	targetTime := time.Date(req.TransitYear, time.Month(req.TransitMonth), tDay, hourInt, minInt, secInt, 0, loc)

	dasha := CalculateVimshottariDasha(
		moonLon, req.Year, req.Month, req.Day, req.HourThai, req.TimeZone, targetTime,
	)

	// 3. คำนวณดาวจรนวเคราะห์ 9 ดวง ณ วันเวลาของเดือนจร
	tJD, _ := ComputeJulianDayAndAyanamsa(
		req.TransitYear, req.TransitMonth, tDay, req.HourThai, req.TimeZone,
	)
	tPlanetRes := CalculateSiderealPlanets(tJD)

	var sunLon float64
	planetByCanonical := make(map[string]PlanetResult, len(tPlanetRes))
	for _, p := range tPlanetRes {
		_, _, canonical := TranslateBodyInfo(p.Name)
		planetByCanonical[canonical] = p
		if canonical == "Sun" {
			sunLon = p.Longitude
		}
	}

	transitHouseByShort := make(map[string]int, len(NavagrahaOrder))
	transitPlanets := make([]TransitPromptBody, 0, len(NavagrahaOrder))
	vedicBodies := make([]TransitPromptBody, 0, len(NavagrahaOrder))

	for _, canonicalName := range NavagrahaOrder {
		p, exists := planetByCanonical[canonicalName]
		if !exists {
			continue
		}

		pShort := ShortPlanetCode(p.Name)
		posDetail := ParsePositionDetail(p.Longitude)
		navamsaRes := CalculateNavamsa(p.Longitude)
		nak, pada, _ := GetNakshatra(p.Longitude)
		dignity, _ := CalculateDignity(p.Name, p.Longitude)
		isCombust := IsPlanetCombust(p.Name, p.Longitude, sunLon, p.IsRetro)

		pSignIdx, _ := ParseLongitude(p.Longitude)
		houseNum := ((pSignIdx - ascSignIdx + 12) % 12) + 1
		transitHouseByShort[pShort] = houseNum

		d1SignShort := ShortSignCode(posDetail.Sign)
		d9SignShort := ShortSignCode(navamsaRes.Sign)

		var flags []string
		if p.IsRetro {
			flags = append(flags, "R")
		}
		if isCombust {
			flags = append(flags, "C")
		}
		if posDetail.Sign != "" && posDetail.Sign == navamsaRes.Sign {
			flags = append(flags, "Varg")
		}
		flagStr := ""
		if len(flags) > 0 {
			flagStr = "[" + strings.Join(flags, ",") + "]"
		}

		d1DigShort := ShortDignityCode(dignity)
		d9SignIdx := GetSignIndexByName(navamsaRes.Sign)
		d9DigEng, _ := CalculateSignDignity(p.Name, d9SignIdx)
		d9DigShort := ShortDignityCode(d9DigEng)

		nakShort := ShortNakshatraName(nak.Name)
		lordShort := ShortPlanetCode(nak.Lord)

		sec := int(math.Round(posDetail.Second))
		min := posDetail.Minute
		deg := posDetail.Degree
		if sec >= 60 {
			sec = 0
			min++
			if min >= 60 {
				min = 0
				deg++
			}
		}
		dmsStr := fmt.Sprintf("%02d°%02d'%02d\"", deg, min, sec)

		formattedLine := fmt.Sprintf("%s%s:%s %s(%s)|H%d|%s(%d,%s)|D9:%s(%s)",
			pShort, flagStr,
			d1SignShort, dmsStr, d1DigShort,
			houseNum,
			nakShort, pada, lordShort,
			d9SignShort, d9DigShort,
		)

		body := TransitPromptBody{
			Code:        pShort,
			House:       houseNum,
			Lon:         p.Longitude,
			D1SignShort: d1SignShort,
			DMS:         dmsStr,
			D1DigShort:  d1DigShort,
			NakShort:    nakShort,
			Pada:        pada,
			LordShort:   lordShort,
			D9SignShort: d9SignShort,
			D9DigShort:  d9DigShort,
			Flags:       flagStr,
			Formatted:   formattedLine,
		}

		transitPlanets = append(transitPlanets, body)
		vedicBodies = append(vedicBodies, body)
	}

	return TransitDashaPromptData{
		DateISO:             dateISO,
		NatalAscSignShort:   ascSignShort,
		NatalMoonSignShort:  moonSignShort,
		NatalMoonHouse:      moonHouse,
		NatalMoonNakName:    dasha.MoonNakshatra,
		Dasha:               dasha,
		NatalHouseByShort:   natalHouseByShort,
		TransitHouseByShort: transitHouseByShort,
		TransitPlanets:      transitPlanets,
		VedicBodies:         vedicBodies,
	}
}
