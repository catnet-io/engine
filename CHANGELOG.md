# Changelog

## [0.1.0] - 2026-06-02

### Added
- pkg/scanner: DeviceInfo, ScanConfig, DefaultConfig, Sanitize
- pkg/scanner: validateIPv4, Ping, ReverseDNS, GetMAC, ScanPorts
- pkg/scanner: ParseRange com suporte a CIDR e dash range
- pkg/scanner: StartScan, StopScan com goroutines e cancelamento por context
- pkg/scanner: build tags separados para Windows (SendARP) e POSIX (arp)
- pkg/exporter: ExportJSON, ExportCSV, ExportXML com sanitização de injeção CSV
- CI: go vet + go test -race no GitHub Actions
- CI: govulncheck semanal

[0.1.0]: https://github.com/mendsec/catnet-core/releases/tag/v0.1.0
