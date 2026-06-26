package export

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"

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

	// ⚡ Bolt Optimization: Reuse string slice and avoid strings.Join overhead
	// Reduces multiple allocations per item to a single allocation for the ports string.
	row := make([]string, 5)
	var portsBuf []byte

	for _, d := range devices {
		status := "Dead"
		if d.Alive {
			status = "Alive"
		}

		portsBuf = portsBuf[:0]
		for i, p := range d.OpenPorts {
			if i > 0 {
				portsBuf = append(portsBuf, ';')
			}
			portsBuf = strconv.AppendInt(portsBuf, int64(p), 10)
		}

		row[0] = d.IP
		row[1] = sanitizeCSVField(d.Hostname)
		row[2] = d.MAC
		row[3] = status
		row[4] = string(portsBuf)

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), nil
}
