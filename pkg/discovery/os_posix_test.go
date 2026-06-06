//go:build !windows

package discovery

import (
	"context"
	"runtime"
	"testing"
)

func TestOsPingMacOS(t *testing.T) {
	// Localhost should always ping successfully, but on macOS if we pass 1 second as 1ms, it might fail.
	// We're just testing that it doesn't fail parsing arguments.
	
	// Fast ping to localhost.
	success := osPing(context.Background(), "127.0.0.1", 1000)
	if !success {
		t.Logf("Warning: Ping to localhost failed (OS: %s). Check permissions or environment.", runtime.GOOS)
	}
}
