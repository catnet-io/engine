package netutil

import "testing"

func TestValidateIPv4(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"Valid IPv4", "192.168.1.1", false},
		{"Valid IPv4 Loopback", "127.0.0.1", false},
		{"Invalid IP String", "not_an_ip", true},
		{"Invalid IPv4 Format", "999.999.999.999", true},
		{"Valid IPv6 (Unsupported)", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"Empty String", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPv4(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPv4() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
