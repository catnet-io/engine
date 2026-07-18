# Roadmap

This document outlines the near-term priorities and non-goals for `catnet-io/engine`. As the foundational engine for the CatNet ecosystem, we prioritize extreme reliability and backwards compatibility over rapid feature expansion.

## Milestones

### ✅ M1: Bootstrapping & Core Primitives
- [x] Basic discovery flow (ICMP, ARP, Reverse DNS).
- [x] Concurrent execution pipeline.
- [x] Target parsing logic.

### ✅ M2: Contracts & Schema Definition
- [x] Structured error taxonomy.
- [x] Versioned JSON schema for scan results.
- [x] Typed scan lifecycle events for UI synchronization.

### ✅ M3: Reliability, Testing, and Performance
- [x] Deterministic cancellation tests (zero goroutine leaks).
- [x] Memory and allocation benchmarks for core hot paths.
- [x] Integration harness with canonical fixtures.

### ✅ M4: Ecosystem Readiness
- [x] Define SemVer compatibility policy.
- [x] Enforce hard concurrency limits.
- [x] Provide ecosystem integration examples.

### 🟡 M5: Consumers & Asynchronous Architecture (Current)
- [ ] **Genuine Event Asynchrony:** Decouple callback latency from engine throughput using an internal async event dispatcher in `pkg/engine`.
- [ ] **Context Cancellation Validation:** Write tests verifying zero goroutine/socket leaks on premature context cancellation.
- [ ] **Consumer Integration:** Deprecate and remove deprecated core modules (`pkg/store` and `pkg/diff`), coordinate with CLI (`catnet-io/catnet`), GUI (`catnet-io/app`), and TUI (`catnet-io/tui`).

### 🔜 M6: Core Network Stack Optimization
- **Privileged Raw Sockets:** Refactor the discovery engine to use raw sockets on Windows/Linux to bypass the OS network stack limit for massive concurrent ARP/ICMP bursts.
- **Port Scanner Parallelization:** Re-write the internal port scanner to be fully asynchronous per-port instead of per-host, avoiding TCP connection timeouts holding back the worker pool.

## Non-Goals (Deferred Indefinitely)

To prevent scope creep, `catnet-io/engine` explicitly will **NOT** implement:
- **Cloud integration natively (AWS/GCP APIs):** The core focuses exclusively on network-level IP packets. Cloud metadata enrichment belongs in the UI layers or a wrapper plugin.
- **Vulnerability exploitation:** We are a discovery and enumeration tool, not a Metasploit alternative.
- **GUI code:** No native bindings (Wails/Electron) will live in this repository. All interfaces must consume the core as a pure Go module.
