package scan

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/catnet-io/engine/pkg/events"
	"github.com/catnet-io/engine/pkg/profile"
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

	go func() { for range ch1 {} }()
	go func() { for range ch2 {} }()

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

func TestNoGoroutineLeakOnPrematureCancel(t *testing.T) {
	ips := make([]string, 100)
	for i := range ips {
		// Just some random IP that might or might not resolve fast
		ips[i] = "10.0.0.1"
	}
	cfg := profile.DefaultProfile()
	cfg.TimeoutMs = 500
	cfg.Concurrency = 10

	before := runtime.NumGoroutine()

	engine := NewEngine()
	// Unbuffered channel to easily trigger blocking and simulate unconsumed events if consumer dies
	ch := make(chan events.Event)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Read a few events to let things start
		count := 0
		for range ch {
			count++
			if count == 2 {
				break
			}
		}
		// Allow some goroutines to spawn
		time.Sleep(10 * time.Millisecond)
		// Cancel abruptly simulating consumer exit/error
		cancel()
	}()

	_ = engine.ScanStream(ctx, ips, cfg, ch)

	// Wait for background routines to settle
	time.Sleep(100 * time.Millisecond)
	runtime.GC()

	after := runtime.NumGoroutine()
	leak := after - before
	if leak > 3 {
		t.Errorf("goroutine leak detected: before=%d, after=%d, leak=%d", before, after, leak)
	}
}
