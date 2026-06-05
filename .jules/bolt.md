## 2024-06-03 - [Concurrent Port Scanning]
**Learning:** Sequential network I/O operations (like `net.DialTimeout`) cause severe performance bottlenecks when waiting for timeouts. Cumulative timeouts block the scan entirely.
**Action:** Always introduce bounded concurrency (e.g., using a channel semaphore) and thread-safe data aggregation (`sync.Mutex`) when handling multiple network operations that may drop packets or time out. Ensure deterministic output by sorting the results before returning.

## 2024-05-24
*   **Optimization:** Removed the overhead of spawning an external shell process (`exec.Command`) for ICMP pings on Windows by leveraging the native Win32 `IcmpSendEcho` API from `iphlpapi.dll`.
*   **Bottleneck Addressed:** High process-creation overhead on Windows. Spawning `ping.exe` for every IP scan is significantly slower than direct DLL calls, especially when scanning large subnets concurrently.
*   **Edge Case / Learning:** The ICMP reply structure on Windows returns the actual status code at offset 4 of the payload byte array. The lazy loading of `iphlpapi.dll` is efficient enough for this context as the network latency dominates, though global loading could save a few nanoseconds. The change eliminates command injection risks while providing a massive performance boost.

2023-10-27: Found infinite loop and OOM vulnerability when parsing CIDRs containing 0.0.0.0 and exceeding limits. Prevented infinite looping by making increment function return overflow indication, and prevented OOMs by capping CIDR generation limit to 65536 IPs.
