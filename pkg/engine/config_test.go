package engine

import "testing"

func TestConfigSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    ScanConfig
		expected ScanConfig
	}{
		{
			name: "Valid defaults should remain unchanged",
			input: ScanConfig{
				MaxThreads:    64,
				PortTimeoutMs: 500,
				PingTimeoutMs: 1000,
			},
			expected: ScanConfig{
				MaxThreads:    64,
				PortTimeoutMs: 500,
				PingTimeoutMs: 1000,
			},
		},
		{
			name: "Pathological limits should be clamped",
			input: ScanConfig{
				MaxThreads:    1000000,
				PortTimeoutMs: -10,
				PingTimeoutMs: 0,
			},
			expected: ScanConfig{
				MaxThreads:    16,
				PortTimeoutMs: 500,
				PingTimeoutMs: 1000,
			},
		},
		{
			name: "Negative limits should be clamped",
			input: ScanConfig{
				MaxThreads:    -1,
				PortTimeoutMs: 999999,
				PingTimeoutMs: -100,
			},
			expected: ScanConfig{
				MaxThreads:    16,
				PortTimeoutMs: 500,
				PingTimeoutMs: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.input
			cfg.Sanitize()

			if cfg.MaxThreads != tt.expected.MaxThreads {
				t.Errorf("Expected MaxThreads %d, got %d", tt.expected.MaxThreads, cfg.MaxThreads)
			}
			if cfg.PortTimeoutMs != tt.expected.PortTimeoutMs {
				t.Errorf("Expected PortTimeoutMs %d, got %d", tt.expected.PortTimeoutMs, cfg.PortTimeoutMs)
			}
			if cfg.PingTimeoutMs != tt.expected.PingTimeoutMs {
				t.Errorf("Expected PingTimeoutMs %d, got %d", tt.expected.PingTimeoutMs, cfg.PingTimeoutMs)
			}
		})
	}
}
