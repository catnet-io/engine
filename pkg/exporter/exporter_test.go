package exporter

import (
	"strconv"
	"strings"
	"testing"

	"github.com/mendsec/catnet-core/pkg/results"
)

func TestExportCSV(t *testing.T) {
	report := results.NewScanReport()
	report.Devices = []results.DeviceInfo{
		{IP: "+malicious", Hostname: "=cmd|' /C calc'!A0", MAC: "-mac", IsAlive: true, OpenPorts: []int{80}},
		{IP: "192.168.1.11", Hostname: "NormalHost", MAC: "CC:DD", IsAlive: false, OpenPorts: nil},
	}

	out, err := ExportCSV(report)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	strOut := string(out)
	if !strings.Contains(strOut, "'=cmd|' /C calc'!A0") {
		t.Errorf("CSV injection not sanitized correctly for Hostname, got: %s", strOut)
	}
	if !strings.Contains(strOut, "'+malicious") {
		t.Errorf("CSV injection not sanitized correctly for IP, got: %s", strOut)
	}
	if !strings.Contains(strOut, "'-mac") {
		t.Errorf("CSV injection not sanitized correctly for MAC, got: %s", strOut)
	}
	if !strings.Contains(strOut, "192.168.1.11") {
		t.Errorf("Expected IP in output")
	}
}

func TestExportXML(t *testing.T) {
	report := results.NewScanReport()
	report.Devices = []results.DeviceInfo{
		{IP: "192.168.1.10", Hostname: "HostA", MAC: "AA:BB", IsAlive: true},
	}

	out, err := ExportXML(report)
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
	report := results.NewScanReport()
	report.Devices = []results.DeviceInfo{
		{IP: "192.168.1.10", Hostname: "HostA", MAC: "AA:BB", IsAlive: true},
	}

	out, err := ExportJSON(report)
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

func getBenchmarkReport() *results.ScanReport {
	report := results.NewScanReport()
	for i := 0; i < 1000; i++ {
		report.Devices = append(report.Devices, results.DeviceInfo{
			IP:        "192.168.1." + strconv.Itoa(i%255),
			MAC:       "AA:BB:CC:DD:EE:FF",
			Hostname:  "host-" + strconv.Itoa(i),
			IsAlive:   i%2 == 0,
			OpenPorts: []int{80, 443, 8080},
		})
	}
	report.Total = 1000
	return report
}

func BenchmarkExportJSON(b *testing.B) {
	report := getBenchmarkReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExportJSON(report)
	}
}

func BenchmarkExportCSV(b *testing.B) {
	report := getBenchmarkReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExportCSV(report)
	}
}

func BenchmarkExportXML(b *testing.B) {
	report := getBenchmarkReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExportXML(report)
	}
}
