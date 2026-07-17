package engine

import "github.com/catnet-io/engine/pkg/results"

// ScanEventType defines the type of event fired during a scan.
type ScanEventType int

const (
	EventLifecycleStart ScanEventType = iota
	EventLifecycleComplete
	EventLifecycleCancel
	EventWarning
	EventResult
	EventProgress
)

// ScanEvent holds the event data, suitable for synchronous callbacks and FFI.
type ScanEvent struct {
	Type     ScanEventType
	Device   *results.DeviceInfo
	Progress float64
	Message  string
}

// EventCallback is the signature of the synchronous callback function.
type EventCallback func(event ScanEvent)
