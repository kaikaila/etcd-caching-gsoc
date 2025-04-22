package proxy

import (
	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// NewStoreObjFromEvent converts an Event into a storeObj snapshot state.
// This is useful when rebuilding snapshot from event logs.
func NewStoreObjFromEvent(ev eventlog.Event) *StoreObj {
    return &StoreObj{
        Key:            ev.Key,
        Value:          ev.Value,
        Revision:      ev.Revision,
        ModRev:         ev.Revision, // or ev.ModRev if future events carry it separately
        EventType:      mvccpb.Event_EventType(ev.Type), // convert to etcd's enum type
    }
}