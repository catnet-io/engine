package ports

import (
	"context"
	"net"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestScanPorts(t *testing.T) {
	// Start a local test server to listen on a random port
	ts := httptest.NewServer(nil)
	defer ts.Close()

	// Extract the port
	_, portStr, err := net.SplitHostPort(ts.Listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to split host/port: %v", err)
	}
	openPort, _ := strconv.Atoi(portStr)

	tests := []struct {
		name      string
		ip        string
		ports     []int
		timeoutMs int
		wantCount int
	}{
		{"Valid IP and Open Port", "127.0.0.1", []int{openPort}, 500, 1},
		{"Valid IP but Closed Port", "127.0.0.1", []int{openPort + 1}, 50, 0},
		{"Invalid IP", "999.999.999.999", []int{80}, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openChan := ScanPorts(context.Background(), tt.ip, tt.ports, tt.timeoutMs)
			count := 0
			for range openChan {
				count++
			}
			if count != tt.wantCount {
				t.Errorf("ScanPorts() returned %d open ports, want %d", count, tt.wantCount)
			}
		})
	}
}
