package fingerprint

import "testing"

func TestVendorFromMAC(t *testing.T) {
	tests := []struct {
		name     string
		mac      string
		expected string
	}{
		{
			name:     "Valid standard format",
			mac:      "B8:27:EB:00:00:00",
			expected: "Raspberry Pi",
		},
		{
			name:     "Valid with hyphens",
			mac:      "DC-A6-32-11-22-33",
			expected: "Raspberry Pi",
		},
		{
			name:     "With spaces and lowercase",
			mac:      " e4:5f:01:xx:yy:zz ",
			expected: "Raspberry Pi",
		},
		{
			name:     "Another vendor",
			mac:      "00:50:f2:aa:bb:cc",
			expected: "Microsoft",
		},
		{
			name:     "Unknown vendor",
			mac:      "FF:FF:FF:00:00:00",
			expected: "",
		},
		{
			name:     "Malformed - Too short",
			mac:      "00:00",
			expected: "",
		},
		{
			name:     "Malformed - Random string",
			mac:      "invalid_mac",
			expected: "",
		},
		{
			name:     "Empty",
			mac:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VendorFromMAC(tt.mac)
			if result != tt.expected {
				t.Errorf("VendorFromMAC(%q) = %q; expected %q", tt.mac, result, tt.expected)
			}
		})
	}
}
