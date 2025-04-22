package proxy

import "go.etcd.io/etcd/api/v3/mvccpb"

// StoreObj holds a single key’s value and metadata in the cache.
type StoreObj struct {
    Key            string
    Value          []byte
    Revision      int64 // global revision: indicates the change's order among all operations
    ModRev         int64
    EventType      mvccpb.Event_EventType
}

// DeepCopy creates a new copy of StoreObj to avoid shared memory.
func (o *StoreObj) DeepCopy() *StoreObj {
    copy := *o
    return &copy
}