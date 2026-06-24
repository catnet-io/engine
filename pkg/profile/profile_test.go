package profile

import "testing"

func TestDefaultProfile(t *testing.T) {
	p := DefaultProfile()
	if p.Concurrency != 64 {
		t.Errorf("Concurrency = %d, want 64", p.Concurrency)
	}
	if p.TimeoutMs != 1000 {
		t.Errorf("TimeoutMs = %d, want 1000", p.TimeoutMs)
	}
	if len(p.DefaultPorts) == 0 {
		t.Errorf("DefaultPorts is empty")
	}
}

func TestSanitizeClampsConcurrency(t *testing.T) {
	p := ScanProfile{Concurrency: 0, TimeoutMs: 500}
	p.Sanitize()
	if p.Concurrency != 16 {
		t.Errorf("Concurrency should be clamped to 16, got %d", p.Concurrency)
	}
}

func TestSanitizeClampsHighConcurrency(t *testing.T) {
	p := ScanProfile{Concurrency: 1000, TimeoutMs: 500}
	p.Sanitize()
	if p.Concurrency != 16 {
		t.Errorf("Concurrency should be clamped to 16, got %d", p.Concurrency)
	}
}

func TestSanitizeClampsTimeout(t *testing.T) {
	p := ScanProfile{Concurrency: 10, TimeoutMs: 0}
	p.Sanitize()
	if p.TimeoutMs != 1000 {
		t.Errorf("TimeoutMs should be clamped to 1000, got %d", p.TimeoutMs)
	}
}

func TestSanitizeClampsHighTimeout(t *testing.T) {
	p := ScanProfile{Concurrency: 10, TimeoutMs: 20000}
	p.Sanitize()
	if p.TimeoutMs != 1000 {
		t.Errorf("TimeoutMs should be clamped to 1000, got %d", p.TimeoutMs)
	}
}

func TestSanitizeValidValues(t *testing.T) {
	p := ScanProfile{Concurrency: 32, TimeoutMs: 2000}
	p.Sanitize()
	if p.Concurrency != 32 {
		t.Errorf("Concurrency changed from 32 to %d", p.Concurrency)
	}
	if p.TimeoutMs != 2000 {
		t.Errorf("TimeoutMs changed from 2000 to %d", p.TimeoutMs)
	}
}
