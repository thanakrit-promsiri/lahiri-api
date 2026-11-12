package calculator

import (
	"math"
	"strings"
	"testing"
)

func TestComputeGephiHouseGraphTriadStructure(t *testing.T) {
	persons := sampleTriadPersons()
	graph := ComputeGephiHouseGraph(GephiReq{Persons: persons})

	// 3 คน × (9 ดาว + 12 ภพเรือน) = 63 โหนด
	if len(graph.Nodes) != 63 {
		t.Fatalf("expected 63 nodes (21 per person), got %d", len(graph.Nodes))
	}

	nodeMap := make(map[string]GephiHouseNode, len(graph.Nodes))
	planetCount, houseCount := 0, 0
	for _, n := range graph.Nodes {
		if _, dup := nodeMap[n.ID]; dup {
			t.Fatalf("duplicate node id %q", n.ID)
		}
		nodeMap[n.ID] = n
		switch n.NodeKind {
		case NodeKindPlanet:
			planetCount++
		case NodeKindHouse:
			houseCount++
		default:
			t.Fatalf("unexpected node kind %q for %s", n.NodeKind, n.ID)
		}
		if n.Rasi < 1 || n.Rasi > 12 {
			t.Errorf("node %s invalid Rasi %d", n.ID, n.Rasi)
		}
		if n.House < 1 || n.House > 12 {
			t.Errorf("node %s invalid House %d", n.ID, n.House)
		}
		if n.NodeSize <= 0 || n.Bala <= 0 {
			t.Errorf("node %s must have positive bala and node size, got bala=%v size=%v", n.ID, n.Bala, n.NodeSize)
		}
	}

	if planetCount != 27 || houseCount != 36 {
		t.Fatalf("expected 27 planet nodes and 36 house nodes, got %d / %d", planetCount, houseCount)
	}

	// ใน Whole Sign ทุกบุคคลมีครบ 12 ราศี -> House Lord = 36, Placed In = 27 * 3 = 81 -> รวม 117 เส้น
	if len(graph.Edges) != 117 {
		t.Fatalf("expected 117 house-flow edges for 3 persons, got %d", len(graph.Edges))
	}

	relCounts := make(map[string]int)
	scopeCounts := make(map[string]int)
	edgeTypeCounts := make(map[string]int)
	for _, e := range graph.Edges {
		src, okSrc := nodeMap[e.Source]
		dst, okDst := nodeMap[e.Target]
		if !okSrc || !okDst {
			t.Fatalf("edge %s references unknown endpoint %s -> %s", e.ID, e.Source, e.Target)
		}
		if e.Type != EdgeDirDirected {
			t.Errorf("edge %s expected directed type, got %q", e.ID, e.Type)
		}
		relCounts[e.Relation]++
		scopeCounts[e.Scope]++
		edgeTypeCounts[e.EdgeType]++

		switch e.Relation {
		case EdgeRelHouseLord:
			if src.NodeKind != NodeKindHouse || dst.NodeKind != NodeKindPlanet || src.PersonID != dst.PersonID {
				t.Errorf("invalid House Lord endpoints on %s: %s (%s) -> %s (%s)", e.ID, src.ID, src.NodeKind, dst.ID, dst.NodeKind)
			}
			if dst.Body != src.SignLord {
				t.Errorf("House Lord mismatch on %s: house sign lord=%s, target=%s", e.ID, src.SignLord, dst.Body)
			}
		case EdgeRelPlacedIn:
			if src.NodeKind != NodeKindPlanet || dst.NodeKind != NodeKindHouse {
				t.Errorf("invalid Placed In endpoints on %s: %s (%s) -> %s (%s)", e.ID, src.ID, src.NodeKind, dst.ID, dst.NodeKind)
			}
			if src.Rasi != dst.Rasi {
				t.Errorf("Placed In sign mismatch on %s: planet rasi=%d, house rasi=%d", e.ID, src.Rasi, dst.Rasi)
			}
			if e.Weight < 4.0 || e.Weight > 10.0 {
				t.Errorf("Placed In weight out of [4,10] on %s: %v", e.ID, e.Weight)
			}
		default:
			t.Errorf("unexpected relation %q on edge %s", e.Relation, e.ID)
		}
	}

	if relCounts[EdgeRelHouseLord] != 36 {
		t.Errorf("expected 36 House Lord edges, got %d", relCounts[EdgeRelHouseLord])
	}
	if relCounts[EdgeRelPlacedIn] != 81 {
		t.Errorf("expected 81 Placed In edges, got %d", relCounts[EdgeRelPlacedIn])
	}
	if scopeCounts[EdgeScopeIntra] == 0 || scopeCounts[EdgeScopeInter] == 0 {
		t.Errorf("expected both intra and inter edges, got %+v", scopeCounts)
	}
	if edgeTypeCounts[EdgeTypeHarmonious] == 0 || edgeTypeCounts[EdgeTypeChallenging] == 0 {
		t.Errorf("expected both Harmonious and Challenging edges, got %+v", edgeTypeCounts)
	}
}

func TestHouseFlowColorsAndHelpers(t *testing.T) {
	// ทดสอบคะแนนมาตรฐานดาวและโหนดราหู/เกตุ
	cases := []struct {
		body     string
		dignity  string
		signLord string
		want     int
	}{
		{"Sun", "Exalted", "Mars", 1},
		{"Moon", "Own Sign", "Moon", 1},
		{"Mars", "Friend", "Sun", 1},
		{"Venus", "Debilitated", "Mercury", -1},
		{"Saturn", "Enemy", "Sun", -1},
		{"Mercury", "Neutral", "Jupiter", 0},
		{"Rahu", "Neutral", "Saturn", 1},
		{"Rahu", "Neutral", "Sun", -1},
		{"Ketu", "Neutral", "Jupiter", 1},
		{"Ketu", "Neutral", "Moon", -1},
	}
	for _, tc := range cases {
		if got := planetSignScore(tc.body, tc.dignity, tc.signLord); got != tc.want {
			t.Errorf("planetSignScore(%q, %q, %q) = %d, want %d", tc.body, tc.dignity, tc.signLord, got, tc.want)
		}
	}

	// ทดสอบ Dusthana (-1) + Debilitated (-1) -> Challenging / Red
	p := GephiHouseNode{ID: "P1_Sun", PersonID: "P1", Body: "Sun", NodeKind: NodeKindPlanet, Rasi: 7, Degree: 10.0, Dignity: "Debilitated"}
	h8 := GephiHouseNode{ID: "P2_H8", PersonID: "P2", Body: "H8", NodeKind: NodeKindHouse, House: 8, Rasi: 7, SignLord: "Venus", Degree: 10.0}
	e := newHouseFlowEdge(1, p, h8, EdgeRelPlacedIn, "Placed In H8", 10.0, true, 0, -2)
	if e.EdgeType != EdgeTypeChallenging || e.ColorLabel != "Red" || e.Scope != EdgeScopeInter || e.Polarity != -1 {
		t.Errorf("expected Red/Challenging inter edge, got %+v", e)
	}

	// ทดสอบ Trikona (+1) + Exalted (+1) -> Harmonious / Blue
	eBlue := newHouseFlowEdge(2, p, h8, EdgeRelPlacedIn, "Placed In H9", 10.0, true, 0, 2)
	if eBlue.EdgeType != EdgeTypeHarmonious || eBlue.ColorLabel != "Blue" || eBlue.Polarity != 1 {
		t.Errorf("expected Blue/Harmonious edge, got %+v", eBlue)
	}

	// ทดสอบสีโหนดภพเรือนสว่างกว่าดาวของคนเดียวกันแต่ยังคงโทนสีประจำบุคคล
	base := GephiColor{R: 37, G: 99, B: 235}
	tinted := tintColorForHouse(base)
	if tinted.R <= base.R || tinted.G <= base.G || tinted.B <= base.B {
		t.Errorf("expected tinted house color brighter than base %+v, got %+v", base, tinted)
	}

	// ทดสอบสเกลน้ำหนัก Placed In กับระยะ Orb 0° vs 30°
	wClose := calcOrbWeight(0.0, 10.0, 4.0, 3.0)
	wFar := calcOrbWeight(30.0, 10.0, 4.0, 3.0)
	if wClose != 10.0 || math.Abs(wFar-4.0) > 1e-6 {
		t.Errorf("unexpected Placed In orb weights: close=%v far=%v", wClose, wFar)
	}

	// ตรวจสอบว่า node P1_H1 มี LabelThai ภาษาไทยครบ
	graph := ComputeGephiHouseGraph(GephiReq{Persons: sampleTriadPersons()[:2]})
	foundH1 := false
	for _, n := range graph.Nodes {
		if n.ID == "P1_H1" {
			foundH1 = true
			if !strings.Contains(n.LabelThai, "ภพ 1") || n.HouseGroup != "Kendra/Trikona" {
				t.Errorf("unexpected P1_H1 metadata: %+v", n)
			}
		}
	}
	if !foundH1 {
		t.Error("expected P1_H1 node in 2-person house graph")
	}
}
