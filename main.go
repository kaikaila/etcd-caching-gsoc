package main

// Java: 相当于 public class Main，main 函数的入口类必须是这个名字

import (

	// Java: 类似 java.time.Duration，用来控制 sleep、超时
	"fmt"

	"github.com/kaikaila/etcd-caching-gsoc/cache"
)

func main() {
	// Java: 相当于 public static void main(String[] args) {}
	
	// 创建一个 etcd 客户端连接（配置 IP、端口、超时时间）
	// cli, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{"localhost:2379"}, 
	// 	// Java: 等价于 new EtcdClient.Builder().setEndpoints("localhost:2379")

	// 	DialTimeout: 5 * time.Second,            
	// 	// Java: 等价于 .setDialTimeout(Duration.ofSeconds(5))
	// })
	// if err != nil {
	// 	// Java: 相当于 if (err != null) throw new RuntimeException(err)
	// 	log.Fatal(err)
	// }

	// defer cli.Close()
	// // Java: 相当于 try-finally 块，在函数结束时自动执行 cli.close()

	// // 启动 watcher，监听 /foo；用 goroutine 启动后台任务，类似 Java 的 new Thread(() -> ...)
	// watcher.WatchKey(cli, "/foo", func(key, val string) {
	// 	fmt.Printf("🔔 收到 etcd 事件：%s = %s\n", key, val)
	
	// 	// 将来你可以在这里加 cache.Set(key, val)
	// })
	// // 阻塞主线程，防止退出。select {} 是 Go 的“永久等待”，类似 Java 的 while (true) {}
	// select {}
	
	// cache starts here
	// 创建一个内存缓存实例
    c := cache.NewMemoryCache()

    // 设置一个 key
    c.Set("foo", "bar")

    // 获取 key
    val, ok := c.Get("foo")
    if ok {
        fmt.Println("Got value:", val) // 预期输出: Got value: bar
    } else {
        fmt.Println("Key not found")
    }
}