# Integration Examples

This document demonstrates how to integrate `catnet-core` into different types of frontends safely and efficiently.

## 1. CLI Consumer (Synchronous)

CLI applications usually block the main thread and print to `stdout` sequentially.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mendsec/catnet-core/pkg/engine"
	"github.com/mendsec/catnet-core/pkg/exporter"
)

func main() {
	ips := []string{"192.168.1.1", "192.168.1.2"}
	cfg := engine.DefaultConfig()

	ctx := context.Background()

	report, err := engine.StartScan(ctx, ips, cfg, func(event engine.ScanEvent) {
		switch event.Type {
		case engine.EventLifecycleStart:
			fmt.Println("Starting scan...")
		case engine.EventResult:
			if event.Device.IsAlive {
				fmt.Printf("[+] Alive: %s (MAC: %s)\n", event.Device.IP, event.Device.MAC)
			}
		}
	})

	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	jsonBytes, _ := exporter.ExportJSON(report)
	os.WriteFile("output.json", jsonBytes, 0644)
	fmt.Println("Scan saved to output.json")
}
```

## 2. TUI Consumer (Bubble Tea / Async Channels)

TUI frameworks like Bubble Tea require events to be sent as messages so the UI can update asynchronously without blocking the render loop.

```go
package tui

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mendsec/catnet-core/pkg/engine"
)

type scanMsg engine.ScanEvent

func startEngineCommand(ips []string, cfg engine.ScanConfig) tea.Cmd {
	return func() tea.Msg {
		ch := make(chan engine.ScanEvent, 100)
		
		go func() {
			engine.StartScan(context.Background(), ips, cfg, func(e engine.ScanEvent) {
				ch <- e // Push event to channel to decouple from worker threads
			})
			close(ch)
		}()
		
		// Return the first event as a message to start the UI loop
		if e, ok := <-ch; ok {
			return scanMsg(e)
		}
		return nil
	}
}
```

## 3. GUI Consumer (Wails/React - Debouncing)

In GUI wrappers like Wails or Electron, sending thousands of progress updates via RPC per second will freeze the webview. You should debounce the `EventProgress`.

```go
package gui

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/mendsec/catnet-core/pkg/engine"
)

type App struct {
	ctx context.Context
}

func (a *App) RunScan(ips []string) {
	cfg := engine.DefaultConfig()
	
	var lastUpdate time.Time

	engine.StartScan(a.ctx, ips, cfg, func(event engine.ScanEvent) {
		if event.Type == engine.EventProgress {
			// Debounce to 10 FPS (100ms)
			if time.Since(lastUpdate) < 100*time.Millisecond {
				return
			}
			lastUpdate = time.Now()
			runtime.EventsEmit(a.ctx, "scan_progress", event.Progress)
		} else if event.Type == engine.EventResult {
			// Results can be sent immediately or batched
			runtime.EventsEmit(a.ctx, "scan_result", event.Device)
		}
	})
}
```
