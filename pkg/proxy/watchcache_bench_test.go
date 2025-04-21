package proxy

import (
	"testing"
)

// BenchmarkHandlePut_String tests the original string-based HandlePut.
func BenchmarkHandlePut_String(b *testing.B) {
    wc := NewWatchCache(nil)
    key := "foo"
    // 模拟 1KB 的文本 payload
    val := string(make([]byte, 1024))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wc.HandlePut(key, val, int64(i))
    }
}

// BenchmarkHandlePut_Bytes tests the new []byte-based HandlePutBytes.
func BenchmarkHandlePut_Bytes(b *testing.B) {
    wc := NewWatchCache(nil)
    key := "foo"
    // 模拟 1KB 的二进制 payload
    val := make([]byte, 1024)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wc.HandlePutBytes(key, val, int64(i))
    }
}