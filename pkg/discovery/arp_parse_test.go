//go:build !windows

package discovery

import (
	"testing"
)

func TestParseProcNetArp(t *testing.T) {
	procNetArpData := []byte(`IP address       HW type     Flags       HW address            Mask     Device
192.168.1.1      0x1         0x2         aa:bb:cc:dd:ee:ff     *        eth0
192.168.1.2      0x1         0x2         00:00:00:00:00:00     *        eth0
10.0.0.5         0x1         0x2         11:22:33:44:55:66     *        eth1
192.168.1.100    0x1         0x2         00:11:22:33:44:55     *        eth0
`)

	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "Existing complete MAC",
			ip:       "192.168.1.1",
			expected: "AA-BB-CC-DD-EE-FF",
		},
		{
			name:     "Existing incomplete MAC (zeroes)",
			ip:       "192.168.1.2",
			expected: "",
		},
		{
			name:     "Existing MAC in middle",
			ip:       "10.0.0.5",
			expected: "11-22-33-44-55-66",
		},
		{
			name:     "Non-existent IP",
			ip:       "192.168.1.99",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseProcNetArp(procNetArpData, tt.ip)
			if result != tt.expected {
				t.Errorf("parseProcNetArp(..., %q) = %q, want %q", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestParseProcNetArp_EdgeCases(t *testing.T) {
	dataNoHeader := []byte(`192.168.1.1 0x1 0x2 aa:bb:cc:dd:ee:ff * eth0
192.168.1.2 0x1 0x2 11:22:33:44:55:66 * eth0`)

	// First line without newline
	res1 := parseProcNetArp(dataNoHeader, "192.168.1.1")
	if res1 != "AA-BB-CC-DD-EE-FF" {
		t.Errorf("Expected AA-BB-CC-DD-EE-FF, got %q", res1)
	}

	// Malformed / incomplete lines
	dataCorrupted := []byte(`IP address       HW type     Flags       HW address            Mask     Device
192.168.1.1 0x1 0x2
192.168.1.2 0x1 0x2 aa:bb:cc:dd:ee:ff * eth0
`)
	res2 := parseProcNetArp(dataCorrupted, "192.168.1.1")
	if res2 != "" {
		t.Errorf("Expected empty string for malformed line, got %q", res2)
	}
	res3 := parseProcNetArp(dataCorrupted, "192.168.1.2")
	if res3 != "AA-BB-CC-DD-EE-FF" {
		t.Errorf("Expected AA-BB-CC-DD-EE-FF, got %q", res3)
	}
}
