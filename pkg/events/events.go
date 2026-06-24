package events

import "github.com/mendsec/catnet-core/pkg/results"

// EventType identifica o tipo de um evento emitido pelo Engine.
type EventType string

const (
	// ScanStarted é emitido uma vez quando o Engine inicia a varredura.
	ScanStarted EventType = "scan_started"
	// HostDiscovered é emitido cada vez que um host é processado.
	HostDiscovered EventType = "host_discovered"
	// ScanProgress é emitido periodicamente com o progresso da varredura.
	ScanProgress EventType = "scan_progress"
	// ScanCompleted é emitido uma vez quando todos os hosts foram processados.
	ScanCompleted EventType = "scan_completed"
)

// Event é a estrutura genérica emitida pelo Engine no channel de eventos.
// O campo Data deve ser convertido para o tipo específico via type assertion,
// usando HostDiscoveredData ou ProgressData conforme o Type.
type Event struct {
	Type EventType
	Data any
}

// HostDiscoveredData é o payload de um evento HostDiscovered.
type HostDiscoveredData struct {
	Host results.HostResult
}

// ProgressData é o payload de um evento ScanProgress.
type ProgressData struct {
	// Ratio é um valor entre 0.0 e 1.0 representando o progresso da varredura.
	Ratio float64
}
