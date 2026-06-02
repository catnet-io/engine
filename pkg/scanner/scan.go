package scanner

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	isScanning atomic.Bool
	scanMu     sync.Mutex
	cancelScan context.CancelFunc
)

// StartScan inicia uma varredura de rede concorrente.
func StartScan(ips []string, cfg ScanConfig, onResult func(DeviceInfo), onProgress func(float64)) error {
	if !isScanning.CompareAndSwap(false, true) {
		return fmt.Errorf("scan already in progress")
	}
	defer isScanning.Store(false)

	ctx, cancel := context.WithCancel(context.Background())
	scanMu.Lock()
	cancelScan = cancel
	scanMu.Unlock()
	defer func() {
		scanMu.Lock()
		cancelScan = nil
		scanMu.Unlock()
		cancel()
	}()

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
					di := DeviceInfo{IP: ip}
					di.IsAlive = Ping(ip, cfg.PingTimeoutMs)
					if di.IsAlive {
						di.Hostname = ReverseDNS(ip)
						di.MAC = GetMAC(ip)
						di.OpenPorts = ScanPorts(ip, cfg.DefaultPorts, cfg.PortTimeoutMs)
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

// StopScan cancela a varredura em andamento.
func StopScan() {
	scanMu.Lock()
	defer scanMu.Unlock()
	if cancelScan != nil {
		cancelScan()
	}
}
