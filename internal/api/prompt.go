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

// FormatRomanCompactPrompt สร้างข้อความดวงชะตาพระเวทแบบ Roman Ultra-Compact ขั้นสูง (พร้อม Orb, Drishti, และ PD)
func FormatRomanCompactPrompt(resp CalculateResp) string {
	var sb strings.Builder

	sb.WriteString("[Vedic Sidereal(Lahiri)/WholeSign]\n")

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
	houseBodies := make(map[int][]bodyDrishtiItem, 12)
	for _, b := range vedicBodies {
		if b.House >= 1 && b.House <= 12 {
			houseBodies[b.House] = append(houseBodies[b.House], b)
		}
	}

	var conjParts []string
	for h := 1; h <= 12; h++ {
		bodies := houseBodies[h]
		if len(bodies) < 2 {
			continue
		}

		// กรณีมี 2 จุด เช่น H1(ASC+Ra[orb 10°53'])
		if len(bodies) == 2 {
			b1, b2 := bodies[0], bodies[1]
			orb := calculator.AngularDistance(b1.Lon, b2.Lon)
			extra := ""
			// เช็คว่ามีดวงอาทิตย์กุมและดับลึกหรือไม่ (<= 5 องศา)
			if (b1.Code == "Su" || b2.Code == "Su") && orb <= 5.0 {
				extra = ",DeepCombust"
			}
			conjParts = append(conjParts, fmt.Sprintf("H%d(%s+%s[orb %s%s])",
				h, b1.Code, b2.Code, formatDMSDiff(orb), extra,
			))
		} else {
			// กรณีมีมากกว่า 2 จุด เช่น H11: คำนวณคู่เทียบกับจุดแรก (หรือดวงอาทิตย์ถ้ามี)
			// ถ้ามีอาทิตย์ ให้ใช้อาทิตย์เป็นแกนเทียบ
			sunIdx := -1
			for idx, b := range bodies {
				if b.Code == "Su" {
					sunIdx = idx
					break
				}
			}

			var pairStrs []string
			if sunIdx >= 0 {
				sBody := bodies[sunIdx]
				for idx, b := range bodies {
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
				// ไม่มีอาทิตย์ ให้จับคู่ตัวแรกกับตัวถัดๆ ไป
				base := bodies[0]
				for idx := 1; idx < len(bodies); idx++ {
					b := bodies[idx]
					orb := calculator.AngularDistance(base.Lon, b.Lon)
					pairStrs = append(pairStrs, fmt.Sprintf("%s+%s[orb %s]",
						base.Code, b.Code, formatDMSDiff(orb),
					))
				}
			}
			conjParts = append(conjParts, fmt.Sprintf("H%d(%s)", h, strings.Join(pairStrs, ", ")))
		}
	}
	if len(conjParts) > 0 {
		sb.WriteString("Conj: " + strings.Join(conjParts, " ") + "\n")
	}

	// 5. สรุปเกณฑ์สายตาของดาว (Drishti: Special & Mutual)
	// ดาวในแต่ละภพ (เฉพาะดาว 9 ดวง ไม่รวม ASC)
	housePlanets := make(map[int][]string, 12)
	planetHouse := make(map[string]int)
	for _, b := range vedicBodies {
		if b.Code == "ASC" {
			continue
		}
		housePlanets[b.House] = append(housePlanets[b.House], b.Code)
		planetHouse[b.Code] = b.House
	}

	// ฟังก์ชันคำนวณภพเป้าหมาย (1-12) ตามระยะนับ
	getTargetHouse := func(fromHouse, count int) int {
		return ((fromHouse - 1 + (count - 1)) % 12) + 1
	}

	// ฟังก์ชันจัดรูปแบบเป้าหมาย เช่น H5(Mo) หรือ H8
	formatTarget := func(targetH int) string {
		occupants := housePlanets[targetH]
		if len(occupants) > 0 {
			return fmt.Sprintf("H%d(%s)", targetH, strings.Join(occupants, ","))
		}
		return fmt.Sprintf("H%d", targetH)
	}

	var specialDrishtiParts []string

	// เสาร์ (Sa): ส่องเกณฑ์ 3, 7, 10
	if h, ok := planetHouse["Sa"]; ok {
		h3 := getTargetHouse(h, 3)
		h7 := getTargetHouse(h, 7)
		h10 := getTargetHouse(h, 10)
		specialDrishtiParts = append(specialDrishtiParts, fmt.Sprintf("Sa->%s,%s,%s",
			formatTarget(h3), formatTarget(h7), formatTarget(h10),
		))
	}

	// อังคาร (Ma): ส่องเกณฑ์ 4, 7, 8
	if h, ok := planetHouse["Ma"]; ok {
		h4 := getTargetHouse(h, 4)
		h7 := getTargetHouse(h, 7)
		h8 := getTargetHouse(h, 8)
		specialDrishtiParts = append(specialDrishtiParts, fmt.Sprintf("Ma->%s,%s,%s",
			formatTarget(h4), formatTarget(h7), formatTarget(h8),
		))
	}

	// พฤหัสบดี (Ju): ส่องเกณฑ์ 5, 7, 9
	if h, ok := planetHouse["Ju"]; ok {
		h5 := getTargetHouse(h, 5)
		h7 := getTargetHouse(h, 7)
		h9 := getTargetHouse(h, 9)
		specialDrishtiParts = append(specialDrishtiParts, fmt.Sprintf("Ju->%s,%s,%s",
			formatTarget(h5), formatTarget(h7), formatTarget(h9),
		))
	}

	// การเล็งกันข้ามภพ 180 องศา (Mutual Aspects: ภพ X <-> ภพ X+6)
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
	if len(drishtiLineParts) > 0 {
		sb.WriteString("Drishti(Special/Mutual): " + strings.Join(drishtiLineParts, " | ") + "\n")
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
