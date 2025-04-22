// pkg/clientlibrary/clientlibrary_test.go
package clientlibrary

import (
	"testing"
	"yourproject/pkg/proxy/fake" // 假造一个 in-memory proxy
)

func TestClientSession_MVP(t *testing.T) {
    // 1. 构造一个 fake proxy，预先放入一些事件
    fp := fake.NewFakeProxy()
    fp.AppendEvent("key1", "val1")  // rev=1
    fp.AppendEvent("key2", "val2")  // rev=2

    cl := NewClientLibrary(fp)
    sess, err := cl.NewSession(1)   // 从 rev=1 开始
    if err != nil {
        t.Fatal(err)
    }
    defer sess.Close()

    // 2. 初始 snapshot 应该看到 rev<=1 的内容
    snap := sess.List()
    if _, ok := snap["key1"]; !ok {
        t.Errorf("expected key1 in snapshot")
    }
    if _, ok := snap["key2"]; ok {  // key2 rev=2 不应出现在 snapshot
        t.Errorf("did not expect key2 in snapshot")
    }

    // 3. Watch 应该能收到 rev>1 的事件
    events := sess.Watch()
    ev := <-events
    if ev.Key != "key2" {
        t.Errorf("expected key2 event, got %v", ev)
    }
}