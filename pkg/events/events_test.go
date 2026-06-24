package events

import (
	"testing"

	"github.com/mendsec/catnet-core/pkg/results"
)

func TestEventTypeValues(t *testing.T) {
	if ScanStarted != "scan_started" {
		t.Errorf("ScanStarted = %q, want %q", ScanStarted, "scan_started")
	}
	if HostDiscovered != "host_discovered" {
		t.Errorf("HostDiscovered = %q, want %q", HostDiscovered, "host_discovered")
	}
	if ScanProgress != "scan_progress" {
		t.Errorf("ScanProgress = %q, want %q", ScanProgress, "scan_progress")
	}
	if ScanCompleted != "scan_completed" {
		t.Errorf("ScanCompleted = %q, want %q", ScanCompleted, "scan_completed")
	}
}

func TestEventDataTypeAssertion(t *testing.T) {
	ev := Event{
		Type: HostDiscovered,
		Data: HostDiscoveredData{Host: results.HostResult{IP: "10.0.0.1", Alive: true}},
	}
	data, ok := ev.Data.(HostDiscoveredData)
	if !ok {
		t.Fatal("type assertion to HostDiscoveredData failed")
	}
	if data.Host.IP != "10.0.0.1" {
		t.Errorf("unexpected IP: %s", data.Host.IP)
	}

	evP := Event{Type: ScanProgress, Data: ProgressData{Ratio: 0.5}}
	pd, ok := evP.Data.(ProgressData)
	if !ok {
		t.Fatal("type assertion to ProgressData failed")
	}
	if pd.Ratio != 0.5 {
		t.Errorf("unexpected Ratio: %f", pd.Ratio)
	}
}
