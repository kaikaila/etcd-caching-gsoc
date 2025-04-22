## Workflow Overview

1. **Initialize ClientLibrary**

   - 用户在程序入口调用 `NewClientLibrary(proxy)`
   - 关联已有的 **WatchCache** 实例（或其他实现了 `proxy.WatchProxy` 的组件）

2. **Create ClientSession**

   - 调用 `ClientLibrary.NewSession(startRevision)`
   - 从 WatchCache 拿到一次全局快照 `Snapshot()`
   - 启动一个独立的 watch 流 `Watch(ctx, startRevision+1)`

3. **使用 ClientSession**

   - `List()`：返回初始快照中的所有数据（读快照）
   - `Watch()`：消费后续所有版本号大于 `startRevision` 的事件

4. **包装为 ClientCacheView（可选）**

   - 将快照 map 转为有序切片并深拷贝，提供分页等只读视图

5. **Clean Up**
   - 用户调用 `Session.Close()` 取消底层 watch 流，释放资源

---

## Relation to Existing Components

- **EventLog & Watcher**
  - 持续往底层日志中追加事件
- **WatchCache**
  - 从 `EventLog` 拉取事件，维护内存缓存
  - 实现 `Snapshot() map[string]*StoreObj` 与 `Watch(ctx, sinceRev) <-chan Event`
- **ClientLibrary Layer**
  - 构建在 WatchCache 之上，管理多用户会话隔离
  - 为每个会话提供独立的快照和事件流

---

## Deep Dive: `cacheview.go`

`cacheview.go` 的职责是为 `ClientSession.List()` 提供一个**只读、可分页、全局有序**的视图层，解耦底层缓存实现。

```go
package clientlibrary                     // 1

import (                                  // 2
    "sort"                               // 2a
    "github.com/kaikaila/etcd-caching-gsoc/pkg/proxy" // 2b
)

// cacheView 实现了 ClientCacheView 接口      // 3
type cacheView struct {
    data []*proxy.StoreObj               // 3a 存放深拷贝并排序后的对象切片
}

// NewCacheView 构造函数                  // 4
// - 接受原始 snapshot map
// - 深拷贝每个 StoreObj 避免共享内存
// - 按 GlobalRev 排序，确保遍历顺序稳定
func NewCacheView(snapshot map[string]*proxy.StoreObj) *cacheView {
    objs := make([]*proxy.StoreObj, 0, len(snapshot)) // 4a 初始化切片
    for _, o := range snapshot {                      // 4b 遍历 map
        objs = append(objs, o.DeepCopy())             // 4c 深拷贝对象
    }
    sort.Slice(objs, func(i, j int) bool {            // 4d 按全局版本排序
        return objs[i].GlobalRev < objs[j].GlobalRev
    })
    return &cacheView{data: objs}                     // 4e 返回封装好的视图
}

// List 返回整个对象切片                      // 5
func (cv *cacheView) List() []*proxy.StoreObj {
    return cv.data
}

// Page 提供分页访问                          // 6
// - page 从 1 开始，size 为每页条目数
// - 越界时返回 nil，保证安全性
func (cv *cacheView) Page(page, size int) []*proxy.StoreObj {
    if page <= 0 || size <= 0 {                      // 6a 非法参数检查
        return nil
    }
    start := (page - 1) * size                       // 6b 计算起始索引
    if start >= len(cv.data) {                       // 6c 起始越界
        return nil
    }
    end := start + size                              // 6d 计算结束索引
    if end > len(cv.data) {                          // 6e 结束越界截断
        end = len(cv.data)
    }
    return cv.data[start:end]                        // 6f 返回子切片
}
```
