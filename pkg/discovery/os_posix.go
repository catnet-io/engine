//go:build !windows

package discovery

import (
	"net"
	"os/exec"
	"strings"
)

// osPing faz ping em sistemas POSIX
func osPing(ip string, timeoutMs int) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	return cmd.Run() == nil
}

// osGetMAC obtém o MAC em sistemas POSIX usando comando arp
func osGetMAC(ip string) string {
	if net.ParseIP(ip) == nil {
		return ""
	}
	cmd := exec.Command("arp", "-an")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	searchParen := "(" + ip + ")"
	searchSpace := " " + ip + " "
	for _, line := range lines {
		if strings.Contains(line, searchParen) || strings.Contains(line, searchSpace) {
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
