package cache

// SnapshotView holds an ordered slice of storeObj for paging.
type SnapshotView struct {
  Data []*storeObj
}

// Page returns items for the given page number (1-based) and page size.
func (sv *SnapshotView) Page(page, size int) []*storeObj {
  start := (page - 1) * size
  if start >= len(sv.Data) {
    return nil
  }
  end := start + size
  if end > len(sv.Data) {
    end = len(sv.Data)
  }
  return sv.Data[start:end]
}