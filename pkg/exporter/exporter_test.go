package exporter

import (
	"github.com/mendsec/catnet-core/pkg/results"
	"strings"
	"testing"
)

func TestExportCSV(t *testing.T) {
	devices := []results.DeviceInfo{
		{IP: "192.168.1.10", Hostname: "=cmd|' /C calc'!A0", MAC: "AA:BB", IsAlive: true, OpenPorts: []int{80}},
		{IP: "192.168.1.11", Hostname: "NormalHost", MAC: "CC:DD", IsAlive: false, OpenPorts: nil},
	}

	out, err := ExportCSV(devices)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	strOut := string(out)
	if !strings.Contains(strOut, "'=cmd|' /C calc'!A0") {
		t.Errorf("CSV injection not sanitized correctly, got: %s", strOut)
	}
	if !strings.Contains(strOut, "192.168.1.10") {
		t.Errorf("Expected IP in output")
	}
}

func TestExportXML(t *testing.T) {
	devices := []results.DeviceInfo{
		{IP: "192.168.1.10", Hostname: "HostA", MAC: "AA:BB", IsAlive: true},
	}

	out, err := ExportXML(devices)
	if err != nil {
		t.Fatalf("ExportXML failed: %v", err)
	}

	strOut := string(out)
	if !strings.Contains(strOut, "<ip>192.168.1.10</ip>") {
		t.Errorf("Expected IP in XML output")
	}
	if !strings.Contains(strOut, "<status>Alive</status>") {
		t.Errorf("Expected Alive status in XML output")
	}
}

func TestExportJSON(t *testing.T) {
	devices := []results.DeviceInfo{
		{IP: "192.168.1.10", Hostname: "HostA", MAC: "AA:BB", IsAlive: true},
	}

	out, err := ExportJSON(devices)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	strOut := string(out)
	if !strings.Contains(strOut, "\"ip\": \"192.168.1.10\"") {
		t.Errorf("Expected IP in JSON output")
	}
	if !strings.Contains(strOut, "\"isAlive\": true") {
		t.Errorf("Expected isAlive in JSON output")
	}
}
