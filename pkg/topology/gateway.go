package topology

import (
	"strings"

	"github.com/catnet-io/engine/pkg/results"
)

// DetectGateway heuristically finds a likely gateway in the scanned devices.
func DetectGateway(devices []results.DeviceInfo) string {
	// Heuristic 1: IP ending in .1 or .254
	for _, dev := range devices {
		if !dev.IsAlive {
			continue
		}
		if strings.HasSuffix(dev.IP, ".1") || strings.HasSuffix(dev.IP, ".254") {
			return dev.IP
		}
	}

	// Heuristic 2: Vendor is router-like
	for _, dev := range devices {
		if !dev.IsAlive {
			continue
		}
		vendor := strings.ToLower(dev.Vendor)
		if strings.Contains(vendor, "ubiquiti") || strings.Contains(vendor, "cisco") || strings.Contains(vendor, "mikrotik") {
			return dev.IP
		}
	}

	// Heuristic 3: Host with ports 443 AND 80 AND 22 open
	for _, dev := range devices {
		if !dev.IsAlive {
			continue
		}
		has443 := false
		has80 := false
		has22 := false
		for _, port := range dev.OpenPorts {
			if port == 443 {
				has443 = true
			} else if port == 80 {
				has80 = true
			} else if port == 22 {
				has22 = true
			}
		}
		if has443 && has80 && has22 {
			return dev.IP
		}
	}

	return ""
}
