package profile

// ScanProfile defines the parameters for a network scan execution.
type ScanProfile struct {
	Name          string   `json:"name"`
	DiscoveryMode string   `json:"discovery_mode"`
	TimeoutMs     int      `json:"timeout_ms"`
	Concurrency   int      `json:"concurrency"`
	ResolveDNS    bool     `json:"resolve_dns"`
	ResolveMAC    bool     `json:"resolve_mac"`
	ExportFormats []string `json:"export_formats"`
	Ports         []int    `json:"ports,omitempty"`
}

// DefaultProfile returns a sensible default scan configuration.
func DefaultProfile() ScanProfile {
	return ScanProfile{
		Name:          "default",
		DiscoveryMode: "icmp+tcp",
		TimeoutMs:     1000,
		Concurrency:   64,
		ResolveDNS:    true,
		ResolveMAC:    true,
		ExportFormats: []string{"json"},
		Ports:         []int{22, 80, 443, 139, 445, 3389},
	}
}
