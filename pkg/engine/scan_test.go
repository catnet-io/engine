package engine

import (
	"context"
	"testing"

	"github.com/mendsec/catnet-core/pkg/results"
)

func TestScanConcurrency(t *testing.T) {
	ips := []string{"192.0.2.1", "192.0.2.2", "192.0.2.3", "192.0.2.4", "192.0.2.5"}
	cfg := DefaultConfig()
	cfg.PingTimeoutMs = 10 // small timeout for test
	cfg.MaxThreads = 2

	var eventDevices []results.DeviceInfo

	report, err := StartScan(context.Background(), ips, cfg, func(event ScanEvent) {
		if event.Type == EventResult && event.Device != nil {
			eventDevices = append(eventDevices, *event.Device)
		}
	})

	if err != nil {
		t.Fatalf("StartScan failed: %v", err)
	}

	if report == nil {
		t.Fatalf("Expected report, got nil")
	}

	if len(report.Devices) != len(ips) {
		t.Errorf("Expected %d results in report, got %d", len(ips), len(report.Devices))
	}

	// Verify all ips are present in report
	ipMap := make(map[string]bool)
	for _, r := range report.Devices {
		ipMap[r.IP] = true
	}
	for _, ip := range ips {
		if !ipMap[ip] {
			t.Errorf("Expected IP %s in report", ip)
		}
	}
	
	if len(eventDevices) != len(ips) {
		t.Errorf("Expected %d results from events, got %d", len(ips), len(eventDevices))
	}
}
