package engine

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mendsec/catnet-core/pkg/coreerr"
	"github.com/mendsec/catnet-core/pkg/discovery"
	"github.com/mendsec/catnet-core/pkg/ports"
	"github.com/mendsec/catnet-core/pkg/results"
)

// StartScan inicia uma varredura de rede concorrente e retorna um relatório completo.
func StartScan(ctx context.Context, ips []string, cfg ScanConfig, onEvent EventCallback) (*results.ScanReport, error) {
	report := results.NewScanReport()
	total := len(ips)
	report.Total = total
	report.Devices = make([]results.DeviceInfo, 0, total)
	
	if onEvent != nil {
		onEvent(ScanEvent{
			Type:    EventLifecycleStart,
			Message: "Scan started",
		})
	}

	if total == 0 {
		report.EndTime = time.Now()
		if onEvent != nil {
			onEvent(ScanEvent{
				Type:    EventLifecycleComplete,
				Message: "Scan completed (empty)",
			})
		}
		return report, nil
	}

	threads := cfg.MaxThreads
	if threads <= 0 {
		threads = 16
	}
	const maxAllowedThreads = 256
	if threads > maxAllowedThreads {
		threads = maxAllowedThreads
	}
	if _, ok := ctx.Deadline(); !ok {
		// Calcula timeout defensivo:
		// Se cada ping/port timeout levar o tempo máximo sequencialmente, 
		// com threads em paralelo. Apenas um fallback.
		maxTimePerHost := time.Duration(cfg.PingTimeoutMs) * time.Millisecond
		if len(cfg.DefaultPorts) > 0 {
			maxTimePerHost += time.Duration(len(cfg.DefaultPorts)) * time.Duration(cfg.PortTimeoutMs) * time.Millisecond
		}
		maxDuration := time.Duration(total) * maxTimePerHost / time.Duration(threads)
		maxDuration += time.Minute // Buffer de segurança
		// Limite fixo absoluto de 2 horas
		if maxDuration > 2*time.Hour {
			maxDuration = 2 * time.Hour
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, maxDuration)
		defer cancel()
	}

	ipChan := make(chan string, total)
	for _, ip := range ips {
		ipChan <- ip
	}
	close(ipChan)

	var wg sync.WaitGroup


	var processed int32
	var mu sync.Mutex

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				select {
				case <-ctx.Done():
					return
				default:
					di := results.DeviceInfo{IP: ip}
					di.IsAlive = discovery.Ping(ip, cfg.PingTimeoutMs)
					if di.IsAlive {
						di.Hostname = discovery.ReverseDNS(ip)
						di.MAC = discovery.GetMAC(ip)
						di.OpenPorts = ports.ScanPorts(ip, cfg.DefaultPorts, cfg.PortTimeoutMs)
						di.OpenPortsCount = len(di.OpenPorts)
					}

					mu.Lock()
					report.Devices = append(report.Devices, di)
					if di.IsAlive {
						report.Alive++
					}
					mu.Unlock()

					curr := atomic.AddInt32(&processed, 1)
					if onEvent != nil {
						onEvent(ScanEvent{
							Type:     EventResult,
							Device:   &di,
							Progress: float64(curr) / float64(total),
						})
						onEvent(ScanEvent{
							Type:     EventProgress,
							Device:   nil,
							Progress: float64(curr) / float64(total),
						})
					}
				}
			}
		}()
	}
	wg.Wait()
	report.EndTime = time.Now()
	
	if ctx.Err() == context.DeadlineExceeded {
		if onEvent != nil {
			onEvent(ScanEvent{Type: EventLifecycleCancel, Message: "Scan timeout"})
		}
		return report, fmt.Errorf("%w: scan reached timeout", coreerr.ErrTimeout)
	} else if ctx.Err() == context.Canceled {
		if onEvent != nil {
			onEvent(ScanEvent{Type: EventLifecycleCancel, Message: "Scan cancelled"})
		}
		return report, fmt.Errorf("%w: scan was cancelled", coreerr.ErrCancelled)
	}

	if onEvent != nil {
		onEvent(ScanEvent{Type: EventLifecycleComplete, Message: "Scan completed successfully"})
	}
	return report, nil
}
