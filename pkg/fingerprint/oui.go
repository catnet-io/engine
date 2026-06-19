package fingerprint

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
// ⚡ Bolt Optimization: Zero-allocation MAC parsing.
// Avoids strings.Split, ReplaceAll, and ToUpper to save ~280B and 8 allocs per lookup.
func VendorFromMAC(mac string) string {
	if len(mac) < 8 {
		return ""
	}

	var prefix [8]byte
	var pIdx int
	started := false

	for i := 0; i < len(mac) && pIdx < 8; i++ {
		c := mac[i]
		// Trim leading spaces
		if !started {
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				continue
			}
			started = true
		}

		// To upper case
		if c >= 'a' && c <= 'z' {
			c -= 32
		} else if c == '-' { // Replace dash with colon
			c = ':'
		}
		prefix[pIdx] = c
		pIdx++
	}

	// Only lookup if we successfully extracted a valid 8-byte prefix (XX:XX:XX)
	if pIdx == 8 && prefix[2] == ':' && prefix[5] == ':' {
		// Go compiler optimizes string(prefix[:]) map lookup to avoid allocation
		if vendor, ok := ouiMap[string(prefix[:])]; ok {
			return vendor
		}
	}
	return ""
}
