package fingerprint

import (
	"testing"
)

func FuzzVendorFromMAC(f *testing.F) {
	f.Add("B8:27:EB:00:00:00")
	f.Add("DC-A6-32-11-22-33")
	f.Add("")
	f.Add("invalid_mac_string_with_unicode_\x00\xff")
	f.Add("00:00:00:00:00:00:00:00:00")
	f.Fuzz(func(t *testing.T, mac string) {
		_ = VendorFromMAC(mac)
	})
}

func FuzzSanitizeBanner(f *testing.F) {
	f.Add("SSH-2.0-OpenSSH_8.9p1 Ubuntu")
	f.Add("\x00\x01\x02\xff\xfe")
	f.Add("\n\r\t")
	f.Add(string(make([]byte, 4096)))
	f.Fuzz(func(t *testing.T, input string) {
		_ = sanitizeBanner(input)
	})
}
