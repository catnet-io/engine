package exporter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mendsec/catnet-core/pkg/coreerr"
	"github.com/mendsec/catnet-core/pkg/results"
	"strconv"
	"strings"
)

// ExportJSON exporta resultados para formato JSON.
func ExportJSON(report *results.ScanReport) ([]byte, error) {
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode JSON: %v", coreerr.ErrExport, err)
	}
	return out, nil
}

// ExportXML exporta resultados para formato XML.
func ExportXML(report *results.ScanReport) ([]byte, error) {
	type XMLDevice struct {
		IP       string `xml:"ip"`
		Hostname string `xml:"hostname"`
		MAC      string `xml:"mac"`
		Status   string `xml:"status"`
	}
	type XMLResults struct {
		XMLName xml.Name    `xml:"results"`
		Devices []XMLDevice `xml:"device"`
	}
	res := XMLResults{}
	for _, d := range report.Devices {
		status := "Dead"
		if d.IsAlive {
			status = "Alive"
		}
		res.Devices = append(res.Devices, XMLDevice{
			IP: d.IP, Hostname: d.Hostname, MAC: d.MAC, Status: status,
		})
	}
	out, err := xml.MarshalIndent(res, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode XML: %v", coreerr.ErrExport, err)
	}
	return append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), out...), nil
}

// sanitizeCSVField limpa caracteres perigosos para prevenção de injeção CSV.
func sanitizeCSVField(field string) string {
	field = strings.ReplaceAll(field, "\n", " ")
	field = strings.ReplaceAll(field, "\r", " ")

	if len(field) == 0 {
		return field
	}

	fc := field[0]
	if fc == '=' || fc == '+' || fc == '-' || fc == '@' || fc == '\t' {
		return "'" + field
	}

	return field
}

// ExportCSV exporta resultados para formato CSV.
func ExportCSV(report *results.ScanReport) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"IP", "Hostname", "MAC", "Status", "Open Ports"}); err != nil {
		return nil, fmt.Errorf("%w: failed to write CSV header: %v", coreerr.ErrExport, err)
	}

	// ⚡ Bolt Optimization: Reuse string slice and avoid strings.Join overhead
	// This reduces allocations per device from ~6 to 1 and speeds up CSV generation
	row := make([]string, 5)
	var portsBuf []byte

	for _, d := range report.Devices {
		row[0] = sanitizeCSVField(d.IP)
		row[1] = sanitizeCSVField(d.Hostname)
		row[2] = sanitizeCSVField(d.MAC)
		if d.IsAlive {
			row[3] = "Alive"
		} else {
			row[3] = "Dead"
		}

		portsBuf = portsBuf[:0]
		for i, p := range d.OpenPorts {
			if i > 0 {
				portsBuf = append(portsBuf, ';')
			}
			portsBuf = strconv.AppendInt(portsBuf, int64(p), 10)
		}
		row[4] = string(portsBuf)

		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("%w: failed to write CSV record: %v", coreerr.ErrExport, err)
		}
	}
	writer.Flush()
	return buf.Bytes(), nil
}
