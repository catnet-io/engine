package netutil

import (
	"fmt"
	"net"
)

// ValidateIPv4 ensures the provided string is a valid IPv4 address.
func ValidateIPv4(ip string) error {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return fmt.Errorf("endereço IP inválido: %q", ip)
	}
	if parsed.To4() == nil {
		return fmt.Errorf("apenas IPv4 é suportado: %q", ip)
	}
	return nil
}
