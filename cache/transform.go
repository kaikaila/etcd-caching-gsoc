package cache

import (
	"github.com/kaikaila/etcd-caching-gsoc/cache/event"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// NewStoreObjFromEvent converts an Event into a storeObj snapshot state.
// This is useful when rebuilding snapshot from event logs.
func NewStoreObjFromEvent(ev event.Event) *StoreObj {
    return &StoreObj{
        Key:            ev.Key,
        Value:          ev.Value,
        KeyRev:         ev.KeyRev,
        GlobalRev:      ev.GlobalRev,
        ModRev:         ev.GlobalRev, // or ev.ModRev if future events carry it separately
        EventType:      mvccpb.Event_EventType(ev.Type), // convert to etcd's enum type
    }
}