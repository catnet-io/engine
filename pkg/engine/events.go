package engine

import "github.com/mendsec/catnet-core/pkg/results"

// ScanEventType define o tipo do evento disparado durante a varredura.
type ScanEventType int

const (
	EventResult ScanEventType = iota
	EventProgress
)

// ScanEvent é a estrutura que engloba dados de evento, amigável para callbacks síncronos e FFI.
type ScanEvent struct {
	Type     ScanEventType
	Device   *results.DeviceInfo
	Progress float64
}

// EventCallback é a assinatura da função de callback síncrono.
type EventCallback func(event ScanEvent)
