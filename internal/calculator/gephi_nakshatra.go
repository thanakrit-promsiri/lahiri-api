package calculator

import (
	"fmt"
	"sort"
	"strings"
)

const (
	NodeKindPoint        = "Point"
	NodeKindNakshatra    = "Nakshatra"
	NakshatraHubPersonID = "NAKSHATRA"

	EdgeRelOccupiedNak = "Occupied Nakshatra"
	EdgeScopeHub       = "hub"
)

// nakshatraHubColor ใช้สี Slate เข้มสำหรับโหนด Hub นักษัตรเพื่อตัดกับชุดสีประจำบุคคลอย่างชัดเจนใน Gephi
var nakshatraHubColor = GephiColor{R: 51, G: 65, B: 85}

// ComputeGephiNakshatraGraph สร้างโครงข่ายกราฟ Bipartite โดยมีนักษัตรเป็น Hub กลาง (เฉพาะนักษัตรที่ถูกใช้งาน 1-27)
// เชื่อมโยงกับลัคนาและดาวเคราะห์ทั้ง 9 ดวงของแต่ละบุคคล (คนละ 10 จุด)
func ComputeGephiNakshatraGraph(req GephiReq) GephiNakshatraGraph {
	pointNodes := make([]GephiNakshatraNode, 0, len(req.Persons)*10)
	occupiedCounts := make(map[int]int, 27)

	for i, person := range req.Persons {
		resolved := resolvePersonMeta(person, i)
		color := personColorPalette[i%len(personColorPalette)]
		baseNodes := buildPersonGephiNodes(resolved, color)
		for _, base := range baseNodes {
			pNode := toNakshatraPointNode(base)
			pointNodes = append(pointNodes, pNode)
			if pNode.Nakshatra >= 1 && pNode.Nakshatra <= 27 {
				occupiedCounts[pNode.Nakshatra]++
			}
		}
	}

	// สร้างโหนด Hub เฉพาะนักษัตรที่ถูกใช้งานจริงเรียงตามลำดับ 1..27
	activeIndices := make([]int, 0, len(occupiedCounts))
	for idx := range occupiedCounts {
		activeIndices = append(activeIndices, idx)
	}
	sort.Ints(activeIndices)

	hubByIndex := make(map[int]GephiNakshatraNode, len(activeIndices))
	hubNodes := make([]GephiNakshatraNode, 0, len(activeIndices))
	for _, nakIdx := range activeIndices {
		hub := buildNakshatraHubNode(nakIdx, occupiedCounts[nakIdx])
		hubByIndex[nakIdx] = hub
		hubNodes = append(hubNodes, hub)
	}

	// รวมโหนดทั้งหมด: จุดของทุกคน + โหนด Hub นักษัตร
	nodes := make([]GephiNakshatraNode, 0, len(pointNodes)+len(hubNodes))
	nodes = append(nodes, pointNodes...)
	nodes = append(nodes, hubNodes...)

	// สร้างเส้นเชื่อมแบบมีทิศทาง Point -> Nakshatra Hub
	edges := make([]GephiNakshatraEdge, 0, len(pointNodes))
	for _, pNode := range pointNodes {
		hub, ok := hubByIndex[pNode.Nakshatra]
		if !ok {
			continue
		}
		edges = append(edges, newNakshatraOccupiedEdge(len(edges)+1, pNode, hub))
	}

	return GephiNakshatraGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

func toNakshatraPointNode(base GephiNode) GephiNakshatraNode {
	namesThai := GetThaiNakshatraNames()
	nakThai := ""
	if base.Nakshatra >= 1 && base.Nakshatra <= len(namesThai) {
		nakThai = namesThai[base.Nakshatra-1]
	}

	// เพิ่มน้ำหนักขนาดเล็กน้อยให้ Moon (จันทรนักษัตร) และ Lagna เพื่อความเด่นในกราฟนักษัตร
	size := base.NodeSize
	if base.Body == "Moon" || base.Body == "Lagna" {
		size += 3.0
	}

	return GephiNakshatraNode{
		ID:                base.ID,
		Label:             base.Label,
		NodeKind:          NodeKindPoint,
		PersonID:          base.PersonID,
		PersonName:        base.PersonName,
		Body:              base.Body,
		BodyThai:          base.BodyThai,
		Rasi:              base.Rasi,
		RasiName:          base.RasiName,
		Degree:            base.Degree,
		Longitude:         base.Longitude,
		Nakshatra:         base.Nakshatra,
		NakshatraName:     base.NakshatraName,
		NakshatraThai:     nakThai,
		NakshatraPada:     base.NakshatraPada,
		NakshatraLord:     base.NakshatraLord,
		NakshatraLordThai: TranslateLordThai(base.NakshatraLord),
		Dignity:           base.Dignity,
		DignityThai:       base.DignityThai,
		Bala:              base.Bala,
		NodeSize:          roundTo2(size),
		House:             base.House,
		IsRetrograde:      base.IsRetrograde,
		IsCombust:         base.IsCombust,
		IsVargottama:      base.IsVargottama,
		D9Rasi:            base.D9Rasi,
		D9RasiName:        base.D9RasiName,
		D9Dignity:         base.D9Dignity,
		Color:             base.Color,
	}
}

func buildNakshatraHubNode(nakIndex, occupiedCount int) GephiNakshatraNode {
	names := GetNakshatraNames()
	namesThai := GetThaiNakshatraNames()
	lords := GetNakshatraLords()

	idx0 := nakIndex - 1
	name := fmt.Sprintf("Nak%d", nakIndex)
	nameThai := name
	if idx0 >= 0 && idx0 < len(names) {
		name = names[idx0]
		nameThai = namesThai[idx0]
	}
	lord := lords[idx0%len(lords)]
	lordThai := TranslateLordThai(lord)

	span := 360.0 / 27.0
	startLon := float64(idx0) * span
	signIdx, degInSign := ParseLongitude(startLon)
	engSigns := GetEnglishSigns()
	rasiName := ""
	if signIdx >= 0 && signIdx < len(engSigns) {
		rasiName = engSigns[signIdx]
	}

	slug := strings.ReplaceAll(name, " ", "_")
	nodeID := fmt.Sprintf("Nak_%02d_%s", nakIndex, slug)
	size := roundTo2(20.0 + float64(occupiedCount)*5.0)
	if size > 55.0 {
		size = 55.0
	}

	return GephiNakshatraNode{
		ID:                nodeID,
		Label:             fmt.Sprintf("%s (%d)", name, nakIndex),
		NodeKind:          NodeKindNakshatra,
		PersonID:          NakshatraHubPersonID,
		PersonName:        "Nakshatra Hub",
		Body:              fmt.Sprintf("Nak_%02d", nakIndex),
		BodyThai:          nameThai,
		Rasi:              signIdx + 1,
		RasiName:          rasiName,
		Degree:            roundTo4(degInSign),
		Longitude:         roundTo4(startLon),
		Nakshatra:         nakIndex,
		NakshatraName:     name,
		NakshatraThai:     nameThai,
		NakshatraLord:     lord,
		NakshatraLordThai: lordThai,
		OccupiedCount:     occupiedCount,
		Dignity:           NodeKindNakshatra,
		DignityThai:       "นักษัตรฤกษ์",
		Bala:              roundTo2(float64(occupiedCount) * 15.0),
		NodeSize:          size,
		Color:             nakshatraHubColor,
	}
}

func newNakshatraOccupiedEdge(seq int, point, hub GephiNakshatraNode) GephiNakshatraEdge {
	weight := nakshatraOccupiedWeight(point)
	score := nakshatraAffinityScore(point, hub.NakshatraLord)
	edgeType, colorLabel, polarity, color := edgeStyleByScore(score)
	label := fmt.Sprintf("Pad %d (%s)", point.NakshatraPada, point.Dignity)

	return GephiNakshatraEdge{
		GephiEdge: GephiEdge{
			ID:         fmt.Sprintf("e%d", seq),
			Source:     point.ID,
			Target:     hub.ID,
			Type:       EdgeDirDirected,
			Weight:     weight,
			Label:      label,
			Relation:   EdgeRelOccupiedNak,
			EdgeType:   edgeType,
			Scope:      EdgeScopeHub,
			Polarity:   polarity,
			ColorLabel: colorLabel,
			Color:      color,
		},
		Pada:          point.NakshatraPada,
		NakshatraLord: hub.NakshatraLord,
	}
}

func nakshatraOccupiedWeight(p GephiNakshatraNode) float64 {
	var base float64
	switch p.Dignity {
	case "Exalted":
		base = 8.5
	case "Moolatrikona", "Own Sign":
		base = 7.5
	case "Lagna":
		base = 8.0
	case "Friend":
		base = 6.5
	case "Neutral":
		base = 5.5
	case "Enemy":
		base = 4.5
	case "Debilitated":
		base = 3.5
	default:
		base = 5.0
	}

	if p.Body == "Moon" {
		base += 1.5
	}
	if p.IsVargottama {
		base += 0.5
	}
	if base > 10.0 {
		base = 10.0
	}
	return roundTo2(base)
}

func nakshatraAffinityScore(p GephiNakshatraNode, nakLord string) int {
	if p.Body == "Lagna" {
		return 1
	}

	score := 0
	switch p.Dignity {
	case "Exalted", "Moolatrikona", "Own Sign", "Friend":
		score++
	case "Enemy", "Debilitated":
		score--
	}

	if p.Body == nakLord {
		score++
	} else {
		score += grahaAttitude(p.Body, nakLord)
	}
	return score
}
