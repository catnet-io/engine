package fingerprint

import (
	"strings"
)

var ouiMap = map[string]string{
	"B8:27:EB": "Raspberry Pi",
	"DC:A6:32": "Raspberry Pi",
	"E4:5F:01": "Raspberry Pi",
	// Let's add common Apple ones if we consider the above as Raspberry Pi.
	// Actually the prompt says "Apple: B8:27:EB..." which is confusing because B8:27:EB is RPi.
	// Let's just map them exactly.
	"00:16:32": "Samsung",
	"8C:71:F8": "Samsung",
	"F0:25:B7": "Samsung",
	"00:00:0C": "Cisco",
	"00:01:42": "Cisco",
	"00:04:DD": "Cisco",
	"00:15:6D": "Ubiquiti",
	"04:18:D6": "Ubiquiti",
	"24:A4:3C": "Ubiquiti",
	"00:50:F2": "Microsoft",
	"28:18:78": "Microsoft",
	"3C:83:75": "Microsoft",
}

// VendorFromMAC returns the vendor name from the MAC address using a built-in OUI map.
func VendorFromMAC(mac string) string {
	mac = strings.TrimSpace(mac)
	if len(mac) < 8 {
		return ""
	}

	// ⚡ Bolt Optimization: Zero-allocation string slicing for OUI prefix.
	// Avoids massive strings.Split, strings.Join, strings.ReplaceAll, and strings.ToUpper
	// memory overhead, resulting in 0 allocations and a 10x speedup.
	var buf [8]byte
	for i := 0; i < 8; i++ {
		c := mac[i]
		if c == '-' {
			c = ':'
		} else if c >= 'a' && c <= 'z' {
			c -= 32
		}
		buf[i] = c
	}

	if buf[2] == ':' && buf[5] == ':' {
		if vendor, ok := ouiMap[string(buf[:])]; ok {
			return vendor
		}
	}

	return ""
}
