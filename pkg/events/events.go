package events

import "github.com/mendsec/catnet-core/pkg/results"

// EventType categorizes the lifecycle event of a scan.
type EventType string

const (
	ScanStarted    EventType = "ScanStarted"
	TargetQueued   EventType = "TargetQueued"
	HostDiscovered EventType = "HostDiscovered"
	TargetFailed   EventType = "TargetFailed"
	ScanProgress   EventType = "ScanProgress"
	ScanCompleted  EventType = "ScanCompleted"
)

// Event is the generic envelope emitted by the scanner engine.
type Event struct {
	Type EventType
	Data interface{}
}

// ProgressData payload for ScanProgress events.
type ProgressData struct {
	Processed int
	Total     int
	Ratio     float64
}

// HostDiscoveredData payload for HostDiscovered events.
type HostDiscoveredData struct {
	Host results.HostResult
}
