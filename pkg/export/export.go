package export

import (
	"github.com/mendsec/catnet-core/pkg/exporter"
	"github.com/mendsec/catnet-core/pkg/results"
)

func ExportJSON(devices []results.HostResult) ([]byte, error) {
	report := results.NewScanReport()
	report.Devices = make([]results.DeviceInfo, len(devices))
	for i, d := range devices {
		report.Devices[i] = results.DeviceInfo(d)
	}
	return exporter.ExportJSON(report)
}

func ExportCSV(devices []results.HostResult) ([]byte, error) {
	report := results.NewScanReport()
	report.Devices = make([]results.DeviceInfo, len(devices))
	for i, d := range devices {
		report.Devices[i] = results.DeviceInfo(d)
	}
	return exporter.ExportCSV(report)
}
