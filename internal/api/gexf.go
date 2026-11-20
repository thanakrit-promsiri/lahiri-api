package api

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	gexfNS          = "http://gexf.net/1.3"
	gexfVizNS       = "http://gexf.net/1.3/viz"
	xmlSchemaInstNS = "http://www.w3.org/2001/XMLSchema-instance"
	gexfSchemaLoc   = "http://gexf.net/1.3 http://gexf.net/1.3/gexf.xsd"
)

type gexfDoc struct {
	XMLName        xml.Name  `xml:"gexf"`
	Xmlns          string    `xml:"xmlns,attr"`
	XmlnsViz       string    `xml:"xmlns:viz,attr"`
	XmlnsXsi       string    `xml:"xmlns:xsi,attr"`
	SchemaLocation string    `xml:"xsi:schemaLocation,attr"`
	Version        string    `xml:"version,attr"`
	Meta           gexfMeta  `xml:"meta"`
	Graph          gexfGraph `xml:"graph"`
}

type gexfMeta struct {
	LastModifiedDate string `xml:"lastmodifieddate,attr"`
	Creator          string `xml:"creator"`
	Description      string `xml:"description"`
}

type gexfGraph struct {
	Mode            string           `xml:"mode,attr"`
	DefaultEdgeType string           `xml:"defaultedgetype,attr"`
	Attributes      []gexfAttributes `xml:"attributes"`
	Nodes           gexfNodes        `xml:"nodes"`
	Edges           gexfEdges        `xml:"edges"`
}

type gexfAttributes struct {
	Class string          `xml:"class,attr"`
	Mode  string          `xml:"mode,attr"`
	Items []gexfAttribute `xml:"attribute"`
}

type gexfAttribute struct {
	ID    string `xml:"id,attr"`
	Title string `xml:"title,attr"`
	Type  string `xml:"type,attr"`
}

type gexfNodes struct {
	Items []gexfNode `xml:"node"`
}

type gexfNode struct {
	ID        string        `xml:"id,attr"`
	Label     string        `xml:"label,attr"`
	Size      gexfVizSize   `xml:"viz:size"`
	Color     gexfVizColor  `xml:"viz:color"`
	AttValues gexfAttValues `xml:"attvalues"`
}

type gexfEdges struct {
	Items []gexfEdge `xml:"edge"`
}

type gexfEdge struct {
	ID        string        `xml:"id,attr"`
	Source    string        `xml:"source,attr"`
	Target    string        `xml:"target,attr"`
	Type      string        `xml:"type,attr"`
	Weight    float64       `xml:"weight,attr"`
	Label     string        `xml:"label,attr"`
	Color     gexfVizColor  `xml:"viz:color"`
	AttValues gexfAttValues `xml:"attvalues"`
}

type gexfVizSize struct {
	Value float64 `xml:"value,attr"`
}

type gexfVizColor struct {
	R int `xml:"r,attr"`
	G int `xml:"g,attr"`
	B int `xml:"b,attr"`
}

type gexfAttValues struct {
	Items []gexfAttValue `xml:"attvalue"`
}

type gexfAttValue struct {
	For   string `xml:"for,attr"`
	Value string `xml:"value,attr"`
}

var nodeAttrDefs = []gexfAttribute{
	{ID: "person", Title: "Person", Type: "string"},
	{ID: "person_name", Title: "PersonName", Type: "string"},
	{ID: "body", Title: "Body", Type: "string"},
	{ID: "body_thai", Title: "BodyThai", Type: "string"},
	{ID: "rasi", Title: "Rasi", Type: "integer"},
	{ID: "rasi_name", Title: "RasiName", Type: "string"},
	{ID: "rasi_thai", Title: "RasiThai", Type: "string"},
	{ID: "degree", Title: "Degree", Type: "double"},
	{ID: "longitude", Title: "Longitude", Type: "double"},
	{ID: "nakshatra", Title: "Nakshatra", Type: "integer"},
	{ID: "nakshatra_name", Title: "NakshatraName", Type: "string"},
	{ID: "nakshatra_pada", Title: "NakshatraPada", Type: "integer"},
	{ID: "nakshatra_lord", Title: "NakshatraLord", Type: "string"},
	{ID: "dignity", Title: "Dignity", Type: "string"},
	{ID: "dignity_thai", Title: "DignityThai", Type: "string"},
	{ID: "bala", Title: "Bala", Type: "double"},
	{ID: "house", Title: "House", Type: "integer"},
	{ID: "is_retrograde", Title: "IsRetrograde", Type: "boolean"},
	{ID: "is_combust", Title: "IsCombust", Type: "boolean"},
	{ID: "is_vargottama", Title: "IsVargottama", Type: "boolean"},
	{ID: "d9_rasi", Title: "D9Rasi", Type: "integer"},
	{ID: "d9_rasi_name", Title: "D9RasiName", Type: "string"},
	{ID: "d9_dignity", Title: "D9Dignity", Type: "string"},

	// P0: Birth Metadata
	{ID: "birth_date", Title: "BirthDate", Type: "string"},
	{ID: "birth_time", Title: "BirthTime", Type: "string"},
	{ID: "time_zone", Title: "TimeZone", Type: "double"},
	{ID: "latitude", Title: "Latitude", Type: "double"},
	{ID: "birth_longitude", Title: "BirthLongitude", Type: "double"},
	{ID: "julian_day", Title: "JulianDay", Type: "double"},
	{ID: "ayanamsa", Title: "Ayanamsa", Type: "string"},

	// P1: Nakshatra Compatibility (Gana, Yoni, Nadi)
	{ID: "nakshatra_gana", Title: "NakshatraGana", Type: "string"},
	{ID: "nakshatra_yoni", Title: "NakshatraYoni", Type: "string"},
	{ID: "nakshatra_nadi", Title: "NakshatraNadi", Type: "string"},

	// P2: Vimshottari Dasha Active Timing
	{ID: "is_dasha_active", Title: "IsDashaActive", Type: "boolean"},
	{ID: "dasha_role", Title: "DashaRole", Type: "string"},
	{ID: "dasha_period", Title: "DashaPeriod", Type: "string"},
}

var houseNodeAttrDefs = []gexfAttribute{
	{ID: "node_kind", Title: "NodeKind", Type: "string"},
	{ID: "person", Title: "Person", Type: "string"},
	{ID: "person_name", Title: "PersonName", Type: "string"},
	{ID: "body", Title: "Body", Type: "string"},
	{ID: "label_thai", Title: "LabelThai", Type: "string"},
	{ID: "house", Title: "House", Type: "integer"},
	{ID: "house_name_thai", Title: "HouseNameThai", Type: "string"},
	{ID: "house_group", Title: "HouseGroup", Type: "string"},
	{ID: "rasi", Title: "Rasi", Type: "integer"},
	{ID: "rasi_name", Title: "RasiName", Type: "string"},
	{ID: "rasi_thai", Title: "RasiThai", Type: "string"},
	{ID: "sign_lord", Title: "SignLord", Type: "string"},
	{ID: "degree", Title: "Degree", Type: "double"},
	{ID: "nakshatra", Title: "Nakshatra", Type: "integer"},
	{ID: "dignity", Title: "Dignity", Type: "string"},
	{ID: "bala", Title: "Bala", Type: "double"},
	{ID: "is_retrograde", Title: "IsRetrograde", Type: "boolean"},
	{ID: "is_combust", Title: "IsCombust", Type: "boolean"},
	{ID: "is_vargottama", Title: "IsVargottama", Type: "boolean"},
}

var houseEdgeAttrDefs = []gexfAttribute{
	{ID: "relation", Title: "Relation", Type: "string"},
	{ID: "edge_type", Title: "EdgeType", Type: "string"},
	{ID: "scope", Title: "Scope", Type: "string"},
	{ID: "polarity", Title: "Polarity", Type: "integer"},
	{ID: "orb", Title: "Orb", Type: "double"},
	{ID: "color_label", Title: "ColorLabel", Type: "string"},
}

var edgeAttrDefs = []gexfAttribute{
	{ID: "relation", Title: "Relation", Type: "string"},
	{ID: "relation_thai", Title: "RelationThai", Type: "string"},
	{ID: "edge_type", Title: "EdgeType", Type: "string"},
	{ID: "scope", Title: "Scope", Type: "string"},
	{ID: "polarity", Title: "Polarity", Type: "integer"},
	{ID: "orb", Title: "Orb", Type: "double"},
	{ID: "color_label", Title: "ColorLabel", Type: "string"},
	{ID: "rasi", Title: "Rasi", Type: "integer"},
	{ID: "rasi_name", Title: "RasiName", Type: "string"},
	{ID: "rasi_thai", Title: "RasiThai", Type: "string"},
	{ID: "drishti_house", Title: "DrishtiAspect", Type: "integer"},
	{ID: "score", Title: "Score", Type: "integer"},
	{ID: "base_weight", Title: "BaseWeight", Type: "double"},
	{ID: "scoring_formula", Title: "ScoringFormula", Type: "string"},
	{ID: "layer", Title: "Layer", Type: "string"},
	{ID: "drishti_rule", Title: "DrishtiRule", Type: "string"},
	{ID: "is_mutual", Title: "IsMutual", Type: "boolean"},
	{ID: "aspect_angle", Title: "AspectAngle", Type: "integer"},
	{ID: "tara_name", Title: "TaraName", Type: "string"},
	{ID: "tara_score", Title: "TaraScore", Type: "integer"},
	{ID: "gana_compat", Title: "GanaCompat", Type: "string"},
	{ID: "yoni_compat", Title: "YoniCompat", Type: "string"},
	{ID: "nadi_compat", Title: "NadiCompat", Type: "string"},
}

var karmaNodeAttrDefs = []gexfAttribute{
	{ID: "person", Title: "Person", Type: "string"},
	{ID: "person_name", Title: "PersonName", Type: "string"},
	{ID: "body", Title: "Body", Type: "string"},
	{ID: "body_thai", Title: "BodyThai", Type: "string"},
	{ID: "graha_type", Title: "Graha_Type", Type: "string"},
	{ID: "dignity_score", Title: "Dignity_Score", Type: "integer"},
	{ID: "force_role", Title: "Force_Role", Type: "string"},
	{ID: "rasi", Title: "Rasi", Type: "integer"},
	{ID: "rasi_name", Title: "RasiName", Type: "string"},
	{ID: "degree", Title: "Degree", Type: "double"},
	{ID: "longitude", Title: "Longitude", Type: "double"},
	{ID: "nakshatra", Title: "Nakshatra", Type: "integer"},
	{ID: "nakshatra_name", Title: "NakshatraName", Type: "string"},
	{ID: "nakshatra_pada", Title: "NakshatraPada", Type: "integer"},
	{ID: "nakshatra_lord", Title: "NakshatraLord", Type: "string"},
	{ID: "dignity", Title: "Dignity", Type: "string"},
	{ID: "house", Title: "House", Type: "integer"},
	{ID: "is_retrograde", Title: "IsRetrograde", Type: "boolean"},
	{ID: "is_combust", Title: "IsCombust", Type: "boolean"},
	{ID: "is_vargottama", Title: "IsVargottama", Type: "boolean"},
	{ID: "d9_rasi", Title: "D9Rasi", Type: "integer"},
	{ID: "d9_rasi_name", Title: "D9RasiName", Type: "string"},
	{ID: "d9_degree", Title: "D9Degree", Type: "double"},
	{ID: "d9_longitude", Title: "D9Longitude", Type: "double"},
	{ID: "d9_dignity", Title: "D9Dignity", Type: "string"},
}

var karmaEdgeAttrDefs = append(append([]gexfAttribute(nil), houseEdgeAttrDefs...),
	gexfAttribute{ID: "drishti_house", Title: "DrishtiAspect", Type: "integer"},
	gexfAttribute{ID: "layer", Title: "Layer", Type: "string"},
	gexfAttribute{ID: "force_role", Title: "Force_Role", Type: "string"},
)

var nakshatraNodeAttrDefs = []gexfAttribute{
	{ID: "node_kind", Title: "NodeKind", Type: "string"},
	{ID: "person", Title: "Person", Type: "string"},
	{ID: "person_name", Title: "PersonName", Type: "string"},
	{ID: "body", Title: "Body", Type: "string"},
	{ID: "body_thai", Title: "BodyThai", Type: "string"},
	{ID: "rasi", Title: "Rasi", Type: "integer"},
	{ID: "rasi_name", Title: "RasiName", Type: "string"},
	{ID: "degree", Title: "Degree", Type: "double"},
	{ID: "longitude", Title: "Longitude", Type: "double"},
	{ID: "nakshatra", Title: "Nakshatra", Type: "integer"},
	{ID: "nakshatra_name", Title: "NakshatraName", Type: "string"},
	{ID: "nakshatra_thai", Title: "NakshatraThai", Type: "string"},
	{ID: "nakshatra_pada", Title: "NakshatraPada", Type: "integer"},
	{ID: "nakshatra_lord", Title: "NakshatraLord", Type: "string"},
	{ID: "nakshatra_lord_thai", Title: "NakshatraLordThai", Type: "string"},
	{ID: "occupied_count", Title: "OccupiedCount", Type: "integer"},
	{ID: "dignity", Title: "Dignity", Type: "string"},
	{ID: "dignity_thai", Title: "DignityThai", Type: "string"},
	{ID: "bala", Title: "Bala", Type: "double"},
	{ID: "house", Title: "House", Type: "integer"},
	{ID: "is_retrograde", Title: "IsRetrograde", Type: "boolean"},
	{ID: "is_combust", Title: "IsCombust", Type: "boolean"},
	{ID: "is_vargottama", Title: "IsVargottama", Type: "boolean"},
	{ID: "d9_rasi", Title: "D9Rasi", Type: "integer"},
	{ID: "d9_rasi_name", Title: "D9RasiName", Type: "string"},
	{ID: "d9_dignity", Title: "D9Dignity", Type: "string"},
}

var nakshatraEdgeAttrDefs = append(append([]gexfAttribute(nil), houseEdgeAttrDefs...),
	gexfAttribute{ID: "pada", Title: "Pada", Type: "integer"},
	gexfAttribute{ID: "nakshatra_lord", Title: "NakshatraLord", Type: "string"},
)

// FormatGEXF แปลงโครงสร้าง GephiGraph เป็นเอกสาร GEXF 1.2 XML ที่พร้อมนำเข้าซอฟต์แวร์ Gephi ทันที
func FormatGEXF(graph GephiGraph) ([]byte, error) {
	return formatGEXFWithDate(graph, time.Now().UTC().Format("2006-01-02"))
}

// FormatHouseGEXF แปลงโครงสร้าง GephiHouseGraph (ภพเรือน + ดาวเคราะห์) เป็นเอกสาร GEXF 1.2 XML
func FormatHouseGEXF(graph GephiHouseGraph) ([]byte, error) {
	return formatHouseGEXFWithDate(graph, time.Now().UTC().Format("2006-01-02"))
}

// FormatKarmaGEXF แปลงโครงสร้าง GephiKarmaGraph (ผังกรรม Topology of Karma) เป็นเอกสาร GEXF 1.2 XML
func FormatKarmaGEXF(graph GephiKarmaGraph) ([]byte, error) {
	return formatKarmaGEXFWithDate(graph, time.Now().UTC().Format("2006-01-02"))
}

// FormatNakshatraGEXF แปลงโครงสร้าง GephiNakshatraGraph (นักษัตรเป็น Hub กลางเชื่อมจุด/ดาวของทุกคน) เป็นเอกสาร GEXF 1.2 XML
func FormatNakshatraGEXF(graph GephiNakshatraGraph) ([]byte, error) {
	return formatNakshatraGEXFWithDate(graph, time.Now().UTC().Format("2006-01-02"))
}

func formatGEXFWithDate(graph GephiGraph, modifiedDate string) ([]byte, error) {
	nodes := make([]gexfNode, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodes = append(nodes, toGEXFNode(n))
	}

	description := "Vedic Astrology Network Analysis (Sidereal Lahiri / Whole Sign)"
	if len(graph.PairwiseMetrics) > 0 {
		var b strings.Builder
		b.WriteString(description)
		b.WriteString(" | Pairwise Synastry:")
		for _, m := range graph.PairwiseMetrics {
			b.WriteString(fmt.Sprintf(" [%s-%s: Harmony=%.1f%%, Ashtakoota=%.1f/36 (%s)]",
				m.PersonA, m.PersonB, m.HarmonyRatio*100, m.AshtakootaScore, m.AshtakootaVerdict))
		}
		description = b.String()
	}

	return marshalGEXFDoc(
		description,
		"mixed",
		nodeAttrDefs,
		edgeAttrDefs,
		nodes,
		toGEXFEdges(graph.Edges),
		modifiedDate,
	)
}

func formatHouseGEXFWithDate(graph GephiHouseGraph, modifiedDate string) ([]byte, error) {
	nodes := make([]gexfNode, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodes = append(nodes, toGEXFHouseNode(n))
	}

	return marshalGEXFDoc(
		"Vedic Astrology House-Planet Bipartite Network (Sidereal Lahiri / Whole Sign)",
		"directed",
		houseNodeAttrDefs,
		houseEdgeAttrDefs,
		nodes,
		toGEXFEdges(graph.Edges),
		modifiedDate,
	)
}

func formatKarmaGEXFWithDate(graph GephiKarmaGraph, modifiedDate string) ([]byte, error) {
	nodes := make([]gexfNode, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodes = append(nodes, toGEXFKarmaNode(n))
	}

	edges := make([]gexfEdge, 0, len(graph.Edges))
	for _, e := range graph.Edges {
		ge := toGEXFEdge(e.GephiEdge)
		ge.AttValues.Items = append(ge.AttValues.Items,
			gexfAttValue{For: "layer", Value: e.Layer},
			gexfAttValue{For: "force_role", Value: e.ForceRole},
		)
		edges = append(edges, ge)
	}

	return marshalGEXFDoc(
		"Vedic Astrology Topology of Karma (Sidereal Lahiri / True Nodes / D1+D9)",
		"mixed",
		karmaNodeAttrDefs,
		karmaEdgeAttrDefs,
		nodes,
		edges,
		modifiedDate,
	)
}

func toGEXFKarmaNode(n GephiKarmaNode) gexfNode {
	return gexfNode{
		ID:    n.ID,
		Label: n.Label,
		Size:  gexfVizSize{Value: n.NodeSize},
		Color: gexfVizColor{R: n.Color.R, G: n.Color.G, B: n.Color.B},
		AttValues: gexfAttValues{
			Items: []gexfAttValue{
				{For: "person", Value: n.PersonID},
				{For: "person_name", Value: n.PersonName},
				{For: "body", Value: n.Body},
				{For: "body_thai", Value: n.BodyThai},
				{For: "graha_type", Value: n.GrahaType},
				{For: "dignity_score", Value: strconv.Itoa(n.DignityScore)},
				{For: "force_role", Value: n.ForceRole},
				{For: "rasi", Value: strconv.Itoa(n.Rasi)},
				{For: "rasi_name", Value: n.RasiName},
				{For: "degree", Value: strconv.FormatFloat(n.Degree, 'f', -1, 64)},
				{For: "longitude", Value: strconv.FormatFloat(n.Longitude, 'f', -1, 64)},
				{For: "nakshatra", Value: strconv.Itoa(n.Nakshatra)},
				{For: "nakshatra_name", Value: n.NakshatraName},
				{For: "nakshatra_pada", Value: strconv.Itoa(n.NakshatraPada)},
				{For: "nakshatra_lord", Value: n.NakshatraLord},
				{For: "dignity", Value: n.Dignity},
				{For: "house", Value: strconv.Itoa(n.House)},
				{For: "is_retrograde", Value: strconv.FormatBool(n.IsRetrograde)},
				{For: "is_combust", Value: strconv.FormatBool(n.IsCombust)},
				{For: "is_vargottama", Value: strconv.FormatBool(n.IsVargottama)},
				{For: "d9_rasi", Value: strconv.Itoa(n.D9Rasi)},
				{For: "d9_rasi_name", Value: n.D9RasiName},
				{For: "d9_degree", Value: strconv.FormatFloat(n.D9Degree, 'f', -1, 64)},
				{For: "d9_longitude", Value: strconv.FormatFloat(n.D9Longitude, 'f', -1, 64)},
				{For: "d9_dignity", Value: n.D9Dignity},
			},
		},
	}
}

func formatNakshatraGEXFWithDate(graph GephiNakshatraGraph, modifiedDate string) ([]byte, error) {
	nodes := make([]gexfNode, 0, len(graph.Nodes))
	for _, n := range graph.Nodes {
		nodes = append(nodes, toGEXFNakshatraNode(n))
	}

	edges := make([]gexfEdge, 0, len(graph.Edges))
	for _, e := range graph.Edges {
		ge := toGEXFEdge(e.GephiEdge)
		ge.AttValues.Items = append(ge.AttValues.Items,
			gexfAttValue{For: "pada", Value: strconv.Itoa(e.Pada)},
			gexfAttValue{For: "nakshatra_lord", Value: e.NakshatraLord},
		)
		edges = append(edges, ge)
	}

	return marshalGEXFDoc(
		"Vedic Astrology Nakshatra Hub Network (Sidereal Lahiri / Occupied Nakshatra)",
		"directed",
		nakshatraNodeAttrDefs,
		nakshatraEdgeAttrDefs,
		nodes,
		edges,
		modifiedDate,
	)
}

func toGEXFNakshatraNode(n GephiNakshatraNode) gexfNode {
	return gexfNode{
		ID:    n.ID,
		Label: n.Label,
		Size:  gexfVizSize{Value: n.NodeSize},
		Color: gexfVizColor{R: n.Color.R, G: n.Color.G, B: n.Color.B},
		AttValues: gexfAttValues{
			Items: []gexfAttValue{
				{For: "node_kind", Value: n.NodeKind},
				{For: "person", Value: n.PersonID},
				{For: "person_name", Value: n.PersonName},
				{For: "body", Value: n.Body},
				{For: "body_thai", Value: n.BodyThai},
				{For: "rasi", Value: strconv.Itoa(n.Rasi)},
				{For: "rasi_name", Value: n.RasiName},
				{For: "degree", Value: strconv.FormatFloat(n.Degree, 'f', -1, 64)},
				{For: "longitude", Value: strconv.FormatFloat(n.Longitude, 'f', -1, 64)},
				{For: "nakshatra", Value: strconv.Itoa(n.Nakshatra)},
				{For: "nakshatra_name", Value: n.NakshatraName},
				{For: "nakshatra_thai", Value: n.NakshatraThai},
				{For: "nakshatra_pada", Value: strconv.Itoa(n.NakshatraPada)},
				{For: "nakshatra_lord", Value: n.NakshatraLord},
				{For: "nakshatra_lord_thai", Value: n.NakshatraLordThai},
				{For: "occupied_count", Value: strconv.Itoa(n.OccupiedCount)},
				{For: "dignity", Value: n.Dignity},
				{For: "dignity_thai", Value: n.DignityThai},
				{For: "bala", Value: strconv.FormatFloat(n.Bala, 'f', -1, 64)},
				{For: "house", Value: strconv.Itoa(n.House)},
				{For: "is_retrograde", Value: strconv.FormatBool(n.IsRetrograde)},
				{For: "is_combust", Value: strconv.FormatBool(n.IsCombust)},
				{For: "is_vargottama", Value: strconv.FormatBool(n.IsVargottama)},
				{For: "d9_rasi", Value: strconv.Itoa(n.D9Rasi)},
				{For: "d9_rasi_name", Value: n.D9RasiName},
				{For: "d9_dignity", Value: n.D9Dignity},
			},
		},
	}
}

func marshalGEXFDoc(
	description string,
	defaultEdgeType string,
	nodeAttrs []gexfAttribute,
	edgeAttrs []gexfAttribute,
	nodes []gexfNode,
	edges []gexfEdge,
	modifiedDate string,
) ([]byte, error) {
	if defaultEdgeType == "" {
		defaultEdgeType = "mixed"
	}
	doc := gexfDoc{
		Xmlns:          gexfNS,
		XmlnsViz:       gexfVizNS,
		XmlnsXsi:       xmlSchemaInstNS,
		SchemaLocation: gexfSchemaLoc,
		Version:        "1.3",
		Meta: gexfMeta{
			LastModifiedDate: modifiedDate,
			Creator:          "lahiri-api",
			Description:      description,
		},
		Graph: gexfGraph{
			Mode:            "static",
			DefaultEdgeType: defaultEdgeType,
			Attributes: []gexfAttributes{
				{Class: "node", Mode: "static", Items: nodeAttrs},
				{Class: "edge", Mode: "static", Items: edgeAttrs},
			},
			Nodes: gexfNodes{Items: nodes},
			Edges: gexfEdges{Items: edges},
		},
	}

	encoded, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, len(xml.Header)+len(encoded)+1)
	out = append(out, xml.Header...)
	out = append(out, encoded...)
	out = append(out, '\n')
	return out, nil
}

func toGEXFEdges(items []GephiEdge) []gexfEdge {
	edges := make([]gexfEdge, 0, len(items))
	for _, e := range items {
		edges = append(edges, toGEXFEdge(e))
	}
	return edges
}

func toGEXFNode(n GephiNode) gexfNode {
	return gexfNode{
		ID:    n.ID,
		Label: n.Label,
		Size:  gexfVizSize{Value: n.NodeSize},
		Color: gexfVizColor{R: n.Color.R, G: n.Color.G, B: n.Color.B},
		AttValues: gexfAttValues{
			Items: []gexfAttValue{
				{For: "person", Value: n.PersonID},
				{For: "person_name", Value: n.PersonName},
				{For: "body", Value: n.Body},
				{For: "body_thai", Value: n.BodyThai},
				{For: "rasi", Value: strconv.Itoa(n.Rasi)},
				{For: "rasi_name", Value: n.RasiName},
				{For: "rasi_thai", Value: n.RasiThai},
				{For: "degree", Value: strconv.FormatFloat(n.Degree, 'f', -1, 64)},
				{For: "longitude", Value: strconv.FormatFloat(n.Longitude, 'f', -1, 64)},
				{For: "nakshatra", Value: strconv.Itoa(n.Nakshatra)},
				{For: "nakshatra_name", Value: n.NakshatraName},
				{For: "nakshatra_pada", Value: strconv.Itoa(n.NakshatraPada)},
				{For: "nakshatra_lord", Value: n.NakshatraLord},
				{For: "dignity", Value: n.Dignity},
				{For: "dignity_thai", Value: n.DignityThai},
				{For: "bala", Value: strconv.FormatFloat(n.Bala, 'f', -1, 64)},
				{For: "house", Value: strconv.Itoa(n.House)},
				{For: "is_retrograde", Value: strconv.FormatBool(n.IsRetrograde)},
				{For: "is_combust", Value: strconv.FormatBool(n.IsCombust)},
				{For: "is_vargottama", Value: strconv.FormatBool(n.IsVargottama)},
				{For: "d9_rasi", Value: strconv.Itoa(n.D9Rasi)},
				{For: "d9_rasi_name", Value: n.D9RasiName},
				{For: "d9_dignity", Value: n.D9Dignity},

				// P0: Birth Metadata
				{For: "birth_date", Value: n.BirthDate},
				{For: "birth_time", Value: n.BirthTime},
				{For: "time_zone", Value: strconv.FormatFloat(n.TimeZone, 'f', -1, 64)},
				{For: "latitude", Value: strconv.FormatFloat(n.Latitude, 'f', -1, 64)},
				{For: "birth_longitude", Value: strconv.FormatFloat(n.BirthLongitude, 'f', -1, 64)},
				{For: "julian_day", Value: strconv.FormatFloat(n.JulianDay, 'f', -1, 64)},
				{For: "ayanamsa", Value: n.Ayanamsa},

				// P1: Nakshatra Compatibility
				{For: "nakshatra_gana", Value: n.NakshatraGana},
				{For: "nakshatra_yoni", Value: n.NakshatraYoni},
				{For: "nakshatra_nadi", Value: n.NakshatraNadi},

				// P2: Vimshottari Dasha Active Timing
				{For: "is_dasha_active", Value: strconv.FormatBool(n.IsDashaActive)},
				{For: "dasha_role", Value: n.DashaRole},
				{For: "dasha_period", Value: n.DashaPeriod},
			},
		},
	}
}

func toGEXFHouseNode(n GephiHouseNode) gexfNode {
	return gexfNode{
		ID:    n.ID,
		Label: n.Label,
		Size:  gexfVizSize{Value: n.NodeSize},
		Color: gexfVizColor{R: n.Color.R, G: n.Color.G, B: n.Color.B},
		AttValues: gexfAttValues{
			Items: []gexfAttValue{
				{For: "node_kind", Value: n.NodeKind},
				{For: "person", Value: n.PersonID},
				{For: "person_name", Value: n.PersonName},
				{For: "body", Value: n.Body},
				{For: "label_thai", Value: n.LabelThai},
				{For: "house", Value: strconv.Itoa(n.House)},
				{For: "house_name_thai", Value: n.HouseNameThai},
				{For: "house_group", Value: n.HouseGroup},
				{For: "rasi", Value: strconv.Itoa(n.Rasi)},
				{For: "rasi_name", Value: n.RasiName},
				{For: "rasi_thai", Value: n.RasiThai},
				{For: "sign_lord", Value: n.SignLord},
				{For: "degree", Value: strconv.FormatFloat(n.Degree, 'f', -1, 64)},
				{For: "nakshatra", Value: strconv.Itoa(n.Nakshatra)},
				{For: "dignity", Value: n.Dignity},
				{For: "bala", Value: strconv.FormatFloat(n.Bala, 'f', -1, 64)},
				{For: "is_retrograde", Value: strconv.FormatBool(n.IsRetrograde)},
				{For: "is_combust", Value: strconv.FormatBool(n.IsCombust)},
				{For: "is_vargottama", Value: strconv.FormatBool(n.IsVargottama)},
			},
		},
	}
}

func toGEXFEdge(e GephiEdge) gexfEdge {
	attrs := make([]gexfAttValue, 0, 20)
	attrs = append(attrs,
		gexfAttValue{For: "relation", Value: e.Relation},
	)
	if e.RelationThai != "" {
		attrs = append(attrs, gexfAttValue{For: "relation_thai", Value: e.RelationThai})
	}
	attrs = append(attrs,
		gexfAttValue{For: "edge_type", Value: e.EdgeType},
		gexfAttValue{For: "scope", Value: e.Scope},
		gexfAttValue{For: "polarity", Value: strconv.Itoa(e.Polarity)},
	)
	if e.HasOrb {
		attrs = append(attrs, gexfAttValue{For: "orb", Value: strconv.FormatFloat(e.Orb, 'f', -1, 64)})
	}
	if e.DrishtiHouse > 0 {
		attrs = append(attrs, gexfAttValue{For: "drishti_house", Value: strconv.Itoa(e.DrishtiHouse)})
	}
	attrs = append(attrs, gexfAttValue{For: "color_label", Value: e.ColorLabel})
	if e.Rasi > 0 {
		attrs = append(attrs,
			gexfAttValue{For: "rasi", Value: strconv.Itoa(e.Rasi)},
			gexfAttValue{For: "rasi_name", Value: e.RasiName},
			gexfAttValue{For: "rasi_thai", Value: e.RasiThai},
		)
	}

	// P0: Scoring Formula
	if e.Score != 0 {
		attrs = append(attrs, gexfAttValue{For: "score", Value: strconv.Itoa(e.Score)})
	}
	if e.BaseWeight > 0 {
		attrs = append(attrs, gexfAttValue{For: "base_weight", Value: strconv.FormatFloat(e.BaseWeight, 'f', -1, 64)})
	}
	if e.ScoringFormula != "" {
		attrs = append(attrs, gexfAttValue{For: "scoring_formula", Value: e.ScoringFormula})
	}

	// P1: Layer
	if e.Layer != "" {
		attrs = append(attrs, gexfAttValue{For: "layer", Value: e.Layer})
	}

	// P1: Graha Drishti Detail
	if e.DrishtiRule != "" {
		attrs = append(attrs, gexfAttValue{For: "drishti_rule", Value: e.DrishtiRule})
	}
	if e.IsMutual {
		attrs = append(attrs, gexfAttValue{For: "is_mutual", Value: strconv.FormatBool(e.IsMutual)})
	}
	if e.AspectAngle > 0 || e.Relation == "Yuti" || e.Relation == "D9 Yuti" {
		attrs = append(attrs, gexfAttValue{For: "aspect_angle", Value: strconv.Itoa(e.AspectAngle)})
	}

	// P1: Nakshatra Compatibility Details
	if e.TaraName != "" {
		attrs = append(attrs, gexfAttValue{For: "tara_name", Value: e.TaraName})
	}
	if e.TaraScore != 0 {
		attrs = append(attrs, gexfAttValue{For: "tara_score", Value: strconv.Itoa(e.TaraScore)})
	}
	if e.GanaCompat != "" {
		attrs = append(attrs, gexfAttValue{For: "gana_compat", Value: e.GanaCompat})
	}
	if e.YoniCompat != "" {
		attrs = append(attrs, gexfAttValue{For: "yoni_compat", Value: e.YoniCompat})
	}
	if e.NadiCompat != "" {
		attrs = append(attrs, gexfAttValue{For: "nadi_compat", Value: e.NadiCompat})
	}

	return gexfEdge{
		ID:        e.ID,
		Source:    e.Source,
		Target:    e.Target,
		Type:      e.Type,
		Weight:    e.Weight,
		Label:     e.Label,
		Color:     gexfVizColor{R: e.Color.R, G: e.Color.G, B: e.Color.B},
		AttValues: gexfAttValues{Items: attrs},
	}
}
