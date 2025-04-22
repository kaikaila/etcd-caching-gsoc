package eventlog

import "github.com/kaikaila/etcd-caching-gsoc/pkg/api"

type Event = api.Event
const (
    EventPut    = api.EventPut
    EventDelete = api.EventDelete
)