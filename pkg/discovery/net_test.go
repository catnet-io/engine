package discovery

import (
	"context"
	"testing"
)

// These tests verify input validation paths only. Network-dependent paths are covered by integration tests.

func TestPingValidation(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{"empty string", ""},
		{"invalid ip", "999.999.999.999"},
		{"not an ip", "not-an-ip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Ping(context.Background(), tt.ip, 1000) != false {
				t.Errorf("Ping(%q) expected false, got true", tt.ip)
			}
		})
	}
}

func TestReverseDNSValidation(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{"empty string", ""},
		{"out of bounds ip", "256.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := ReverseDNS(context.Background(), tt.ip); res != "" {
				t.Errorf("ReverseDNS(%q) expected empty string, got %q", tt.ip, res)
			}
		})
	}
}

func TestGetMACValidation(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{"empty string", ""},
		{"ipv6 not supported", "::1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := GetMAC(context.Background(), tt.ip); res != "" {
				t.Errorf("GetMAC(%q) expected empty string, got %q", tt.ip, res)
			}
		})
	}
}
