## 2024-06-03 - [Concurrent Port Scanning]
**Learning:** Sequential network I/O operations (like `net.DialTimeout`) cause severe performance bottlenecks when waiting for timeouts. Cumulative timeouts block the scan entirely.
**Action:** Always introduce bounded concurrency (e.g., using a channel semaphore) and thread-safe data aggregation (`sync.Mutex`) when handling multiple network operations that may drop packets or time out. Ensure deterministic output by sorting the results before returning.

2023-10-27: Found infinite loop and OOM vulnerability when parsing CIDRs containing 0.0.0.0 and exceeding limits. Prevented infinite looping by making increment function return overflow indication, and prevented OOMs by capping CIDR generation limit to 65536 IPs.
