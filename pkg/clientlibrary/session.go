// session 实现 api.ClientSession。
// 提供 snapshot 读取、事件订阅流以及生命周期管理（Start/Stop）。
package clientlibrary

import (
	"context"
	"fmt"
	"strings"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/proxy"
)

// session 实现 api.ClientSession
type session struct {
    cache           proxy.WatchCacheInterface
    log             eventlog.EventLog
    startRevision   int64
    initialSnapshot []api.KV
    eventsCh        <-chan eventlog.Event
    cancelWatch     context.CancelFunc
}

func newSession(cache proxy.WatchCacheInterface, log eventlog.EventLog, rv int64) api.ClientSession {
    // 1. 获取初始快照（Snapshot）
    snaps := cache.Snapshot()
    ssdata,_ := snaps.List("")
    // 2. 启动 watch
    ctx, cancel := context.WithCancel(context.Background())
    events, _ := log.Watch(ctx, rv+1)
    return &session{
        cache:           cache,
        log:             log,
        startRevision:   rv,
        initialSnapshot: ssdata,
        eventsCh:        events,
        cancelWatch:     cancel,
    }
}

// List 返回 snapshot
func (s *session) List() []api.KV {
    return s.initialSnapshot
}

// WatchSingle subscribes to changes on a single key
func (s *session) Watch(key string, fromRev int64) (<-chan api.Event, error) {
	ctx, _ := context.WithCancel(context.Background())
	events, err := s.log.Watch(ctx, fromRev)
	if err != nil {
		return nil, err
	}
	out := make(chan api.Event)
	go func() {
		defer close(out)
		for ev := range events {
			if ev.Key == key {
				out <- ev
			}
		}
	}()
	// 暂不追踪 substream cancel，可加字段做管理
	return out, nil
}

// WatchPrefix subscribes to changes on a key prefix
func (s *session) WatchPrefix(prefix string, fromRev int64) (<-chan api.Event, error) {
	ctx, _ := context.WithCancel(context.Background())
	events, err := s.log.Watch(ctx, fromRev)
	if err != nil {
		return nil, err
	}
	out := make(chan api.Event)
	go func() {
		defer close(out)
		for ev := range events {
			if strings.HasPrefix(ev.Key, prefix) {
				out <- ev
			}
		}
	}()
	return out, nil
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

// CacheView returns a read-only snapshot view of the cache.
func (s *session) CacheView() api.SnapshotView {
    return s.cache.Snapshot()
}