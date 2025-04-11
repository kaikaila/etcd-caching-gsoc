package cache

// 建一个以 _test.go 结尾的文件，写函数名是 TestXxx(t *testing.T)，Go 会自动识别成测试用例。
import "testing"

func TestMemoryCache_GetSet(t *testing.T) {
    c := NewMemoryCache()

    c.Set("foo", "bar")

    val, ok := c.Get("foo")
    if !ok {
        t.Fatal("Expected key 'foo' to be found")
    }
    if val != "bar" {
        t.Fatalf("Expected value 'bar', got '%s'", val)
    }
}