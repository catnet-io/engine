package results

import "time"

// ScanReport encapsula o resultado completo de uma varredura.
type ScanReport struct {
	SchemaVersion string       `json:"schemaVersion"`
	StartTime     time.Time    `json:"startTime"`
	EndTime       time.Time    `json:"endTime"`
	Total     int          `json:"total"`
	Alive     int          `json:"alive"`
	Devices   []DeviceInfo `json:"devices"`
}

// NewScanReport cria um novo relatório de varredura.
func NewScanReport() *ScanReport {
	return &ScanReport{
		SchemaVersion: "1.0.0",
		StartTime:     time.Now(),
		Devices:       []DeviceInfo{},
	}
}
