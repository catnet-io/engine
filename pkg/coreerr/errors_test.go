package coreerr

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorWrapping(t *testing.T) {
	err := fmt.Errorf("scan failed: %w", ErrTimeout)
	
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("Expected err to wrap ErrTimeout, but errors.Is returned false")
	}
	
	if errors.Is(err, ErrCancelled) {
		t.Errorf("Expected err to not wrap ErrCancelled")
	}
}
