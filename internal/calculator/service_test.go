package calculator

import "testing"

func sampleCalculateReq() CalculateReq {
	return CalculateReq{
		Year:        1988,
		Month:       12,
		Day:         23,
		HourThai:    10.35,
		TimeZone:    7.0,
		Latitude:    13.7563,
		Longitude:   100.5018,
		HouseSystem: "",
	}
}

func sampleTransitReq() TransitReq {
	return TransitReq{
		Year:         1988,
		Month:        12,
		Day:          23,
		HourThai:     10.35,
		TimeZone:     7.0,
		Latitude:     13.7563,
		Longitude:    100.5018,
		HouseSystem:  "W",
		TransitYear:  2026,
		TransitMonth: 12,
	}
}

func TestSiderealWrappers(t *testing.T) {
	jd, ayanamsa := ComputeJulianDayAndAyanamsa(1988, 12, 23, 10.35, 7.0)

	planets := CalculateSiderealPlanets(jd)
	if len(planets) < 9 {
		t.Fatalf("expected at least 9 planets, got %d", len(planets))
	}

	houses := CalculateSiderealHouses(jd, 13.7563, 100.5018, int('W'), ayanamsa)
	if len(houses) != 4 || houses[0].Name != "Ascendant (ASC)" {
		t.Fatalf("expected 4 cardinal angles starting with Ascendant (ASC), got %+v", houses)
	}
}

func TestComputeHoroscope(t *testing.T) {
	resp := ComputeHoroscope(sampleCalculateReq())

	if resp.Metadata.BirthDate != "1988-12-23" || resp.Metadata.BirthTime != "10:21" {
		t.Errorf("unexpected birth metadata: date=%q time=%q", resp.Metadata.BirthDate, resp.Metadata.BirthTime)
	}
	if resp.Metadata.AyanamsaName != "Lahiri" {
		t.Errorf("expected AyanamsaName 'Lahiri', got %q", resp.Metadata.AyanamsaName)
	}
	if len(resp.Planets) == 0 || len(resp.HousesAndAngles) != 4 || len(resp.WholeSignHouses) != 12 {
		t.Fatalf("unexpected horoscope slices: planets=%d houses=%d wholeSign=%d",
			len(resp.Planets), len(resp.HousesAndAngles), len(resp.WholeSignHouses))
	}
	for _, ws := range resp.WholeSignHouses {
		if ws.LordPlacedInHouse < 1 || ws.LordPlacedInHouse > 12 {
			t.Errorf("house %d has invalid LordPlacedInHouse=%d", ws.House, ws.LordPlacedInHouse)
		}
	}
	if resp.VimshottariDasha.MoonNakshatra != "Ardra" {
		t.Errorf("expected MoonNakshatra 'Ardra', got %q", resp.VimshottariDasha.MoonNakshatra)
	}
}

func TestComputeDashaVariants(t *testing.T) {
	t.Run("CustomTransitYearAndMonth", func(t *testing.T) {
		req := sampleTransitReq()
		d := ComputeDasha(req)
		if d.TargetDate != "23/12/2026" {
			t.Errorf("expected TargetDate '23/12/2026', got %q", d.TargetDate)
		}
		if d.CurrentDasha.MahadashaLord != "Saturn" || d.CurrentDasha.AntardashaLord != "Mercury" {
			t.Errorf("unexpected dasha at 2026-12-23: %+v", d.CurrentDasha)
		}
	})

	t.Run("CustomTransitYearWithoutMonthDefaultsToBirthMonth", func(t *testing.T) {
		req := sampleTransitReq()
		req.TransitMonth = 0
		d := ComputeDasha(req)
		if d.TargetDate != "23/12/2026" {
			t.Errorf("expected TargetDate '23/12/2026', got %q", d.TargetDate)
		}
	})

	t.Run("NoCustomTargetUsesCurrentTime", func(t *testing.T) {
		req := sampleTransitReq()
		req.TransitYear = 0
		req.TransitMonth = 0
		d := ComputeDasha(req)
		if len(d.MahadashaList) != 9 || d.CurrentDasha.MahadashaLord == "" {
			t.Errorf("expected populated MahadashaList and CurrentDasha, got %+v", d)
		}
	})
}

func TestComputeTransitServices(t *testing.T) {
	req := sampleTransitReq()

	ascIdx := ComputeNatalAscSignIdx(req)
	if ascIdx != 10 { // Aquarius
		t.Errorf("expected natal ASC Aquarius (10), got %d", ascIdx)
	}

	annual := ComputeAnnualTransit(req)
	if len(annual) != 12 {
		t.Fatalf("expected 12 monthly transits, got %d", len(annual))
	}

	monthOnly := ComputeMonthTransit(req)
	if monthOnly.Month != 12 || monthOnly.Date != "23/12/2026" || monthOnly.Dasha != nil {
		t.Errorf("unexpected ComputeMonthTransit result: month=%d date=%q dasha=%v",
			monthOnly.Month, monthOnly.Date, monthOnly.Dasha)
	}
	if len(monthOnly.Planets) != 9 {
		t.Errorf("expected 9 Navagraha planets, got %d", len(monthOnly.Planets))
	}

	monthWithDasha := ComputeMonthTransitWithDasha(req)
	if monthWithDasha.Dasha == nil || monthWithDasha.Dasha.MahadashaLord != "Saturn" {
		t.Errorf("unexpected ComputeMonthTransitWithDasha.Dasha: %+v", monthWithDasha.Dasha)
	}

	promptData := ComputeTransitDashaPromptData(req)
	if promptData.DateISO != "2026-12-23" || promptData.NatalBirthDate != "1988-12-23" || promptData.NatalBirthTime != "10:21" {
		t.Errorf("unexpected promptData dates: %+v", promptData)
	}
	if promptData.NatalAscSignShort != "Aq" || promptData.NatalMoonSignShort != "Ge" || promptData.NatalMoonNakName != "Ardra" {
		t.Errorf("unexpected NatalRef: asc=%s moon=%s nak=%s",
			promptData.NatalAscSignShort, promptData.NatalMoonSignShort, promptData.NatalMoonNakName)
	}
	if len(promptData.TransitPlanets) != 9 || len(promptData.VedicBodies) != 9 {
		t.Errorf("expected 9 transit bodies, got %d / %d", len(promptData.TransitPlanets), len(promptData.VedicBodies))
	}
	if promptData.NatalHouseByShort["Su"] == 0 || promptData.TransitHouseByShort["Su"] == 0 {
		t.Errorf("expected non-zero NatalHouseByShort/TransitHouseByShort entries")
	}
}
