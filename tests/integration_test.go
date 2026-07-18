package tests

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/catnet-io/engine/pkg/engine"
	"github.com/catnet-io/engine/pkg/exporter"
	"github.com/catnet-io/engine/pkg/results"
	"github.com/catnet-io/engine/pkg/targets"
)

func TestEndToEndScanAndExport(t *testing.T) {
	// Read targets fixture
	targetsPath := filepath.Join("..", "testdata", "targets_fixture.txt")
	content, err := os.ReadFile(targetsPath)
	if err != nil {
		t.Fatalf("Failed to read targets fixture: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var allIPs []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// We only parse the subset we want to test to keep the test fast
		// e.g. the single IP and the dash range
		if strings.Contains(line, "/") {
			// Skip /16 or /30 from fixture to keep integration test extremely fast
			continue
		}

		parsed, err := targets.ParseRange(line)
		if err != nil {
			t.Fatalf("Failed to parse target %s: %v", line, err)
		}
		allIPs = append(allIPs, parsed...)
	}

	if len(allIPs) == 0 {
		t.Fatalf("No IPs parsed from fixture")
	}

	// Start scan
	cfg := engine.DefaultConfig()
	cfg.PingTimeoutMs = 100 // Short timeout for integration test
	cfg.DefaultPorts = []int{80, 443}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mu sync.Mutex
	var startEventReceived, completeEventReceived bool
	var resultsReceived int

	report, err := engine.StartScan(ctx, allIPs, cfg, func(event engine.ScanEvent) { //nolint:staticcheck // integration test tests the deprecated engine.StartScan API
		mu.Lock()
		defer mu.Unlock()
		switch event.Type {
		case engine.EventLifecycleStart:
			startEventReceived = true
		case engine.EventLifecycleComplete:
			completeEventReceived = true
		case engine.EventResult:
			resultsReceived++
		}
	})

	if err != nil {
		t.Fatalf("StartScan failed: %v", err)
	}

	// Validate report
	if report.SchemaVersion != "2.0.0" {
		t.Errorf("Expected SchemaVersion 2.0.0, got %s", report.SchemaVersion)
	}
	if report.Total != len(allIPs) {
		t.Errorf("Expected Total %d, got %d", len(allIPs), report.Total)
	}

	// Validate events
	if !startEventReceived {
		t.Errorf("Did not receive EventLifecycleStart")
	}
	if !completeEventReceived {
		t.Errorf("Did not receive EventLifecycleComplete")
	}
	if resultsReceived != len(allIPs) {
		t.Errorf("Expected %d result events, got %d", len(allIPs), resultsReceived)
	}

	// Test Export
	jsonBytes, err := exporter.ExportJSON(report)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	// Verify exported structure unmarshals correctly
	var exportedReport results.ScanReport
	if err := json.Unmarshal(jsonBytes, &exportedReport); err != nil {
		t.Fatalf("Failed to unmarshal exported JSON: %v", err)
	}

	if exportedReport.Total != report.Total {
		t.Errorf("Exported JSON Total mismatch: got %d, want %d", exportedReport.Total, report.Total)
	}
}
