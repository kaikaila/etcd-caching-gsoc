# WatchCache 架构

## Watcher

### 来源：etcd Client Watch Stream

### 输出：Event

### 方法：

#### - Start()

#### - Stop()

#### - Chan() <-chan Event

## EventLog

### 接收自：Watcher

### 输出至：WatchCache

### 作用：按顺序记录事件（Append、Replay）

### 方法：

#### - Append(event Event)

#### - Replay(sinceRevision int64) []Event

#### - Compact(revision int64)

## WatchCache

### 接收自：EventLog

### 输出至：外部用户

### 作用：存储当前快照，支持读写缓存

### 方法：

#### - Get(key string) \*storeObj

#### - List() []\*storeObj

#### - Snapshot() map[string]\*storeObj

#### - Compact(revision int64)

## 外部用户

### 来源：K8s / Cilium / Calico 等

### 调用接口：

#### - WatchCache.Get()

#### - WatchCache.List()

#### - WatchCache.Snapshot()
