## 2026-06-25 - Prevent goroutine leaks in context-cancelled channel sends
**Learning:** Worker goroutines can leak when sending on an unbuffered or partially read channel if the consumer prematurely stops reading due to context cancellation.
**Action:** Wrap channel send operations in a `select` statement with `ctx.Done()` to prevent goroutines from blocking indefinitely when the consumer aborts.
## 2026-06-25 - Prevent goroutine leaks in context-cancelled channel sends
**Learning:** Worker goroutines can leak when sending on an unbuffered or partially read channel if the consumer prematurely stops reading due to context cancellation.
**Action:** Wrap channel send operations in a `select` statement with `ctx.Done()` to prevent goroutines from blocking indefinitely when the consumer aborts.
