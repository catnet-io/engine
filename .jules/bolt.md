## 2024-06-03 - [Concurrent Port Scanning]
**Learning:** Sequential network I/O operations (like `net.DialTimeout`) cause severe performance bottlenecks when waiting for timeouts. Cumulative timeouts block the scan entirely.
**Action:** Always introduce bounded concurrency (e.g., using a channel semaphore) and thread-safe data aggregation (`sync.Mutex`) when handling multiple network operations that may drop packets or time out. Ensure deterministic output by sorting the results before returning.
## 2026-06-04 - Optimize IP Range Allocations
**Learning:** Expanding large IP ranges (like CIDR /16 or full class B dash ranges) causes massive memory allocation pressure in Go because the slice backing array constantly resizes and reallocates, and creating a new `net.IP` in every iteration creates garbage collection overhead.
**Action:** When parsing target ranges, always pre-calculate the total number of IPs (using CIDR mask size or end-start+1) and initialize the slice with `make([]string, 0, capacity)`. Additionally, reuse a single `net.IP` buffer across loop iterations when generating sequential IPs.
