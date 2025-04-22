// pkg/clientlibrary/session.go
package clientlibrary

import (
	"context"
	"fmt"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/proxy"
)

// session 实现 api.ClientSession
type session struct {
    cache           proxy.WatchCache
    log             eventlog.EventLog
    startRevision   int64
    initialSnapshot map[string]*proxy.StoreObj
    eventsCh        <-chan eventlog.Event
    cancelWatch     context.CancelFunc
}

func newSession(cache proxy.WatchCache, log eventlog.EventLog, rv int64) api.ClientSession {
    // 1. 获取初始快照（Snapshot）
    snap := cache.Snapshot()
    // 2. 启动 watch
    ctx, cancel := context.WithCancel(context.Background())
    events, _ := log.Watch(ctx, rv+1)
    return &session{
        cache:           cache,
        log:             log,
        startRevision:   rv,
        initialSnapshot: snap,
        eventsCh:        events,
        cancelWatch:     cancel,
    }
}

// List 返回 snapshot
func (s *session) List() map[string]*proxy.StoreObj {
    return s.initialSnapshot
}

// Watch 返回一个事件 channel
func (s *session) Watch() <-chan eventlog.Event {
    return s.eventsCh
}

// Close 取消 watch
func (s *session) Close() {
    s.cancelWatch()
}

// Start starts the session; placeholder if any initialization is needed.
func (s *session) Start() error {
    // No additional startup logic for now.
    return nil
}

// Stop stops the session and releases resources.
func (s *session) Stop() error {
    s.Close()
    return nil
}

// ID returns a unique identifier for this session.
func (s *session) ID() string {
    return fmt.Sprintf("%p", s)
}

// CacheView returns a read-only view of the initial snapshot.
func (s *session) CacheView() api.ClientCacheView {
    return NewCacheView(s.initialSnapshot)
}