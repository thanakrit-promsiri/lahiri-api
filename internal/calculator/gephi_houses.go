package calculator

import (
	"fmt"
	"math"
)

const (
	NodeKindPlanet = "Planet"
	NodeKindHouse  = "House"

	EdgeRelHouseLord = "House Lord"
	EdgeRelPlacedIn  = "Placed In"
)

var houseGroupTable = map[int]struct {
	Group      string
	BaseBala   float64
	LordWeight float64
	Polarity   int
}{
	1:  {Group: "Kendra/Trikona", BaseBala: 70, LordWeight: 8.0, Polarity: 1},
	2:  {Group: "Dhana", BaseBala: 50, LordWeight: 6.0, Polarity: 0},
	3:  {Group: "Upachaya", BaseBala: 50, LordWeight: 6.0, Polarity: 0},
	4:  {Group: "Kendra", BaseBala: 60, LordWeight: 7.0, Polarity: 0},
	5:  {Group: "Trikona", BaseBala: 65, LordWeight: 8.0, Polarity: 1},
	6:  {Group: "Dusthana", BaseBala: 40, LordWeight: 5.0, Polarity: -1},
	7:  {Group: "Kendra", BaseBala: 60, LordWeight: 7.0, Polarity: 0},
	8:  {Group: "Dusthana", BaseBala: 40, LordWeight: 5.0, Polarity: -1},
	9:  {Group: "Trikona", BaseBala: 65, LordWeight: 8.0, Polarity: 1},
	10: {Group: "Kendra", BaseBala: 60, LordWeight: 7.0, Polarity: 0},
	11: {Group: "Labha", BaseBala: 60, LordWeight: 7.0, Polarity: 1},
	12: {Group: "Dusthana", BaseBala: 40, LordWeight: 5.0, Polarity: -1},
}

// ComputeGephiHouseGraph สร้างโครงข่าย Bipartite ระหว่างโหนดดาวเคราะห์ (9 ดวง/คน) กับโหนดภพเรือน (12 ภพ/คน)
// ตามความสัมพันธ์ House Lord Flow (ภพ -> เจ้าเรือนเกษตร และ ดาวเคราะห์ -> ภพที่ไปสถิตทั้งในดวงและข้ามดวง)
func ComputeGephiHouseGraph(req GephiReq) GephiHouseGraph {
	planetNodes := make([]GephiHouseNode, 0, len(req.Persons)*9)
	houseNodes := make([]GephiHouseNode, 0, len(req.Persons)*12)
	planetByKey := make(map[string]GephiHouseNode, len(req.Persons)*9)

	for i, person := range req.Persons {
		resolved := resolvePersonMeta(person, i)
		pColor := personColorPalette[i%len(personColorPalette)]
		hColor := tintColorForHouse(pColor)

		pNodes, hNodes := buildPersonHouseGraphNodes(resolved, pColor, hColor)
		for _, pn := range pNodes {
			planetNodes = append(planetNodes, pn)
			planetByKey[pn.ID] = pn
		}
		houseNodes = append(houseNodes, hNodes...)
	}

	nodes := make([]GephiHouseNode, 0, len(planetNodes)+len(houseNodes))
	nodes = append(nodes, planetNodes...)
	nodes = append(nodes, houseNodes...)

	edges := buildHouseFlowEdges(planetNodes, houseNodes, planetByKey)

	return GephiHouseGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

func buildPersonHouseGraphNodes(person GephiPersonReq, pColor, hColor GephiColor) ([]GephiHouseNode, []GephiHouseNode) {
	horoscope := ComputeHoroscope(person.CalculateReq)

	var ascDegInSign float64
	if len(horoscope.HousesAndAngles) > 0 {
		_, ascDegInSign = ParseLongitude(horoscope.HousesAndAngles[0].Degree)
	}

	planetByBody := make(map[string]ApiPlanet, len(horoscope.Planets))
	for _, p := range horoscope.Planets {
		_, _, canonical := TranslateBodyInfo(p.Name)
		planetByBody[canonical] = p
	}

	planetNodes := make([]GephiHouseNode, 0, len(NavagrahaOrder))
	planetBalaByBody := make(map[string]float64, len(NavagrahaOrder))
	planetDigByBody := make(map[string]string, len(NavagrahaOrder))

	for _, body := range NavagrahaOrder {
		p, ok := planetByBody[body]
		if !ok {
			continue
		}
		pn := buildHouseGraphPlanetNode(person, body, p, pColor)
		planetNodes = append(planetNodes, pn)
		planetBalaByBody[body] = pn.Bala
		planetDigByBody[body] = pn.Dignity
	}

	houseNodes := make([]GephiHouseNode, 0, len(horoscope.WholeSignHouses))
	for _, ws := range horoscope.WholeSignHouses {
		lordBala := planetBalaByBody[ws.Lord]
		lordDig := planetDigByBody[ws.Lord]
		hn := buildHouseGraphHouseNode(person, ws, ascDegInSign, lordBala, lordDig, hColor)
		houseNodes = append(houseNodes, hn)
	}

	return planetNodes, houseNodes
}

func buildHouseGraphPlanetNode(person GephiPersonReq, body string, p ApiPlanet, color GephiColor) GephiHouseNode {
	signIdx, degInSign := ParseLongitude(p.Longitude)
	nak, _, _ := GetNakshatra(p.Longitude)
	bala, size := dignityToBalaAndSize(p.Dignity)
	nodeID := person.ID + "_" + body

	return GephiHouseNode{
		ID:            nodeID,
		Label:         nodeID,
		NodeKind:      NodeKindPlanet,
		PersonID:      person.ID,
		PersonName:    person.Name,
		Body:          body,
		LabelThai:     p.NameThai,
		House:         p.House,
		HouseNameThai: p.HouseNameThai,
		HouseGroup:    houseGroupTable[p.House].Group,
		Rasi:          signIdx + 1,
		RasiName:      p.Position.Sign,
		RasiThai:      p.Position.SignThai,
		SignLord:      GetSignLord(signIdx),
		Degree:        roundTo4(degInSign),
		Nakshatra:     nak.Index,
		Dignity:       p.Dignity,
		Bala:          bala,
		NodeSize:      size,
		IsRetrograde:  p.IsRetrograde,
		IsCombust:     p.IsCombust,
		IsVargottama:  p.Position.Sign != "" && p.Position.Sign == p.Navamsa.Sign,
		Color:         color,
	}
}

func buildHouseGraphHouseNode(person GephiPersonReq, ws ApiWholeSign, ascDegInSign, lordBala float64, lordDignity string, color GephiColor) GephiHouseNode {
	signIdx := GetSignIndexByName(ws.Sign)
	meta := houseGroupTable[ws.House]
	if lordBala <= 0 {
		lordBala = 50
	}
	bala := roundTo2((meta.BaseBala + lordBala) / 2.0)
	size := roundTo2(8.0 + (bala/100.0)*20.0)

	bhavaLon := float64(signIdx)*30.0 + ascDegInSign
	nak, _, _ := GetNakshatra(bhavaLon)

	nodeID := fmt.Sprintf("%s_H%d", person.ID, ws.House)
	label := fmt.Sprintf("%s (%s)", nodeID, ws.HouseNameThai)

	return GephiHouseNode{
		ID:            nodeID,
		Label:         label,
		NodeKind:      NodeKindHouse,
		PersonID:      person.ID,
		PersonName:    person.Name,
		Body:          fmt.Sprintf("H%d", ws.House),
		LabelThai:     fmt.Sprintf("ภพ %d (%s)", ws.House, ws.HouseNameThai),
		House:         ws.House,
		HouseNameThai: ws.HouseNameThai,
		HouseGroup:    meta.Group,
		Rasi:          signIdx + 1,
		RasiName:      ws.Sign,
		RasiThai:      ws.SignThai,
		SignLord:      ws.Lord,
		Degree:        roundTo4(ascDegInSign),
		Nakshatra:     nak.Index,
		Dignity:       lordDignity,
		Bala:          bala,
		NodeSize:      size,
		Color:         color,
	}
}

func buildHouseFlowEdges(planetNodes, houseNodes []GephiHouseNode, planetByKey map[string]GephiHouseNode) []GephiEdge {
	edges := make([]GephiEdge, 0, len(houseNodes)+len(planetNodes)*4)

	// 1. House -> LordPlanet (เจ้าเรือนภพภายในดวงชะตา)
	for _, h := range houseNodes {
		lordNode, ok := planetByKey[h.PersonID+"_"+h.SignLord]
		if !ok {
			continue
		}
		meta := houseGroupTable[h.House]
		score := meta.Polarity + planetSignScore(lordNode.Body, lordNode.Dignity, lordNode.SignLord)
		label := fmt.Sprintf("Lord (%s)", lordNode.Body)
		edges = append(edges, newHouseFlowEdge(len(edges)+1, h, lordNode, EdgeRelHouseLord, label, meta.LordWeight, false, 0, score))
	}

	// 2. Planet -> TargetHouse (ดาวเคราะห์ไปสถิตในภพทั้งในดวงตัวเองและข้ามบุคคล)
	housesByRasi := make(map[int][]GephiHouseNode, 12)
	for _, h := range houseNodes {
		housesByRasi[h.Rasi] = append(housesByRasi[h.Rasi], h)
	}

	for _, p := range planetNodes {
		for _, h := range housesByRasi[p.Rasi] {
			orb := math.Abs(p.Degree - h.Degree)
			weight := calcOrbWeight(orb, 10.0, 4.0, 3.0)
			meta := houseGroupTable[h.House]
			score := meta.Polarity + planetSignScore(p.Body, p.Dignity, h.SignLord)
			label := fmt.Sprintf("Placed In H%d (%.2f°)", h.House, roundTo2(orb))
			edges = append(edges, newHouseFlowEdge(len(edges)+1, p, h, EdgeRelPlacedIn, label, weight, true, orb, score))
		}
	}

	return edges
}

func newHouseFlowEdge(seq int, src, dst GephiHouseNode, relation, label string, weight float64, hasOrb bool, orb float64, score int) GephiEdge {
	edgeType, colorLabel, polarity, color := edgeStyleByScore(score)
	scope := EdgeScopeInter
	if src.PersonID == dst.PersonID {
		scope = EdgeScopeIntra
	}

	return GephiEdge{
		ID:         fmt.Sprintf("e%d", seq),
		Source:     src.ID,
		Target:     dst.ID,
		Type:       EdgeDirDirected,
		Weight:     roundTo2(weight),
		Label:      label,
		Relation:   relation,
		EdgeType:   edgeType,
		Scope:      scope,
		Polarity:   polarity,
		HasOrb:     hasOrb,
		Orb:        roundTo2(orb),
		ColorLabel: colorLabel,
		Color:      color,
	}
}

func planetSignScore(body, dignity, signLord string) int {
	switch dignity {
	case "Exalted", "Moolatrikona", "Own Sign", "Friend":
		return 1
	case "Enemy", "Debilitated":
		return -1
	default:
		if body == "Rahu" || body == "Ketu" {
			return grahaAttitude(body, signLord)
		}
		return 0
	}
}

// tintColorForHouse ปรับโทนสีให้โหนดภพเรือนสว่าง/นุ่มกว่าโหนดดาวของบุคคลเดียวกันเพื่อแยกแยะได้ง่ายใน Gephi
func tintColorForHouse(c GephiColor) GephiColor {
	return GephiColor{
		R: c.R + (255-c.R)/2,
		G: c.G + (255-c.G)/2,
		B: c.B + (255-c.B)/2,
	}
}
