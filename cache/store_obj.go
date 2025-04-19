// storeobj.go
package cache

import "go.etcd.io/etcd/api/v3/mvccpb"

// storeObj holds a single key’s value and metadata in the cache.
type storeObj struct {
    Key         string
    Value       []byte
    Revision    int64
    ModRevision int64
    EventType   mvccpb.Event_EventType
}

// DeepCopy creates a new copy of storeObj to avoid shared memory.
func (o *storeObj) DeepCopy() *storeObj {
    copy := *o
    // If there were slice or map fields, copy them here as well.
    return &copy
}