package topology

import (
	"encoding/json"
)

type exportGraph struct {
	Nodes   []TopologyNode `json:"nodes"`
	Links   []TopologyEdge `json:"links"` // D3 uses "links" instead of "edges"
	Gateway string         `json:"gateway,omitempty"`
}

// ExportD3JSON serializes the topology graph to a D3.js compatible JSON structure.
func ExportD3JSON(graph *TopologyGraph) ([]byte, error) {
	if graph == nil {
		return json.Marshal(exportGraph{
			Nodes: []TopologyNode{},
			Links: []TopologyEdge{},
		})
	}

	eg := exportGraph{
		Nodes:   graph.Nodes,
		Links:   graph.Edges,
		Gateway: graph.Gateway,
	}
	if eg.Nodes == nil {
		eg.Nodes = []TopologyNode{}
	}
	if eg.Links == nil {
		eg.Links = []TopologyEdge{}
	}

	return json.Marshal(eg)
}
