package event

// EventLog defines the interface for an append-only event history.
// 不同的实现可以是内存环形缓冲区、WAL 文件、etcd API 等。目前先做了memory版本，WAL,etcd历史接口的版本以后再说
type EventLog interface {
    Append(ev Event) error                      // 追加一条事件
    Replay(fromRev int64) ([]Event, error)      // 从某个 global revision 开始 replay
    LatestRevision() int64                      // 返回当前已记录的最大 global revision
}
