//go:build windows

package discovery

import (
	"context"
	"testing"
)

func TestOsPingWindows(t *testing.T) {
	ctx := context.Background()
	// Just test if it compiles
	_ = osPing(ctx, "127.0.0.1", 1000)
}
