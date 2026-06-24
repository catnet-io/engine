package results

// HostResult is a compatibility alias for DeviceInfo.
type HostResult = DeviceInfo

// DeviceInfo representa o resultado da varredura de um único host.
type DeviceInfo struct {
	IP         string `json:"ip"`
	IsAlive    bool   `json:"isAlive"`
	Hostname   string `json:"hostname"`
	MAC        string `json:"mac"`
	OS         string `json:"os,omitempty"`
	OSFamily   string `json:"osFamily,omitempty"`
	DeviceType string `json:"deviceType,omitempty"`
	Vendor     string `json:"vendor,omitempty"`
	OpenPorts  []int  `json:"openPorts"`
}

// PortCount retorna a quantidade de portas abertas.
func (d DeviceInfo) PortCount() int {
	return len(d.OpenPorts)
}
