package scanner

// DeviceInfo representa o resultado da varredura de um único host.
type DeviceInfo struct {
	IP             string `json:"ip"`
	IsAlive        bool   `json:"isAlive"`
	Hostname       string `json:"hostname"`
	MAC            string `json:"mac"`
	OpenPortsCount int    `json:"openPortsCount"`
	OpenPorts      []int  `json:"openPorts"`
}
