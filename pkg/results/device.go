package results

// DeviceInfo representa o resultado da varredura de um único host.
type DeviceInfo struct {
	IP        string `json:"ip"`
	IsAlive   bool   `json:"isAlive"`
	Hostname  string `json:"hostname"`
	MAC       string `json:"mac"`
	OpenPorts []int  `json:"openPorts"`
}

// PortCount retorna a quantidade de portas abertas.
func (d DeviceInfo) PortCount() int {
	return len(d.OpenPorts)
}
