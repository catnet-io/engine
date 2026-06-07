package topology

type NodeRole string

const (
	RoleGateway NodeRole = "gateway"
	RoleHost    NodeRole = "host"
	RoleServer  NodeRole = "server"
	RoleUnknown NodeRole = "unknown"
)

type TopologyNode struct {
	ID         string   `json:"id"`    // IP address
	Label      string   `json:"label"` // hostname or IP
	Role       NodeRole `json:"role"`
	DeviceType string   `json:"deviceType"`
	IsAlive    bool     `json:"isAlive"`
	OpenPorts  []int    `json:"openPorts"`
}

type TopologyEdge struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"` // 1.0 = same subnet, 0.5 = cross-subnet
}

type TopologyGraph struct {
	Nodes   []TopologyNode `json:"nodes"`
	Edges   []TopologyEdge `json:"edges"`
	Gateway string         `json:"gateway,omitempty"`
}
