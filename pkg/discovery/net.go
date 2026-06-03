package discovery

import (
	"net"
	"strings"

	"github.com/mendsec/catnet-core/internal/netutil"
)

// Ping realiza um ping ICMP na máquina com timeout em milissegundos.
func Ping(ip string, timeoutMs int) bool {
	if err := netutil.ValidateIPv4(ip); err != nil {
		return false
	}
	return osPing(ip, timeoutMs)
}

// ReverseDNS resolve o nome do host do endereço IP dado.
func ReverseDNS(ip string) string {
	if err := netutil.ValidateIPv4(ip); err != nil {
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
	if err := netutil.ValidateIPv4(ip); err != nil {
		return ""
	}
	return osGetMAC(ip)
}


