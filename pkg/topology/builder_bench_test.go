package topology

import (
	"fmt"
	"github.com/catnet-io/engine/pkg/results"
	"testing"
)

func BenchmarkBuildGraph(b *testing.B) {
	report := &results.ScanReport{
		Devices: make([]results.DeviceInfo, 1000),
	}
	for i := 0; i < 1000; i++ {
		report.Devices[i] = results.DeviceInfo{
			IP:        fmt.Sprintf("192.168.1.%d", i%254+1), // all in same /24
			Hostname:  "host",
			IsAlive:   true,
			OpenPorts: []int{80, 443},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildGraph(report)
	}
}
