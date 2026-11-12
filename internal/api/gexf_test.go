package api

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleExportGephiTriadGEXF(t *testing.T) {
	payload := `{
  "persons": [
    {
      "id": "P1",
      "name": "Alice",
      "year": 1988,
      "month": 12,
      "day": 23,
      "hour_thai": 10.35,
      "time_zone": 7.0,
      "latitude": 13.7563,
      "longitude": 100.5018,
      "house_system": "W"
    },
    {
      "id": "P2",
      "name": "Bob",
      "year": 1990,
      "month": 5,
      "day": 15,
      "hour_thai": 14.5,
      "time_zone": 7.0,
      "latitude": 13.7563,
      "longitude": 100.5018
    },
    {
      "id": "P3",
      "name": "Carol",
      "year": 1992,
      "month": 8,
      "day": 9,
      "hour_thai": 8.25,
      "time_zone": 7.0,
      "latitude": 18.7883,
      "longitude": 98.9853
    }
  ]
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleExportGephi(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", res.StatusCode, rec.Body.String())
	}

	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/gexf+xml") {
		t.Errorf("expected Content-Type application/gexf+xml, got %q", ct)
	}
	if cd := res.Header.Get("Content-Disposition"); !strings.Contains(cd, ".gexf") {
		t.Errorf("expected .gexf Content-Disposition, got %q", cd)
	}

	rawXML := rec.Body.String()
	for _, want := range []string{
		`xmlns="http://www.gexf.net/1.2draft"`,
		`xmlns:viz="http://www.gexf.net/1.2draft/viz"`,
		`defaultedgetype="directed"`,
		`<viz:size value="`,
		`<viz:color r="`,
		`id="P1_Lagna"`,
		`id="P2_Sun"`,
		`id="P3_Saturn"`,
		`type="directed"`,
		`type="undirected"`,
	} {
		if !strings.Contains(rawXML, want) {
			t.Errorf("expected GEXF output to contain %q", want)
		}
	}

	var parsed gexfDoc
	if err := xml.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse generated GEXF XML: %v", err)
	}

	if len(parsed.Graph.Nodes.Items) != 30 {
		t.Fatalf("expected 30 GEXF nodes for 3 persons, got %d", len(parsed.Graph.Nodes.Items))
	}
	if len(parsed.Graph.Edges.Items) < 50 {
		t.Fatalf("expected at least 50 GEXF edges, got %d", len(parsed.Graph.Edges.Items))
	}
	if len(parsed.Graph.Attributes) != 2 {
		t.Fatalf("expected 2 attribute blocks (node & edge), got %d", len(parsed.Graph.Attributes))
	}
}

func TestHandleExportGephiWrapperAndBuddhistEra(t *testing.T) {
	payload := `{
  "gephiReq": {
    "persons": [
      {
        "year": 2531,
        "month": 12,
        "day": 23,
        "hour_thai": 10.35,
        "time_zone": 7.0,
        "latitude": 13.7563,
        "longitude": 100.5018
      },
      {
        "year": 2533,
        "month": 5,
        "day": 15,
        "hour_thai": 14.5,
        "time_zone": 7.0,
        "latitude": 13.7563,
        "longitude": 100.5018
      }
    ]
  }
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	HandleExportGephi(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Result().StatusCode, rec.Body.String())
	}

	var parsed gexfDoc
	if err := xml.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse GEXF XML: %v", err)
	}
	if len(parsed.Graph.Nodes.Items) != 20 {
		t.Fatalf("expected 20 GEXF nodes for 2 persons, got %d", len(parsed.Graph.Nodes.Items))
	}
	if parsed.Graph.Nodes.Items[0].ID != "P1_Lagna" || parsed.Graph.Nodes.Items[10].ID != "P2_Lagna" {
		t.Errorf("expected auto-assigned P1_Lagna and P2_Lagna, got %q and %q",
			parsed.Graph.Nodes.Items[0].ID, parsed.Graph.Nodes.Items[10].ID)
	}
}

func TestHandleExportGephiHousesBipartiteGEXF(t *testing.T) {
	payload := `{
  "persons": [
    {
      "id": "P1",
      "name": "Male (Loei 1988)",
      "year": 1988,
      "month": 12,
      "day": 23,
      "hour_thai": 10.35,
      "time_zone": 7.0,
      "latitude": 17.59812056211902,
      "longitude": 101.72141606527812,
      "house_system": "W"
    },
    {
      "id": "P2",
      "name": "Female (Khon Kaen 1991)",
      "year": 1991,
      "month": 4,
      "day": 19,
      "hour_thai": 2.15,
      "time_zone": 7.0,
      "latitude": 16.46773511946009,
      "longitude": 102.82996133788332,
      "house_system": "W"
    },
    {
      "id": "P3",
      "name": "Female (Kamphaeng Phet 2003)",
      "year": 2003,
      "month": 4,
      "day": 25,
      "hour_thai": 9.7,
      "time_zone": 7.0,
      "latitude": 16.463037262887433,
      "longitude": 99.52840747732445,
      "house_system": "W"
    },
    {
      "id": "P4",
      "name": "Female (Bangkok 1990)",
      "year": 1990,
      "month": 9,
      "day": 26,
      "hour_thai": 4.816666666666667,
      "time_zone": 7.0,
      "latitude": 13.909166021400969,
      "longitude": 100.61788433453725,
      "house_system": "W"
    }
  ]
}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_houses", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	HandleExportGephiHouses(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", res.StatusCode, rec.Body.String())
	}

	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/gexf+xml") {
		t.Errorf("expected Content-Type application/gexf+xml, got %q", ct)
	}
	if cd := res.Header.Get("Content-Disposition"); !strings.Contains(cd, "vedic-houses-synastry.gexf") {
		t.Errorf("expected vedic-houses-synastry.gexf Content-Disposition, got %q", cd)
	}

	rawXML := rec.Body.String()
	for _, want := range []string{
		`id="P1_H1"`,
		`id="P4_H12"`,
		`id="P2_Sun"`,
		`for="node_kind" value="House"`,
		`for="node_kind" value="Planet"`,
		`for="relation" value="House Lord"`,
		`for="relation" value="Placed In"`,
	} {
		if !strings.Contains(rawXML, want) {
			t.Errorf("expected House GEXF output to contain %q", want)
		}
	}

	var parsed gexfDoc
	if err := xml.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse generated House GEXF XML: %v", err)
	}

	// 4 คน × 21 โหนด = 84 โหนด, House Lord = 48 + Placed In = 36*4 = 144 -> 192 เส้น
	if len(parsed.Graph.Nodes.Items) != 84 {
		t.Fatalf("expected 84 GEXF nodes for 4 persons, got %d", len(parsed.Graph.Nodes.Items))
	}
	if len(parsed.Graph.Edges.Items) != 192 {
		t.Fatalf("expected 192 GEXF edges for 4 persons, got %d", len(parsed.Graph.Edges.Items))
	}
}

func TestHandleExportGephiHousesErrors(t *testing.T) {
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/hora/gephi_houses", nil)
	recGet := httptest.NewRecorder()
	HandleExportGephiHouses(recGet, reqGet)
	if recGet.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed, got %d", recGet.Result().StatusCode)
	}

	badPayload := `{"persons": [{"year": 1988, "month": 12, "day": 23, "hour_thai": 10.35, "time_zone": 7.0, "latitude": 13.75, "longitude": 100.5}]}`
	reqBad := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_houses", bytes.NewBufferString(badPayload))
	recBad := httptest.NewRecorder()
	HandleExportGephiHouses(recBad, reqBad)
	if recBad.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request when < 2 persons, got %d", recBad.Result().StatusCode)
	}
}

func TestHandleExportGephiErrorPaths(t *testing.T) {
	cases := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"MethodNotAllowed", http.MethodGet, "", http.StatusMethodNotAllowed},
		{"InvalidJSON", http.MethodPost, "{bad-json", http.StatusBadRequest},
		{"LessThanTwoPersons", http.MethodPost, `{"persons":[{"year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7}]}`, http.StatusBadRequest},
		{"DuplicateID", http.MethodPost, `{"persons":[{"id":"A","year":1988,"month":12,"day":23,"hour_thai":10.35,"time_zone":7},{"id":"A","year":1990,"month":5,"day":15,"hour_thai":14.5,"time_zone":7}]}`, http.StatusBadRequest},
		{"InvalidBirthMonth", http.MethodPost, `{"persons":[{"year":1988,"month":13,"day":23,"hour_thai":10.35,"time_zone":7},{"year":1990,"month":5,"day":15,"hour_thai":14.5,"time_zone":7}]}`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/v1/hora/gephi", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			HandleExportGephi(rec, req)

			if rec.Result().StatusCode != tc.wantStatus {
				t.Errorf("%s: expected status %d, got %d (%s)", tc.name, tc.wantStatus, rec.Result().StatusCode, rec.Body.String())
			}
		})
	}
}

func TestHandleExportGephiKarmaFormatsAndErrors(t *testing.T) {
	payload := `{
  "persons": [
    {
      "id": "P1",
      "name": "Male (Loei 1988)",
      "year": 1988,
      "month": 12,
      "day": 23,
      "hour_thai": 10.35,
      "time_zone": 7.0,
      "latitude": 17.59812056211902,
      "longitude": 101.72141606527812,
      "house_system": "W"
    },
    {
      "id": "P2",
      "name": "Female (Khon Kaen 1991)",
      "year": 1991,
      "month": 4,
      "day": 19,
      "hour_thai": 2.15,
      "time_zone": 7.0,
      "latitude": 16.46773511946009,
      "longitude": 102.82996133788332,
      "house_system": "W"
    },
    {
      "id": "P3",
      "name": "Female (Kamphaeng Phet 2003)",
      "year": 2003,
      "month": 4,
      "day": 25,
      "hour_thai": 9.7,
      "time_zone": 7.0,
      "latitude": 16.463037262887433,
      "longitude": 99.52840747732445,
      "house_system": "W"
    },
    {
      "id": "P4",
      "name": "Female (Bangkok 1990)",
      "year": 1990,
      "month": 9,
      "day": 26,
      "hour_thai": 4.816666666666667,
      "time_zone": 7.0,
      "latitude": 13.909166021400969,
      "longitude": 100.61788433453725,
      "house_system": "W"
    }
  ]
}`

	// 1. Default GEXF
	req := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_karma", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()
	HandleExportGephiKarma(rec, req)
	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Result().StatusCode, rec.Body.String())
	}
	if cd := rec.Result().Header.Get("Content-Disposition"); !strings.Contains(cd, "vedic-karma-topology.gexf") {
		t.Errorf("unexpected Content-Disposition: %q", cd)
	}
	var parsed gexfDoc
	if err := xml.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse karma GEXF XML: %v", err)
	}
	if len(parsed.Graph.Nodes.Items) != 40 {
		t.Fatalf("expected 40 karma nodes for 4 persons, got %d", len(parsed.Graph.Nodes.Items))
	}
	if len(parsed.Graph.Edges.Items) == 0 {
		t.Fatal("expected karma edges in GEXF output")
	}

	// 2. ?format=zip -> verify 2 CSV files inside ZIP archive
	reqZip := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_karma?format=zip", bytes.NewBufferString(payload))
	recZip := httptest.NewRecorder()
	HandleExportGephiKarma(recZip, reqZip)
	if recZip.Result().StatusCode != http.StatusOK {
		t.Fatalf("zip format failed: %d", recZip.Result().StatusCode)
	}
	zr, err := zip.NewReader(bytes.NewReader(recZip.Body.Bytes()), int64(recZip.Body.Len()))
	if err != nil {
		t.Fatalf("invalid zip payload: %v", err)
	}
	if len(zr.File) != 2 {
		t.Fatalf("expected 2 CSV files in zip, got %d", len(zr.File))
	}

	// 3. ?format=nodes_csv -> verify 41 rows (header + 40 nodes) and exact columns
	reqNodes := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_karma?format=nodes_csv", bytes.NewBufferString(payload))
	recNodes := httptest.NewRecorder()
	HandleExportGephiKarma(recNodes, reqNodes)
	nodeRows, err := csv.NewReader(bytes.NewReader(recNodes.Body.Bytes())).ReadAll()
	if err != nil {
		t.Fatalf("failed parsing nodes csv: %v", err)
	}
	if len(nodeRows) != 41 || nodeRows[0][0] != "Id" || nodeRows[0][6] != "Graha_Type" || nodeRows[0][7] != "Dignity_Score" {
		t.Fatalf("unexpected nodes csv header/rows: rows=%d header=%+v", len(nodeRows), nodeRows[0])
	}

	// 4. ?format=edges_csv -> verify header and Directed/Undirected column
	reqEdges := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_karma?format=edges_csv", bytes.NewBufferString(payload))
	recEdges := httptest.NewRecorder()
	HandleExportGephiKarma(recEdges, reqEdges)
	edgeRows, err := csv.NewReader(bytes.NewReader(recEdges.Body.Bytes())).ReadAll()
	if err != nil || len(edgeRows) < 2 {
		t.Fatalf("failed parsing edges csv: err=%v rows=%d", err, len(edgeRows))
	}
	if edgeRows[0][1] != "Source" || edgeRows[0][2] != "Target" || edgeRows[0][3] != "Type" || edgeRows[0][4] != "Weight" {
		t.Errorf("unexpected edges csv header: %+v", edgeRows[0])
	}

	// 5. ?format=json -> verify structured JSON + CSV strings
	reqJSON := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_karma?format=json", bytes.NewBufferString(payload))
	recJSON := httptest.NewRecorder()
	HandleExportGephiKarma(recJSON, reqJSON)
	var jsonOut struct {
		Status   string          `json:"status"`
		Graph    GephiKarmaGraph `json:"graph"`
		NodesCSV string          `json:"nodes_csv"`
		EdgesCSV string          `json:"edges_csv"`
	}
	if err := json.Unmarshal(recJSON.Body.Bytes(), &jsonOut); err != nil {
		t.Fatalf("failed parsing karma JSON: %v", err)
	}
	if jsonOut.Status != "success" || len(jsonOut.Graph.Nodes) != 40 || jsonOut.NodesCSV == "" || jsonOut.EdgesCSV == "" {
		t.Errorf("unexpected karma JSON output: status=%q nodes=%d", jsonOut.Status, len(jsonOut.Graph.Nodes))
	}

	// 6. Error cases: unsupported format (400), method not allowed (405)
	reqBadFmt := httptest.NewRequest(http.MethodPost, "/api/v1/hora/gephi_karma?format=parquet", bytes.NewBufferString(payload))
	recBadFmt := httptest.NewRecorder()
	HandleExportGephiKarma(recBadFmt, reqBadFmt)
	if recBadFmt.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for unsupported format, got %d", recBadFmt.Result().StatusCode)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/hora/gephi_karma", nil)
	recGet := httptest.NewRecorder()
	HandleExportGephiKarma(recGet, reqGet)
	if recGet.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 for GET, got %d", recGet.Result().StatusCode)
	}
}
