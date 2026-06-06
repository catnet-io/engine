## 2024-06-03 - [Concurrent Port Scanning]
**Learning:** Sequential network I/O operations (like `net.DialTimeout`) cause severe performance bottlenecks when waiting for timeouts. Cumulative timeouts block the scan entirely.
**Action:** Always introduce bounded concurrency (e.g., using a channel semaphore) and thread-safe data aggregation (`sync.Mutex`) when handling multiple network operations that may drop packets or time out. Ensure deterministic output by sorting the results before returning.

## 2026-06-04 - Optimize IP Range Allocations
**Learning:** Expanding large IP ranges (like CIDR /16 or full class B dash ranges) causes massive memory allocation pressure in Go because the slice backing array constantly resizes and reallocates, and creating a new `net.IP` in every iteration creates garbage collection overhead.
**Action:** When parsing target ranges, always pre-calculate the total number of IPs (using CIDR mask size or end-start+1) and initialize the slice with `make([]string, 0, capacity)`. Additionally, reuse a single `net.IP` buffer across loop iterations when generating sequential IPs.

## 2024-05-24
*   **Optimization:** Removed the overhead of spawning an external shell process (`exec.Command`) for ICMP pings on Windows by leveraging the native Win32 `IcmpSendEcho` API from `iphlpapi.dll`.
*   **Bottleneck Addressed:** High process-creation overhead on Windows. Spawning `ping.exe` for every IP scan is significantly slower than direct DLL calls, especially when scanning large subnets concurrently.
*   **Edge Case / Learning:** The ICMP reply structure on Windows returns the actual status code at offset 4 of the payload byte array. The lazy loading of `iphlpapi.dll` is efficient enough for this context as the network latency dominates, though global loading could save a few nanoseconds. The change eliminates command injection risks while providing a massive performance boost.

2023-10-27: Found infinite loop and OOM vulnerability when parsing CIDRs containing 0.0.0.0 and exceeding limits. Prevented infinite looping by making increment function return overflow indication, and prevented OOMs by capping CIDR generation limit to 65536 IPs.

## 2026-06-05 - [MAC Address Resolution Fast Path]
**Learning:** Spawning an external process (`fork`/`exec` via `exec.Command`) to read the local ARP table (`arp -an`) creates massive performance penalties and OS scheduling contention when scaled across many concurrent scanning threads. For example, spawning a process for 100 IPs can take ~300ms, whereas reading a file takes ~2.5ms.
**Action:** When gathering MAC addresses on POSIX systems (specifically Linux), implement a fast path that reads directly from `/proc/net/arp`. Only fallback to `exec.Command` if the file doesn't exist or isn't mounted (e.g. macOS/BSD).
## 2023-10-27 - Atomic Index Loop Over Buffered Channel
**Learning:** For distributing pre-known arrays of work (like IP slices up to 65536 items) across worker threads, creating a buffered channel and pushing all items into it introduces massive setup overhead and memory allocation (O(N)).
**Action:** Use a pre-allocated array and a shared `int32` atomic index (`atomic.AddInt32(&index, 1)`) inside the worker thread loop to read the slice dynamically without synchronization channels. This gives ~3x speedup.
