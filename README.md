# catnet-io/engine

[![CI Status](https://img.shields.io/github/actions/workflow/status/catnet-io/engine/ci.yml?branch=main&style=flat-square)](https://github.com/catnet-io/engine/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/catnet-io/engine?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/catnet-io/engine?style=flat-square)](https://goreportcard.com/report/github.com/catnet-io/engine)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![pkg.go.dev](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/catnet-io/engine)
[![Latest Release](https://img.shields.io/github/v/release/catnet-io/engine?style=flat-square)](https://github.com/catnet-io/engine/releases)

Shared Go scanning engine for the CatNet ecosystem. Pure Go library — no binary, no CGO, no UI code.

## Documentation

- [API Stability Policy](docs/contracts/api-stability.md)
- [Compatibility Policy](docs/contracts/compatibility.md)
- [Event System Contract](docs/contracts/events.md)
- [Integration Examples](docs/examples/integration_examples.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Roadmap](ROADMAP.md)

## Packages

| Package | Key Exports | Description | Status |
|---|---|---|---|
| `pkg/engine` | `StartScan`, `ScanConfig`, `DefaultConfig`, `EventCallback` | Callback-based scan orchestrator. Use for CLI and scripting. | Stable |
| `pkg/scan` | `Engine`, `NewEngine`, `ScanStream`, `Stop` | Channel-based scan orchestrator. **Preferred for GUI and TUI consumers.** | Stable |
| `pkg/events` | `Event`, `EventType`, `HostDiscoveredData`, `ProgressData` | Async event types for the channel API. String-based `EventType` for safe serialization. | Stable |
| `pkg/profile` | `ScanProfile`, `DefaultProfile`, `Sanitize` | Scan configuration for the channel API. | Stable |
| `pkg/results` | `DeviceInfo`, `ScanReport`, `HostResult` | Core domain types used across the ecosystem. | Stable |
| `pkg/discovery` | `Ping`, `ReverseDNS`, `GetMAC` | Host liveness and resolution primitives. | Stable |
| `pkg/ports` | `ScanPorts` | TCP port scanner. Returns `<-chan int`. | Stable |
| `pkg/targets` | `ParseRange` | IP range parsing — CIDR and dash range. | Stable |
| `pkg/fingerprint` | `Fingerprint`, `FingerprintWithConfig`, `GrabBanners`, `VendorFromMAC` | Heuristic OS/device/vendor detection via TTL, banner, OUI, and RDP probe. | Stable · Experimental results |
| `pkg/topology` | `BuildGraph`, `ExportD3JSON`, `DetectGateway` | Network topology graph builder. | Stable · Experimental results |
| `pkg/export` | `ExportJSON`, `ExportCSV` | Export for `[]results.HostResult` with CSV sanitization. | Stable |
| `pkg/exporter` | `ExportJSON`, `ExportCSV`, `ExportXML` | Export for `*results.ScanReport`. | Stable |
| `pkg/coreerr` | `ErrTimeout`, `ErrCancelled`, `ErrInvalidInput` | Typed error taxonomy for `errors.Is`. | Stable |
| `pkg/store` | `ScanStore`, `NewSQLiteStore` | SQLite scan history. | Deprecated — moving to consumers |
| `pkg/diff` | `Compare`, `HostDiff` | Scan comparison. | Deprecated — moving to consumers |

## Quickstart

### CLI and scripting — callback API

Use `pkg/engine.StartScan` when you want a synchronous callback invoked per event.
Suitable for CLI tools and scripts where the callback is lightweight.

```go
package main

import (
	"context"
	"fmt"

	"github.com/catnet-io/engine/pkg/engine"
)

func main() {
	ips := []string{"192.168.1.1", "192.168.1.2"}
	cfg := engine.DefaultConfig()

	report, err := engine.StartScan(context.Background(), ips, cfg, func(ev engine.ScanEvent) {
		if ev.Type == engine.EventResult && ev.Device != nil && ev.Device.IsAlive {
			fmt.Printf("Found: %s (%s)\n", ev.Device.IP, ev.Device.Hostname)
		}
	})
	if err != nil {
		fmt.Printf("Scan failed: %v\n", err)
		return
	}
	fmt.Printf("Scanned: %d  Alive: %d\n", report.Total, report.Alive)
}
```

### GUI and TUI — channel API

Use `pkg/scan.Engine.ScanStream` when the consumer is a GUI or TUI that renders events
asynchronously. Events are delivered via a Go channel, decoupling scan workers from
rendering latency.

```go
package main

import (
	"context"
	"fmt"

	"github.com/catnet-io/engine/pkg/events"
	"github.com/catnet-io/engine/pkg/profile"
	"github.com/catnet-io/engine/pkg/scan"
)

func main() {
	ips := []string{"192.168.1.1", "192.168.1.2"}
	prof := profile.DefaultProfile()
	eng := scan.NewEngine()

	eventChan := make(chan events.Event)
	go func() {
		for ev := range eventChan {
			if ev.Type == events.HostDiscovered {
				data := ev.Data.(events.HostDiscoveredData)
				fmt.Printf("Found: %s\n", data.Host.IP)
			}
		}
	}()

	if err := eng.ScanStream(context.Background(), ips, prof, eventChan); err != nil {
		fmt.Printf("Scan failed: %v\n", err)
	}
}
```

## Fingerprinting

`pkg/fingerprint` identifies OS, device type, and vendor using three passive/active signal sources:

| Signal | Source | Method |
|---|---|---|
| TTL | ICMP response | Passive — inferred from hop count |
| OUI | MAC address prefix | Passive — offline vendor lookup |
| Banner | Open TCP ports | Active — protocol-specific probes |

**Active probes** send a minimal payload to elicit a service response:

| Port | Protocol | Probe | Notes |
|---|---|---|---|
| 80, 8080 | HTTP | `HEAD / HTTP/1.0` | — |
| 445 | SMB | Negotiate request | Opt-in — `BannerGrabConfig.AggressiveSMB = true`. May trigger IDS/IPS. |
| 3389 | RDP | TPKT + X.224 CR | Identifies MS-RDP from Connection Confirm flags. |

Use `FingerprintWithConfig` to control probe behavior:

```go
cfg := fingerprint.BannerGrabConfig{
	AggressiveSMB: false, // set true only if permitted by engagement rules
	Concurrency:   5,
}
result := fingerprint.FingerprintWithConfig(ctx, ip, mac, ttl, ports, timeoutMs, cfg)
```

## Security and CI

- `go vet`, `go test -race`, and dependency verification on every push
- Fuzz testing — `go test -fuzz=FuzzParseRange` for target parsing
- Linting via `golangci-lint` (see `.golangci.yml`)
- Vulnerability scanning with `govulncheck`
- Dependency updates via Dependabot
- GPG-signed automated merges from `develop` → `main`
- Security reporting guidance in [`SECURITY.md`](SECURITY.md)
- Full DevSecOps guidance in [`docs/devsecops.md`](docs/devsecops.md)

## Status

Current version: **v0.3.0** — see [CHANGELOG.md](CHANGELOG.md) for details.

## License

MIT — see [LICENSE](LICENSE).

## Part of the CatNet ecosystem

| | Repository | Role |
|---|---|---|
| ⚙️ | [catnet-io/engine](https://github.com/catnet-io/engine) | Shared Go scanning engine |
| 💻 | [catnet-io/catnet](https://github.com/catnet-io/catnet) | CLI |
| 🖥️ | [catnet-io/app](https://github.com/catnet-io/app) | Desktop app |
| 📟 | [catnet-io/tui](https://github.com/catnet-io/tui) | Terminal UI |