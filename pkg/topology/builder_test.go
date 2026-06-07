package topology

import (
	"testing"

	"github.com/mendsec/catnet-core/pkg/results"
)

func TestDetectGateway(t *testing.T) {
	devices := []results.DeviceInfo{
		{IP: "192.168.1.10", IsAlive: true},
		{IP: "192.168.1.254", IsAlive: true},
	}

	gw := DetectGateway(devices)
	if gw != "192.168.1.254" {
		t.Errorf("Expected 192.168.1.254, got %s", gw)
	}
}

func TestBuildGraph_EmptyReport(t *testing.T) {
	graph := BuildGraph(&results.ScanReport{})
	if len(graph.Nodes) != 0 || len(graph.Edges) != 0 {
		t.Errorf("Expected empty graph")
	}
}

func TestBuildGraph_SingleHost(t *testing.T) {
	report := &results.ScanReport{
		Devices: []results.DeviceInfo{
			{IP: "10.0.0.5", IsAlive: true},
		},
	}
	graph := BuildGraph(report)
	if len(graph.Nodes) != 1 {
		t.Errorf("Expected 1 node")
	}
	if len(graph.Edges) != 0 {
		t.Errorf("Expected 0 edges")
	}
}

func TestBuildGraph_WithGateway(t *testing.T) {
	report := &results.ScanReport{
		Devices: []results.DeviceInfo{
			{IP: "192.168.1.1", IsAlive: true}, // gateway
			{IP: "192.168.1.10", IsAlive: true, OpenPorts: []int{80}},
			{IP: "192.168.1.11", IsAlive: true, OpenPorts: []int{80}},
			{IP: "192.168.1.12", IsAlive: true},
		},
	}

	graph := BuildGraph(report)
	if graph.Gateway != "192.168.1.1" {
		t.Errorf("Expected gateway 192.168.1.1, got %s", graph.Gateway)
	}
	if len(graph.Nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(graph.Nodes))
	}

	if len(graph.Edges) != 4 {
		t.Errorf("Expected 4 edges, got %d", len(graph.Edges))
	}
}
