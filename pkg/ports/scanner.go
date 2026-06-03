package ports

import (
	"net"
	"strconv"
	"time"

	"github.com/mendsec/catnet-core/internal/netutil"
)

// ScanPorts varre uma lista de portas em um IP e retorna as abertas.
func ScanPorts(ip string, ports []int, timeoutMs int) []int {
	if err := netutil.ValidateIPv4(ip); err != nil {
		return nil
	}
	var openPorts []int
	timeout := time.Duration(timeoutMs) * time.Millisecond
	for _, port := range ports {
		address := net.JoinHostPort(ip, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			conn.Close()
			openPorts = append(openPorts, port)
		}
	}
	return openPorts
}
