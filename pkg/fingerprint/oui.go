package fingerprint

import (
	"strings"
)

var ouiMap = map[string]string{
	"B8:27:EB": "Raspberry Pi",
	"DC:A6:32": "Raspberry Pi",
	"E4:5F:01": "Raspberry Pi",
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
// ⚡ Bolt Optimization: Zero-allocation MAC prefix extraction using fixed-size byte array.
// Avoids strings.ToUpper, strings.ReplaceAll, strings.Split, and strings.Join memory allocations.
func VendorFromMAC(mac string) string {
	mac = strings.TrimSpace(mac)
	if len(mac) < 8 {
		return ""
	}

	var prefix [8]byte
	j := 0
	for i := 0; i < len(mac) && j < 8; i++ {
		c := mac[i]
		if c == '-' {
			c = ':'
		}
		if c >= 'a' && c <= 'z' {
			c -= 32 // ToUpper
		}
		prefix[j] = c
		j++
	}

	if j == 8 {
		// Go compiler optimizes string(byteSlice) as map key to avoid allocation
		if vendor, ok := ouiMap[string(prefix[:])]; ok {
			return vendor
		}
	}
	return ""
}
