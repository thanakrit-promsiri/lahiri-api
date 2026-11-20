package calculator

import (
	"math"
	"strings"
	"testing"
)

func sampleTriadPersons() []GephiPersonReq {
	return []GephiPersonReq{
		{
			ID:   "P1",
			Name: "Person One",
			CalculateReq: CalculateReq{
				Year:        1988,
				Month:       12,
				Day:         23,
				HourThai:    10.35,
				TimeZone:    7.0,
				Latitude:    13.7563,
				Longitude:   100.5018,
				HouseSystem: "W",
			},
		},
		{
			ID:   "P2",
			Name: "Person Two",
			CalculateReq: CalculateReq{
				Year:        1990,
				Month:       5,
				Day:         15,
				HourThai:    14.5,
				TimeZone:    7.0,
				Latitude:    13.7563,
				Longitude:   100.5018,
				HouseSystem: "W",
			},
		},
		{
			ID:   "P3",
			Name: "Person Three",
			CalculateReq: CalculateReq{
				Year:        1992,
				Month:       8,
				Day:         9,
				HourThai:    8.25,
				TimeZone:    7.0,
				Latitude:    18.7883,
				Longitude:   98.9853,
				HouseSystem: "W",
			},
		},
	}
}

func TestGephiReqValidation(t *testing.T) {
	valid := sampleTriadPersons()

	if err := (GephiReq{Persons: valid[:1]}).Validate(); err == nil {
		t.Error("expected error when fewer than 2 persons provided")
	}

	if err := (GephiReq{Persons: valid[:2]}).Validate(); err != nil {
		t.Errorf("expected 2 valid persons to pass validation, got %v", err)
	}

	if err := (GephiReq{Persons: valid}).Validate(); err != nil {
		t.Errorf("expected 3 valid persons to pass validation, got %v", err)
	}

	// ชน ID อัตโนมัติ P1 กับ ID ที่กำหนดเอง "P1"
	dupAuto := []GephiPersonReq{
		{CalculateReq: valid[0].CalculateReq},
		{ID: "P1", CalculateReq: valid[1].CalculateReq},
	}
	if err := (GephiReq{Persons: dupAuto}).Validate(); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("expected duplicate person id error, got %v", err)
	}

	// ID มีช่องว่าง
	spaceID := []GephiPersonReq{
		{ID: "P 1", CalculateReq: valid[0].CalculateReq},
		{ID: "P2", CalculateReq: valid[1].CalculateReq},
	}
	if err := (GephiReq{Persons: spaceID}).Validate(); err == nil {
		t.Error("expected error for whitespace in person id")
	}

	// ข้อมูลวันเกิดของคนที่ 2 ไม่ถูกต้อง
	badChild := []GephiPersonReq{
		valid[0],
		{ID: "P2", CalculateReq: CalculateReq{Year: 1990, Month: 13, Day: 15, HourThai: 10, TimeZone: 7}},
	}
	if err := (GephiReq{Persons: badChild}).Validate(); err == nil || !strings.Contains(err.Error(), "person 2 (P2)") {
		t.Errorf("expected wrapped child validation error, got %v", err)
	}
}

func TestGephi4PersonsJupiterConjunction(t *testing.T) {
	req := GephiReq{
		Persons: []GephiPersonReq{
			{
				ID: "P1", Name: "Male 1988",
				CalculateReq: CalculateReq{Year: 1988, Month: 12, Day: 23, HourThai: 10.35, TimeZone: 7, Latitude: 17.598, Longitude: 101.721, HouseSystem: "W"},
			},
			{
				ID: "P2", Name: "Female 1991",
				CalculateReq: CalculateReq{Year: 1991, Month: 4, Day: 19, HourThai: 2.15, TimeZone: 7, Latitude: 16.468, Longitude: 102.83, HouseSystem: "W"},
			},
			{
				ID: "P3", Name: "Female 2003",
				CalculateReq: CalculateReq{Year: 2003, Month: 4, Day: 25, HourThai: 9.7, TimeZone: 7, Latitude: 16.463, Longitude: 99.528, HouseSystem: "W"},
			},
			{
				ID: "P4", Name: "Female 1990",
				CalculateReq: CalculateReq{Year: 1990, Month: 9, Day: 26, HourThai: 4.816667, TimeZone: 7, Latitude: 13.909, Longitude: 100.618, HouseSystem: "W"},
			},
		},
	}

	graph := ComputeGephiGraph(req)
	var foundP2P3, foundP2P4, foundP3P4 bool
	for _, e := range graph.Edges {
		if e.Relation != EdgeRelYuti {
			continue
		}
		pair := e.Source + "-" + e.Target
		if pair == "P2_Jupiter-P3_Jupiter" || pair == "P3_Jupiter-P2_Jupiter" {
			foundP2P3 = true
			if e.Rasi != 4 || e.RasiName != "Cancer" || e.RelationThai != "กุม (Yuti)" {
				t.Errorf("unexpected P2-P3 Jupiter edge attributes: %+v", e)
			}
		}
		if pair == "P2_Jupiter-P4_Jupiter" || pair == "P4_Jupiter-P2_Jupiter" {
			foundP2P4 = true
		}
		if pair == "P3_Jupiter-P4_Jupiter" || pair == "P4_Jupiter-P3_Jupiter" {
			foundP3P4 = true
		}
	}

	if !foundP2P3 || !foundP2P4 || !foundP3P4 {
		t.Fatalf("expected all 3 Jupiter-Jupiter conjunction edges in Cancer, got P2-P3=%v, P2-P4=%v, P3-P4=%v",
			foundP2P3, foundP2P4, foundP3P4)
	}
}

func TestComputeGephiGraphTriad(t *testing.T) {
	req := GephiReq{Persons: sampleTriadPersons()}
	graph := ComputeGephiGraph(req)

	if len(graph.Nodes) != 30 {
		t.Fatalf("expected 30 nodes for 3 persons, got %d", len(graph.Nodes))
	}

	seenNodeIDs := make(map[string]bool, 30)
	for _, n := range graph.Nodes {
		if seenNodeIDs[n.ID] {
			t.Errorf("duplicate node id: %s", n.ID)
		}
		seenNodeIDs[n.ID] = true

		if n.Rasi < 1 || n.Rasi > 12 {
			t.Errorf("node %s invalid rasi: %d", n.ID, n.Rasi)
		}
		if n.Degree < 0 || n.Degree >= 30 {
			t.Errorf("node %s invalid degree: %f", n.ID, n.Degree)
		}
		if n.Nakshatra < 1 || n.Nakshatra > 27 {
			t.Errorf("node %s invalid nakshatra: %d", n.ID, n.Nakshatra)
		}
		if n.Dignity == "" || n.Bala <= 0 || n.NodeSize <= 0 {
			t.Errorf("node %s invalid dignity/bala/size: %q %f %f", n.ID, n.Dignity, n.Bala, n.NodeSize)
		}
		if n.House < 1 || n.House > 12 {
			t.Errorf("node %s invalid house: %d", n.ID, n.House)
		}
	}

	if !seenNodeIDs["P1_Lagna"] || !seenNodeIDs["P2_Sun"] || !seenNodeIDs["P3_Saturn"] {
		t.Errorf("missing expected standard node IDs in graph")
	}

	if len(graph.Edges) < 50 {
		t.Fatalf("expected substantial edges for 30 nodes (inter+intra), got %d", len(graph.Edges))
	}

	var hasDirected, hasUndirected, hasInter, hasIntra bool
	seenEdgeIDs := make(map[string]bool, len(graph.Edges))
	for _, e := range graph.Edges {
		if seenEdgeIDs[e.ID] {
			t.Errorf("duplicate edge id: %s", e.ID)
		}
		seenEdgeIDs[e.ID] = true

		if e.Source == e.Target {
			t.Errorf("unexpected self-loop edge on %s", e.Source)
		}
		if e.Type == EdgeDirDirected {
			hasDirected = true
		} else if e.Type == EdgeDirUndirected {
			hasUndirected = true
		}
		if e.Scope == EdgeScopeInter {
			hasInter = true
		} else if e.Scope == EdgeScopeIntra {
			hasIntra = true
		}
		if e.Weight <= 0 {
			t.Errorf("edge %s has non-positive weight: %f", e.ID, e.Weight)
		}
	}

	if !hasDirected || !hasUndirected || !hasInter || !hasIntra {
		t.Errorf("expected mixed directed/undirected and inter/intra edges, got directed=%v undirected=%v inter=%v intra=%v",
			hasDirected, hasUndirected, hasInter, hasIntra)
	}

	// P0: ตรวจสอบ Birth Metadata ในโหนด
	if graph.Nodes[0].BirthDate == "" || graph.Nodes[0].BirthTime == "" || graph.Nodes[0].JulianDay <= 0 {
		t.Errorf("expected birth metadata populated in node, got %+v", graph.Nodes[0])
	}

	// P1: ตรวจสอบ Nakshatra Compatibility Attributes (Gana, Yoni, Nadi)
	if graph.Nodes[0].NakshatraGana == "" || graph.Nodes[0].NakshatraYoni == "" || graph.Nodes[0].NakshatraNadi == "" {
		t.Errorf("expected nakshatra attributes in node, got gana=%q yoni=%q nadi=%q",
			graph.Nodes[0].NakshatraGana, graph.Nodes[0].NakshatraYoni, graph.Nodes[0].NakshatraNadi)
	}

	// P2: ตรวจสอบ Vimshottari Dasha Active Timing
	var hasActiveDasha bool
	for _, n := range graph.Nodes {
		if n.IsDashaActive {
			hasActiveDasha = true
			if n.DashaRole == "" || n.DashaPeriod == "" {
				t.Errorf("node %s has active dasha but empty role/period", n.ID)
			}
		}
	}
	if !hasActiveDasha {
		t.Error("expected at least one active dasha planet per person")
	}

	// P2: ตรวจสอบ PairwiseSynastryMetrics
	if len(graph.PairwiseMetrics) != 3 {
		t.Fatalf("expected 3 pairwise synastry metrics for 3 persons, got %d", len(graph.PairwiseMetrics))
	}
	m := graph.PairwiseMetrics[0]
	if m.PersonA != "P1" || m.PersonB != "P2" || m.AshtakootaScore <= 0 || m.AshtakootaVerdict == "" {
		t.Errorf("unexpected pairwise metric values: %+v", m)
	}

	// ตรวจสอบว่าเส้นทั้งหมดเป็น D1 (ราศีจักร) เท่านั้น และไม่มี D9 Yuti
	var hasNakshatraCompat bool
	for _, e := range graph.Edges {
		if e.Layer != "D1" {
			t.Errorf("expected D1 edge only, got edge %s layer=%s", e.ID, e.Layer)
		}
		if e.Relation == EdgeRelD9Yuti {
			t.Errorf("expected no D9 Yuti edges in D1 graph, got %s", e.ID)
		}
		if e.Relation == EdgeRelNakshatra {
			hasNakshatraCompat = true
			if e.TaraName == "" || e.TaraScore == 0 {
				t.Errorf("expected tara attributes on nakshatra edge, got %+v", e)
			}
		}
	}
	if !hasNakshatraCompat {
		t.Error("expected cross-chart Nakshatra Compatibility edges between Moons")
	}
}

func TestGephiYutiEdgesAndOrbWeights(t *testing.T) {
	nodes := []GephiNode{
		{ID: "P1_Moon", PersonID: "P1", Body: "Moon", Rasi: 2, RasiName: "Taurus", RasiThai: "พฤษภ", Degree: 5.0},
		{ID: "P2_Venus", PersonID: "P2", Body: "Venus", Rasi: 2, RasiName: "Taurus", RasiThai: "พฤษภ", Degree: 6.5},      // orb 1.5 -> W = 10
		{ID: "P3_Mercury", PersonID: "P3", Body: "Mercury", Rasi: 2, RasiName: "Taurus", RasiThai: "พฤษภ", Degree: 21.5}, // orb vs P1 = 16.5 -> W = 6.5
		{ID: "P1_Sun", PersonID: "P1", Body: "Sun", Rasi: 3, RasiName: "Gemini", RasiThai: "เมถุน", Degree: 10.0},
	}

	edges := appendYutiEdges(nil, nodes)
	if len(edges) != 3 {
		t.Fatalf("expected 3 Yuti edges in Rasi 2, got %d", len(edges))
	}

	if edges[0].Source != "P1_Moon" || edges[0].Target != "P2_Venus" || edges[0].Weight != 10.0 || edges[0].Type != EdgeDirUndirected {
		t.Errorf("unexpected tight Yuti edge: %+v", edges[0])
	}
	if edges[0].Rasi != 2 || edges[0].RasiName != "Taurus" || edges[0].RasiThai != "พฤษภ" {
		t.Errorf("expected Rasi 2 / Taurus / พฤษภ on Yuti edge, got: %+v", edges[0])
	}
	if edges[0].RelationThai != "กุม (Yuti)" {
		t.Errorf("expected RelationThai 'กุม (Yuti)', got: %s", edges[0].RelationThai)
	}
	if !strings.Contains(edges[0].Label, "Taurus") {
		t.Errorf("expected label to contain sign Taurus, got %s", edges[0].Label)
	}
	if edges[0].ScoringFormula == "" || edges[0].BaseWeight <= 0 {
		t.Errorf("expected scoring formula and base weight on Yuti edge: %+v", edges[0])
	}
	if edges[1].Source != "P1_Moon" || edges[1].Target != "P3_Mercury" || math.Abs(edges[1].Weight-6.5) > 0.01 {
		t.Errorf("unexpected mid-orb Yuti edge (expected W=6.5): %+v", edges[1])
	}

	if got := calcOrbWeight(30.0, 10.0, 3.0, 3.0); got != 3.0 {
		t.Errorf("expected minW 3.0 at 30 deg orb, got %f", got)
	}
}

func TestGephiDrishtiAspectsAndWeights(t *testing.T) {
	nodes := []GephiNode{
		{ID: "P1_Mars", PersonID: "P1", Body: "Mars", Rasi: 1, Degree: 10.0},      // 4->Cn(4), 7->Li(7), 8->Sc(8)
		{ID: "P1_Jupiter", PersonID: "P1", Body: "Jupiter", Rasi: 1, Degree: 10.0}, // 5->Le(5), 7->Li(7), 9->Sg(9)
		{ID: "P1_Saturn", PersonID: "P1", Body: "Saturn", Rasi: 1, Degree: 10.0},   // 3->Ge(3), 7->Li(7), 10->Cp(10)
		{ID: "P1_Rahu", PersonID: "P1", Body: "Rahu", Rasi: 1, Degree: 10.0},       // 5, 7, 9
		{ID: "P1_Ketu", PersonID: "P1", Body: "Ketu", Rasi: 1, Degree: 10.0},       // 5, 7, 9
		{ID: "P1_Sun", PersonID: "P1", Body: "Sun", Rasi: 1, Degree: 10.0},         // 7 only
		{ID: "P1_Lagna", PersonID: "P1", Body: "Lagna", Rasi: 1, Degree: 10.0},     // sends NO drishti
		{ID: "P2_Gem", PersonID: "P2", Body: "Moon", Rasi: 3, Degree: 11.0},        // orb 1 -> W=7.0
		{ID: "P2_Lagna", PersonID: "P2", Body: "Lagna", Rasi: 4, Degree: 26.5},     // orb 16.5 -> W=5.5
		{ID: "P2_Leo", PersonID: "P2", Body: "Sun", Rasi: 5, Degree: 10.0},
		{ID: "P2_Lib", PersonID: "P2", Body: "Venus", Rasi: 7, Degree: 10.0},
		{ID: "P2_Sco", PersonID: "P2", Body: "Mercury", Rasi: 8, Degree: 10.0},
		{ID: "P2_Sag", PersonID: "P2", Body: "Jupiter", Rasi: 9, Degree: 10.0},
		{ID: "P2_Cap", PersonID: "P2", Body: "Saturn", Rasi: 10, Degree: 10.0},
	}

	edges := appendDrishtiEdges(nil, nodes)

	aspectHit := make(map[string]GephiEdge)
	for _, e := range edges {
		if e.Type != EdgeDirDirected || e.Relation != EdgeRelDrishti {
			t.Errorf("unexpected edge metadata in drishti: %+v", e)
		}
		if strings.HasPrefix(e.Source, "P1_Lagna") {
			t.Errorf("Lagna must not cast Drishti: %+v", e)
		}
		aspectHit[e.Source+"->"+e.Target] = e
	}

	// Mars -> Cn(4), Li(7), Sc(8)
	if _, ok := aspectHit["P1_Mars->P2_Lagna"]; !ok {
		t.Error("expected Mars 4th aspect on P2_Lagna in Cancer")
	}
	if _, ok := aspectHit["P1_Mars->P2_Lib"]; !ok {
		t.Error("expected Mars 7th aspect on Libra")
	}
	if _, ok := aspectHit["P1_Mars->P2_Sco"]; !ok {
		t.Error("expected Mars 8th aspect on Scorpio")
	}
	if _, ok := aspectHit["P1_Mars->P2_Leo"]; ok {
		t.Error("Mars should not aspect Leo (5th)")
	}

	// Saturn -> Ge(3), Li(7), Cp(10)
	if e, ok := aspectHit["P1_Saturn->P2_Gem"]; !ok || e.Weight != 7.0 || e.DrishtiHouse != 3 {
		t.Errorf("expected Saturn 3rd aspect with W=7.0, got %+v (ok=%v)", e, ok)
	}
	if _, ok := aspectHit["P1_Saturn->P2_Cap"]; !ok {
		t.Error("expected Saturn 10th aspect on Capricorn")
	}

	// Jupiter, Rahu, Ketu -> Le(5), Li(7), Sg(9)
	for _, src := range []string{"P1_Jupiter", "P1_Rahu", "P1_Ketu"} {
		if _, ok := aspectHit[src+"->P2_Leo"]; !ok {
			t.Errorf("expected %s 5th aspect on Leo", src)
		}
		if _, ok := aspectHit[src+"->P2_Sag"]; !ok {
			t.Errorf("expected %s 9th aspect on Sagittarius", src)
		}
	}

	// Sun -> 7th only
	if _, ok := aspectHit["P1_Sun->P2_Lib"]; !ok {
		t.Error("expected Sun 7th aspect on Libra")
	}
	if _, ok := aspectHit["P1_Sun->P2_Leo"]; ok {
		t.Error("Sun should not cast 5th aspect")
	}

	// Weight check for orb 16.5 -> 5.5
	if e := aspectHit["P1_Mars->P2_Lagna"]; math.Abs(e.Weight-5.5) > 0.01 {
		t.Errorf("expected orb 16.5 drishti weight 5.5, got %f", e.Weight)
	}
}

func TestGephiParivartanaEdges(t *testing.T) {
	nodes := []GephiNode{
		{ID: "P1_Sun", PersonID: "P1", Body: "Sun", Rasi: 4},    // Cancer (Moon)
		{ID: "P3_Moon", PersonID: "P3", Body: "Moon", Rasi: 5},  // Leo (Sun) -> Inter Parivartana!
		{ID: "P1_Mars", PersonID: "P1", Body: "Mars", Rasi: 2},  // Taurus (Venus)
		{ID: "P1_Venus", PersonID: "P1", Body: "Venus", Rasi: 1}, // Aries (Mars) -> Intra Parivartana!
		{ID: "P2_Mars", PersonID: "P2", Body: "Mars", Rasi: 8},  // Scorpio (Mars) -> Same body, not Parivartana
		{ID: "P2_Rahu", PersonID: "P2", Body: "Rahu", Rasi: 5},  // Non-sign ruler
	}

	edges := appendParivartanaEdges(nil, nodes)
	if len(edges) != 2 {
		t.Fatalf("expected exactly 2 Parivartana edges, got %d: %+v", len(edges), edges)
	}

	if edges[0].Source != "P1_Sun" || edges[0].Target != "P3_Moon" || edges[0].Weight != 8.0 || edges[0].Type != EdgeDirUndirected || edges[0].Scope != EdgeScopeInter || edges[0].HasOrb {
		t.Errorf("unexpected inter Parivartana edge: %+v", edges[0])
	}
	if edges[1].Source != "P1_Mars" || edges[1].Target != "P1_Venus" || edges[1].Weight != 8.0 || edges[1].Scope != EdgeScopeIntra {
		t.Errorf("unexpected intra Parivartana edge: %+v", edges[1])
	}
}

func TestGephiGrahaMaitriClassification(t *testing.T) {
	cases := []struct {
		name      string
		u, v      GephiNode
		wantType  string
		wantColor string
		wantPol   int
	}{
		{
			name:      "JupiterToMoon_HarmoniousBlue",
			u:         GephiNode{Body: "Jupiter", Rasi: 1},
			v:         GephiNode{Body: "Moon", Rasi: 4},
			wantType:  EdgeTypeHarmonious,
			wantColor: "Blue",
			wantPol:   1,
		},
		{
			name:      "SaturnToSun_ChallengingRed",
			u:         GephiNode{Body: "Saturn", Rasi: 10},
			v:         GephiNode{Body: "Sun", Rasi: 5},
			wantType:  EdgeTypeChallenging,
			wantColor: "Red",
			wantPol:   -1,
		},
		{
			name:      "MarsToVenus_NeutralGray",
			u:         GephiNode{Body: "Mars", Rasi: 1},
			v:         GephiNode{Body: "Venus", Rasi: 7},
			wantType:  EdgeTypeNeutral,
			wantColor: "Gray",
			wantPol:   0,
		},
		{
			name:      "MoonToMercury_NeutralGray",
			u:         GephiNode{Body: "Moon", Rasi: 4},
			v:         GephiNode{Body: "Mercury", Rasi: 3},
			wantType:  EdgeTypeNeutral,
			wantColor: "Gray",
			wantPol:   0,
		},
		{
			name:      "RahuVenus_Harmonious",
			u:         GephiNode{Body: "Rahu", Rasi: 2},
			v:         GephiNode{Body: "Venus", Rasi: 2},
			wantType:  EdgeTypeHarmonious,
			wantColor: "Blue",
			wantPol:   1,
		},
		{
			name:      "RahuJupiter_Challenging",
			u:         GephiNode{Body: "Rahu", Rasi: 9},
			v:         GephiNode{Body: "Jupiter", Rasi: 9},
			wantType:  EdgeTypeChallenging,
			wantColor: "Red",
			wantPol:   -1,
		},
		{
			name:      "KetuJupiter_Harmonious",
			u:         GephiNode{Body: "Ketu", Rasi: 8},
			v:         GephiNode{Body: "Jupiter", Rasi: 12},
			wantType:  EdgeTypeHarmonious,
			wantColor: "Blue",
			wantPol:   1,
		},
		{
			name:      "KetuSun_Challenging",
			u:         GephiNode{Body: "Ketu", Rasi: 5},
			v:         GephiNode{Body: "Sun", Rasi: 5},
			wantType:  EdgeTypeChallenging,
			wantColor: "Red",
			wantPol:   -1,
		},
		{
			name:      "RahuKetu_Neutral",
			u:         GephiNode{Body: "Rahu", Rasi: 1},
			v:         GephiNode{Body: "Ketu", Rasi: 7},
			wantType:  EdgeTypeNeutral,
			wantColor: "Gray",
			wantPol:   0,
		},
		{
			name:      "JupiterToCancerLagna_HarmoniousViaLagneshMoon",
			u:         GephiNode{Body: "Jupiter", Rasi: 12},
			v:         GephiNode{Body: "Lagna", Rasi: 4}, // Moon
			wantType:  EdgeTypeHarmonious,
			wantColor: "Blue",
			wantPol:   1,
		},
		{
			name:      "SaturnToLeoLagna_ChallengingViaLagneshSun",
			u:         GephiNode{Body: "Saturn", Rasi: 11},
			v:         GephiNode{Body: "Lagna", Rasi: 5}, // Sun
			wantType:  EdgeTypeChallenging,
			wantColor: "Red",
			wantPol:   -1,
		},
		{
			name:      "SameBeneficMoonMoon_Harmonious",
			u:         GephiNode{Body: "Moon", Rasi: 4},
			v:         GephiNode{Body: "Moon", Rasi: 4},
			wantType:  EdgeTypeHarmonious,
			wantColor: "Blue",
			wantPol:   1,
		},
		{
			name:      "SameMaleficSaturnSaturn_Neutral",
			u:         GephiNode{Body: "Saturn", Rasi: 10},
			v:         GephiNode{Body: "Saturn", Rasi: 10},
			wantType:  EdgeTypeNeutral,
			wantColor: "Gray",
			wantPol:   0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotType, gotColor, gotPol, _ := classifyGrahaMaitri(tc.u, tc.v)
			if gotType != tc.wantType || gotColor != tc.wantColor || gotPol != tc.wantPol {
				t.Errorf("got (%s, %s, %d); want (%s, %s, %d)",
					gotType, gotColor, gotPol, tc.wantType, tc.wantColor, tc.wantPol)
			}
		})
	}
}

func TestGephiHelpers(t *testing.T) {
	// ทดสอบ fallback ของ dignityToBalaAndSize และ effectivePlanet กรณี Rasi นอกช่วง
	bala, size := dignityToBalaAndSize("UnknownDignity")
	if bala != 50 || size != 15 {
		t.Errorf("expected default bala=50 size=15, got %f %f", bala, size)
	}
	if eff := effectivePlanet(GephiNode{Body: "Lagna", Rasi: 99}); eff != "Lagna" {
		t.Errorf("expected fallback 'Lagna' for invalid Rasi, got %q", eff)
	}
	// ทดสอบ default ID/Name ใน resolvePersonMeta
	meta := resolvePersonMeta(GephiPersonReq{}, 0)
	if meta.ID != "P1" || meta.Name != "P1" {
		t.Errorf("expected default P1/P1, got %q/%q", meta.ID, meta.Name)
	}
}

func TestAshtakootaCompatibilityAndTables(t *testing.T) {
	// Ashwini (1) to Rohini (4):
	// Count from 1 to 4 = 4. 4 % 9 = 4 -> Kshema (Well-being)
	tara := ComputeTaraBala(1, 4)
	if !strings.Contains(tara.Name, "Kshema") || tara.Score != 1 || !tara.IsAuspicious {
		t.Errorf("expected Kshema score 1 auspicious, got %+v", tara)
	}

	// Gana compatibility: Deva and Manushya -> Harmonious (score 5)
	ganaLabel, ganaScore := ComputeGanaCompat(GanaDeva, GanaManushya)
	if !strings.Contains(ganaLabel, "Harmonious") || ganaScore != 5.0 {
		t.Errorf("expected Deva-Manushya to be Harmonious score 5, got label=%s score=%f", ganaLabel, ganaScore)
	}

	// Yoni compatibility: Horse and Elephant
	yoniLabel, yoniScore := ComputeYoniCompat(YoniHorse, YoniElephant)
	if yoniLabel == "" || yoniScore < 0 {
		t.Errorf("expected valid yoni score, got label=%s score=%f", yoniLabel, yoniScore)
	}

	// Nadi compatibility: Different nadis (Adi vs Madhya) -> Harmonious score 8
	nadiLabel, nadiScore := ComputeNadiCompat(NadiAdi, NadiMadhya)
	if !strings.Contains(nadiLabel, "Different") || nadiScore != 8.0 {
		t.Errorf("expected Nadi Harmonious score 8, got label=%s score=%f", nadiLabel, nadiScore)
	}

	// Full Ashtakoota
	moonA := GephiNode{Rasi: 1, Nakshatra: 1, NakshatraGana: GanaDeva, NakshatraYoni: YoniHorse, NakshatraNadi: NadiAdi}
	moonB := GephiNode{Rasi: 2, Nakshatra: 4, NakshatraGana: GanaManushya, NakshatraYoni: YoniSerpent, NakshatraNadi: NadiAntya}
	koota := ComputeAshtakoota(moonA, moonB)
	if koota.TotalScore <= 0 || koota.Verdict == "" {
		t.Errorf("expected positive total score and verdict, got %+v", koota)
	}

	// Boundary / fallback tests
	if g := GetNakshatraGana(99); g != GanaDeva {
		t.Errorf("expected default Deva for invalid nakshatra, got %s", g)
	}
	if y := GetNakshatraYoni(99); y != YoniHorse {
		t.Errorf("expected default YoniHorse for invalid yoni, got %s", y)
	}
	if n := GetNakshatraNadi(99); n != NadiAdi {
		t.Errorf("expected default NadiAdi for invalid nadi, got %s", n)
	}
}

func TestGephiD9SynastryAndScoring(t *testing.T) {
	nodes := []GephiNode{
		{
			ID: "P1_Moon", PersonID: "P1", Body: "Moon", Rasi: 2, D9Rasi: 5,
		},
		{
			ID: "P2_Venus", PersonID: "P2", Body: "Venus", Rasi: 8, D9Rasi: 5, // D9 Yuti in Leo!
		},
		{
			ID: "P1_Jupiter", PersonID: "P1", Body: "Jupiter", Rasi: 3, D9Rasi: 5, // Intra with P1_Moon, cross-chart with P2_Venus
		},
		{
			ID: "P2_Mars", PersonID: "P2", Body: "Mars", Rasi: 1, D9Rasi: 9, // Different D9 sign
		},
	}

	edges := appendD9SynastryEdges(nil, nodes)
	if len(edges) != 2 {
		t.Fatalf("expected 2 cross-chart D9 Yuti edges, got %d", len(edges))
	}
	if edges[0].Source != "P1_Moon" || edges[0].Target != "P2_Venus" || edges[0].Type != EdgeDirUndirected {
		t.Errorf("unexpected D9 edge: %+v", edges[0])
	}
	if edges[0].Layer != "D9" || edges[0].Relation != EdgeRelD9Yuti || edges[0].ScoringFormula == "" {
		t.Errorf("unexpected D9 edge attributes: %+v", edges[0])
	}
}

func TestGephiPairwiseMetricsComputation(t *testing.T) {
	req := GephiReq{
		Persons: []GephiPersonReq{
			{
				ID: "P1", Name: "Alpha",
				CalculateReq: CalculateReq{Year: 1988, Month: 12, Day: 23, HourThai: 10.35, TimeZone: 7.0, Latitude: 13.75, Longitude: 100.5},
			},
			{
				ID: "P2", Name: "Beta",
				CalculateReq: CalculateReq{Year: 1991, Month: 4, Day: 19, HourThai: 2.15, TimeZone: 7.0, Latitude: 16.46, Longitude: 102.82},
			},
		},
	}
	graph := ComputeGephiGraph(req)
	if len(graph.PairwiseMetrics) != 1 {
		t.Fatalf("expected 1 pairwise metric for 2 persons, got %d", len(graph.PairwiseMetrics))
	}
	pm := graph.PairwiseMetrics[0]
	if pm.PersonA != "P1" || pm.PersonB != "P2" {
		t.Errorf("unexpected pairwise persons: %s - %s", pm.PersonA, pm.PersonB)
	}
	if pm.TotalEdges == 0 || pm.AshtakootaScore <= 0 {
		t.Errorf("expected positive total edges and Ashtakoota score, got %+v", pm)
	}
}
