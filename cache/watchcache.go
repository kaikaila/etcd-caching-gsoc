package cache

import "sync"

type WatchCache struct {
	mu        sync.RWMutex
	state     map[string]string // 当前最新 key-value 状态副本
	revision  int64             // 当前 revision，确保顺序一致
	eventSink EventSink         // 下沉接口（观察者模式）
}

func (w *WatchCache) HandlePut(key, val string, rev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if rev <= w.revision {
		return // 忽略陈旧事件
	}

	w.state[key] = val
	w.revision = rev
	if w.eventSink != nil {
		w.eventSink.HandlePut(key, val)
	}
}