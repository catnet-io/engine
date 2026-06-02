package results

// HostResult represents the outcome of scanning a single host.
type HostResult struct {
	IP          string `json:"ip"`
	Hostname    string `json:"hostname,omitempty"`
	MAC         string `json:"mac,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	RTTMs       int    `json:"rtt_ms,omitempty"`
	Alive       bool   `json:"alive"`
	Source      string `json:"source,omitempty"`
	LastSeenUTC string `json:"last_seen_utc,omitempty"`
	OpenPorts   []int  `json:"open_ports,omitempty"`
}

// ScanStats holds aggregated statistics of a scan run.
type ScanStats struct {
	TotalTargets int   `json:"total_targets"`
	Responded    int   `json:"responded"`
	DurationMs   int64 `json:"duration_ms"`
	Errors       int   `json:"errors"`
}

// ScanResult contains the full payload of a completed scan.
type ScanResult struct {
	ScanID        string       `json:"scan_id"`
	StartedAtUTC  string       `json:"started_at_utc"`
	FinishedAtUTC string       `json:"finished_at_utc"`
	ProfileName   string       `json:"profile_name"`
	Targets       []string     `json:"targets"`
	Hosts         []HostResult `json:"hosts"`
	Stats         ScanStats    `json:"stats"`
}
