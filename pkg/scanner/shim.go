// Package scanner is deprecated: Use pkg/engine, pkg/discovery, pkg/ports, pkg/targets, and pkg/results instead.
// This package is maintained temporarily for backward compatibility.
package scanner

import (
	"context"
	"sync"

	"github.com/mendsec/catnet-core/pkg/discovery"
	"github.com/mendsec/catnet-core/pkg/engine"
	"github.com/mendsec/catnet-core/pkg/ports"
	"github.com/mendsec/catnet-core/pkg/results"
	"github.com/mendsec/catnet-core/pkg/targets"
)

// DeviceInfo is deprecated: Use results.DeviceInfo instead.
type DeviceInfo = results.DeviceInfo

// ScanConfig is deprecated: Use engine.ScanConfig instead.
type ScanConfig = engine.ScanConfig

// DefaultConfig is deprecated: Use engine.DefaultConfig instead.
func DefaultConfig() ScanConfig {
	return engine.DefaultConfig()
}

// ParseRange is deprecated: Use targets.ParseRange instead.
func ParseRange(input string) ([]string, error) {
	return targets.ParseRange(input)
}

// Ping is deprecated: Use discovery.Ping instead.
func Ping(ip string, timeoutMs int) bool {
	return discovery.Ping(context.Background(), ip, timeoutMs)
}

// ReverseDNS is deprecated: Use discovery.ReverseDNS instead.
func ReverseDNS(ip string) string {
	return discovery.ReverseDNS(ip)
}

// GetMAC is deprecated: Use discovery.GetMAC instead.
func GetMAC(ip string) string {
	return discovery.GetMAC(ip)
}

// ScanPorts is deprecated: Use ports.ScanPorts instead.
func ScanPorts(ip string, portsList []int, timeoutMs int) []int {
	return ports.ScanPorts(context.Background(), ip, portsList, timeoutMs)
}

var (
	shimMu     sync.Mutex
	shimCancel context.CancelFunc
)

// StartScan is deprecated: Use engine.StartScan with a context.Context instead.
func StartScan(ips []string, cfg ScanConfig, onResult func(DeviceInfo), onProgress func(float64)) error {
	ctx, cancel := context.WithCancel(context.Background())
	shimMu.Lock()
	shimCancel = cancel
	shimMu.Unlock()

	defer func() {
		shimMu.Lock()
		shimCancel = nil
		shimMu.Unlock()
		cancel()
	}()

	_, err := engine.StartScan(ctx, ips, cfg, func(event engine.ScanEvent) {
		if event.Type == engine.EventResult && event.Device != nil && onResult != nil {
			onResult(*event.Device)
		} else if event.Type == engine.EventProgress && onProgress != nil {
			onProgress(event.Progress)
		}
	})
	return err
}

// StopScan is deprecated: Caller should pass a cancelable context to engine.StartScan instead.
func StopScan() {
	shimMu.Lock()
	defer shimMu.Unlock()
	if shimCancel != nil {
		shimCancel()
	}
}
