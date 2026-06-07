# catnet-core

[![CI Status](https://img.shields.io/github/actions/workflow/status/mendsec/catnet-core/ci.yml?branch=main&style=flat-square)](https://github.com/mendsec/catnet-core/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mendsec/catnet-core?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/mendsec/catnet-core?style=flat-square)](https://goreportcard.com/report/github.com/mendsec/catnet-core)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![pkg.go.dev](https://img.shields.io/badge/pkg.go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/mendsec/catnet-core)
[![Latest Release](https://img.shields.io/github/v/release/mendsec/catnet-core?style=flat-square)](https://github.com/mendsec/catnet-core/releases)

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
| `pkg/results` | `DeviceInfo` | Core models used across the ecosystem. |
| `pkg/targets` | `ParseRange` | Target parsing and CIDR utilities. |
| `pkg/discovery` | `Ping`, `ReverseDNS`, `GetMAC` | Host liveness and resolution primitives. |
| `pkg/ports` | `ScanPorts` | Port scanning utilities. |
| `pkg/exporter` | `ExportJSON`, `ExportXML`, `ExportCSV` | Safe result export functions. JSON is the canonical schema reference format. |
| `pkg/scanner` | *Deprecated* | Temporary compatibility wrapper. Do not use for new code. |

## Quickstart

```go
package main

import (
	"context"
	"fmt"

	"github.com/mendsec/catnet-core/pkg/engine"
	"github.com/mendsec/catnet-core/pkg/results"
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
| [`catnet-core`](https://github.com/mendsec/catnet-core) | Shared Go engine — no GUI |
| [`catnet-scanner`](https://github.com/mendsec/catnet-scanner) | Desktop frontend (Raygui) — Planned evolution to Wails + React |
| [`catnet`](https://github.com/mendsec/catnet) | Scriptable Go CLI |
| [`catnet-tui`](https://github.com/mendsec/catnet-tui) | Interactive TUI (Go + Bubble Tea) |

## Status

Current version: v0.2.0
See [CHANGELOG.md](CHANGELOG.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
