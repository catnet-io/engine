# catnet-core

catnet-core is the shared Go engine behind the CatNet ecosystem.

It provides reusable packages for target parsing, host discovery, scan execution, result modeling, export, profiles, and scan lifecycle events. This repository contains no GUI, no TUI, and no presentation-layer code.

## Goals
- Provide a stable scanning engine for the ecosystem.
- Keep domain logic independent from interface layers.
- Standardize results, events, errors, and export formats.
- Support CLI, TUI, and desktop frontends without code duplication.

## Scope
- Target parsing
- Discovery workflows
- Scan orchestration
- Result models
- JSON and CSV export
- Profiles and execution events

## Non-goals
- Desktop UI
- Web UI
- Terminal rendering
- Wails, React, or Bubble Tea application code

## Status
Bootstrap phase. Public contracts may evolve until the first stable release.
