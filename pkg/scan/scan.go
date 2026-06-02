package scan

import (
	"catnet-core/pkg/errors"
	"catnet-core/pkg/events"
	"catnet-core/pkg/profile"
	"catnet-core/pkg/results"
	"context"
	"sync"
	"sync/atomic"
)

// Engine is the main scanner orchestrator.
type Engine struct {
	mu         sync.Mutex
	isScanning atomic.Bool
	cancelScan context.CancelFunc
}

// NewEngine creates a new Engine instance.
func NewEngine() *Engine {
	return &Engine{}
}

// ScanStream executes a parallel scan over the provided IPs, emitting events to the eventChan.
func (e *Engine) ScanStream(ctx context.Context, ips []string, prof profile.ScanProfile, eventChan chan<- events.Event) error {
	if !e.isScanning.CompareAndSwap(false, true) {
		return errors.ErrScanInProgress
	}
	defer e.isScanning.Store(false)

	ctx, cancel := context.WithCancel(ctx)
	e.mu.Lock()
	e.cancelScan = cancel
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.cancelScan = nil
		e.mu.Unlock()
		cancel()
	}()

	total := len(ips)
	if total == 0 {
		return nil
	}

	if eventChan != nil {
		eventChan <- events.Event{Type: events.ScanStarted, Data: total}
	}

	ipChan := make(chan string, total)
	for _, ip := range ips {
		ipChan <- ip
	}
	close(ipChan)

	var wg sync.WaitGroup
	threads := prof.Concurrency
	if threads <= 0 {
		threads = 16
	}

	const maxAllowedThreads = 256
	if threads > maxAllowedThreads {
		threads = maxAllowedThreads
	}

	var processed int32

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				select {
				case <-ctx.Done():
					return
				default:
					if eventChan != nil {
						eventChan <- events.Event{Type: events.TargetQueued, Data: ip}
					}

					di := results.HostResult{IP: ip}
					di.Alive = Ping(ip, prof.TimeoutMs)

					if di.Alive {
						if prof.ResolveDNS {
							di.Hostname = ReverseDNS(ip)
						}
						if prof.ResolveMAC {
							di.MAC = GetMAC(ip)
						}
						if len(prof.Ports) > 0 {
							di.OpenPorts = ScanPorts(ip, prof.Ports, prof.TimeoutMs)
						}
					}

					if eventChan != nil {
						eventChan <- events.Event{Type: events.HostDiscovered, Data: events.HostDiscoveredData{Host: di}}
					}

					curr := atomic.AddInt32(&processed, 1)
					if eventChan != nil {
						eventChan <- events.Event{
							Type: events.ScanProgress,
							Data: events.ProgressData{
								Processed: int(curr),
								Total:     total,
								Ratio:     float64(curr) / float64(total),
							},
						}
					}
				}
			}
		}()
	}

	wg.Wait()
	if eventChan != nil {
		eventChan <- events.Event{Type: events.ScanCompleted, Data: nil}
	}
	return nil
}

// Stop cancels any ongoing scan gracefully.
func (e *Engine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cancelScan != nil {
		e.cancelScan()
	}
}
