package cache

import "sync"

// memoryCache 是 Cache 接口的内存实现，使用 map 存储数据，并用互斥锁确保并发安全
// Java: 相当于 class MemoryCache implements Cache { Map store; synchronized(lock) {...} }
type memoryCache struct {
    store map[string]string // 存储实际的 key-value 对
    mu    sync.Mutex        // 用于加锁保护 map 的并发访问
}

// NewMemoryCache 是构造函数，返回一个新的 Cache 实例
// Java: 相当于 new MemoryCache()
func NewMemoryCache() Cache {
    return &memoryCache{
        store: make(map[string]string), // 初始化 map
    }
}

// Get 实现了 Cache 接口的 Get 方法，带互斥锁保护
func (m *memoryCache) Get(key string) (string, bool) {
    m.mu.Lock()         // 加锁，保证 map 并发读写安全
    defer m.mu.Unlock() // 延迟释放锁，函数退出时自动执行

    val, ok := m.store[key] // 从 map 中读取 key
    return val, ok          // 返回 value 和是否存在
}

// Set 实现了 Cache 接口的 Set 方法，带互斥锁保护
func (m *memoryCache) Set(key string, value string) {
    m.mu.Lock()         // 加锁，防止其他 goroutine 同时修改 map
    defer m.mu.Unlock() // 延迟释放锁

    m.store[key] = value // 写入或更新 key 的值
}
