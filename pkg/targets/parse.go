package targets

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/mendsec/catnet-core/pkg/coreerr"
)

// ParseRange interpreta uma string e retorna uma lista de endereços IP correspondentes.
func ParseRange(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("%w: empty input", coreerr.ErrInvalidInput)
	}
	if strings.Contains(input, "/") {
		return parseCIDR(input)
	}
	if strings.Contains(input, "-") {
		return parseDashRange(input)
	}
	parsed := net.ParseIP(input)
	if parsed == nil {
		return nil, fmt.Errorf("%w: invalid IP format", coreerr.ErrInvalidInput)
	}
	return []string{parsed.String()}, nil
}

func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

func parseDashRange(dashStr string) ([]string, error) {
	parts := strings.Split(dashStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: invalid range format, expected start-end", coreerr.ErrInvalidInput)
	}
	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])
	if !strings.Contains(endStr, ".") {
		lastDot := strings.LastIndex(startStr, ".")
		if lastDot != -1 {
			endStr = startStr[:lastDot+1] + endStr
		}
	}
	startIP := net.ParseIP(startStr).To4()
	endIP := net.ParseIP(endStr).To4()
	if startIP == nil || endIP == nil {
		return nil, fmt.Errorf("%w: invalid IP in range", coreerr.ErrInvalidInput)
	}
	start := binary.BigEndian.Uint32(startIP)
	end := binary.BigEndian.Uint32(endIP)
	if start > end {
		return nil, fmt.Errorf("%w: start IP is greater than end IP", coreerr.ErrInvalidInput)
	}
	if end-start > 65536 {
		return nil, fmt.Errorf("%w: range too large (max 65536)", coreerr.ErrInvalidInput)
	}
	var ips []string
	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
