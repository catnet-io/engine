# catnet-core

Shared Go engine for the CatNet scanning ecosystem.

## Documentation & Contracts

Before adopting `catnet-core`, please review our [API Stability Policy](docs/contracts/api-stability.md).

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
| [`catnet-scanner`](https://github.com/mendsec/catnet-scanner) | Desktop frontend (Go + Wails + React) |
| [`catnet`](https://github.com/mendsec/catnet) | Scriptable Go CLI |
| [`catnet-tui`](https://github.com/mendsec/catnet-tui) | Interactive TUI (Go + Bubble Tea) |

## Status

Current version: v0.1.0
See [CHANGELOG.md](CHANGELOG.md) for details.
