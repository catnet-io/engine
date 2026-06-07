package targets

import (
	"testing"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr bool
	}{
		{"Single IP valid", "10.0.0.1", 1, false},
		{"CIDR valid", "192.168.1.0/30", 2, false}, // 192.168.1.1 and 192.168.1.2
		{"CIDR too large", "0.0.0.0/0", 0, true},
		{"CIDR max limit", "10.0.0.0/16", 65534, false},
		{"Dash range full valid", "192.168.1.10-192.168.1.12", 3, false},
		{"Dash range shorthand valid", "192.168.1.10-12", 3, false},
		{"Invalid format", "invalid", 0, true},
		{"Empty input", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("ParseRange() got %v items, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func BenchmarkParseRange(b *testing.B) {
	input := "10.0.0.0/16" // 65534 IPs
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseRange(input)
	}
}

func FuzzParseRange(f *testing.F) {
	f.Add("10.0.0.1")
	f.Add("192.168.1.0/24")
	f.Add("192.168.1.1-192.168.1.10")
	f.Add("10.0.0.1-255")
	f.Add("invalid")
	f.Add("10.0.0.0/0")
	f.Add("10.0.0.0/33")

	f.Fuzz(func(t *testing.T, input string) {
		// The goal of fuzzing here is to ensure ParseRange never panics on malformed input.
		ParseRange(input)
	})
}
