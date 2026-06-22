package ports

import (
	"context"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mendsec/catnet-core/internal/netutil"
)

// ScanConcurrency define o número máximo de conexões TCP simultâneas por IP.
const ScanConcurrency = 10

// ScanPorts varre uma lista de portas em um IP concorrentemente e retorna os resultados via canal.
// O canal é fechado automaticamente quando todas as portas foram testadas.
func ScanPorts(ctx context.Context, ip string, ports []int, timeoutMs int) <-chan int {
	out := make(chan int, len(ports))
	if err := netutil.ValidateIPv4(ip); err != nil {
		close(out)
		return out
	}
	
	timeout := time.Duration(timeoutMs) * time.Millisecond

	var wg sync.WaitGroup
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
					out <- p
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
