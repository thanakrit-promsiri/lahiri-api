package api

import (
	"fmt"
	"math"
	"strings"

	"lahiri-api/internal/calculator"
)

// formatDMSInt จัดรูปแบบองศา ลิปดา ฟิลิปดา (ปัดเศษฟิลิปดาเป็นจำนวนเต็ม) เช่น 07°49'03"
func formatDMSInt(pos calculator.PositionDetail) string {
	sec := int(math.Round(pos.Second))
	min := pos.Minute
	deg := pos.Degree
	if sec >= 60 {
		sec = 0
		min++
		if min >= 60 {
			min = 0
			deg++
		}
	}
	return fmt.Sprintf("%02d°%02d'%02d\"", deg, min, sec)
}

// formatDMSDiff แปลงระยะห่างเชิงมุม (องศาทศนิยม) เป็นรูปแบบ deg°min'
func formatDMSDiff(diffDeg float64) string {
	deg := int(diffDeg)
	min := int(math.Round((diffDeg - float64(deg)) * 60.0))
	if min >= 60 {
		min = 0
		deg++
	}
	return fmt.Sprintf("%d°%02d'", deg, min)
}

// bodyDrishtiItem เก็บข้อมูลดาว/จุดสำหรับคำนวณเกณฑ์สายตา
type bodyDrishtiItem struct {
	Code  string
	House int
	Lon   float64
}

// formatConjunctions สรุปดาวกุมกันพร้อม Orb และ DeepCombust
func formatConjunctions(bodies []bodyDrishtiItem) string {
	houseBodies := make(map[int][]bodyDrishtiItem, 12)
	for _, b := range bodies {
		if b.House >= 1 && b.House <= 12 {
			houseBodies[b.House] = append(houseBodies[b.House], b)
		}
	}

	var conjParts []string
	for h := 1; h <= 12; h++ {
		bList := houseBodies[h]
		if len(bList) < 2 {
			continue
		}

		if len(bList) == 2 {
			b1, b2 := bList[0], bList[1]
			orb := calculator.AngularDistance(b1.Lon, b2.Lon)
			extra := ""
			if (b1.Code == "Su" || b2.Code == "Su") && orb <= 5.0 {
				extra = ",DeepCombust"
			}
			conjParts = append(conjParts, fmt.Sprintf("H%d(%s+%s[orb %s%s])",
				h, b1.Code, b2.Code, formatDMSDiff(orb), extra,
			))
		} else {
			sunIdx := -1
			for idx, b := range bList {
				if b.Code == "Su" {
					sunIdx = idx
					break
				}
			}

			var pairStrs []string
			if sunIdx >= 0 {
				sBody := bList[sunIdx]
				for idx, b := range bList {
					if idx == sunIdx {
						continue
					}
					orb := calculator.AngularDistance(sBody.Lon, b.Lon)
					extra := ""
					if orb <= 5.0 {
						extra = ",DeepCombust"
					}
					pairStrs = append(pairStrs, fmt.Sprintf("%s+%s[orb %s%s]",
						sBody.Code, b.Code, formatDMSDiff(orb), extra,
					))
				}
			} else {
				base := bList[0]
				for idx := 1; idx < len(bList); idx++ {
					b := bList[idx]
					orb := calculator.AngularDistance(base.Lon, b.Lon)
					pairStrs = append(pairStrs, fmt.Sprintf("%s+%s[orb %s]",
						base.Code, b.Code, formatDMSDiff(orb),
					))
				}
			}
			conjParts = append(conjParts, fmt.Sprintf("H%d(%s)", h, strings.Join(pairStrs, ", ")))
		}
	}
	return strings.Join(conjParts, " ")
}

// formatDrishti สรุปเกณฑ์สายตาของดาว (Special & Mutual)
func formatDrishti(bodies []bodyDrishtiItem) string {
	housePlanets := make(map[int][]string, 12)
	planetHouse := make(map[string]int, len(bodies))
	for _, b := range bodies {
		if b.Code == "ASC" {
			continue
		}
		housePlanets[b.House] = append(housePlanets[b.House], b.Code)
		planetHouse[b.Code] = b.House
	}

	getTargetHouse := func(fromHouse, count int) int {
		return ((fromHouse - 1 + (count - 1)) % 12) + 1
	}

	formatTarget := func(targetH int) string {
		occupants := housePlanets[targetH]
		if len(occupants) > 0 {
			return fmt.Sprintf("H%d(%s)", targetH, strings.Join(occupants, ","))
		}
		return fmt.Sprintf("H%d", targetH)
	}

	var specialDrishtiParts []string
	if h, ok := planetHouse["Sa"]; ok {
		specialDrishtiParts = append(specialDrishtiParts, fmt.Sprintf("Sa->%s,%s,%s",
			formatTarget(getTargetHouse(h, 3)),
			formatTarget(getTargetHouse(h, 7)),
			formatTarget(getTargetHouse(h, 10)),
		))
	}
	if h, ok := planetHouse["Ma"]; ok {
		specialDrishtiParts = append(specialDrishtiParts, fmt.Sprintf("Ma->%s,%s,%s",
			formatTarget(getTargetHouse(h, 4)),
			formatTarget(getTargetHouse(h, 7)),
			formatTarget(getTargetHouse(h, 8)),
		))
	}
	if h, ok := planetHouse["Ju"]; ok {
		specialDrishtiParts = append(specialDrishtiParts, fmt.Sprintf("Ju->%s,%s,%s",
			formatTarget(getTargetHouse(h, 5)),
			formatTarget(getTargetHouse(h, 7)),
			formatTarget(getTargetHouse(h, 9)),
		))
	}

	var mutualParts []string
	checkedPairs := make(map[int]bool)
	for h := 1; h <= 6; h++ {
		oppH := h + 6
		pList1 := housePlanets[h]
		pList2 := housePlanets[oppH]
		if len(pList1) > 0 && len(pList2) > 0 && !checkedPairs[h] {
			checkedPairs[h] = true
			mutualParts = append(mutualParts, fmt.Sprintf("H%d(%s)<->H%d(%s)",
				h, strings.Join(pList1, ","),
				oppH, strings.Join(pList2, ","),
			))
		}
	}

	var drishtiLineParts []string
	if len(specialDrishtiParts) > 0 {
		drishtiLineParts = append(drishtiLineParts, strings.Join(specialDrishtiParts, " | "))
	}
	if len(mutualParts) > 0 {
		drishtiLineParts = append(drishtiLineParts, "Mutual: "+strings.Join(mutualParts, ", "))
	}
	return strings.Join(drishtiLineParts, " | ")
}

// FormatRomanCompactPrompt สร้างข้อความดวงชะตาพระเวทแบบ Roman Ultra-Compact ขั้นสูง (พร้อม Orb, Drishti, และ PD)
func FormatRomanCompactPrompt(resp CalculateResp) string {
	var sb strings.Builder

	if resp.Metadata.BirthDate != "" {
		tzStr := ""
		if resp.Metadata.TimeZone != 0 {
			tzStr = fmt.Sprintf(" (TZ%+g)", resp.Metadata.TimeZone)
		}
		if resp.Metadata.BirthTime != "" {
			sb.WriteString(fmt.Sprintf("[Vedic Sidereal(Lahiri)/WholeSign | Birth:%s %s%s]\n",
				resp.Metadata.BirthDate, resp.Metadata.BirthTime, tzStr))
		} else {
			sb.WriteString(fmt.Sprintf("[Vedic Sidereal(Lahiri)/WholeSign | Birth:%s%s]\n",
				resp.Metadata.BirthDate, tzStr))
		}
	} else {
		sb.WriteString("[Vedic Sidereal(Lahiri)/WholeSign]\n")
	}

	// 1. บรรทัดลัคนาและมุมหลัก (ASC, MC, DSC, IC)
	if len(resp.HousesAndAngles) > 0 {
		asc := resp.HousesAndAngles[0]
		ascSignShort := calculator.ShortSignCode(asc.Position.Sign)
		ascD9Short := calculator.ShortSignCode(asc.Navamsa.Sign)
		ascNakShort := calculator.ShortNakshatraName(asc.Nakshatra)
		ascLordShort := calculator.ShortPlanetCode(asc.Lord)

		sb.WriteString(fmt.Sprintf("ASC:%s %s|D9:%s|Nak:%s(%d,%s)",
			ascSignShort, formatDMSInt(asc.Position),
			ascD9Short, ascNakShort, asc.Pada, ascLordShort,
		))

		if len(resp.HousesAndAngles) >= 4 {
			mc := resp.HousesAndAngles[1]
			dsc := resp.HousesAndAngles[2]
			ic := resp.HousesAndAngles[3]
			sb.WriteString(fmt.Sprintf(" | MC:%s %s(H%d) DSC:%s %s(H%d) IC:%s %s(H%d)",
				calculator.ShortSignCode(mc.Position.Sign), formatDMSInt(mc.Position), mc.House,
				calculator.ShortSignCode(dsc.Position.Sign), formatDMSInt(dsc.Position), dsc.House,
				calculator.ShortSignCode(ic.Position.Sign), formatDMSInt(ic.Position), ic.House,
			))
		}
		sb.WriteString("\n")
	}

	// 2. ตารางดาวเคราะห์ (D1 Sign, DMS, Dignity | House | Nakshatra | D9 Sign & Dignity)
	sb.WriteString("Planets(D1:Sign,DMS,Dig|House|Nak(Pada,Lord)|D9:Sign,Dig):\n")

	// รายชื่อดาวหลักที่พิจารณาในโหราศาสตร์พระเวท
	vedicBodies := []bodyDrishtiItem{}
	if len(resp.HousesAndAngles) > 0 {
		vedicBodies = append(vedicBodies, bodyDrishtiItem{
			Code:  "ASC",
			House: resp.HousesAndAngles[0].House,
			Lon:   resp.HousesAndAngles[0].Degree,
		})
	}

	for _, p := range resp.Planets {
		pShort := calculator.ShortPlanetCode(p.Name)
		d1SignShort := calculator.ShortSignCode(p.Position.Sign)
		d9SignShort := calculator.ShortSignCode(p.Navamsa.Sign)

		// แท็กสถานะพิเศษ [R] = พักร์, [C] = ดับ, [Varg] = วรโคตมนวางศ์
		var flags []string
		if p.IsRetrograde {
			flags = append(flags, "R")
		}
		if p.IsCombust {
			flags = append(flags, "C")
		}
		if p.Position.Sign != "" && p.Position.Sign == p.Navamsa.Sign {
			flags = append(flags, "Varg")
		}
		flagStr := ""
		if len(flags) > 0 {
			flagStr = "[" + strings.Join(flags, ",") + "]"
		}

		d1DigShort := calculator.ShortDignityCode(p.Dignity)

		// คำนวณมาตรฐานดาวในกระดานนวางศ์ (D9)
		d9SignIdx := calculator.GetSignIndexByName(p.Navamsa.Sign)
		d9DigEng, _ := calculator.CalculateSignDignity(p.Name, d9SignIdx)
		d9DigShort := calculator.ShortDignityCode(d9DigEng)

		nakShort := calculator.ShortNakshatraName(p.Nakshatra)
		lordShort := calculator.ShortPlanetCode(p.Lord)

		sb.WriteString(fmt.Sprintf("%s%s:%s %s(%s)|H%d|%s(%d,%s)|D9:%s(%s)\n",
			pShort, flagStr,
			d1SignShort, formatDMSInt(p.Position), d1DigShort,
			p.House,
			nakShort, p.Pada, lordShort,
			d9SignShort, d9DigShort,
		))

		// เฉพาะดาวพระเวท 9 ดวง
		switch pShort {
		case "Su", "Mo", "Ma", "Me", "Ju", "Ve", "Sa", "Ra", "Ke":
			vedicBodies = append(vedicBodies, bodyDrishtiItem{
				Code:  pShort,
				House: p.House,
				Lon:   p.Longitude,
			})
		}
	}

	// 3. ความสัมพันธ์เจ้าเรือนภพทั้ง 12 (Lords)
	var lordParts []string
	for _, ws := range resp.WholeSignHouses {
		signShort := calculator.ShortSignCode(ws.Sign)
		lordShort := calculator.ShortPlanetCode(ws.Lord)
		lordParts = append(lordParts, fmt.Sprintf("%d:%s,%s->H%d",
			ws.House, signShort, lordShort, ws.LordPlacedInHouse,
		))
	}
	if len(lordParts) > 0 {
		sb.WriteString("Lords(H:Sign,Lord->InH): " + strings.Join(lordParts, " ") + "\n")
	}

	// 4. สรุปดาวกุมกันพร้อม Orb และ DeepCombust (Conj)
	conjStr := formatConjunctions(vedicBodies)
	if conjStr != "" {
		sb.WriteString("Conj: " + conjStr + "\n")
	}

	// 5. สรุปเกณฑ์สายตาของดาว (Drishti: Special & Mutual)
	drishtiStr := formatDrishti(vedicBodies)
	if drishtiStr != "" {
		sb.WriteString("Drishti(Special/Mutual): " + drishtiStr + "\n")
	}

	// 6. วิมโสตรีทศา (Vimshottari Dasha พร้อม Pratyantardasha PD)
	d := resp.VimshottariDasha
	if d.MoonNakshatra != "" {
		mdShort := calculator.ShortPlanetCode(d.CurrentDasha.MahadashaLord)
		adShort := calculator.ShortPlanetCode(d.CurrentDasha.AntardashaLord)
		pdShort := calculator.ShortPlanetCode(d.CurrentDasha.PratyantardashaLord)
		if pdShort == "" {
			pdShort = adShort
		}

		mdRange := d.CurrentDasha.MahadashaStart + ".." + d.CurrentDasha.MahadashaEnd
		adRange := d.CurrentDasha.AntardashaStart + ".." + d.CurrentDasha.AntardashaEnd

		sb.WriteString(fmt.Sprintf("Dasha(Moon:%s): MD:%s(%s)/AD:%s(%s)/PD:%s",
			calculator.ShortNakshatraName(d.MoonNakshatra),
			mdShort, mdRange, adShort, adRange, pdShort,
		))
	}

	return strings.TrimSpace(sb.String())
}

// FormatTransitDashaCompactPrompt สร้างข้อความสรุปดาวจรในเดือนและทศา/อนุทศา ณ เวลาดาวจร แบบ Roman Ultra-Compact สำหรับป้อน AI
func FormatTransitDashaCompactPrompt(req TransitReq) string {
	data := calculator.ComputeTransitDashaPromptData(req)

	var sb strings.Builder

	// 1. Header & NatalRef
	sb.WriteString(fmt.Sprintf("[Vedic Transit(Lahiri)/WholeSign | Date:%s]\n", data.DateISO))
	birthRef := ""
	if data.NatalBirthDate != "" {
		tzStr := ""
		if data.NatalTimeZone != 0 {
			tzStr = fmt.Sprintf(" (TZ%+g)", data.NatalTimeZone)
		}
		if data.NatalBirthTime != "" {
			birthRef = fmt.Sprintf(" | Birth:%s %s%s", data.NatalBirthDate, data.NatalBirthTime, tzStr)
		} else {
			birthRef = fmt.Sprintf(" | Birth:%s%s", data.NatalBirthDate, tzStr)
		}
	}
	sb.WriteString(fmt.Sprintf("NatalRef: ASC:%s(H1) | Moon:%s(H%d,%s)%s\n",
		data.NatalAscSignShort, data.NatalMoonSignShort, data.NatalMoonHouse, calculator.ShortNakshatraName(data.NatalMoonNakName),
		birthRef,
	))

	// 2. Dasha@Date
	dasha := data.Dasha
	mdShort := calculator.ShortPlanetCode(dasha.CurrentDasha.MahadashaLord)
	adShort := calculator.ShortPlanetCode(dasha.CurrentDasha.AntardashaLord)
	pdShort := calculator.ShortPlanetCode(dasha.CurrentDasha.PratyantardashaLord)
	if pdShort == "" {
		pdShort = adShort
	}

	sb.WriteString(fmt.Sprintf("Dasha@%s: MD:%s(%s..%s)[nH%d->tH%d] / AD:%s(%s..%s)[nH%d->tH%d] / PD:%s(%s..%s)[nH%d->tH%d]\n",
		data.DateISO,
		mdShort, dasha.CurrentDasha.MahadashaStart, dasha.CurrentDasha.MahadashaEnd, data.NatalHouseByShort[mdShort], data.TransitHouseByShort[mdShort],
		adShort, dasha.CurrentDasha.AntardashaStart, dasha.CurrentDasha.AntardashaEnd, data.NatalHouseByShort[adShort], data.TransitHouseByShort[adShort],
		pdShort, dasha.CurrentDasha.PratyantardashaStart, dasha.CurrentDasha.PratyantardashaEnd, data.NatalHouseByShort[pdShort], data.TransitHouseByShort[pdShort],
	))

	// 3. TransitPlanets table
	sb.WriteString("TransitPlanets(D1:Sign,DMS,Dig|House|Nak(Pada,Lord)|D9:Sign,Dig):\n")
	for _, p := range data.TransitPlanets {
		sb.WriteString(p.Formatted + "\n")
	}

	// 4. Convert VedicBodies for Conj & Drishti
	vedicBodies := make([]bodyDrishtiItem, len(data.VedicBodies))
	for i, b := range data.VedicBodies {
		vedicBodies[i] = bodyDrishtiItem{
			Code:  b.Code,
			House: b.House,
			Lon:   b.Lon,
		}
	}

	// 5. TransitConj
	conjStr := formatConjunctions(vedicBodies)
	if conjStr != "" {
		sb.WriteString("TransitConj: " + conjStr + "\n")
	}

	// 6. TransitDrishti
	drishtiStr := formatDrishti(vedicBodies)
	if drishtiStr != "" {
		sb.WriteString("TransitDrishti(Special/Mutual): " + drishtiStr)
	}

	return strings.TrimSpace(sb.String())
}
