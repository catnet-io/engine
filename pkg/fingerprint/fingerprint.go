package fingerprint

import (
	"context"
)

// Fingerprint orchestrates the OS, device type, and vendor detection process.
func Fingerprint(ctx context.Context, ip, mac string, ttl int, openPorts []int, timeoutMs int) FingerprintResult {
	// 1. Gather inputs
	ttlResult := GuessOSFromTTL(ttl)

	banners := GrabBanners(ctx, ip, openPorts, timeoutMs)
	bannerResult := OsFromBanners(banners)

	vendor := VendorFromMAC(mac)

	// 2. Combine logic based on priority
	var final FingerprintResult

	if bannerResult.Confidence > 70 {
		final = bannerResult
	} else if ttlResult.Confidence > 50 {
		final = ttlResult
	} else {
		final = bannerResult // Fallback
	}

	// 3. Set Vendor and complement DeviceType
	final.Vendor = vendor
	if final.DeviceType == DeviceUnknown || final.DeviceType == "" {
		if vendor == "Ubiquiti" || vendor == "Cisco" {
			final.DeviceType = DeviceRouter
		} else if vendor == "Raspberry Pi" {
			final.DeviceType = DeviceIoT
		}
	}

	// Make sure fields are somewhat initialized
	if final.OS == "" {
		final.OS = "Unknown"
	}
	if final.OSFamily == "" {
		final.OSFamily = "unknown"
	}
	if final.DeviceType == "" {
		final.DeviceType = DeviceUnknown
	}

	return final
}
