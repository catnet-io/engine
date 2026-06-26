package scan

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mendsec/catnet-core/pkg/events"
	"github.com/mendsec/catnet-core/pkg/profile"
)

func TestEngineScansAllIPs(t *testing.T) {
	ips := []string{"192.0.2.1", "192.0.2.2", "192.0.2.3"}
	cfg := profile.DefaultProfile()
	cfg.TimeoutMs = 10
	cfg.Concurrency = 2

	engine := NewEngine()
	ch := make(chan events.Event, 64)

	var wg sync.WaitGroup
	wg.Add(1)

	var discovered int
	go func() {
		defer wg.Done()
		for ev := range ch {
			if ev.Type == events.HostDiscovered {
				discovered++
			}
		}
	}()

	err := engine.ScanStream(context.Background(), ips, cfg, ch)
	close(ch)
	wg.Wait()

	if err != nil {
		t.Fatalf("ScanStream failed: %v", err)
	}
	if discovered != len(ips) {
		t.Errorf("expected %d HostDiscovered events, got %d", len(ips), discovered)
	}
}

func TestEngineStop(t *testing.T) {
	ips := make([]string, 50)
	for i := range ips {
		ips[i] = "192.0.2.1"
	}
	cfg := profile.DefaultProfile()
	cfg.TimeoutMs = 500
	cfg.Concurrency = 2

	engine := NewEngine()
	ch := make(chan events.Event, 256)

	go func() {
		for range ch {
		}
	}()

	done := make(chan error, 1)
	go func() {
		done <- engine.ScanStream(context.Background(), ips, cfg, ch)
	}()

	time.Sleep(20 * time.Millisecond)
	engine.Stop()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("ScanStream did not return after Stop within timeout")
	}
	close(ch)
}

func TestEngineRejectsParallelScan(t *testing.T) {
	ips := []string{"192.0.2.1", "192.0.2.2"}
	cfg := profile.DefaultProfile()
	cfg.TimeoutMs = 200
	cfg.Concurrency = 1

	engine := NewEngine()
	ch1 := make(chan events.Event, 64)
	ch2 := make(chan events.Event, 64)

	go func() {
		for range ch1 {
		}
	}()
	go func() {
		for range ch2 {
		}
	}()

	var firstErr, secondErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		firstErr = engine.ScanStream(context.Background(), ips, cfg, ch1)
	}()

	time.Sleep(10 * time.Millisecond)
	secondErr = engine.ScanStream(context.Background(), ips, cfg, ch2)

	wg.Wait()
	close(ch1)
	close(ch2)

	if secondErr == nil {
		t.Error("expected error on parallel ScanStream call, got nil")
	}
	if firstErr != nil {
		t.Errorf("first ScanStream failed unexpectedly: %v", firstErr)
	}
}
