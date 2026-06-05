package engine

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mendsec/catnet-core/pkg/coreerr"
	"github.com/mendsec/catnet-core/pkg/results"
)

func TestScanConcurrency(t *testing.T) {
	ips := []string{"192.0.2.1", "192.0.2.2", "192.0.2.3", "192.0.2.4", "192.0.2.5"}
	cfg := DefaultConfig()
	cfg.PingTimeoutMs = 10 // small timeout for test
	cfg.MaxThreads = 2

	var mu sync.Mutex
	var eventDevices []results.DeviceInfo

	report, err := StartScan(context.Background(), ips, cfg, func(event ScanEvent) {
		if event.Type == EventResult && event.Device != nil {
			mu.Lock()
			eventDevices = append(eventDevices, *event.Device)
			mu.Unlock()
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

func TestScanCancellation(t *testing.T) {
	ips := make([]string, 100)
	for i := 0; i < 100; i++ {
		ips[i] = "192.0.2." + strconv.Itoa(i)
	}
	cfg := DefaultConfig()
	cfg.PingTimeoutMs = 1000
	cfg.MaxThreads = 2

	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancel almost immediately
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	_, err := StartScan(ctx, ips, cfg, nil)
	if err == nil {
		t.Fatalf("Expected cancellation error, got nil")
	}

	if !errors.Is(err, coreerr.ErrCancelled) {
		t.Errorf("Expected coreerr.ErrCancelled, got %v", err)
	}
}

func TestParallelScans(t *testing.T) {
	ips1 := []string{"192.0.2.1", "192.0.2.2"}
	ips2 := []string{"192.0.2.3", "192.0.2.4"}
	cfg := DefaultConfig()
	cfg.PingTimeoutMs = 10
	cfg.MaxThreads = 2

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := StartScan(context.Background(), ips1, cfg, nil)
		if err != nil {
			t.Errorf("Scan 1 failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		_, err := StartScan(context.Background(), ips2, cfg, nil)
		if err != nil {
			t.Errorf("Scan 2 failed: %v", err)
		}
	}()

	wg.Wait()
}

func BenchmarkStartScan(b *testing.B) {
	ips := make([]string, 100)
	for i := 0; i < 100; i++ {
		ips[i] = "127.0.0.1" // Localhost to avoid dropping packets
	}
	cfg := DefaultConfig()
	cfg.PingTimeoutMs = 10
	cfg.MaxThreads = 10
	cfg.DefaultPorts = []int{} // No port scanning for baseline

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = StartScan(context.Background(), ips, cfg, nil)
	}
}
