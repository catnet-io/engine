# API Stability Policy

`catnet-core` is currently in a pre-`v1.0.0` stage. However, as the shared engine for the CatNet ecosystem, we are establishing clear expectations around API stability to allow downstream consumers (CLI, TUI, GUI) to build upon it reliably.

## Stable Packages

The following packages form the core stable surface of the ecosystem. Breaking changes to these packages will be avoided whenever possible, and if required, they will be communicated via explicit version bumps and detailed changelogs.

- **`pkg/engine`**: The primary orchestrator. `StartScan` and `ScanConfig` are the standard way to run scans.
- **`pkg/results`**: Defines the `DeviceInfo` model and other schemas. This is the canonical schema used across the ecosystem.
- **`pkg/targets`**: Target parsing utilities like `ParseRange`.
- **`pkg/discovery`**: Discovery tools for Liveness, DNS, and MAC resolution.
- **`pkg/ports`**: Port scanning utilities.
- **`pkg/exporter`**: Exporting capabilities (JSON, CSV, XML). JSON is the canonical format.

## Experimental / Internal Packages

Any package under `internal/` is strictly internal and offers **zero stability guarantees**. It may be changed or removed at any time without warning. Consumers must not import `internal/` packages.

## Deprecated Packages

- **`pkg/scanner`**: **DEPRECATED**. This package previously housed all scanning logic. It is now maintained solely as a temporary compatibility wrapper to ease the migration to the new package structure. It will be removed in a future release. Downstream consumers should migrate to `pkg/engine`, `pkg/discovery`, `pkg/ports`, `pkg/targets`, and `pkg/results`.

## Execution Model: Context

The preferred and official execution model uses `context.Context` to manage timeouts, deadlines, and cancellation. Legacy cancellation methods such as `StopScan()` are deprecated and should not be used in new integrations.
