# catnet-core

Shared Go engine for the CatNet scanning ecosystem.

## Packages

| Package | Funções e Tipos Públicos | Descrição |
|---|---|---|
| `pkg/scanner` | `DeviceInfo`, `ScanConfig`, `DefaultConfig`, `StartScan`, `StopScan`, `Ping`, `ReverseDNS`, `GetMAC`, `ScanPorts`, `ParseRange` | Lógica central de varredura e primitivas de rede. |
| `pkg/exporter` | `ExportJSON`, `ExportXML`, `ExportCSV` | Funções para exportação segura de resultados (prevenção contra CSV injection). |

## Quickstart

```go
package main

import (
	"fmt"
	"github.com/mendsec/catnet-core/pkg/scanner"
)

func main() {
	ips := []string{"192.168.1.1"}
	cfg := scanner.DefaultConfig()
	cfg.Sanitize()

	err := scanner.StartScan(ips, cfg, func(d scanner.DeviceInfo) {
		if d.IsAlive {
			fmt.Printf("Encontrado: %s (%s)\n", d.IP, d.MAC)
		}
	}, func(progress float64) {
		fmt.Printf("Progresso: %.2f%%\n", progress*100)
	})

	if err != nil {
		panic(err)
	}
}
```

## Ecosystem

| Repositório | Papel |
|---|---|
| [`catnet-core`](https://github.com/mendsec/catnet-core) | Engine Go compartilhada — sem GUI |
| [`catnet-scanner`](https://github.com/mendsec/catnet-scanner) | Frontend desktop Go + Wails + React |
| [`catnet`](https://github.com/mendsec/catnet) | CLI Go scriptável |
| [`catnet-tui`](https://github.com/mendsec/catnet-tui) | TUI interativa Go + Bubble Tea |

## Status

Versão atual: v0.1.0
Veja o [CHANGELOG.md](CHANGELOG.md) para detalhes.
