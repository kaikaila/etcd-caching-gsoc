package proxy

import (
	"fmt"
	"strings"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
)

// SnapshotView holds an ordered slice of storeObj for paging.
type CacheSnapshotView struct {
  data []*StoreObj
  index map[string]*StoreObj
  revision int64
}

// the constructor is in watchcache.go

// Page returns items for the given page number (1-based) and page size.
// It may return a non-nil error in future if, for example, the snapshot has expired,
// the page/size parameters are invalid (e.g., negative or zero), or the snapshot is compacted.
// Currently, it always returns a nil error for simplicity.
func (sv *CacheSnapshotView) Page(page, size int) ([]api.KV, error) {
    start := (page - 1) * size
    if start >= len(sv.data) {
        return nil, fmt.Errorf("page start index %d out of bounds (total %d items)", start, len(sv.data))
    }
    end := start + size
    if end > len(sv.data) {
        end = len(sv.data)
    }
    var result []api.KV
    for _, obj := range sv.data[start:end] {
        valCopy := append([]byte(nil), obj.Value...)
        result = append(result, api.KV{
            Key:      obj.Key,
            Value:    valCopy,
            Revision: obj.Revision,
        })
    }
    //
    return result, nil
}

// Get returns the StoreObj for a single key if present.
func (sv *CacheSnapshotView) Get(key string) (api.KV, bool) {
    if obj, ok := sv.index[key]; ok {
        valCopy := append([]byte(nil), obj.Value...)
        return api.KV{Key: obj.Key, Value: valCopy, Revision: obj.Revision}, true
    }
    return api.KV{}, false
}

// List returns all StoreObj items whose key has the given prefix.
// May return error in future if snapshot is expired, compacted, or invalid.
func (sv *CacheSnapshotView) List(prefix string) ([]api.KV, error) {
    var result []api.KV
    for _, obj := range sv.data {
        if strings.HasPrefix(obj.Key, prefix) {
            valCopy := append([]byte(nil), obj.Value...)
            result = append(result, api.KV{
                Key:      obj.Key,
                Value:    valCopy,
                Revision: obj.Revision,
            })
        }
    }
    // Currently always returns nil error, but structured to support future error cases
    // (e.g., snapshot expiration, compaction, invalid parameters).
    return result, nil
}

// Revision returns the highest Revision in this view.
func (sv *CacheSnapshotView) Revision() int64 {
    return sv.revision
}