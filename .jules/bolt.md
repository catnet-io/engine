
## 2024-05-24
*   **Optimization:** Removed the overhead of spawning an external shell process (`exec.Command`) for ICMP pings on Windows by leveraging the native Win32 `IcmpSendEcho` API from `iphlpapi.dll`.
*   **Bottleneck Addressed:** High process-creation overhead on Windows. Spawning `ping.exe` for every IP scan is significantly slower than direct DLL calls, especially when scanning large subnets concurrently.
*   **Edge Case / Learning:** The ICMP reply structure on Windows returns the actual status code at offset 4 of the payload byte array. The lazy loading of `iphlpapi.dll` is efficient enough for this context as the network latency dominates, though global loading could save a few nanoseconds. The change eliminates command injection risks while providing a massive performance boost.
