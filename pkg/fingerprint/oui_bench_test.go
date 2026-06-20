package fingerprint

import "testing"

func BenchmarkVendorFromMAC(b *testing.B) {
	mac := "00:16:32:aa:bb:cc"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VendorFromMAC(mac)
	}
}
