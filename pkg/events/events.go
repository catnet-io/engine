package events

type Type int

const (
	ScanStarted Type = iota
	HostDiscovered
	ScanProgress
	ScanCompleted
)

type Event struct {
	Type Type
	Data interface{}
}

type HostDiscoveredData struct {
	Host string
}

type ProgressData struct {
	Ratio float64
}
