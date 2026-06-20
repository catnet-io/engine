package fingerprint

import (
	"context"
	"testing"
)

func TestGuessOSFromTTL(t *testing.T) {
	tests := []struct {
		ttl        int
		osFamily   string
		deviceType DeviceType
	}{
		{64, "linux", ""}, // usually empty device type for basic linux unless set
		{128, "windows", ""},
		{255, "unix", ""},
		{1, "unknown", DeviceRouter},
	}

	for _, tt := range tests {
		res := GuessOSFromTTL(tt.ttl)
		if res.OSFamily != tt.osFamily {
			t.Errorf("ttl %d: expected OSFamily %s, got %s", tt.ttl, tt.osFamily, res.OSFamily)
		}
		if tt.deviceType != "" && res.DeviceType != tt.deviceType {
			t.Errorf("ttl %d: expected DeviceType %s, got %s", tt.ttl, tt.deviceType, res.DeviceType)
		}
	}
}

func TestOsFromBanners(t *testing.T) {
	tests := []struct {
		banners    map[int]string
		os         string
		osFamily   string
		deviceType DeviceType
	}{
		{
			banners:  map[int]string{22: "SSH-2.0-OpenSSH_8.9p1 Ubuntu"},
			os:       "Ubuntu",
			osFamily: "linux",
		},
		{
			banners:  map[int]string{80: "Server: IIS"},
			os:       "Windows",
			osFamily: "windows",
		},
		{
			banners:  map[int]string{22: "OpenSSH_for_Windows"},
			os:       "Windows",
			osFamily: "windows",
		},
	}

	for _, tt := range tests {
		res := OsFromBanners(tt.banners)
		if res.OSFamily != tt.osFamily {
			t.Errorf("expected %s, got %s", tt.osFamily, res.OSFamily)
		}
		if res.OS != tt.os {
			t.Errorf("expected %s, got %s", tt.os, res.OS)
		}
	}
}

func TestFingerprint(t *testing.T) {
	// A basic test to make sure it runs without crashing
	ctx := context.Background()
	res := Fingerprint(ctx, "127.0.0.1", "B8:27:EB:11:22:33", 64, []int{}, 100)
	if res.Vendor != "Raspberry Pi" {
		t.Errorf("expected vendor Raspberry Pi, got %s", res.Vendor)
	}
}
