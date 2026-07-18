# Scan Lifecycle Events

The `catnet-core` engine exposes a stream of typed events during scan execution to allow consumers (CLI, TUI, GUI) to react to progress and lifecycle changes without coupling to internal logic.

## Event Signature

Events are dispatched synchronously via the `engine.EventCallback` function:
```go
type EventCallback func(event ScanEvent)
```

## ScanEvent Struct

```go
type ScanEvent struct {
	Type     ScanEventType
	Device   *results.DeviceInfo
	Progress float64
	Message  string
}
```

## Event Types

| Type | Description | Included Data |
|------|-------------|---------------|
| `EventLifecycleStart` | Fired once when `StartScan` is invoked. | `Message` |
| `EventLifecycleComplete` | Fired once when `StartScan` completes successfully. | `Message` |
| `EventLifecycleCancel` | Fired once when the scan is cancelled or times out. | `Message` |
| `EventWarning` | Fired when a non-fatal error occurs (e.g. permission issues). | `Message` |
| `EventProgress` | Fired periodically to indicate scan progress. | `Progress` (0.0 to 1.0) |
| `EventResult` | Fired when a device scan completes (whether alive or dead). | `Device` |

## Consumer Responsibilities

1. **Asynchronous dispatch safety:** Since v0.5.0, the engine uses an internal asynchronous dispatcher (`asyncDispatcher`) with a buffered channel to decouple scan worker execution from callback execution. However, consumers should still keep callback code efficient to prevent backpressure from eventually filling the buffer limit (512) and blocking.
2. **Thread safety:** The `Device` pointer in `EventResult` points to a copied value, but consumers should avoid mutating it.
