# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2026-06-06

### Added
- Added package-level documentation (`doc.go`) to all public packages and `internal/netutil`.
- Added mockable, deterministic unit tests for input validation in `pkg/discovery`.
- Added `macos-latest` to GitHub Actions CI workflow to validate POSIX builds.
- Added fuzz testing (`FuzzParseRange`) to target parsing logic to ensure resilience against malformed inputs.
- Added `os_posix_test.go` with OS-specific validations for fallback ping execution.

### Changed
- **BREAKING CHANGE**: Removed `OpenPortsCount` field from `DeviceInfo` struct to eliminate redundancy and potential state corruption. Use the new `PortCount()` method instead.
- Improved defensive timeout calculation in `StartScan` by considering concurrent port scan batches, avoiding overestimated timeouts.
- Updated minimum Go version requirement to `1.26.4` in `go.mod` and `1.26.x` in CI workflows.
- Improved `sanitizeCSVField` to block newline (`\n` / `\r`) injections and control characters within the body of CSV fields.
- Context is now explicitly propagated to blocking I/O calls (`discovery.Ping` and `ports.ScanPorts`) for immediate hard-cancellation support.

### Fixed
- Fixed pointer memory aliasing bug in `EventCallback` within `StartScan` loop, ensuring consumers receive a safe, distinct copy of `DeviceInfo`.
- Fixed `osPing` in both Windows and POSIX implementations to handle `timeoutMs <= 0` gracefully with a safe default (1000ms).
- Fixed POSIX `osPing` ignoring `timeoutMs` parameter, properly converting to whole seconds for the `-W` flag on Linux, and using milliseconds on macOS.

## [0.1.0] - 2026-06-02

### Added
- `pkg/scanner`: `DeviceInfo`, `ScanConfig`, `DefaultConfig`, `Sanitize`.
- `pkg/scanner`: `validateIPv4`, `Ping`, `ReverseDNS`, `GetMAC`, `ScanPorts`.
- `pkg/scanner`: `ParseRange` com suporte a CIDR e dash range.
- `pkg/scanner`: `StartScan`, `StopScan` com goroutines e cancelamento por context.
- `pkg/scanner`: build tags separados para Windows (`SendARP`) e POSIX (`arp`).
- `pkg/exporter`: `ExportJSON`, `ExportCSV`, `ExportXML` com sanitização de injeção CSV.
- CI: `go vet` + `go test -race` no GitHub Actions.
- CI: `govulncheck` semanal.

[Unreleased]: https://github.com/mendsec/catnet-core/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/mendsec/catnet-core/releases/tag/v0.1.1
[0.1.0]: https://github.com/mendsec/catnet-core/releases/tag/v0.1.0
