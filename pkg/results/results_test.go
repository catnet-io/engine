package results

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHostResultJSON(t *testing.T) {
	h := HostResult{
		IP: "192.168.1.1", Alive: true,
		Hostname: "router", MAC: "AA-BB-CC-DD-EE-FF",
		OpenPorts: []int{80, 443},
	}
	data, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var h2 HostResult
	if err := json.Unmarshal(data, &h2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if h2.IP != h.IP || h2.Alive != h.Alive || len(h2.OpenPorts) != len(h.OpenPorts) {
		t.Errorf("round-trip mismatch: got %+v", h2)
	}
}

func TestToDeviceInfo(t *testing.T) {
	h := HostResult{
		IP: "10.0.0.1", Alive: true,
		Hostname: "test", MAC: "AA:BB:CC:DD:EE:FF",
		OpenPorts: []int{22, 80},
	}
	d := h.ToDeviceInfo()
	if d.IP != h.IP {
		t.Errorf("IP mismatch: %s != %s", d.IP, h.IP)
	}
	if d.IsAlive != h.Alive {
		t.Errorf("IsAlive mismatch: %v != %v", d.IsAlive, h.Alive)
	}
	if d.Hostname != h.Hostname {
		t.Errorf("Hostname mismatch: %s != %s", d.Hostname, h.Hostname)
	}
	if d.MAC != h.MAC {
		t.Errorf("MAC mismatch: %s != %s", d.MAC, h.MAC)
	}
	if len(d.OpenPorts) != len(h.OpenPorts) {
		t.Errorf("OpenPorts length mismatch: %d != %d", len(d.OpenPorts), len(h.OpenPorts))
	}
}

func TestHostResultJSONTags(t *testing.T) {
	h := HostResult{IP: "192.168.1.1", Alive: true}
	data, _ := json.Marshal(h)
	s := string(data)
	if !strings.Contains(s, `"alive":`) {
		t.Errorf("expected JSON key 'alive', got: %s", s)
	}
	if !strings.Contains(s, `"open_ports":`) {
		t.Errorf("expected JSON key 'open_ports', got: %s", s)
	}
}
