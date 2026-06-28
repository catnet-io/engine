# catnet-core

[![CI Status](https://img.shields.io/github/actions/workflow/status/catnet-io/engine/ci.yml?branch=main&style=flat-square)](https://github.com/catnet-io/engine/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/catnet-io/engine?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/catnet-io/engine?style=flat-square)](https://goreportcard.com/report/github.com/catnet-io/engine)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![pkg.go.dev](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/catnet-io/engine)
[![Latest Release](https://img.shields.io/github/v/release/catnet-io/engine?style=flat-square)](https://github.com/catnet-io/engine/releases)

Shared Go engine for the CatNet scanning ecosystem.

## Documentation & Contracts

Before adopting `catnet-core`, please review our documentation:
- [API Stability Policy](docs/contracts/api-stability.md)
- [Compatibility Policy](docs/contracts/compatibility.md)
- [Integration Examples](docs/examples/integration_examples.md)
- [Project Roadmap](ROADMAP.md)

## Packages

| Package | Public Types & Functions | Description |
|---|---|---|
| `pkg/engine` | `ScanConfig`, `DefaultConfig`, `StartScan` | Main scan orchestrator using `context.Context`. |
| `pkg/results` | `DeviceInfo`, `ScanReport`, `HostResult` | Core models used across the ecosystem. `HostResult` is the canonical type for the event-driven API. |
| `pkg/targets` | `ParseRange` | Target parsing and CIDR utilities. |
| `pkg/discovery` | `Ping`, `ReverseDNS`, `GetMAC` | Host liveness and resolution primitives. |
| `pkg/ports` | `ScanPorts` | Port scanning utilities. |
| `pkg/exporter` | `ExportJSON`, `ExportXML`, `ExportCSV` | Safe result export functions (`DeviceInfo`-based). |
| `pkg/export` | `ExportJSON`, `ExportCSV` | Export for `[]results.HostResult` with CSV sanitization. |
| `pkg/events` | `Event`, `EventType`, `HostDiscoveredData`, `ProgressData` | Async event system via Go channel. String-based `EventType` for Wails/TUI serialization. |
| `pkg/profile` | `ScanProfile`, `DefaultProfile`, `Sanitize` | Scan configuration with concurrency and timeout. |
| `pkg/scan` | `Engine`, `NewEngine`, `ScanStream`, `Stop`, `Ping`, `ReverseDNS`, `GetMAC`, `ScanPorts` | Event-driven orchestrator. Main entry point for frontends. |
| `pkg/fingerprint` | `Fingerprint`, `GrabBanners`, `VendorFromMAC` | Heuristic OS/device detection. **Experimental.** |
| `pkg/topology` | `BuildGraph`, `ExportD3JSON`, `DetectGateway` | Network topology graph builder. **Experimental.** |
| `pkg/coreerr` | `ErrTimeout`, `ErrCancelled`, `ErrInvalidInput`, ... | Structured error taxonomy for `errors.Is`. |

## Quickstart

```go
package main

import (
	"context"
	"fmt"

	"github.com/catnet-io/engine/pkg/engine"
	"github.com/catnet-io/engine/pkg/results"
)

func main() {
	ips := []string{"192.168.1.1"}
	cfg := engine.DefaultConfig()
	cfg.Sanitize()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	report, err := engine.StartScan(ctx, ips, cfg, func(event engine.ScanEvent) {
		switch event.Type {
		case engine.EventResult:
			if event.Device != nil && event.Device.IsAlive {
				fmt.Printf("Found: %s (%s)\n", event.Device.IP, event.Device.MAC)
			}
		case engine.EventProgress:
			fmt.Printf("Progress: %.2f%%\n", event.Progress*100)
		case engine.EventLifecycleStart:
			fmt.Println("Scan started...")
		case engine.EventLifecycleComplete:
			fmt.Println("Scan completed!")
		}
	})

	if err != nil {
		fmt.Printf("Scan failed: %v\n", err)
	} else {
		fmt.Printf("Total Scanned: %d, Alive: %d\n", report.Total, report.Alive)
	}
}
```

## Ecosystem

| Repository | Role |
|---|---|
| [`catnet-core`](https://github.com/catnet-io/engine) | Shared Go engine Ã¢â‚¬â€ no GUI |
| [`app`](https://github.com/catnet-io/app) | Desktop frontend (Raygui) Ã¢â‚¬â€ Planned evolution to Wails + React |
| [`catnet`](https://github.com/catnet-io/catnet) | Scriptable Go CLI |
| [`tui`](https://github.com/catnet-io/tui) | Interactive TUI (Go + Bubble Tea) |

## Status

Current version: v0.3.0
See [CHANGELOG.md](CHANGELOG.md) for details.

## Security and CI

This repository follows DevSecOps practices by integrating quality and security checks into the development workflow:

- GitHub Actions CI on `main` and `develop`
- `go vet`, `go test -race`, and dependency verification
- `go test -fuzz=FuzzParseRange` fuzzing of target parsing
- linting via `golangci-lint`
- vulnerability scanning with `govulncheck`
- dependency updates via Dependabot
- security reporting guidance in `SECURITY.md`
- full DevSecOps guidance in [docs/devsecops.md](docs/devsecops.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
