# Performance Decision: HandlePut Parameter Type

We benchmarked two implementations of WatchCache.HandlePut:

| Version                        | ns/op | B/op | allocs/op |
| ------------------------------ | ----- | ---- | --------- |
| HandlePut(key,string,rev)      | 107.8 | 1088 | 2         |
| HandlePutBytes(key,[]byte,rev) | 6.986 | 0    | 0         |

**Analysis**:

- Passing raw `[]byte` avoids two conversions (string→[]byte and string allocation).
- Eliminates memory allocations and GC overhead.
- Per-op latency reduced by ~15×.

**Decision**:

- **Use `HandlePutBytes(key string, val []byte, rev int64)`** in the high-throughput path.
- Keep `HandlePut(key, val string, rev)` around as a backward-compatible wrapper if needed.
