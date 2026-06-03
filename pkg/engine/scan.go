package engine

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/mendsec/catnet-core/pkg/discovery"
	"github.com/mendsec/catnet-core/pkg/ports"
	"github.com/mendsec/catnet-core/pkg/results"
)

// StartScan inicia uma varredura de rede concorrente.
func StartScan(ctx context.Context, ips []string, cfg ScanConfig, onResult func(results.DeviceInfo), onProgress func(float64)) error {
	total := len(ips)
	if total == 0 {
		return nil
	}

	ipChan := make(chan string, total)
	for _, ip := range ips {
		ipChan <- ip
	}
	close(ipChan)

	var wg sync.WaitGroup
	threads := cfg.MaxThreads
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
					di := results.DeviceInfo{IP: ip}
					di.IsAlive = discovery.Ping(ip, cfg.PingTimeoutMs)
					if di.IsAlive {
						di.Hostname = discovery.ReverseDNS(ip)
						di.MAC = discovery.GetMAC(ip)
						di.OpenPorts = ports.ScanPorts(ip, cfg.DefaultPorts, cfg.PortTimeoutMs)
						di.OpenPortsCount = len(di.OpenPorts)
					}
					if onResult != nil {
						onResult(di)
					}
					curr := atomic.AddInt32(&processed, 1)
					if onProgress != nil {
						onProgress(float64(curr) / float64(total))
					}
				}
			}
		}()
	}
	wg.Wait()
	return nil
}
