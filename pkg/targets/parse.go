package targets

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// ParseRange interpreta uma string e retorna uma lista de endereços IP correspondentes.
func ParseRange(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	if strings.Contains(input, "/") {
		return parseCIDR(input)
	}
	if strings.Contains(input, "-") {
		return parseDashRange(input)
	}
	parsed := net.ParseIP(input)
	if parsed == nil {
		return nil, fmt.Errorf("invalid IP format")
	}
	return []string{parsed.String()}, nil
}

func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	// ⚡ Bolt Optimization: Pre-allocate slice to prevent dynamic resizing overhead.
	ones, bits := ipnet.Mask.Size()
	// Avoid massive instant allocations on huge CIDRs (e.g. /0 or /8)
	if bits - ones > 16 { // Cap at 65536 hosts
		return nil, fmt.Errorf("CIDR range too large (max 65536 hosts)")
	}
	capacity := 1 << (bits - ones)
	ips := make([]string, 0, capacity)
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) > 2 { // Skip network and broadcast addresses
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

func parseDashRange(dashStr string) ([]string, error) {
	parts := strings.Split(dashStr, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range format, expected start-end")
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
		return nil, fmt.Errorf("invalid IP in range")
	}
	start := binary.BigEndian.Uint32(startIP)
	end := binary.BigEndian.Uint32(endIP)
	if start > end {
		return nil, fmt.Errorf("start IP is greater than end IP")
	}
	if end-start > 65536 {
		return nil, fmt.Errorf("range too large (max 65536)")
	}
	// ⚡ Bolt Optimization: Pre-allocate slice and reuse IP buffer to reduce memory allocations.
	capacity := end - start + 1
	ips := make([]string, 0, capacity)
	ip := make(net.IP, 4)
	for i := start; i <= end; i++ {
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
