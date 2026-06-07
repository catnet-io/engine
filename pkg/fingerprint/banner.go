package fingerprint

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// GrabBanners attempts to read banners from open ports.
func GrabBanners(ctx context.Context, ip string, openPorts []int, timeoutMs int) map[int]string {
	banners := make(map[int]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	timeout := time.Duration(timeoutMs) * time.Millisecond

	for _, port := range openPorts {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			banner := grabBannerFromPort(ctx, ip, p, timeout)
			if banner != "" {
				mu.Lock()
				banners[p] = banner
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return banners
}

func grabBannerFromPort(ctx context.Context, ip string, port int, timeout time.Duration) string {
	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return ""
	}
	defer conn.Close()

	stop := context.AfterFunc(ctx, func() {
		conn.Close()
	})
	defer stop()

	conn.SetDeadline(time.Now().Add(timeout))

	if port == 80 || port == 8080 {
		req := "HEAD / HTTP/1.0\r\n\r\n"
		conn.Write([]byte(req))
	} else if port == 445 {
		// Basic SMB negotiate request
		smbReq := []byte{
			0x00, 0x00, 0x00, 0x2f, 0xff, 0x53, 0x4d, 0x42,
			0x72, 0x00, 0x00, 0x00, 0x00, 0x08, 0x01, 0x40,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x0c,
			0x00, 0x00, 0x01, 0x00, 0x00, 0x0c, 0x00, 0x02,
			0x4e, 0x54, 0x20, 0x4c, 0x4d, 0x20, 0x30, 0x2e,
			0x31, 0x32, 0x00,
		}
		conn.Write(smbReq)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err == nil && n > 0 {
		return string(buf[:n])
	}
	return ""
}

// OsFromBanners infers the OS and device type from collected banners.
func OsFromBanners(banners map[int]string) FingerprintResult {
	res := FingerprintResult{
		OS:         "Unknown",
		OSFamily:   "unknown",
		DeviceType: DeviceUnknown,
		Confidence: 0,
	}

	hasSSH := false
	if _, ok := banners[22]; ok {
		hasSSH = true
	}

	for port, banner := range banners {
		lowerBanner := strings.ToLower(banner)

		if port == 22 {
			if strings.Contains(lowerBanner, "ubuntu") || strings.Contains(lowerBanner, "debian") || strings.Contains(lowerBanner, "centos") {
				res.OSFamily = "linux"
				res.OS = "Linux"
				if strings.Contains(lowerBanner, "ubuntu") {
					res.OS = "Ubuntu"
				} else if strings.Contains(lowerBanner, "debian") {
					res.OS = "Debian"
				} else if strings.Contains(lowerBanner, "centos") {
					res.OS = "CentOS"
				}
				res.DeviceType = DeviceServer
				res.Confidence = 90
				return res
			}
			if strings.Contains(lowerBanner, "openssh_for_windows") {
				res.OSFamily = "windows"
				res.OS = "Windows"
				res.DeviceType = DeviceServer
				res.Confidence = 90
				return res
			}
		}

		if port == 80 || port == 8080 {
			if strings.Contains(lowerBanner, "apache") || strings.Contains(lowerBanner, "nginx") {
				res.OSFamily = "linux"
				res.OS = "Linux"
				res.DeviceType = DeviceServer
				res.Confidence = 80
			} else if strings.Contains(lowerBanner, "iis") {
				res.OSFamily = "windows"
				res.OS = "Windows"
				res.DeviceType = DeviceServer
				res.Confidence = 80
				return res
			} else if strings.Contains(lowerBanner, "lighttpd") {
				res.DeviceType = DeviceIoT
				res.Confidence = 70
			}
		}

		if port == 445 {
			// Windows SMB detect
			if len(lowerBanner) > 0 {
				res.OSFamily = "windows"
				res.OS = "Windows"
				res.Confidence = 70
			}
		}
	}

	// IoT heuristics: lighttpd OR (port 8888/80 without SSH)
	_, has80 := banners[80]
	_, has8888 := banners[8888]
	hasLighttpd := false
	for _, b := range banners {
		if strings.Contains(strings.ToLower(b), "lighttpd") {
			hasLighttpd = true
			break
		}
	}

	if res.DeviceType == DeviceUnknown {
		if hasLighttpd || ((has80 || has8888) && !hasSSH) {
			res.DeviceType = DeviceIoT
			res.Confidence = 60
		}
	}

	return res
}
