package calculator

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	EdgeDirDirected   = "directed"
	EdgeDirMutual     = "undirected" // Gephi recognizes only "directed" and "undirected"
	EdgeDirUndirected = "undirected"

	EdgeRelYuti        = "Yuti"
	EdgeRelDrishti     = "Graha Drishti"
	EdgeRelParivartana = "Parivartana"
	EdgeRelD9Yuti      = "D9 Yuti"
	EdgeRelNakshatra   = "Nakshatra Compatibility"

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

// ComputeGephiGraph สร้างโครงข่ายโหนด (คนละ 10 โหนด) และเส้นเชื่อมความสัมพันธ์ (Yuti, Graha Drishti, Parivartana, D9 Synastry, Nakshatra)
func ComputeGephiGraph(req GephiReq) GephiGraph {
	nodes := make([]GephiNode, 0, len(req.Persons)*10)
	for i, person := range req.Persons {
		resolved := resolvePersonMeta(person, i)
		color := personColorPalette[i%len(personColorPalette)]
		nodes = append(nodes, buildPersonGephiNodes(resolved, color)...)
	}

	edges := buildGephiEdges(nodes)
	pairwiseMetrics := ComputePairwiseMetrics(nodes, edges)

	return GephiGraph{
		Nodes:           nodes,
		Edges:           edges,
		PairwiseMetrics: pairwiseMetrics,
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
	meta := horoscope.Metadata
	dasha := horoscope.VimshottariDasha.CurrentDasha

	// 1. โหนดลัคนา (Lagna)
	if len(horoscope.HousesAndAngles) > 0 {
		asc := horoscope.HousesAndAngles[0]
		nodes = append(nodes, buildLagnaNode(person, asc, meta, color))
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
		nodes = append(nodes, buildPlanetNode(person, body, p, meta, dasha, color))
	}

	return nodes
}

func buildLagnaNode(person GephiPersonReq, asc ApiHouse, meta ApiMetadata, color GephiColor) GephiNode {
	signIdx, degInSign := ParseLongitude(asc.Degree)
	nak, pada, _ := GetNakshatra(asc.Degree)
	d9SignIdx := GetSignIndexByName(asc.Navamsa.Sign)
	bala, size := dignityToBalaAndSize("Lagna")
	nodeID := person.ID + "_Lagna"

	return GephiNode{
		ID:             nodeID,
		Label:          nodeID,
		PersonID:       person.ID,
		PersonName:     person.Name,
		Body:           "Lagna",
		BodyThai:       asc.NameThai,
		Rasi:           signIdx + 1,
		RasiName:       asc.Position.Sign,
		RasiThai:       asc.Position.SignThai,
		Degree:         roundTo4(degInSign),
		Longitude:      roundTo4(asc.Degree),
		Nakshatra:      nak.Index,
		NakshatraName:  asc.Nakshatra,
		NakshatraPada:  pada,
		NakshatraLord:  asc.Lord,
		Dignity:        "Lagna",
		DignityThai:    "ลัคนา",
		Bala:           bala,
		NodeSize:       size,
		House:          asc.House,
		IsVargottama:   asc.Position.Sign != "" && asc.Position.Sign == asc.Navamsa.Sign,
		D9Rasi:         d9SignIdx + 1,
		D9RasiName:     asc.Navamsa.Sign,
		D9Dignity:      "Lagna",
		Color:          color,
		BirthDate:      meta.BirthDate,
		BirthTime:      meta.BirthTime,
		TimeZone:       meta.TimeZone,
		Latitude:       meta.Location.Latitude,
		BirthLongitude: meta.Location.Longitude,
		JulianDay:      roundTo4(meta.JulianDay),
		Ayanamsa:       meta.AyanamsaDegrees,
		NakshatraGana:  GetNakshatraGana(nak.Index),
		NakshatraYoni:  GetNakshatraYoni(nak.Index),
		NakshatraNadi:  GetNakshatraNadi(nak.Index),
	}
}

func buildPlanetNode(person GephiPersonReq, body string, p ApiPlanet, meta ApiMetadata, dasha CurrentDasha, color GephiColor) GephiNode {
	signIdx, degInSign := ParseLongitude(p.Longitude)
	nak, pada, _ := GetNakshatra(p.Longitude)
	d9SignIdx := GetSignIndexByName(p.Navamsa.Sign)
	d9Dig, _ := CalculateSignDignity(p.Name, d9SignIdx)
	bala, size := dignityToBalaAndSize(p.Dignity)
	nodeID := person.ID + "_" + body

	isDashaActive := false
	dashaRole := ""
	dashaPeriod := ""
	if dasha.MahadashaLord == body {
		isDashaActive = true
		dashaRole = "Mahadasha"
		dashaPeriod = fmt.Sprintf("MD: %s (%s..%s)", dasha.MahadashaLord, dasha.MahadashaStart, dasha.MahadashaEnd)
	} else if dasha.AntardashaLord == body {
		isDashaActive = true
		dashaRole = "Antardasha"
		dashaPeriod = fmt.Sprintf("AD: %s (%s..%s)", dasha.AntardashaLord, dasha.AntardashaStart, dasha.AntardashaEnd)
	} else if dasha.PratyantardashaLord == body {
		isDashaActive = true
		dashaRole = "Pratyantardasha"
		dashaPeriod = fmt.Sprintf("PD: %s (%s..%s)", dasha.PratyantardashaLord, dasha.PratyantardashaStart, dasha.PratyantardashaEnd)
	}

	return GephiNode{
		ID:             nodeID,
		Label:          nodeID,
		PersonID:       person.ID,
		PersonName:     person.Name,
		Body:           body,
		BodyThai:       p.NameThai,
		Rasi:           signIdx + 1,
		RasiName:       p.Position.Sign,
		RasiThai:       p.Position.SignThai,
		Degree:         roundTo4(degInSign),
		Longitude:      roundTo4(p.Longitude),
		Nakshatra:      nak.Index,
		NakshatraName:  p.Nakshatra,
		NakshatraPada:  pada,
		NakshatraLord:  p.Lord,
		Dignity:        p.Dignity,
		DignityThai:    p.DignityThai,
		Bala:           bala,
		NodeSize:       size,
		House:          p.House,
		IsRetrograde:   p.IsRetrograde,
		IsCombust:      p.IsCombust,
		IsVargottama:   p.Position.Sign != "" && p.Position.Sign == p.Navamsa.Sign,
		D9Rasi:         d9SignIdx + 1,
		D9RasiName:     p.Navamsa.Sign,
		D9Dignity:      d9Dig,
		Color:          color,
		BirthDate:      meta.BirthDate,
		BirthTime:      meta.BirthTime,
		TimeZone:       meta.TimeZone,
		Latitude:       meta.Location.Latitude,
		BirthLongitude: meta.Location.Longitude,
		JulianDay:      roundTo4(meta.JulianDay),
		Ayanamsa:       meta.AyanamsaDegrees,
		NakshatraGana:  GetNakshatraGana(nak.Index),
		NakshatraYoni:  GetNakshatraYoni(nak.Index),
		NakshatraNadi:  GetNakshatraNadi(nak.Index),
		IsDashaActive:  isDashaActive,
		DashaRole:      dashaRole,
		DashaPeriod:    dashaPeriod,
	}
}

func buildGephiEdges(nodes []GephiNode) []GephiEdge {
	edges := make([]GephiEdge, 0, len(nodes)*6)
	edges = appendYutiEdges(edges, nodes)
	edges = appendDrishtiEdges(edges, nodes)
	edges = appendParivartanaEdges(edges, nodes)
	edges = appendNakshatraSynastryEdges(edges, nodes)
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
			label := fmt.Sprintf("Yuti in %s / กุมใน%s (%.2f°)", u.RasiName, u.RasiThai, roundTo2(orb))
			edge := newGephiEdge(len(edges)+1, u, v, EdgeDirMutual, EdgeRelYuti, label, weight, true, orb, 0)
			edge.BaseWeight = 10.0
			edge.Layer = "D1"
			edge.AspectAngle = 0
			edge.Rasi = u.Rasi
			edge.RasiName = u.RasiName
			edge.RasiThai = u.RasiThai
			edge.ScoringFormula = "W = 10.0 - 7.0*(orb-3.0)/27.0 (orb<=3°: 10.0, orb>=30°: 3.0); Polarity = Attitude(A,B) + Attitude(B,A)"
			edges = append(edges, edge)
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
				rule, angle := getDrishtiRule(u.Body, aspectHouse)

				isMutual := false
				if vAspects, ok := drishtiAspects[v.Body]; ok {
					revHouse := ((u.Rasi - v.Rasi + 12) % 12) + 1
					for _, h := range vAspects {
						if h == revHouse {
							isMutual = true
							break
						}
					}
				}

				label := fmt.Sprintf("Drishti %d (%.2f°)", aspectHouse, roundTo2(orb))
				if isMutual {
					label = fmt.Sprintf("Mutual Drishti %d (%.2f°)", aspectHouse, roundTo2(orb))
				}

				edge := newGephiEdge(len(edges)+1, u, v, EdgeDirDirected, EdgeRelDrishti, label, weight, true, orb, aspectHouse)
				edge.BaseWeight = 7.0
				edge.Layer = "D1"
				edge.DrishtiRule = rule
				edge.AspectAngle = angle
				edge.IsMutual = isMutual
				edge.ScoringFormula = "W = 7.0 - 3.0*(orb-3.0)/27.0 (orb<=3°: 7.0, orb>=30°: 4.0); Polarity = Attitude(A,B) + Attitude(B,A)"
				edges = append(edges, edge)
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
				edge := newGephiEdge(len(edges)+1, u, v, EdgeDirMutual, EdgeRelParivartana, "Parivartana", 8.0, false, 0, 0)
				edge.BaseWeight = 8.0
				edge.Layer = "D1"
				edge.ScoringFormula = "W = 8.0 (Mutual Sign Exchange); Polarity = Attitude(A,B) + Attitude(B,A)"
				edges = append(edges, edge)
			}
		}
	}
	return edges
}

func appendD9SynastryEdges(edges []GephiEdge, nodes []GephiNode) []GephiEdge {
	for i := 0; i < len(nodes); i++ {
		u := nodes[i]
		if u.D9Rasi <= 0 {
			continue
		}
		for j := i + 1; j < len(nodes); j++ {
			v := nodes[j]
			if v.D9Rasi <= 0 || u.PersonID == v.PersonID {
				continue
			}
			if u.D9Rasi == v.D9Rasi && u.Rasi != v.Rasi {
				d9Orb := math.Abs(calcNavamsaDegreeInSign(u.Degree) - calcNavamsaDegreeInSign(v.Degree))
				weight := calcOrbWeight(d9Orb, 8.0, 4.0, 3.0)
				label := fmt.Sprintf("D9 Yuti in %s (%.2f°)", u.D9RasiName, roundTo2(d9Orb))
				edge := newGephiEdge(len(edges)+1, u, v, EdgeDirMutual, EdgeRelD9Yuti, label, weight, true, d9Orb, 0)
				edge.BaseWeight = 8.0
				edge.Layer = "D9"
				edge.AspectAngle = 0
				edge.ScoringFormula = "D9 Navamsa Conjunction: W = 8.0 - 4.0*(d9_orb-3.0)/27.0; Polarity = Attitude(A,B) + Attitude(B,A)"
				edges = append(edges, edge)
			}
		}
	}
	return edges
}

func appendNakshatraSynastryEdges(edges []GephiEdge, nodes []GephiNode) []GephiEdge {
	var moons []GephiNode
	for _, n := range nodes {
		if n.Body == "Moon" {
			moons = append(moons, n)
		}
	}

	for i := 0; i < len(moons); i++ {
		for j := i + 1; j < len(moons); j++ {
			mA, mB := moons[i], moons[j]
			if mA.PersonID == mB.PersonID {
				continue
			}
			ashtakoota := ComputeAshtakoota(mA, mB)
			weight := roundTo2(3.0 + (ashtakoota.TotalScore / 36.0) * 7.0)
			if weight > 10.0 {
				weight = 10.0
			}

			score := 0
			edgeType := EdgeTypeNeutral
			colorLabel := "Gray"
			color := GephiColor{R: 156, G: 163, B: 175}
			polarity := 0

			if ashtakoota.TotalScore >= 18.0 {
				score = 1
				edgeType = EdgeTypeHarmonious
				colorLabel = "Blue"
				color = GephiColor{R: 37, G: 99, B: 235}
				polarity = 1
			} else {
				score = -1
				edgeType = EdgeTypeChallenging
				colorLabel = "Red"
				color = GephiColor{R: 220, G: 38, B: 38}
				polarity = -1
			}

			taraAB := ComputeTaraBala(mA.Nakshatra, mB.Nakshatra)
			label := fmt.Sprintf("Ashtakoota %.1f/36 (%s)", ashtakoota.TotalScore, ashtakoota.Verdict)

			edges = append(edges, GephiEdge{
				ID:             fmt.Sprintf("e%d", len(edges)+1),
				Source:         mA.ID,
				Target:         mB.ID,
				Type:           EdgeDirMutual,
				Weight:         weight,
				Label:          label,
				Relation:       EdgeRelNakshatra,
				EdgeType:       edgeType,
				Scope:          EdgeScopeInter,
				Polarity:       polarity,
				ColorLabel:     colorLabel,
				Color:          color,
				Score:          score,
				BaseWeight:     10.0,
				ScoringFormula: "Ashtakoota 36 Points: Varna(1)+Vashya(2)+Tara(3)+Yoni(4)+Maitri(5)+Gana(6)+Bhakoot(7)+Nadi(8)",
				Layer:          "D1",
				TaraName:       taraAB.Name,
				TaraScore:      taraAB.Score,
				GanaCompat:     ashtakoota.GanaDescription,
				YoniCompat:     ashtakoota.YoniDescription,
				NadiCompat:     ashtakoota.NadiDescription,
			})
		}
	}
	return edges
}

func getDrishtiRule(body string, house int) (rule string, angle int) {
	switch body {
	case "Mars":
		switch house {
		case 4:
			return "Mars 4th Aspect (Chaturtha - 100% Special)", 90
		case 7:
			return "Mars 7th Aspect (Saptama - Full 180°)", 180
		case 8:
			return "Mars 8th Aspect (Ashtama - 100% Special)", 210
		}
	case "Jupiter":
		switch house {
		case 5:
			return "Jupiter 5th Aspect (Panchama / Trikona - 100% Special)", 120
		case 7:
			return "Jupiter 7th Aspect (Saptama - Full 180°)", 180
		case 9:
			return "Jupiter 9th Aspect (Navama / Trikona - 100% Special)", 240
		}
	case "Saturn":
		switch house {
		case 3:
			return "Saturn 3rd Aspect (Tritiya - 100% Special)", 60
		case 7:
			return "Saturn 7th Aspect (Saptama - Full 180°)", 180
		case 10:
			return "Saturn 10th Aspect (Dashama - 100% Special)", 270
		}
	case "Rahu":
		switch house {
		case 5:
			return "Rahu 5th Aspect (Trine 120°)", 120
		case 7:
			return "Rahu 7th Aspect (Opposition 180°)", 180
		case 9:
			return "Rahu 9th Aspect (Trine 240°)", 240
		}
	case "Ketu":
		switch house {
		case 5:
			return "Ketu 5th Aspect (Trine 120°)", 120
		case 7:
			return "Ketu 7th Aspect (Opposition 180°)", 180
		case 9:
			return "Ketu 9th Aspect (Trine 240°)", 240
		}
	default:
		return fmt.Sprintf("%s 7th Aspect (Saptama - Full 180°)", body), 180
	}
	return fmt.Sprintf("%s %dth Aspect", body, house), house * 30
}

func ComputePairwiseMetrics(nodes []GephiNode, edges []GephiEdge) []PairwiseSynastryMetric {
	personNames := make(map[string]string)
	for _, n := range nodes {
		personNames[n.PersonID] = n.PersonName
	}
	persons := make([]string, 0, len(personNames))
	for pid := range personNames {
		persons = append(persons, pid)
	}
	sort.Strings(persons)

	metrics := make([]PairwiseSynastryMetric, 0, len(persons)*(len(persons)-1)/2)
	nodeMap := make(map[string]GephiNode, len(nodes))
	for _, n := range nodes {
		nodeMap[n.ID] = n
	}

	for i := 0; i < len(persons); i++ {
		pA := persons[i]
		for j := i + 1; j < len(persons); j++ {
			pB := persons[j]

			total, harm, neu, chal := 0, 0, 0, 0
			for _, e := range edges {
				src, okS := nodeMap[e.Source]
				dst, okD := nodeMap[e.Target]
				if !okS || !okD {
					continue
				}
				if (src.PersonID == pA && dst.PersonID == pB) || (src.PersonID == pB && dst.PersonID == pA) {
					total++
					switch e.EdgeType {
					case EdgeTypeHarmonious:
						harm++
					case EdgeTypeChallenging:
						chal++
					default:
						neu++
					}
				}
			}

			ratio := 0.0
			if total > 0 {
				ratio = roundTo2(float64(harm) / float64(total))
			}

			moonA, okA := nodeMap[pA+"_Moon"]
			moonB, okB := nodeMap[pB+"_Moon"]
			ash := AshtakootaResult{Verdict: "N/A"}
			if okA && okB {
				ash = ComputeAshtakoota(moonA, moonB)
			}

			metrics = append(metrics, PairwiseSynastryMetric{
				PersonA:           pA,
				PersonB:           pB,
				PersonAName:       personNames[pA],
				PersonBName:       personNames[pB],
				TotalEdges:        total,
				HarmoniousEdges:   harm,
				NeutralEdges:      neu,
				ChallengingEdges:  chal,
				HarmonyRatio:      ratio,
				MoonTaraBala:      ash.TaraDescription,
				MoonGanaCompat:    ash.GanaDescription,
				MoonYoniCompat:    ash.YoniDescription,
				MoonNadiCompat:    ash.NadiDescription,
				AshtakootaScore:   ash.TotalScore,
				AshtakootaVerdict: ash.Verdict,
			})
		}
	}
	return metrics
}

func newGephiEdge(seq int, u, v GephiNode, direction, relation, label string, weight float64, hasOrb bool, orb float64, drishtiHouse int) GephiEdge {
	edgeType, colorLabel, polarity, color := classifyGrahaMaitri(u, v)
	scope := EdgeScopeInter
	if u.PersonID == v.PersonID {
		scope = EdgeScopeIntra
	}

	score := calculateMaitriRawScore(u, v)

	relationThai := relation
	switch relation {
	case EdgeRelYuti:
		relationThai = "กุม (Yuti)"
	case EdgeRelDrishti:
		relationThai = "เกณฑ์สายตา (Drishti)"
	case EdgeRelParivartana:
		relationThai = "สลับเรือนเกษตร (Parivartana)"
	case EdgeRelNakshatra:
		relationThai = "สมพงศ์นักษัตร (Nakshatra)"
	}

	return GephiEdge{
		ID:           fmt.Sprintf("e%d", seq),
		Source:       u.ID,
		Target:       v.ID,
		Type:         direction,
		Weight:       roundTo2(weight),
		Label:        label,
		Relation:     relation,
		RelationThai: relationThai,
		EdgeType:     edgeType,
		Scope:        scope,
		Polarity:     polarity,
		HasOrb:       hasOrb,
		Orb:          roundTo2(orb),
		DrishtiHouse: drishtiHouse,
		ColorLabel:   colorLabel,
		Color:        color,
		Score:        score,
	}
}

func calculateMaitriRawScore(u, v GephiNode) int {
	a := effectivePlanet(u)
	b := effectivePlanet(v)
	if a == b {
		if harmoniousSamePlanets[a] {
			return 1
		}
		return 0
	}
	return grahaAttitude(a, b) + grahaAttitude(b, a)
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
