package results

import "time"

// ScanReport encapsula o resultado completo de uma varredura.
type ScanReport struct {
	StartTime time.Time    `json:"startTime"`
	EndTime   time.Time    `json:"endTime"`
	Total     int          `json:"total"`
	Alive     int          `json:"alive"`
	Devices   []DeviceInfo `json:"devices"`
}

// NewScanReport cria um novo relatório de varredura.
func NewScanReport() *ScanReport {
	return &ScanReport{
		StartTime: time.Now(),
		Devices:   []DeviceInfo{},
	}
}
