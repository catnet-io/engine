package topology

import (
	"fmt"
	"testing"

	"github.com/catnet-io/engine/pkg/results"
)

func TestBuildGraph_EdgeLimit(t *testing.T) {
	devices := make([]results.DeviceInfo, 300)
	for i := range devices {
		devices[i] = results.DeviceInfo{
			IP:        fmt.Sprintf("10.0.0.%d", i%254+1),
			IsAlive:   true,
			OpenPorts: []int{80},
		}
	}
	report := &results.ScanReport{Devices: devices}
	graph := BuildGraph(report)

	hostEdges := 0
	for _, e := range graph.Edges {
		if e.Weight == 0.3 {
			hostEdges++
		}
	}
	if hostEdges > 200 {
		t.Errorf("Expected <= 200 host-to-host edges, got %d", hostEdges)
	}
}
