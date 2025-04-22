// pkg/clientlibrary/clientlibrary.go
package clientlibrary

import (
	"errors"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/proxy"
)

// ClientLibrary 实现 api.ClientLibrary
type clientLibrary struct {
    cache proxy.WatchCacheInterface
    log   eventlog.EventLog
}

// NewClientLibrary 构造
func NewClientLibrary(cache proxy.WatchCacheInterface, log eventlog.EventLog) api.ClientLibrary {
    return &clientLibrary{cache: cache, log: log}
}

// NewSession 创建一个新会话
func (cl *clientLibrary) NewSession(clientID string) (api.ClientSession, error) {
    if cl.log == nil {
        return nil, errors.New("event log is nil")
    }
    rv := cl.log.LatestRevision()
    sess := newSession(cl.cache, cl.log, rv)
    return sess, nil
}

// BroadcastUpdate broadcasts an event into the EventLog.
func (cl *clientLibrary) BroadcastUpdate(ev api.Event) {
    // Append into the underlying EventLog.
    _ = cl.log.Append(ev)
}