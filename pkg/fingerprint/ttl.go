package fingerprint

// GuessOSFromTTL attempts to detect the OS and device type based on the IP TTL value.
func GuessOSFromTTL(ttl int) FingerprintResult {
	res := FingerprintResult{
		OS:         "Unknown",
		OSFamily:   "unknown",
		DeviceType: DeviceUnknown,
		Confidence: 0,
	}

	if ttl <= 0 {
		return res
	}

	if ttl < 30 {
		res.DeviceType = DeviceRouter
		res.Confidence = 50
	} else if ttl >= 60 && ttl <= 64 {
		res.OSFamily = "linux"
		res.OS = "Linux/Unix"
		res.Confidence = 80
	} else if ttl >= 120 && ttl <= 128 {
		res.OSFamily = "windows"
		res.OS = "Windows"
		res.Confidence = 80
	} else if ttl >= 250 && ttl <= 255 {
		res.OSFamily = "unix"
		res.OS = "Cisco/macOS"
		res.Confidence = 60
	}

	return res
}
