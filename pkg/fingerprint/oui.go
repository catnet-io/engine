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
	mac = strings.ToUpper(strings.TrimSpace(mac))
	mac = strings.ReplaceAll(mac, "-", ":")

	colonsFound := 0
	for i := 0; i < len(mac); i++ {
		if mac[i] == ':' {
			colonsFound++
			if colonsFound == 3 {
				prefix := mac[:i]
				if vendor, ok := ouiMap[prefix]; ok {
					return vendor
				}
				return ""
			}
		}
	}
	if colonsFound == 2 {
		if vendor, ok := ouiMap[mac]; ok {
			return vendor
		}
	}
	return ""
}
