package results

// HostResult é o tipo canônico de resultado de varredura para a API event-driven.
// Usado pelos pacotes scan, events e export como formato de intercâmbio.
type HostResult struct {
	IP        string `json:"ip"`
	Alive     bool   `json:"alive"`
	Hostname  string `json:"hostname"`
	MAC       string `json:"mac"`
	OpenPorts []int  `json:"open_ports"`
}

// ToDeviceInfo converte HostResult para DeviceInfo, preservando campos comuns.
func (h HostResult) ToDeviceInfo() DeviceInfo {
	return DeviceInfo{
		IP:        h.IP,
		IsAlive:   h.Alive,
		Hostname:  h.Hostname,
		MAC:       h.MAC,
		OpenPorts: h.OpenPorts,
	}
}

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
