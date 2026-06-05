# Compatibility Policy

As the shared engine of the CatNet ecosystem, `catnet-core` strictly defines how downstream repositories must track versions and handle breaking changes. This reduces coordination costs and ensures stability across CLI, TUI, and GUI consumers.

## Versioning Strategy

`catnet-core` uses [Semantic Versioning (SemVer)](https://semver.org/).

Given that the project is currently in a pre-`v1.0.0` state (`0.x.x`), the SemVer specification dictates that **anything may change at any time** and the public API should not be considered stable. 

However, to maintain ecosystem sanity, we apply the following internal rules during the `0.x.x` phase:
- **Patch (`0.0.x`)**: Bug fixes, performance improvements, and non-breaking internal refactors.
- **Minor (`0.x.0`)**: New features, new fields in `ScanReport`, or breaking changes to public contracts (like `ScanConfig` or `ScanEvent`).

## Ecosystem Consumer Requirements

All official consumers (`catnet`, `catnet-tui`, `catnet-scanner`) **MUST**:
1. **Pin Exact Minor Versions:** Consumers must pin their `go.mod` to a specific minor version (e.g., `v0.1.x`) and avoid using `@latest` blindly.
2. **Handle Unknown JSON Fields:** Consumers reading the JSON export must not crash if new fields appear in `ScanReport` (Forward Compatibility).
3. **Respect `SchemaVersion`:** If the major version of `SchemaVersion` changes, the consumer must either abort the import or update their internal models to handle the new schema.

## Breaking Changes Communication

Before any breaking change is merged into `develop` (such as modifying `ScanEventType` or `ScanConfig` fields):
1. An issue must be created and tagged with `contracts`.
2. The PR must detail the migration path for CLI, TUI, and GUI.
3. The `CHANGELOG.md` must clearly highlight the breaking change under a `### Changed` header.
