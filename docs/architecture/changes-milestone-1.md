# Milestone 1: Core Architecture Hardening

## Overview

As part of Milestone 1, `catnet-core` underwent significant architectural hardening. The goal was to remove the monolithic `pkg/scanner` and split responsibilities into focused domain packages, establish a context-aware execution model, and set strict API stability guidelines for downstream consumers.

## Package Boundaries

### Before
All scanning, parsing, network, and execution logic was housed inside a single package: `pkg/scanner`.
The entry point `StartScan` relied on global variables (`sync.Mutex`, `atomic.Bool`) preventing multiple concurrent scan jobs from running in the same process.

### After
The responsibilities were divided into the following domain packages:
- **`pkg/engine`**: Contains `ScanConfig`, `DefaultConfig`, and the context-aware `StartScan` orchestrator.
- **`pkg/results`**: Defines `DeviceInfo` as the standard ecosystem result schema.
- **`pkg/targets`**: Contains `ParseRange` and internal CIDR/dash range parsing tools.
- **`pkg/discovery`**: Contains `Ping` and OS-specific liveness logic, `ReverseDNS`, and `GetMAC`.
- **`pkg/ports`**: Contains `ScanPorts` logic.

### Deprecation Shim
A new temporary file `pkg/scanner/shim.go` was introduced. `pkg/scanner` is now officially **deprecated**. The shim aliases the structs, functions, and variables to their new respective packages, maintaining backward compatibility while consumers like `catnet-scanner` and `catnet-tui` migrate.

## API Changes

### Context-Aware Execution
- `engine.StartScan` now accepts a `context.Context` as its first parameter instead of relying on a global package lock.
- `scanner.StopScan` is deprecated. While the shim maintains it using a hidden package-level context, all new integrations should manage cancellation directly through the passed context.

### Global State Removal
- Global variables `isScanning`, `scanMu`, and `cancelScan` were completely removed from `pkg/engine`. 

## File Operations

- **Created Directories**: `pkg/engine`, `pkg/results`, `pkg/targets`, `pkg/discovery`, `pkg/ports`, `docs/contracts`, `docs/architecture`.
- **Moved/Renamed**:
  - `pkg/scanner/scan.go` -> `pkg/engine/scan.go`
  - `pkg/scanner/scan_test.go` -> `pkg/engine/scan_test.go`
  - `pkg/scanner/config.go` -> `pkg/engine/config.go`
  - `pkg/scanner/device.go` -> `pkg/results/device.go`
  - `pkg/scanner/utils.go` -> `pkg/targets/parse.go`
  - `pkg/scanner/utils_test.go` -> `pkg/targets/parse_test.go`
  - `pkg/scanner/os_windows.go` -> `pkg/discovery/os_windows.go`
  - `pkg/scanner/os_posix.go` -> `pkg/discovery/os_posix.go`
  - `pkg/scanner/net.go` -> `pkg/discovery/net.go`
  - `pkg/scanner/net_test.go` -> `pkg/discovery/net_test.go`
- **Created**: 
  - `pkg/ports/scanner.go` (Extracted `ScanPorts`)
  - `pkg/scanner/shim.go` (Deprecation wrapper)
  - `docs/contracts/api-stability.md`
  - `docs/architecture/changes-milestone-1.md`

## Test Enhancements
All tests were adjusted to target the new domain boundaries. Removed unused imports and verified that `go test ./...` and `go build ./...` succeed seamlessly on both Windows and POSIX environments.
