package ports

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func validateIPv4(ip string) error {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return fmt.Errorf("endereço IP inválido: %q", ip)
	}
	if parsed.To4() == nil {
		return fmt.Errorf("apenas IPv4 é suportado: %q", ip)
	}
	return nil
}

// ScanPorts varre uma lista de portas em um IP e retorna as abertas.
func ScanPorts(ip string, ports []int, timeoutMs int) []int {
	if err := validateIPv4(ip); err != nil {
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
