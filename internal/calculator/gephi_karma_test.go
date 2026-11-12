package calculator

import (
	"math"
	"testing"
)

func TestTrueNodeDiffersFromMeanNodeAndRemainsOpposite(t *testing.T) {
	req := sampleTriadPersons()[0].CalculateReq
	meanHoroscope := ComputeHoroscope(req)
	trueHoroscope := ComputeHoroscopeTrueNode(req)

	findNodeLon := func(h CalculateResp, body string) float64 {
		for _, p := range h.Planets {
			_, _, canonical := TranslateBodyInfo(p.Name)
			if canonical == body {
				return p.Longitude
			}
		}
		t.Fatalf("planet %s not found", body)
		return 0
	}

	meanRahu := findNodeLon(meanHoroscope, "Rahu")
	trueRahu := findNodeLon(trueHoroscope, "Rahu")
	trueKetu := findNodeLon(trueHoroscope, "Ketu")

	if math.Abs(meanRahu-trueRahu) < 1e-4 {
		t.Errorf("expected True Node Rahu (%f) to differ from Mean Node Rahu (%f)", trueRahu, meanRahu)
	}
	expectedKetu := NormalizeAngle(trueRahu + 180.0)
	if math.Abs(trueKetu-expectedKetu) > 1e-6 {
		t.Errorf("expected True Ketu (%f) to be 180° opposite True Rahu (%f)", trueKetu, expectedKetu)
	}
}

func TestComputeGephiKarmaGraphTriadStructure(t *testing.T) {
	graph := ComputeGephiKarmaGraph(GephiReq{Persons: sampleTriadPersons()})

	if len(graph.Nodes) != 30 {
		t.Fatalf("expected 30 nodes for 3 persons, got %d", len(graph.Nodes))
	}

	nodeMap := make(map[string]GephiKarmaNode, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodeMap[n.ID] = n
		if n.Rasi < 1 || n.Rasi > 12 {
			t.Errorf("node %s invalid Rasi %d", n.ID, n.Rasi)
		}
		if n.D9Rasi < 1 || n.D9Rasi > 12 {
			t.Errorf("node %s invalid D9Rasi %d", n.ID, n.D9Rasi)
		}
		if n.D9Degree < 0 || n.D9Degree >= 30.0 {
			t.Errorf("node %s D9Degree out of [0,30): %v", n.ID, n.D9Degree)
		}
		if n.DignityScore < -5 || n.DignityScore > 5 {
			t.Errorf("node %s DignityScore out of [-5,5]: %d", n.ID, n.DignityScore)
		}
		switch n.Body {
		case "Rahu", "Ketu":
			if n.GrahaType != GrahaTypeKarmicAxis || n.ForceRole != ForceRoleRepel {
				t.Errorf("node %s expected Karmic Axis / Repel, got %s / %s", n.ID, n.GrahaType, n.ForceRole)
			}
		case "Sun", "Mars", "Saturn":
			if n.GrahaType != GrahaTypeMalefic || n.ForceRole != ForceRoleRepel {
				t.Errorf("node %s expected Malefic / Repel, got %s / %s", n.ID, n.GrahaType, n.ForceRole)
			}
		default:
			if n.GrahaType != GrahaTypeBenefic || n.ForceRole != ForceRoleAttract {
				t.Errorf("node %s expected Benefic / Attract, got %s / %s", n.ID, n.GrahaType, n.ForceRole)
			}
		}
	}

	if len(graph.Edges) == 0 {
		t.Fatal("expected karma edges for 3 persons")
	}

	relCounts := make(map[string]int)
	for _, e := range graph.Edges {
		src, okS := nodeMap[e.Source]
		dst, okT := nodeMap[e.Target]
		if !okS || !okT {
			t.Fatalf("edge %s references missing node %s -> %s", e.ID, e.Source, e.Target)
		}
		if src.PersonID == dst.PersonID {
			t.Errorf("expected cross-person karma edge on %s, got same person %s", e.ID, src.PersonID)
		}
		relCounts[e.Relation]++
	}

	if relCounts[KarmaRelAxis] == 0 {
		t.Errorf("expected Karmic Axis Alignment edges in triad graph, got %+v", relCounts)
	}
	if relCounts[KarmaRelTaskmaster] == 0 {
		t.Errorf("expected Karmic Taskmaster edges in triad graph, got %+v", relCounts)
	}
	if relCounts[KarmaRelDispositor] == 0 {
		t.Errorf("expected Nakshatra Dispositor edges in triad graph, got %+v", relCounts)
	}
}

func TestKarmaEdgeRulesUnit(t *testing.T) {
	// 1. Karmic Axis Alignment: D1 Conjunct, D1 Opposite, D9 Conjunct
	p1Rahu := GephiKarmaNode{ID: "P1_Rahu", PersonID: "P1", Body: "Rahu", GrahaType: GrahaTypeKarmicAxis, Rasi: 2, Degree: 10, D9Rasi: 4, D9Degree: 15}
	p2SunConj := GephiKarmaNode{ID: "P2_Sun", PersonID: "P2", Body: "Sun", GrahaType: GrahaTypeMalefic, Rasi: 2, Degree: 12, D9Rasi: 4, D9Degree: 16}
	p2MoonOpp := GephiKarmaNode{ID: "P2_Moon", PersonID: "P2", Body: "Moon", GrahaType: GrahaTypeBenefic, Rasi: 8, Degree: 11, D9Rasi: 1}
	p2VenD9Only := GephiKarmaNode{ID: "P2_Venus", PersonID: "P2", Body: "Venus", GrahaType: GrahaTypeBenefic, Rasi: 5, Degree: 5, D9Rasi: 4, D9Degree: 18}

	axisEdges := appendKarmicAxisEdges(nil, []GephiKarmaNode{p1Rahu, p2SunConj, p2MoonOpp, p2VenD9Only})
	if len(axisEdges) != 3 {
		t.Fatalf("expected 3 karmic axis edges (D1+D9, D1 opp, D9 only), got %d: %+v", len(axisEdges), axisEdges)
	}
	if axisEdges[0].Weight != 10.0 || axisEdges[0].Layer != "D1+D9" {
		t.Errorf("expected D1+D9 Weight=10, got Weight=%v Layer=%q", axisEdges[0].Weight, axisEdges[0].Layer)
	}
	if axisEdges[1].Weight != 10.0 || axisEdges[1].Layer != "D1" {
		t.Errorf("expected D1 opposite Weight=10, got Weight=%v Layer=%q", axisEdges[1].Weight, axisEdges[1].Layer)
	}
	if axisEdges[2].Weight != 8.0 || axisEdges[2].Layer != "D9" {
		t.Errorf("expected D9-only Weight=8, got Weight=%v Layer=%q", axisEdges[2].Weight, axisEdges[2].Layer)
	}

	// 2. Karmic Taskmaster: Saturn at Rasi 1 hits Rasi 3 (H3), 7 (H7), 10 (H10), ignores Rasi 2 (H2)
	p1Sat := GephiKarmaNode{ID: "P1_Saturn", PersonID: "P1", Body: "Saturn", Rasi: 1, Degree: 10}
	p2H3 := GephiKarmaNode{ID: "P2_Sun", PersonID: "P2", Body: "Sun", Rasi: 3, Degree: 10}
	p2H7 := GephiKarmaNode{ID: "P2_Moon", PersonID: "P2", Body: "Moon", Rasi: 7, Degree: 25}
	p2H10 := GephiKarmaNode{ID: "P2_Mars", PersonID: "P2", Body: "Mars", Rasi: 10, Degree: 12}
	p2H2 := GephiKarmaNode{ID: "P2_Mercury", PersonID: "P2", Body: "Mercury", Rasi: 2, Degree: 10}

	satEdges := appendKarmicTaskmasterEdges(nil, []GephiKarmaNode{p1Sat, p2H3, p2H7, p2H10, p2H2})
	if len(satEdges) != 3 {
		t.Fatalf("expected 3 Saturn Drishti edges (3,7,10), got %d", len(satEdges))
	}
	for _, e := range satEdges {
		if e.Type != EdgeDirDirected || e.Source != "P1_Saturn" || e.Weight < 7.0 || e.Weight > 9.0 {
			t.Errorf("unexpected Saturn taskmaster edge: %+v", e)
		}
	}

	// 3. Nakshatra Dispositor: P1_Moon in Bharani (Lord Venus), P2_Venus Exalted, P3_Venus Debilitated
	p1Moon := GephiKarmaNode{ID: "P1_Moon", PersonID: "P1", Body: "Moon", NakshatraName: "Bharani", NakshatraLord: "Venus"}
	p2VenExalted := GephiKarmaNode{ID: "P2_Venus", PersonID: "P2", Body: "Venus", GrahaType: GrahaTypeBenefic, Dignity: "Exalted", D9Dignity: "Own Sign", IsVargottama: true}
	p3VenDeb := GephiKarmaNode{ID: "P3_Venus", PersonID: "P3", Body: "Venus", GrahaType: GrahaTypeBenefic, Dignity: "Debilitated", D9Dignity: "Enemy"}

	dispEdges := appendNakshatraDispositorEdges(nil, []GephiKarmaNode{p1Moon, p2VenExalted, p3VenDeb})
	if len(dispEdges) != 1 {
		t.Fatalf("expected 1 Nakshatra Dispositor edge from prominent P2_Venus only, got %d: %+v", len(dispEdges), dispEdges)
	}
	if dispEdges[0].Source != "P2_Venus" || dispEdges[0].Target != "P1_Moon" || dispEdges[0].Weight != 9.0 || dispEdges[0].ForceRole != ForceRoleAttract {
		t.Errorf("unexpected Dispositor edge: %+v", dispEdges[0])
	}

	// Score table sanity check
	if dignityToKarmaScore("Exalted") != 5 || dignityToKarmaScore("Own Sign") != 4 || dignityToKarmaScore("Debilitated") != -5 {
		t.Error("dignityToKarmaScore mismatch against +5/+4/-5 specification")
	}
}
