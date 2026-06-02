package scanner

import (
	"sync"
	"testing"
)

func TestScanConcurrency(t *testing.T) {
	ips := []string{"192.0.2.1", "192.0.2.2", "192.0.2.3", "192.0.2.4", "192.0.2.5"}
	cfg := DefaultConfig()
	cfg.PingTimeoutMs = 10 // small timeout for test
	cfg.MaxThreads = 2

	var mu sync.Mutex
	var results []DeviceInfo

	err := StartScan(ips, cfg, func(d DeviceInfo) {
		mu.Lock()
		defer mu.Unlock()
		results = append(results, d)
	}, nil)

	if err != nil {
		t.Fatalf("StartScan failed: %v", err)
	}

	if len(results) != len(ips) {
		t.Errorf("Expected %d results, got %d", len(ips), len(results))
	}

	// Verify all ips are present
	ipMap := make(map[string]bool)
	for _, r := range results {
		ipMap[r.IP] = true
	}
	for _, ip := range ips {
		if !ipMap[ip] {
			t.Errorf("Expected IP %s in results", ip)
		}
	}
}
