package ports

import (
	"context"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/mendsec/catnet-core/internal/netutil"
)

// ScanConcurrency define o número máximo de conexões TCP simultâneas por IP.
const ScanConcurrency = 10

// ScanPorts varre uma lista de portas em um IP e retorna as abertas.
// ⚡ Bolt Optimization: Concurrently scan ports to prevent cumulative timeouts from blocking the scan.
func ScanPorts(ctx context.Context, ip string, ports []int, timeoutMs int) []int {
	if err := netutil.ValidateIPv4(ip); err != nil {
		return nil
	}
	var openPorts []int
	timeout := time.Duration(timeoutMs) * time.Millisecond

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Limit concurrent connections per IP to prevent FD exhaustion
	sem := make(chan struct{}, ScanConcurrency)

	for _, port := range ports {
		wg.Add(1)
		sem <- struct{}{}

		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			
			if ctx.Err() != nil {
				return
			}

			address := net.JoinHostPort(ip, strconv.Itoa(p))
			dialer := net.Dialer{Timeout: timeout}
			conn, err := dialer.DialContext(ctx, "tcp", address)
			if err == nil {
				conn.Close()
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()

	// Keep output deterministic
	sort.Ints(openPorts)

	return openPorts
}
