//go:build !windows

package scanner

import (
	"os/exec"
	"strings"
)

// osPing faz ping em sistemas POSIX
func osPing(ip string, timeoutMs int) bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	return cmd.Run() == nil
}

// osGetMAC obtém o MAC em sistemas POSIX usando comando arp
func osGetMAC(ip string) string {
	cmd := exec.Command("arp", "-n", ip)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) {
			parts := strings.Fields(line)
			for _, p := range parts {
				if strings.Contains(p, ":") && len(p) == 17 {
					return strings.ToUpper(strings.ReplaceAll(p, ":", "-"))
				}
			}
		}
	}
	return ""
}
