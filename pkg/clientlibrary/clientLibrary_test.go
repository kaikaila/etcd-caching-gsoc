// pkg/clientlibrary/clientlibrary_test.go
package clientlibrary

import (
	"testing"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/proxy" // 假造一个 in-memory proxy
)

func TestClientSession_MVP(t *testing.T) {
    // 1. 构造一个 fake proxy，预先放入一些事件
    log := eventlog.NewMemoryEventLog(5)
	fp := proxy.NewWatchCacheWithLog(nil, log)
    
    ev1 := api.Event{
        Type:     api.EventPut,
        Key:      "key1",
        Value:    []byte("Alice"),
        Revision: 1,   // 全局单调版本号
        ModRev:   100, // etcd 原生 ModRevision
    }

    // 第二个事件：一次 DELETE 操作
    ev2 := api.Event{
        Type:     api.EventDelete,
        Key:      "key2",
        Value:    nil, // DELETE 时没有新值
        Revision: 2,
        ModRev:   101,
    }

    fp.AddEvent(ev1)  
    fp.AddEvent(ev2) 

    el := eventlog.NewMemoryEventLog(10)
    cl := NewClientLibrary(fp, el)
    sess, err := cl.NewSession("test-client")
    if err != nil {
        t.Fatal(err)
    }
    defer sess.Stop()
    view := sess.CacheView()
    // 2. 初始 snapshot 应该看到 rev<=1 的内容
    snap := sess.List()
    if _, ok := snap["key1"]; !ok {
        t.Errorf("expected key1 in snapshot")
    }
    if _, ok := snap["key2"]; ok {  // key2 rev=2 不应出现在 snapshot
        t.Errorf("did not expect key2 in snapshot")
    }

    // 3. Watch 应该能收到 rev>1 的事件
    events := view.Watch()
    ev := <-events
    if ev.Key != "key2" {
        t.Errorf("expected key2 event, got %v", ev)
    }
}