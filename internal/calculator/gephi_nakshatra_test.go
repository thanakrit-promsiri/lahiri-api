package calculator

import (
	"strings"
	"testing"
)

func TestComputeGephiNakshatraGraphTriadStructure(t *testing.T) {
	graph := ComputeGephiNakshatraGraph(GephiReq{Persons: sampleTriadPersons()})

	pointCount, hubCount := 0, 0
	hubDegree := make(map[string]int)
	nodeMap := make(map[string]GephiNakshatraNode, len(graph.Nodes))

	for _, n := range graph.Nodes {
		if _, dup := nodeMap[n.ID]; dup {
			t.Fatalf("duplicate node id %q", n.ID)
		}
		nodeMap[n.ID] = n
		switch n.NodeKind {
		case NodeKindPoint:
			pointCount++
			if n.Nakshatra < 1 || n.Nakshatra > 27 || n.NakshatraPada < 1 || n.NakshatraPada > 4 {
				t.Errorf("point %s invalid nakshatra/pada: %d / %d", n.ID, n.Nakshatra, n.NakshatraPada)
			}
		case NodeKindNakshatra:
			hubCount++
			hubDegree[n.ID] = 0
			if n.PersonID != NakshatraHubPersonID || n.OccupiedCount <= 0 {
				t.Errorf("hub %s must have PersonID=%s and positive OccupiedCount, got %q / %d",
					n.ID, NakshatraHubPersonID, n.PersonID, n.OccupiedCount)
			}
		default:
			t.Fatalf("unexpected NodeKind %q on %s", n.NodeKind, n.ID)
		}
	}

	// 3 คน × 10 จุด = 30 Point nodes, และมีแต่ Hub นักษัตรที่ถูกใช้งาน (ต้องไม่เกิน 27 และมากกว่า 0)
	if pointCount != 30 {
		t.Fatalf("expected 30 Point nodes for 3 persons, got %d", pointCount)
	}
	if hubCount == 0 || hubCount > 27 {
		t.Fatalf("expected 1..27 active Nakshatra hub nodes, got %d", hubCount)
	}

	// ทุก Point 1 จุดต้องมี 1 Directed Edge ไปยัง Nakshatra Hub ของตนเองพอดี
	if len(graph.Edges) != pointCount {
		t.Fatalf("expected %d Point -> Nakshatra edges, got %d", pointCount, len(graph.Edges))
	}

	for _, e := range graph.Edges {
		src, okS := nodeMap[e.Source]
		dst, okT := nodeMap[e.Target]
		if !okS || !okT {
			t.Fatalf("edge %s references missing node %s -> %s", e.ID, e.Source, e.Target)
		}
		if src.NodeKind != NodeKindPoint || dst.NodeKind != NodeKindNakshatra {
			t.Errorf("expected Point -> Nakshatra on %s, got %s -> %s", e.ID, src.NodeKind, dst.NodeKind)
		}
		if src.Nakshatra != dst.Nakshatra {
			t.Errorf("edge %s nakshatra mismatch: point=%d hub=%d", e.ID, src.Nakshatra, dst.Nakshatra)
		}
		if e.Relation != EdgeRelOccupiedNak || e.Type != EdgeDirDirected || e.Pada != src.NakshatraPada {
			t.Errorf("unexpected edge metadata on %s: %+v", e.ID, e)
		}
		if e.Weight < 3.5 || e.Weight > 10.0 {
			t.Errorf("edge %s weight out of [3.5, 10]: %v", e.ID, e.Weight)
		}
		hubDegree[dst.ID]++
	}

	// ห้ามมีโหนด Hub นักษัตรลอยว่างเปล่า (in-degree ต้องตรงกับ OccupiedCount)
	for hubID, deg := range hubDegree {
		if deg == 0 || deg != nodeMap[hubID].OccupiedCount {
			t.Errorf("hub %s expected in-degree == OccupiedCount (%d), got %d", hubID, nodeMap[hubID].OccupiedCount, deg)
		}
	}
}

func TestNakshatraWeightAndAffinityHelpers(t *testing.T) {
	// Moon + Exalted + Vargottama ควรถูกจำกัดเพดานไม่เกิน 10.0
	wMoonExa := nakshatraOccupiedWeight(GephiNakshatraNode{Body: "Moon", Dignity: "Exalted", IsVargottama: true})
	if wMoonExa != 10.0 {
		t.Errorf("expected Moon Exalted Vargottama weight=10.0, got %v", wMoonExa)
	}

	// Moon Neutral ควรมีน้ำหนักมากกว่า Mars Neutral (7.0 vs 5.5)
	wMoonNeu := nakshatraOccupiedWeight(GephiNakshatraNode{Body: "Moon", Dignity: "Neutral"})
	wMarsNeu := nakshatraOccupiedWeight(GephiNakshatraNode{Body: "Mars", Dignity: "Neutral"})
	if wMoonNeu <= wMarsNeu {
		t.Errorf("expected Moon (%v) > Mars (%v) with same Neutral dignity", wMoonNeu, wMarsNeu)
	}

	// Compatibility score กับ Nakshatra Lord
	if s := nakshatraAffinityScore(GephiNakshatraNode{Body: "Lagna"}, "Saturn"); s <= 0 {
		t.Errorf("expected Lagna positive affinity, got %d", s)
	}
	if s := nakshatraAffinityScore(GephiNakshatraNode{Body: "Jupiter", Dignity: "Exalted"}, "Jupiter"); s <= 0 {
		t.Errorf("expected Exalted Jupiter in own star lord positive affinity, got %d", s)
	}
	if s := nakshatraAffinityScore(GephiNakshatraNode{Body: "Saturn", Dignity: "Debilitated"}, "Sun"); s >= 0 {
		t.Errorf("expected Debilitated Saturn in enemy Sun star lord negative affinity, got %d", s)
	}

	// Hub label & Thai name check
	hub := buildNakshatraHubNode(6, 4)
	if hub.ID != "Nak_06_Ardra" || !strings.Contains(hub.Label, "Ardra") || hub.NakshatraThai != "อารทรา" || hub.NakshatraLord != "Rahu" {
		t.Errorf("unexpected Nakshatra 6 hub metadata: %+v", hub)
	}
}
