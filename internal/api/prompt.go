package api

import (
	"fmt"
	"strings"

	"lahiri-api/internal/calculator"
)

// formatCompactDMS จัดรูปแบบองศา ลิปดา ฟิลิปดา แบบกระชับ (เช่น 07°49'03.39")
func formatCompactDMS(pos calculator.PositionDetail) string {
	return fmt.Sprintf("%02d°%02d'%05.2f\"", pos.Degree, pos.Minute, pos.Second)
}

// FormatRomanCompactPrompt สร้างข้อความดวงชะตาพระเวทแบบ Roman Ultra-Compact (ละเอียดถึงฟิลิปดา) เพื่อป้อนให้ AI พยากรณ์
func FormatRomanCompactPrompt(resp CalculateResp) string {
	var sb strings.Builder

	sb.WriteString("[Vedic Sidereal(Lahiri)/WholeSign]\n")

	// 1. บรรทัดลัคนาและมุมหลัก (ASC, MC, DSC, IC) ละเอียดถึงฟิลิปดา
	if len(resp.HousesAndAngles) > 0 {
		asc := resp.HousesAndAngles[0]
		ascSignShort := calculator.ShortSignCode(asc.Position.Sign)
		ascD9Short := calculator.ShortSignCode(asc.Navamsa.Sign)
		ascNakShort := calculator.ShortNakshatraName(asc.Nakshatra)
		ascLordShort := calculator.ShortPlanetCode(asc.Lord)

		sb.WriteString(fmt.Sprintf("ASC:%s %s|D9:%s|Nak:%s(%d,%s)",
			ascSignShort, formatCompactDMS(asc.Position),
			ascD9Short, ascNakShort, asc.Pada, ascLordShort,
		))

		if len(resp.HousesAndAngles) >= 4 {
			mc := resp.HousesAndAngles[1]
			dsc := resp.HousesAndAngles[2]
			ic := resp.HousesAndAngles[3]
			sb.WriteString(fmt.Sprintf(" | MC:%s %s(H%d) DSC:%s %s(H%d) IC:%s %s(H%d)",
				calculator.ShortSignCode(mc.Position.Sign), formatCompactDMS(mc.Position), mc.House,
				calculator.ShortSignCode(dsc.Position.Sign), formatCompactDMS(dsc.Position), dsc.House,
				calculator.ShortSignCode(ic.Position.Sign), formatCompactDMS(ic.Position), ic.House,
			))
		}
		sb.WriteString("\n")
	}

	// 2. ตารางดาวเคราะห์ (D1 Sign, DMS ละเอียดถึงฟิลิปดา, Dignity | House | Nakshatra | D9 Sign & Dignity)
	sb.WriteString("Planets(D1:Sign,DMS,Dig|House|Nak(Pada,Lord)|D9:Sign,Dig):\n")

	// เก็บรายชื่อดาวในแต่ละภพเพื่อสรุป Conjunctions
	houseOccupants := make(map[int][]string, 12)
	if len(resp.HousesAndAngles) > 0 && resp.HousesAndAngles[0].House > 0 {
		houseOccupants[resp.HousesAndAngles[0].House] = append(houseOccupants[resp.HousesAndAngles[0].House], "ASC")
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
			d1SignShort, formatCompactDMS(p.Position), d1DigShort,
			p.House,
			nakShort, p.Pada, lordShort,
			d9SignShort, d9DigShort,
		))

		if p.House >= 1 && p.House <= 12 {
			houseOccupants[p.House] = append(houseOccupants[p.House], pShort)
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

	// 4. สรุปดาวกุมกันในภพเดียวกัน (Conj) เฉพาะภพที่มีตั้งแต่ 2 จุดขึ้นไป
	var conjParts []string
	for h := 1; h <= 12; h++ {
		if len(houseOccupants[h]) >= 2 {
			conjParts = append(conjParts, fmt.Sprintf("H%d(%s)", h, strings.Join(houseOccupants[h], "+")))
		}
	}
	if len(conjParts) > 0 {
		sb.WriteString("Conj: " + strings.Join(conjParts, " ") + "\n")
	}

	// 5. วิมโสตรีทศา (Vimshottari Dasha)
	d := resp.VimshottariDasha
	if d.MoonNakshatra != "" {
		balShort := formatCompactBalance(d.BalanceAtBirth)
		mdShort := calculator.ShortPlanetCode(d.CurrentDasha.MahadashaLord)
		adShort := calculator.ShortPlanetCode(d.CurrentDasha.AntardashaLord)
		mdRange := d.CurrentDasha.MahadashaStart + ".." + d.CurrentDasha.MahadashaEnd
		adRange := d.CurrentDasha.AntardashaStart + ".." + d.CurrentDasha.AntardashaEnd

		sb.WriteString(fmt.Sprintf("Dasha(Moon:%s,Bal:%s): Now MD:%s(%s)/AD:%s(%s)",
			calculator.ShortNakshatraName(d.MoonNakshatra), balShort,
			mdShort, mdRange, adShort, adRange,
		))
	}

	return strings.TrimSpace(sb.String())
}

// formatCompactBalance แปลง "Rahu 17 ปี 11 เดือน 11 วัน" เป็น "Ra 17y11m11d"
func formatCompactBalance(balanceThai string) string {
	var lord string
	var y, m, d int
	n, err := fmt.Sscanf(balanceThai, "%s %d ปี %d เดือน %d วัน", &lord, &y, &m, &d)
	if err == nil && n == 4 {
		return fmt.Sprintf("%s %dy%dm%dd", calculator.ShortPlanetCode(lord), y, m, d)
	}
	return balanceThai
}
