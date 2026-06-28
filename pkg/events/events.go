package events

import "github.com/catnet-io/engine/pkg/results"

type EventType string

const (
	ScanStarted    EventType = "scan_started"
	HostDiscovered EventType = "host_discovered"
	ScanProgress   EventType = "scan_progress"
	ScanCompleted  EventType = "scan_completed"
)

type Event struct {
	Type EventType
	Data any
}

type HostDiscoveredData struct {
	Host results.HostResult
}

type ProgressData struct {
	Ratio float64
}
