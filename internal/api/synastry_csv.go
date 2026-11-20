package api

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"strconv"
)

// FormatGephiNodesCSV แปลงรายการโหนด Synastry เป็นไฟล์ Nodes Table (.csv) สำหรับ Cosmograph / Gephi Data Laboratory
func FormatGephiNodesCSV(nodes []GephiNode) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{
		"id", "label", "person", "person_name", "body", "body_thai",
		"rasi", "rasi_name", "rasi_thai", "degree", "longitude",
		"nakshatra", "nakshatra_name", "nakshatra_pada", "nakshatra_lord",
		"dignity", "dignity_thai", "bala", "size", "color", "house",
		"is_retrograde", "is_combust", "is_vargottama",
		"d9_rasi", "d9_rasi_name", "d9_dignity",
		"is_dasha_active", "dasha_role", "dasha_period",
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
			strconv.Itoa(n.Rasi),
			n.RasiName,
			n.RasiThai,
			strconv.FormatFloat(n.Degree, 'f', -1, 64),
			strconv.FormatFloat(n.Longitude, 'f', -1, 64),
			strconv.Itoa(n.Nakshatra),
			n.NakshatraName,
			strconv.Itoa(n.NakshatraPada),
			n.NakshatraLord,
			n.Dignity,
			n.DignityThai,
			strconv.FormatFloat(n.Bala, 'f', -1, 64),
			strconv.FormatFloat(n.NodeSize, 'f', -1, 64),
			formatRGBHex(n.Color),
			strconv.Itoa(n.House),
			strconv.FormatBool(n.IsRetrograde),
			strconv.FormatBool(n.IsCombust),
			strconv.FormatBool(n.IsVargottama),
			strconv.Itoa(n.D9Rasi),
			n.D9RasiName,
			n.D9Dignity,
			strconv.FormatBool(n.IsDashaActive),
			n.DashaRole,
			n.DashaPeriod,
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

// FormatGephiLinksCSV แปลงรายการเส้นเชื่อม Synastry เป็นไฟล์ Links Table (.csv) สำหรับ Cosmograph / Gephi Data Laboratory
func FormatGephiLinksCSV(edges []GephiEdge) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write([]string{
		"source", "target", "id", "label", "type", "weight",
		"relation", "relation_thai", "edge_type", "scope", "polarity",
		"orb", "rasi", "rasi_name", "rasi_thai", "drishti_house", "color",
		"scoring_formula",
	}); err != nil {
		return nil, err
	}

	for _, e := range edges {
		orbStr := ""
		if e.HasOrb {
			orbStr = strconv.FormatFloat(e.Orb, 'f', -1, 64)
		}
		drishtiStr := ""
		if e.DrishtiHouse > 0 {
			drishtiStr = strconv.Itoa(e.DrishtiHouse)
		}
		rasiStr := ""
		if e.Rasi > 0 {
			rasiStr = strconv.Itoa(e.Rasi)
		}

		if err := w.Write([]string{
			e.Source,
			e.Target,
			e.ID,
			e.Label,
			e.Type,
			strconv.FormatFloat(e.Weight, 'f', -1, 64),
			e.Relation,
			e.RelationThai,
			e.EdgeType,
			e.Scope,
			strconv.Itoa(e.Polarity),
			orbStr,
			rasiStr,
			e.RasiName,
			e.RasiThai,
			drishtiStr,
			formatRGBHex(e.Color),
			e.ScoringFormula,
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

// FormatGephiCSVZip รวมไฟล์ nodes.csv และ links.csv ลงในไฟล์ ZIP เดียวสำหรับนำเข้า Cosmograph (run.cosmograph.app) หรือ Gephi
func FormatGephiCSVZip(graph GephiGraph) ([]byte, error) {
	nodesCSV, err := FormatGephiNodesCSV(graph.Nodes)
	if err != nil {
		return nil, err
	}
	linksCSV, err := FormatGephiLinksCSV(graph.Edges)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, entry := range []struct {
		name string
		data []byte
	}{
		{name: "nodes.csv", data: nodesCSV},
		{name: "links.csv", data: linksCSV},
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
