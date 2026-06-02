package export

import (
	"bytes"
	"github.com/mendsec/catnet-core/pkg/errors"
	"github.com/mendsec/catnet-core/pkg/results"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ExportJSON marshals a list of HostResults into pretty JSON.
func ExportJSON(hosts []results.HostResult) ([]byte, error) {
	out, err := json.MarshalIndent(hosts, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode JSON: %v", errors.ErrExportFailed, err)
	}
	return out, nil
}

// sanitizeCSVField prevents basic CSV injection attacks.
func sanitizeCSVField(field string) string {
	if len(field) > 0 {
		firstChar := field[0]
		if firstChar == '=' || firstChar == '+' || firstChar == '-' || firstChar == '@' || firstChar == '\t' || firstChar == '\r' {
			return "'" + field
		}
	}
	return field
}

// ExportCSV marshals a list of HostResults into CSV bytes.
func ExportCSV(hosts []results.HostResult) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"IP", "Hostname", "MAC", "Vendor", "Status", "Open Ports", "RTT_ms"})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrExportFailed, err)
	}

	for _, h := range hosts {
		status := "Dead"
		if h.Alive {
			status = "Alive"
		}

		var strPorts []string
		for _, p := range h.OpenPorts {
			strPorts = append(strPorts, strconv.Itoa(p))
		}
		ports := strings.Join(strPorts, ";")

		hostname := sanitizeCSVField(h.Hostname)
		vendor := sanitizeCSVField(h.Vendor)
		rtt := strconv.Itoa(h.RTTMs)

		err = writer.Write([]string{h.IP, hostname, h.MAC, vendor, status, ports, rtt})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrExportFailed, err)
		}
	}
	writer.Flush()
	return buf.Bytes(), nil
}
