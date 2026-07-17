package engine

import "context"

// ScanConfig defines the parameters of a scan.
type ScanConfig struct {
	// DefaultPorts is the list of TCP ports scanned on each alive host.
	// Default: [22 (SSH), 80 (HTTP), 443 (HTTPS), 139 (NetBIOS), 445 (SMB), 3389 (RDP)].
	DefaultPorts []int `json:"defaultPorts"`

	// PortTimeoutMs is the timeout in milliseconds for each TCP connection attempt.
	// Default: 500ms. Range: 1–10000ms (Sanitize clamps values out of range).
	PortTimeoutMs int `json:"portTimeoutMs"`

	// PingTimeoutMs is the timeout in milliseconds for the liveness ICMP ping.
	// Default: 1000ms. Range: 1–10000ms.
	PingTimeoutMs int `json:"pingTimeoutMs"`

	// MaxThreads defines the scan concurrency level (simultaneous goroutines).
	// Default: 64. Absolute maximum limit: 256 (enforced by engine in StartScan).
	// Values above 256 or below 1 are silently clamped to 16.
	MaxThreads int `json:"maxThreads"`

	// FingerprintProvider allows injecting custom fingerprinting logic.
	// If nil, the engine uses pkg/fingerprint with default TTL, banner, and OUI heuristics.
	// Useful for testing (mocking) or extending detection capabilities.
	FingerprintProvider FingerprintProvider `json:"-"`
}

// FingerprintData contains the detection results.
type FingerprintData struct {
	OS         string
	OSFamily   string
	DeviceType string
	Vendor     string
}

// FingerprintProvider defines the contract for OS and device detection heuristics.
type FingerprintProvider interface {
	Fingerprint(ctx context.Context, ip, mac string, ttl int, ports []int, timeoutMs int) FingerprintData
}

// DefaultConfig returns a ScanConfig with conservative default values.
func DefaultConfig() ScanConfig {
	return ScanConfig{
		DefaultPorts:        []int{22, 80, 443, 139, 445, 3389},
		PortTimeoutMs:       500,
		PingTimeoutMs:       1000,
		MaxThreads:          64,
		FingerprintProvider: nil, // Will use default in StartScan if nil
	}
}

// Sanitize corrects values outside safe limits.
// The engine itself executes this sanitization defensively during StartScan,
// but it can be invoked manually to reflect limits in client interfaces.
func (c *ScanConfig) Sanitize() {
	if c.MaxThreads <= 0 || c.MaxThreads > 256 {
		c.MaxThreads = 16
	}
	if c.PortTimeoutMs <= 0 || c.PortTimeoutMs > 10000 {
		c.PortTimeoutMs = 500
	}
	if c.PingTimeoutMs <= 0 || c.PingTimeoutMs > 10000 {
		c.PingTimeoutMs = 1000
	}
}

