package ports

import (
	"context"
	"net"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
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

	// ⚡ Bolt Optimization: Replace channel semaphore with a lock-free atomic index counter.
	// Bypasses the O(N) goroutine allocation and channel operations upfront,
	// drastically reducing memory overhead and improving throughput.
	var index int32 = -1

	// Limit concurrent connections per IP to prevent FD exhaustion
	workers := ScanConcurrency
	if len(ports) < workers {
		workers = len(ports)
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if ctx.Err() != nil {
					return
				}

				idx := int(atomic.AddInt32(&index, 1))
				if idx >= len(ports) {
					return
				}
				p := ports[idx]

				address := net.JoinHostPort(ip, strconv.Itoa(p))
				dialer := net.Dialer{Timeout: timeout}
				conn, err := dialer.DialContext(ctx, "tcp", address)
				if err == nil {
					conn.Close()
					mu.Lock()
					openPorts = append(openPorts, p)
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	// Keep output deterministic
	sort.Ints(openPorts)

	return openPorts
}
