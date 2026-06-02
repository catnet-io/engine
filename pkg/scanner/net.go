package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
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

// Ping realiza um ping ICMP na máquina com timeout em milissegundos.
func Ping(ip string, timeoutMs int) bool {
	if err := validateIPv4(ip); err != nil {
		return false
	}
	return osPing(ip, timeoutMs)
}

// ReverseDNS resolve o nome do host do endereço IP dado.
func ReverseDNS(ip string) string {
	if err := validateIPv4(ip); err != nil {
		return ""
	}
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], ".")
	}
	return ""
}

// GetMAC tenta obter o endereço MAC da máquina alvo.
func GetMAC(ip string) string {
	if err := validateIPv4(ip); err != nil {
		return ""
	}
	return osGetMAC(ip)
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
