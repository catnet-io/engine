package export

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mendsec/catnet-core/pkg/results"
)

func ExportJSON(devices []results.HostResult) ([]byte, error) {
	out, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}
	return out, nil
}

func sanitizeCSVField(field string) string {
	if len(field) > 0 {
		fc := field[0]
		if fc == '=' || fc == '+' || fc == '-' || fc == '@' || fc == '\t' || fc == '\r' {
			return "'" + field
		}
	}
	return field
}

func ExportCSV(devices []results.HostResult) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"IP", "Hostname", "MAC", "Status", "Open Ports"}); err != nil {
		return nil, err
	}
	for _, d := range devices {
		status := "Dead"
		if d.Alive {
			status = "Alive"
		}
		var strPorts []string
		for _, p := range d.OpenPorts {
			strPorts = append(strPorts, strconv.Itoa(p))
		}
		if err := writer.Write([]string{
			d.IP,
			sanitizeCSVField(d.Hostname),
			d.MAC,
			status,
			strings.Join(strPorts, ";"),
		}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), nil
}
