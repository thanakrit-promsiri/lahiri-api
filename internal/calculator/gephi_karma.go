package calculator

import (
	"fmt"
	"math"
)

const (
	GrahaTypeBenefic    = "Benefic"
	GrahaTypeMalefic    = "Malefic"
	GrahaTypeKarmicAxis = "Karmic Axis"

	KarmaRelAxis       = "Karmic Axis Alignment"
	KarmaRelTaskmaster = "Karmic Taskmaster"
	KarmaRelDispositor = "Nakshatra Dispositor"

	ForceRoleAttract = "Attract"
	ForceRoleRepel   = "Repel"
)

var grahaTypeTable = map[string]string{
	"Lagna":   GrahaTypeBenefic,
	"Jupiter": GrahaTypeBenefic,
	"Venus":   GrahaTypeBenefic,
	"Moon":    GrahaTypeBenefic,
	"Mercury": GrahaTypeBenefic,
	"Sun":     GrahaTypeMalefic,
	"Mars":    GrahaTypeMalefic,
	"Saturn":  GrahaTypeMalefic,
	"Rahu":    GrahaTypeKarmicAxis,
	"Ketu":    GrahaTypeKarmicAxis,
}

var dignityScoreTable = map[string]int{
	"Exalted":      5,
	"Moolatrikona": 4,
	"Own Sign":     4,
	"Lagna":        3,
	"Friend":       2,
	"Neutral":      0,
	"Enemy":        -2,
	"Debilitated":  -5,
}

// ComputeGephiKarmaGraph สร้างผังกรรม (Topology of Karma) โดยใช้ตำแหน่ง True Node ของราหู/เกตุ
// พร้อมข้อมูลชั้น D-1 และ D-9 Navamsa และเส้นเชื่อม 3 ระบบ: Karmic Axis Alignment, Karmic Taskmaster และ Nakshatra Dispositor
func ComputeGephiKarmaGraph(req GephiReq) GephiKarmaGraph {
	nodes := make([]GephiKarmaNode, 0, len(req.Persons)*10)
	for i, person := range req.Persons {
		resolved := resolvePersonMeta(person, i)
		color := personColorPalette[i%len(personColorPalette)]
		baseNodes := buildPersonGephiNodesWithNode(resolved, color, true)
		for _, base := range baseNodes {
			nodes = append(nodes, toKarmaNode(base))
		}
	}

	edges := buildKarmaEdges(nodes)
	return GephiKarmaGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

func toKarmaNode(base GephiNode) GephiKarmaNode {
	grahaType := classifyGrahaType(base.Body)
	score := dignityToKarmaScore(base.Dignity)
	nodeSize := roundTo2(18.0 + float64(score)*2.0)

	forceRole := ForceRoleRepel
	if grahaType == GrahaTypeBenefic {
		forceRole = ForceRoleAttract
	}

	d9Deg := calcNavamsaDegreeInSign(base.Degree)
	d9Lon := 0.0
	if base.D9Rasi >= 1 {
		d9Lon = roundTo4(float64(base.D9Rasi-1)*30.0 + d9Deg)
	}

	return GephiKarmaNode{
		ID:            base.ID,
		Label:         base.Label,
		PersonID:      base.PersonID,
		PersonName:    base.PersonName,
		Body:          base.Body,
		BodyThai:      base.BodyThai,
		GrahaType:     grahaType,
		DignityScore:  score,
		ForceRole:     forceRole,
		NodeSize:      nodeSize,
		Rasi:          base.Rasi,
		RasiName:      base.RasiName,
		Degree:        base.Degree,
		Longitude:     base.Longitude,
		Nakshatra:     base.Nakshatra,
		NakshatraName: base.NakshatraName,
		NakshatraPada: base.NakshatraPada,
		NakshatraLord: base.NakshatraLord,
		Dignity:       base.Dignity,
		House:         base.House,
		IsRetrograde:  base.IsRetrograde,
		IsCombust:     base.IsCombust,
		IsVargottama:  base.IsVargottama,
		D9Rasi:        base.D9Rasi,
		D9RasiName:    base.D9RasiName,
		D9Degree:      roundTo4(d9Deg),
		D9Longitude:   d9Lon,
		D9Dignity:     base.D9Dignity,
		Color:         base.Color,
	}
}

func calcNavamsaDegreeInSign(degInSign float64) float64 {
	navamsaSpan := 30.0 / 9.0
	rem := math.Mod(degInSign, navamsaSpan)
	if rem < 0 {
		rem += navamsaSpan
	}
	return (rem / navamsaSpan) * 30.0
}

func classifyGrahaType(body string) string {
	if gt, ok := grahaTypeTable[body]; ok {
		return gt
	}
	return GrahaTypeBenefic
}

func dignityToKarmaScore(dignity string) int {
	if s, ok := dignityScoreTable[dignity]; ok {
		return s
	}
	return 0
}

func buildKarmaEdges(nodes []GephiKarmaNode) []GephiKarmaEdge {
	edges := make([]GephiKarmaEdge, 0, len(nodes)*5)
	edges = appendKarmicAxisEdges(edges, nodes)
	edges = appendKarmicTaskmasterEdges(edges, nodes)
	edges = appendNakshatraDispositorEdges(edges, nodes)
	return edges
}

// 1. Karmic Axis Alignment: ดาว/ลัคนาไปกุมหรือเล็ง (180°) ราหูหรือเกตุของอีกคน -> Weight = 10.0
// พร้อมตรวจชั้นนวางศ์จักร D-9 หากยังไม่เชื่อมใน D-1 แต่กุมราหู/เกตุใน D-9 (Weight = 8.0)
func appendKarmicAxisEdges(edges []GephiKarmaEdge, nodes []GephiKarmaNode) []GephiKarmaEdge {
	for i := 0; i < len(nodes); i++ {
		u := nodes[i]
		for j := i + 1; j < len(nodes); j++ {
			v := nodes[j]
			if u.PersonID == v.PersonID {
				continue
			}
			if u.GrahaType != GrahaTypeKarmicAxis && v.GrahaType != GrahaTypeKarmicAxis {
				continue
			}

			d1Conj := u.Rasi == v.Rasi
			d1Opp := (v.Rasi-u.Rasi+12)%12 == 6
			d9Conj := u.D9Rasi > 0 && u.D9Rasi == v.D9Rasi

			if d1Conj || d1Opp {
				orb := math.Abs(u.Degree - v.Degree)
				mode := "Conjunct"
				if d1Opp {
					mode = "Opposite 180°"
				}
				layer := "D1"
				if d9Conj {
					layer = "D1+D9"
				}
				label := fmt.Sprintf("Karmic Axis %s (%.2f°)", mode, roundTo2(orb))
				edges = append(edges, newKarmaEdge(
					len(edges)+1, u, v, EdgeDirUndirected, KarmaRelAxis,
					label, 10.0, layer, ForceRoleRepel, -1, true, orb, 0,
				))
				continue
			}

			if d9Conj {
				orb := math.Abs(u.D9Degree - v.D9Degree)
				label := fmt.Sprintf("D9 Karmic Axis Conjunct (%.2f°)", roundTo2(orb))
				edges = append(edges, newKarmaEdge(
					len(edges)+1, u, v, EdgeDirUndirected, KarmaRelAxis,
					label, 8.0, "D9", ForceRoleRepel, -1, true, orb, 0,
				))
			}
		}
	}
	return edges
}

// 2. Karmic Taskmaster: Directed Edge จาก Saturn ของทุกคนไปยังจุด/ดาวของคนอื่นที่โดนเกณฑ์ 3, 7, 10
func appendKarmicTaskmasterEdges(edges []GephiKarmaEdge, nodes []GephiKarmaNode) []GephiKarmaEdge {
	satAspects := map[int]bool{3: true, 7: true, 10: true}
	for _, sat := range nodes {
		if sat.Body != "Saturn" {
			continue
		}
		for _, dst := range nodes {
			if sat.PersonID == dst.PersonID {
				continue
			}
			houseDist := ((dst.Rasi - sat.Rasi + 12) % 12) + 1
			if !satAspects[houseDist] {
				continue
			}
			orb := math.Abs(sat.Degree - dst.Degree)
			weight := calcOrbWeight(orb, 9.0, 7.0, 3.0)
			label := fmt.Sprintf("Saturn Drishti H%d (%.2f°)", houseDist, roundTo2(orb))
			edges = append(edges, newKarmaEdge(
				len(edges)+1, sat, dst, EdgeDirDirected, KarmaRelTaskmaster,
				label, weight, "D1", ForceRoleRepel, -1, true, orb, houseDist,
			))
		}
	}
	return edges
}

// 3. Nakshatra Dispositor: หากอีกบุคคลมีดาว Star Lord ที่โดดเด่น (Exalted/Moolatrikona/Own Sign/Friend ใน D1 หรือ D9 หรือ Vargottama)
// ลาก Directed Edge จาก Star Lord ของอีกคน -> ดาว/ลัคนาเจ้าของนักษัตร (น้ำหนัก 6.0–9.0)
func appendNakshatraDispositorEdges(edges []GephiKarmaEdge, nodes []GephiKarmaNode) []GephiKarmaEdge {
	byPersonBody := make(map[string][]GephiKarmaNode, 10)
	for _, n := range nodes {
		byPersonBody[n.Body] = append(byPersonBody[n.Body], n)
	}

	for _, target := range nodes {
		if target.NakshatraLord == "" {
			continue
		}
		for _, lordNode := range byPersonBody[target.NakshatraLord] {
			if lordNode.PersonID == target.PersonID {
				continue
			}
			weight, prominent := nakshatraLordProminenceWeight(lordNode)
			if !prominent {
				continue
			}
			score := 1
			forceRole := ForceRoleAttract
			if lordNode.GrahaType != GrahaTypeBenefic {
				score = -1
				forceRole = ForceRoleRepel
			}
			label := fmt.Sprintf("Nak Lord (%s: %s)", target.NakshatraName, lordNode.Dignity)
			edges = append(edges, newKarmaEdge(
				len(edges)+1, lordNode, target, EdgeDirDirected, KarmaRelDispositor,
				label, weight, "Nakshatra/D9", forceRole, score, false, 0, 0,
			))
		}
	}
	return edges
}

func nakshatraLordProminenceWeight(n GephiKarmaNode) (float64, bool) {
	best := math.Max(dignityProminenceScore(n.Dignity), dignityProminenceScore(n.D9Dignity))
	if best == 0 && !n.IsVargottama {
		return 0, false
	}
	if best == 0 && n.IsVargottama {
		best = 6.5
	}
	if n.IsVargottama {
		best += 0.5
	}
	if dignityProminenceScore(n.D9Dignity) >= 7.5 && dignityProminenceScore(n.Dignity) >= 6.0 {
		best += 0.5
	}
	if best > 9.0 {
		best = 9.0
	}
	return roundTo2(best), true
}

func dignityProminenceScore(dignity string) float64 {
	switch dignity {
	case "Exalted":
		return 8.5
	case "Moolatrikona", "Own Sign":
		return 7.5
	case "Friend":
		return 6.0
	default:
		return 0
	}
}

func newKarmaEdge(
	seq int,
	src, dst GephiKarmaNode,
	direction, relation, label string,
	weight float64,
	layer, forceRole string,
	score int,
	hasOrb bool,
	orb float64,
	drishtiHouse int,
) GephiKarmaEdge {
	edgeType, colorLabel, polarity, color := edgeStyleByScore(score)
	return GephiKarmaEdge{
		GephiEdge: GephiEdge{
			ID:           fmt.Sprintf("e%d", seq),
			Source:       src.ID,
			Target:       dst.ID,
			Type:         direction,
			Weight:       roundTo2(weight),
			Label:        label,
			Relation:     relation,
			EdgeType:     edgeType,
			Scope:        EdgeScopeInter,
			Polarity:     polarity,
			HasOrb:       hasOrb,
			Orb:          roundTo2(orb),
			DrishtiHouse: drishtiHouse,
			ColorLabel:   colorLabel,
			Color:        color,
		},
		Layer:     layer,
		ForceRole: forceRole,
	}
}
