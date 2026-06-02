package targets

// TargetSet holds the parsed boundaries of a scan.
type TargetSet struct {
	RawInputs []string
	CIDRs     []string
	IPs       []string
	Hostnames []string
}
