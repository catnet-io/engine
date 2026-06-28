package engine

import "github.com/catnet-io/engine/pkg/results"

// ScanEventType define o tipo do evento disparado durante a varredura.
type ScanEventType int

const (
	EventLifecycleStart ScanEventType = iota
	EventLifecycleComplete
	EventLifecycleCancel
	EventWarning
	EventResult
	EventProgress
)

// ScanEvent Ã© a estrutura que engloba dados de evento, amigÃ¡vel para callbacks sÃ­ncronos e FFI.
type ScanEvent struct {
	Type     ScanEventType
	Device   *results.DeviceInfo
	Progress float64
	Message  string
}

// EventCallback Ã© a assinatura da funÃ§Ã£o de callback sÃ­ncrono.
type EventCallback func(event ScanEvent)
