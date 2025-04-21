/*
Package eventlog provides modular backends for persisting and replaying etcd change events.

This package defines:
- Event: the internal representation of a cache-level change event, including key, value, revision info.
- EventLog: the core interface for event sinks that store or process historical events.
- MemoryEventLog: an in-memory circular buffer implementation of EventLog.
- EtcdLog and WALLog: planned extensions for persistent event log storage.

This subpackage enables features such as replay, audit, snapshot recovery, and diff-based views
by providing a standardized event stream across caching layers.
*/
package eventlog