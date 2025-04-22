# ClientLibrary Architecture

## ClientLibrary

### Source: WatchCache (proxy.WatchProxy)

### Purpose: Manage multiple client sessions

### Methods:

#### - NewSession(resourceVersion int64) ClientSession

## ClientSession

### Source: ClientLibrary.NewSession

### Purpose: Provide isolated view per client

### Attributes:

#### - startRevision int64

#### - initialSnapshot map[string]\*proxy.StoreObj

#### - eventsCh <-chan proxy.Event

### Methods:

#### - List() map[string]\*proxy.StoreObj

#### - Watch() <-chan proxy.Event

#### - Close()

## ClientCacheView

### Source: Initial snapshot in Session

### Purpose: Read‑only, sorted deep‑copy of snapshot data

### Methods:

#### - List() []\*proxy.StoreObj

#### - Page(page, size int) []\*proxy.StoreObj

## Relation to Existing Components

### WatchCache

#### - Provides Snapshot() map[string]\*proxy.StoreObj

#### - Provides Watch(ctx, sinceRev int64) <-chan proxy.Event

### EventLog & Watcher

#### - Emit events into WatchCache

### ClientLibrary

#### - Builds on top of WatchCache for multi-client isolation

## External Client

### Usage Flow:

#### 1. Call ClientLibrary.NewSession(startRevision)

#### 2. Use Session.List() to get initial view

#### 3. Use Session.Watch() to receive incremental events

#### 4. Optionally wrap snapshot with ClientCacheView for paging
