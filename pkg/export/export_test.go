package export

import (
	"strings"
	"testing"

	"github.com/mendsec/catnet-core/pkg/results"
)

func TestExportJSON(t *testing.T) {
	devices := []results.HostResult{
		{IP: "192.168.1.1", Alive: true, Hostname: "router", OpenPorts: []int{80, 443}},
	}
	out, err := ExportJSON(devices)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, `"ip": "192.168.1.1"`) {
		t.Errorf("missing IP in JSON: %s", s)
	}
	if !strings.Contains(s, `"alive": true`) {
		t.Errorf("missing alive field in JSON: %s", s)
	}
	if !strings.Contains(s, `"open_ports"`) {
		t.Errorf("missing open_ports field in JSON: %s", s)
	}
}

func TestExportCSVInjection(t *testing.T) {
	devices := []results.HostResult{
		{IP: "10.0.0.1", Alive: true, Hostname: "=cmd|' /C calc'!A0", OpenPorts: nil},
	}
	out, err := ExportCSV(devices)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "=cmd") && !strings.Contains(s, "'=cmd") {
		t.Errorf("CSV injection not sanitized: %s", s)
	}
}

func TestExportCSVStatus(t *testing.T) {
	devices := []results.HostResult{
		{IP: "10.0.0.1", Alive: true},
		{IP: "10.0.0.2", Alive: false},
	}
	out, err := ExportCSV(devices)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "Alive") {
		t.Errorf("missing Alive status in CSV")
	}
	if !strings.Contains(s, "Dead") {
		t.Errorf("missing Dead status in CSV")
	}
}
