//go:build !windows

package discovery

import (
	"testing"
)

func FuzzParseProcNetArp(f *testing.F) {
	f.Add([]byte("192.168.1.1 0x1 0x2 aa:bb:cc:dd:ee:ff * eth0\n"), "192.168.1.1")
	f.Add([]byte("\n\n\n"), "192.168.1.1")
	f.Add([]byte{}, "10.0.0.1")
	f.Add([]byte("\n\n192.168.1.1 0x1 0x2 aa:bb:cc:dd:ee:ff * eth0\n"), "192.168.1.1")
	f.Add([]byte("incomplete"), "10.0.0.1")
	f.Fuzz(func(t *testing.T, data []byte, ip string) {
		_ = parseProcNetArp(data, ip)
	})
}
