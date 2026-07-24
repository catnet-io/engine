package scan

import (
	"context"
	"fmt"
	"sync"

	"github.com/catnet-io/engine/pkg/discovery"
	"github.com/catnet-io/engine/pkg/engine"
	"github.com/catnet-io/engine/pkg/events"
	"github.com/catnet-io/engine/pkg/ports"
	"github.com/catnet-io/engine/pkg/profile"
	"github.com/catnet-io/engine/pkg/results"
)

type Engine struct {
	mu     sync.Mutex
	cancel context.CancelFunc
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) ScanStream(ctx context.Context, ips []string, cfg profile.ScanProfile, eventChan chan<- events.Event) error {
	e.mu.Lock()
	if e.cancel != nil {
		e.mu.Unlock()
		return fmt.Errorf("scan already in progress")
	}
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.cancel = nil
		e.mu.Unlock()
		cancel()
	}()

	scanCfg := engine.ScanConfig{
		DefaultPorts: cfg.DefaultPorts,
		MaxThreads:   cfg.Concurrency,
		PingTimeoutMs: cfg.TimeoutMs,
		PortTimeoutMs: cfg.TimeoutMs,
	}
	scanCfg.Sanitize()

	onEvent := func(se engine.ScanEvent) {
		var ev events.Event
		switch se.Type {
		case engine.EventLifecycleStart:
			ev = events.Event{Type: events.ScanStarted}
		case engine.EventResult:
			hr := results.HostResult{}
			if se.Device != nil {
				hr = results.HostResult{
					IP:        se.Device.IP,
					Alive:     se.Device.IsAlive,
					Hostname:  se.Device.Hostname,
					MAC:       se.Device.MAC,
					OpenPorts: se.Device.OpenPorts,
				}
			}
			ev = events.Event{
				Type: events.HostDiscovered,
				Data: events.HostDiscoveredData{Host: hr},
			}
		case engine.EventProgress:
			ev = events.Event{
				Type: events.ScanProgress,
				Data: events.ProgressData{Ratio: se.Progress},
			}
		case engine.EventLifecycleComplete:
			ev = events.Event{Type: events.ScanCompleted}
		case engine.EventLifecycleCancel:
			return
		}
		select {
		case eventChan <- ev:
		case <-ctx.Done():
		}
	}

	_, err := engine.StartScan(ctx, ips, scanCfg, onEvent) //nolint:staticcheck // pkg/scan wraps pkg/engine.StartScan internally
	return err
}

func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cancel != nil {
		e.cancel()
	}
}

func Ping(ip string, timeoutMs int) bool {
	return discovery.Ping(context.Background(), ip, timeoutMs)
}

func ReverseDNS(ctx context.Context, ip string) string {
	return discovery.ReverseDNS(ctx, ip)
}

func GetMAC(ctx context.Context, ip string) string {
	return discovery.GetMAC(ctx, ip)
}

func ScanPorts(ip string, portsList []int, timeoutMs int) []int {
	var result []int
	for p := range ports.ScanPorts(context.Background(), ip, portsList, timeoutMs) {
		result = append(result, p)
	}
	return result
}
