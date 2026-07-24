package discovery

import (
	"context"
	"net"
	"strings"

	"github.com/catnet-io/engine/internal/netutil"
)

// Ping realiza um ping ICMP na mÃ¡quina com timeout em milissegundos.
func Ping(ctx context.Context, ip string, timeoutMs int) bool {
	if err := netutil.ValidateIPv4(ip); err != nil {
		return false
	}
	return osPing(ctx, ip, timeoutMs)
}

// ReverseDNS resolve o nome do host do endereÃ§o IP dado.
func ReverseDNS(ctx context.Context, ip string) string {
	if err := netutil.ValidateIPv4(ip); err != nil {
		return ""
	}
	names, err := net.DefaultResolver.LookupAddr(ctx, ip)
	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], ".")
	}
	return ""
}

// GetMAC tenta obter o endereÃ§o MAC da mÃ¡quina alvo.
func GetMAC(ctx context.Context, ip string) string {
	if err := netutil.ValidateIPv4(ip); err != nil {
		return ""
	}
	return osGetMAC(ctx, ip)
}
