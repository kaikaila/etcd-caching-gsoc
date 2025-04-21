/*
Package cache implements a generic, pluggable caching proxy layer for etcd.

It is designed to provide enhanced watch scalability, low-latency reads, and revision-aware snapshotting,
while remaining decoupled from Kubernetes-specific assumptions.

Core components:

- memoryCache: an in-memory key-value store implementing the Cache interface.
- WatchCache: a higher-level cache with snapshot and compaction support.
- EventSink: an interface for observing change events (used for replay, metrics, or replication).
- StoreObj and SnapshotView: internal data models for consistent snapshotting and versioning.

This package serves as the foundation of a generic watch cache proxy, enabling downstream systems
to build client libraries and adapters on top of it.
*/
package proxy