package calculator

import (
	"fmt"
	"math"
	"strings"
)

const (
	EdgeDirDirected   = "directed"
	EdgeDirUndirected = "undirected"

	EdgeRelYuti        = "Yuti"
	EdgeRelDrishti     = "Graha Drishti"
	EdgeRelParivartana = "Parivartana"

	EdgeTypeHarmonious  = "Harmonious"
	EdgeTypeNeutral     = "Neutral"
	EdgeTypeChallenging = "Challenging"

	EdgeScopeInter = "inter"
	EdgeScopeIntra = "intra"
)

type dignityMetric struct {
	Bala     float64
	NodeSize float64
}

var dignityBalaTable = map[string]dignityMetric{
	"Exalted":      {Bala: 100, NodeSize: 32},
	"Moolatrikona": {Bala: 85, NodeSize: 27},
	"Own Sign":     {Bala: 75, NodeSize: 24},
	"Lagna":        {Bala: 65, NodeSize: 22},
	"Friend":       {Bala: 60, NodeSize: 19},
	"Neutral":      {Bala: 50, NodeSize: 15},
	"Enemy":        {Bala: 35, NodeSize: 11},
	"Debilitated":  {Bala: 20, NodeSize: 8},
}

var personColorPalette = []GephiColor{
	{R: 245, G: 158, B: 11}, // P1: Amber
	{R: 16, G: 185, B: 129}, // P2: Emerald
	{R: 139, G: 92, B: 246}, // P3: Violet
	{R: 6, G: 182, B: 212},  // P4: Cyan
	{R: 236, G: 72, B: 153}, // P5: Pink
	{R: 249, G: 115, B: 22}, // P6: Orange
}

// drishtiAspects กำหนดระยะเกณฑ์สายตาของดาวพระเวท (นับรวมราศีที่ดาวสถิตเป็น 1)
var drishtiAspects = map[string][]int{
	"Sun":     {7},
	"Moon":    {7},
	"Mercury": {7},
	"Venus":   {7},
	"Mars":    {4, 7, 8},
	"Jupiter": {5, 7, 9},
	"Saturn":  {3, 7, 10},
	"Rahu":    {5, 7, 9},
	"Ketu":    {5, 7, 9},
}

// signRulingPlanets ระบุเฉพาะดาวเจ้าเรือนเกษตร 7 ดวงสำหรับตรวจจับ Parivartana
var signRulingPlanets = map[string]bool{
	"Sun":     true,
	"Moon":    true,
	"Mars":    true,
	"Mercury": true,
	"Jupiter": true,
	"Venus":   true,
	"Saturn":  true,
}

var harmoniousSamePlanets = map[string]bool{
	"Sun":     true,
	"Moon":    true,
	"Mercury": true,
	"Jupiter": true,
	"Venus":   true,
}

// ความสัมพันธ์เพิ่มเติมสำหรับราหูและเกตุ เพื่อเสริมตาราง Naisargika Maitri ใน dignity.go
var extraNodeFriends = map[string]map[string]bool{
	"Rahu":    {"Mercury": true, "Venus": true, "Saturn": true},
	"Ketu":    {"Mars": true, "Jupiter": true},
	"Mercury": {"Rahu": true},
	"Venus":   {"Rahu": true},
	"Saturn":  {"Rahu": true},
	"Mars":    {"Ketu": true},
	"Jupiter": {"Ketu": true},
}

var extraNodeEnemies = map[string]map[string]bool{
	"Rahu":    {"Sun": true, "Moon": true, "Mars": true},
	"Ketu":    {"Sun": true, "Moon": true, "Venus": true},
	"Sun":     {"Rahu": true, "Ketu": true},
	"Moon":    {"Rahu": true, "Ketu": true},
	"Mars":    {"Rahu": true},
	"Jupiter": {"Rahu": true},
	"Venus":   {"Ketu": true},
}

// ComputeGephiGraph สร้างโครงข่ายโหนด (คนละ 10 โหนด) และเส้นเชื่อมความสัมพันธ์ (Yuti, Graha Drishti, Parivartana)
func ComputeGephiGraph(req GephiReq) GephiGraph {
	nodes := make([]GephiNode, 0, len(req.Persons)*10)
	for i, person := range req.Persons {
		resolved := resolvePersonMeta(person, i)
		color := personColorPalette[i%len(personColorPalette)]
		nodes = append(nodes, buildPersonGephiNodes(resolved, color)...)
	}

	edges := buildGephiEdges(nodes)
	return GephiGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

func resolvePersonMeta(person GephiPersonReq, idx int) GephiPersonReq {
	id := strings.TrimSpace(person.ID)
	if id == "" {
		id = fmt.Sprintf("P%d", idx+1)
	}
	name := strings.TrimSpace(person.Name)
	if name == "" {
		name = id
	}
	person.ID = id
	person.Name = name
	return person
}

func buildPersonGephiNodes(person GephiPersonReq, color GephiColor) []GephiNode {
	return buildPersonGephiNodesWithNode(person, color, false)
}

func buildPersonGephiNodesWithNode(person GephiPersonReq, color GephiColor, useTrueNode bool) []GephiNode {
	var horoscope CalculateResp
	if useTrueNode {
		horoscope = ComputeHoroscopeTrueNode(person.CalculateReq)
	} else {
		horoscope = ComputeHoroscope(person.CalculateReq)
	}
	nodes := make([]GephiNode, 0, 10)

	// 1. โหนดลัคนา (Lagna)
	if len(horoscope.HousesAndAngles) > 0 {
		asc := horoscope.HousesAndAngles[0]
		nodes = append(nodes, buildLagnaNode(person, asc, color))
	}

	// 2. โหนดดาวนวเคราะห์ 9 ดวงตามลำดับมาตรฐาน NavagrahaOrder
	planetByBody := make(map[string]ApiPlanet, len(horoscope.Planets))
	for _, p := range horoscope.Planets {
		_, _, canonical := TranslateBodyInfo(p.Name)
		planetByBody[canonical] = p
	}

	for _, body := range NavagrahaOrder {
		p, ok := planetByBody[body]
		if !ok {
			continue
		}
		nodes = append(nodes, buildPlanetNode(person, body, p, color))
	}

	return nodes
}

func buildLagnaNode(person GephiPersonReq, asc ApiHouse, color GephiColor) GephiNode {
	signIdx, degInSign := ParseLongitude(asc.Degree)
	nak, pada, _ := GetNakshatra(asc.Degree)
	d9SignIdx := GetSignIndexByName(asc.Navamsa.Sign)
	bala, size := dignityToBalaAndSize("Lagna")
	nodeID := person.ID + "_Lagna"

	return GephiNode{
		ID:            nodeID,
		Label:         nodeID,
		PersonID:      person.ID,
		PersonName:    person.Name,
		Body:          "Lagna",
		BodyThai:      asc.NameThai,
		Rasi:          signIdx + 1,
		RasiName:      asc.Position.Sign,
		RasiThai:      asc.Position.SignThai,
		Degree:        roundTo4(degInSign),
		Longitude:     roundTo4(asc.Degree),
		Nakshatra:     nak.Index,
		NakshatraName: asc.Nakshatra,
		NakshatraPada: pada,
		NakshatraLord: asc.Lord,
		Dignity:       "Lagna",
		DignityThai:   "ลัคนา",
		Bala:          bala,
		NodeSize:      size,
		House:         asc.House,
		IsVargottama:  asc.Position.Sign != "" && asc.Position.Sign == asc.Navamsa.Sign,
		D9Rasi:        d9SignIdx + 1,
		D9RasiName:    asc.Navamsa.Sign,
		D9Dignity:     "Lagna",
		Color:         color,
	}
}

func buildPlanetNode(person GephiPersonReq, body string, p ApiPlanet, color GephiColor) GephiNode {
	signIdx, degInSign := ParseLongitude(p.Longitude)
	nak, pada, _ := GetNakshatra(p.Longitude)
	d9SignIdx := GetSignIndexByName(p.Navamsa.Sign)
	d9Dig, _ := CalculateSignDignity(p.Name, d9SignIdx)
	bala, size := dignityToBalaAndSize(p.Dignity)
	nodeID := person.ID + "_" + body

	return GephiNode{
		ID:            nodeID,
		Label:         nodeID,
		PersonID:      person.ID,
		PersonName:    person.Name,
		Body:          body,
		BodyThai:      p.NameThai,
		Rasi:          signIdx + 1,
		RasiName:      p.Position.Sign,
		RasiThai:      p.Position.SignThai,
		Degree:        roundTo4(degInSign),
		Longitude:     roundTo4(p.Longitude),
		Nakshatra:     nak.Index,
		NakshatraName: p.Nakshatra,
		NakshatraPada: pada,
		NakshatraLord: p.Lord,
		Dignity:       p.Dignity,
		DignityThai:   p.DignityThai,
		Bala:          bala,
		NodeSize:      size,
		House:         p.House,
		IsRetrograde:  p.IsRetrograde,
		IsCombust:     p.IsCombust,
		IsVargottama:  p.Position.Sign != "" && p.Position.Sign == p.Navamsa.Sign,
		D9Rasi:        d9SignIdx + 1,
		D9RasiName:    p.Navamsa.Sign,
		D9Dignity:     d9Dig,
		Color:         color,
	}
}

func buildGephiEdges(nodes []GephiNode) []GephiEdge {
	edges := make([]GephiEdge, 0, len(nodes)*4)
	edges = appendYutiEdges(edges, nodes)
	edges = appendDrishtiEdges(edges, nodes)
	edges = appendParivartanaEdges(edges, nodes)
	return edges
}

func appendYutiEdges(edges []GephiEdge, nodes []GephiNode) []GephiEdge {
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			u, v := nodes[i], nodes[j]
			if u.Rasi != v.Rasi {
				continue
			}
			orb := math.Abs(u.Degree - v.Degree)
			weight := calcOrbWeight(orb, 10.0, 3.0, 3.0)
			label := fmt.Sprintf("Yuti (%.2f°)", roundTo2(orb))
			edges = append(edges, newGephiEdge(len(edges)+1, u, v, EdgeDirUndirected, EdgeRelYuti, label, weight, true, orb, 0))
		}
	}
	return edges
}

func appendDrishtiEdges(edges []GephiEdge, nodes []GephiNode) []GephiEdge {
	nodesByRasi := make(map[int][]GephiNode, 12)
	for _, n := range nodes {
		nodesByRasi[n.Rasi] = append(nodesByRasi[n.Rasi], n)
	}

	for _, u := range nodes {
		aspects, canAspect := drishtiAspects[u.Body]
		if !canAspect {
			continue
		}
		for _, aspectHouse := range aspects {
			targetRasi := ((u.Rasi - 1 + aspectHouse - 1) % 12) + 1
			for _, v := range nodesByRasi[targetRasi] {
				orb := math.Abs(u.Degree - v.Degree)
				weight := calcOrbWeight(orb, 7.0, 4.0, 3.0)
				label := fmt.Sprintf("Drishti %d (%.2f°)", aspectHouse, roundTo2(orb))
				edges = append(edges, newGephiEdge(len(edges)+1, u, v, EdgeDirDirected, EdgeRelDrishti, label, weight, true, orb, aspectHouse))
			}
		}
	}
	return edges
}

func appendParivartanaEdges(edges []GephiEdge, nodes []GephiNode) []GephiEdge {
	for i := 0; i < len(nodes); i++ {
		u := nodes[i]
		if !signRulingPlanets[u.Body] {
			continue
		}
		lordU := GetSignLord(u.Rasi - 1)
		for j := i + 1; j < len(nodes); j++ {
			v := nodes[j]
			if !signRulingPlanets[v.Body] || u.Body == v.Body {
				continue
			}
			lordV := GetSignLord(v.Rasi - 1)
			if lordU == v.Body && lordV == u.Body {
				edges = append(edges, newGephiEdge(len(edges)+1, u, v, EdgeDirUndirected, EdgeRelParivartana, "Parivartana", 8.0, false, 0, 0))
			}
		}
	}
	return edges
}

func newGephiEdge(seq int, u, v GephiNode, direction, relation, label string, weight float64, hasOrb bool, orb float64, drishtiHouse int) GephiEdge {
	edgeType, colorLabel, polarity, color := classifyGrahaMaitri(u, v)
	scope := EdgeScopeInter
	if u.PersonID == v.PersonID {
		scope = EdgeScopeIntra
	}

	return GephiEdge{
		ID:           fmt.Sprintf("e%d", seq),
		Source:       u.ID,
		Target:       v.ID,
		Type:         direction,
		Weight:       roundTo2(weight),
		Label:        label,
		Relation:     relation,
		EdgeType:     edgeType,
		Scope:        scope,
		Polarity:     polarity,
		HasOrb:       hasOrb,
		Orb:          roundTo2(orb),
		DrishtiHouse: drishtiHouse,
		ColorLabel:   colorLabel,
		Color:        color,
	}
}

// effectivePlanet แปลงโหนดลัคนาให้เป็นดาวเจ้าเรือนลัคนา (Lagnesh) สำหรับประเมิน Graha Maitri
func effectivePlanet(n GephiNode) string {
	if n.Body == "Lagna" {
		if lord := GetSignLord(n.Rasi - 1); lord != "" {
			return lord
		}
	}
	return n.Body
}

func grahaAttitude(from, to string) int {
	if naturalFriends[from][to] || extraNodeFriends[from][to] {
		return 1
	}
	if naturalEnemies[from][to] || extraNodeEnemies[from][to] {
		return -1
	}
	return 0
}

func classifyGrahaMaitri(u, v GephiNode) (edgeType, colorLabel string, polarity int, color GephiColor) {
	a := effectivePlanet(u)
	b := effectivePlanet(v)

	if a == b {
		if harmoniousSamePlanets[a] {
			return edgeStyleByScore(1)
		}
		return edgeStyleByScore(0)
	}

	return edgeStyleByScore(grahaAttitude(a, b) + grahaAttitude(b, a))
}

func edgeStyleByScore(score int) (edgeType, colorLabel string, polarity int, color GephiColor) {
	switch {
	case score > 0:
		return EdgeTypeHarmonious, "Blue", 1, GephiColor{R: 37, G: 99, B: 235}
	case score < 0:
		return EdgeTypeChallenging, "Red", -1, GephiColor{R: 220, G: 38, B: 38}
	default:
		return EdgeTypeNeutral, "Gray", 0, GephiColor{R: 156, G: 163, B: 175}
	}
}

// calcOrbWeight คำนวณน้ำหนักแปรผกผันกับระยะห่างองศา (Orb)
// เมื่อ orb <= tightOrb จะได้ maxW และค่อยๆ ลดลงเชิงเส้นจนถึง minW ที่ 30°
func calcOrbWeight(orb, maxW, minW, tightOrb float64) float64 {
	if orb <= tightOrb {
		return maxW
	}
	if orb >= 30.0 {
		return minW
	}
	w := maxW - (maxW-minW)*((orb-tightOrb)/(30.0-tightOrb))
	return roundTo2(w)
}

func dignityToBalaAndSize(dignity string) (bala, size float64) {
	if m, ok := dignityBalaTable[dignity]; ok {
		return m.Bala, m.NodeSize
	}
	return 50, 15
}

func roundTo2(val float64) float64 {
	return math.Round(val*100) / 100
}

func roundTo4(val float64) float64 {
	return math.Round(val*10000) / 10000
}
