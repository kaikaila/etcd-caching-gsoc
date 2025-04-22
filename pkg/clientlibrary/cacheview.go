package clientlibrary

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/proxy"
)

// cacheView implements the ClientCacheView interface.
// It holds a read-only, ordered snapshot of StoreObj items.
type cacheView struct {
	data     []*proxy.StoreObj
	revision int64
}

// NewCacheView creates a new cache view from the given snapshot map.
// It deep-copies each object and sorts them by Revision for deterministic ordering.
func NewCacheView(snapshot map[string]*proxy.StoreObj) *cacheView {
	var maxRev int64
	for _, o := range snapshot {
		if o.Revision > maxRev {
			maxRev = o.Revision
		}
	}
	objs := make([]*proxy.StoreObj, 0, len(snapshot))
	for _, o := range snapshot {
		objs = append(objs, o.DeepCopy())
	}
	sort.Slice(objs, func(i, j int) bool {
		return objs[i].Revision < objs[j].Revision
	})
	return &cacheView{data: objs, revision: maxRev}
}

// Get returns the KV for a single key if present.
func (cv *cacheView) Get(key string) (api.KV, bool) {
	for _, o := range cv.data {
		if o.Key == key {
			return api.KV{Key: o.Key, Value: o.Value, Revision: o.Revision}, true
		}
	}
	return api.KV{}, false
}

// List returns all KVs whose key has the given prefix.
func (cv *cacheView) List(prefix string) ([]api.KV, error) {
	result := make([]api.KV, 0)
	for _, o := range cv.data {
		if strings.HasPrefix(o.Key, prefix) {
			result = append(result, api.KV{Key: o.Key, Value: o.Value, Revision: o.Revision})
		}
	}
	return result, nil
}

// Watch returns an error because cache view does not support watching.
func (cv *cacheView) Watch(key string, fromRevision int64) (<-chan api.Event, error) {
	return nil, fmt.Errorf("Watch not supported in cache view")
}

// ResourceVersion returns the highest Revision in the snapshot.
func (cv *cacheView) ResourceVersion() int64 {
	return cv.revision
}

// Page returns a slice of items for the given 1-based page number and size.
// If the page or size is out of range, it returns an empty slice.
func (cv *cacheView) Page(page, size int) []*proxy.StoreObj {
	if page <= 0 || size <= 0 {
		return nil
	}
	start := (page - 1) * size
	if start >= len(cv.data) {
		return nil
	}
	end := start + size
	if end > len(cv.data) {
		end = len(cv.data)
	}
	return cv.data[start:end]
}
