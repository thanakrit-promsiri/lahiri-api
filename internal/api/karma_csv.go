package api

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
)

// FormatKarmaNodesCSV แปลงรายการโหนดผังกรรมเป็นไฟล์ Nodes Table (.csv) สำหรับนำเข้า Gephi Data Laboratory
func FormatKarmaNodesCSV(nodes []GephiKarmaNode) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{
		"Id", "Label", "Person", "Person_Name", "Body", "Body_Thai",
		"Graha_Type", "Dignity_Score", "Force_Role", "Node_Size",
		"Rasi", "Rasi_Name", "Degree", "Longitude",
		"Nakshatra", "Nakshatra_Name", "Nakshatra_Pada", "Nakshatra_Lord",
		"Dignity", "House", "Is_Retrograde", "Is_Combust", "Is_Vargottama",
		"D9_Rasi", "D9_Rasi_Name", "D9_Degree", "D9_Longitude", "D9_Dignity", "Color",
	}); err != nil {
		return nil, err
	}

	for _, n := range nodes {
		if err := w.Write([]string{
			n.ID,
			n.Label,
			n.PersonID,
			n.PersonName,
			n.Body,
			n.BodyThai,
			n.GrahaType,
			strconv.Itoa(n.DignityScore),
			n.ForceRole,
			strconv.FormatFloat(n.NodeSize, 'f', -1, 64),
			strconv.Itoa(n.Rasi),
			n.RasiName,
			strconv.FormatFloat(n.Degree, 'f', -1, 64),
			strconv.FormatFloat(n.Longitude, 'f', -1, 64),
			strconv.Itoa(n.Nakshatra),
			n.NakshatraName,
			strconv.Itoa(n.NakshatraPada),
			n.NakshatraLord,
			n.Dignity,
			strconv.Itoa(n.House),
			strconv.FormatBool(n.IsRetrograde),
			strconv.FormatBool(n.IsCombust),
			strconv.FormatBool(n.IsVargottama),
			strconv.Itoa(n.D9Rasi),
			n.D9RasiName,
			strconv.FormatFloat(n.D9Degree, 'f', -1, 64),
			strconv.FormatFloat(n.D9Longitude, 'f', -1, 64),
			n.D9Dignity,
			formatRGBHex(n.Color),
		}); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FormatKarmaEdgesCSV แปลงรายการเส้นเชื่อมผังกรรมเป็นไฟล์ Edges Table (.csv) สำหรับนำเข้า Gephi Data Laboratory
func FormatKarmaEdgesCSV(edges []GephiKarmaEdge) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{
		"Id", "Source", "Target", "Type", "Weight", "Label",
		"Karma_Relation", "Layer", "Force_Role", "Edge_Type",
		"Scope", "Polarity", "Orb", "Drishti_House", "Color_Label", "Color",
	}); err != nil {
		return nil, err
	}

	for _, e := range edges {
		edgeDir := "Directed"
		if e.Type == "undirected" {
			edgeDir = "Undirected"
		}
		orbStr := ""
		if e.HasOrb {
			orbStr = strconv.FormatFloat(e.Orb, 'f', -1, 64)
		}
		drishtiStr := ""
		if e.DrishtiHouse > 0 {
			drishtiStr = strconv.Itoa(e.DrishtiHouse)
		}

		if err := w.Write([]string{
			e.ID,
			e.Source,
			e.Target,
			edgeDir,
			strconv.FormatFloat(e.Weight, 'f', -1, 64),
			e.Label,
			e.Relation,
			e.Layer,
			e.ForceRole,
			e.EdgeType,
			e.Scope,
			strconv.Itoa(e.Polarity),
			orbStr,
			drishtiStr,
			e.ColorLabel,
			formatRGBHex(e.Color),
		}); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FormatKarmaCSVZip รวมไฟล์ karma_nodes.csv และ karma_edges.csv ลงในไฟล์ ZIP เดียวสำหรับ Gephi Data Laboratory
func FormatKarmaCSVZip(graph GephiKarmaGraph) ([]byte, error) {
	nodesCSV, err := FormatKarmaNodesCSV(graph.Nodes)
	if err != nil {
		return nil, err
	}
	edgesCSV, err := FormatKarmaEdgesCSV(graph.Edges)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, entry := range []struct {
		name string
		data []byte
	}{
		{name: "karma_nodes.csv", data: nodesCSV},
		{name: "karma_edges.csv", data: edgesCSV},
	} {
		fw, err := zw.Create(entry.name)
		if err != nil {
			return nil, err
		}
		if _, err := fw.Write(entry.data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatRGBHex(c GephiColor) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
