package main

// Java: 相当于 public class Main，main 函数的入口类必须是这个名字

import (

	// Java: 类似 java.time.Duration，用来控制 sleep、超时
	"fmt"
	"log"
	"time"

	"github.com/kaikaila/etcd-caching-gsoc/cache"
	"github.com/kaikaila/etcd-caching-gsoc/watcher"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main_temp() {
	// Java: 相当于 public static void main(String[] args) {}
	
	// 创建一个内存缓存实例
    c := cache.NewMemoryCache()


	// 创建一个 etcd 客户端连接（配置 IP、端口、超时时间）
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, 
		// Java: 等价于 new EtcdClient.Builder().setEndpoints("localhost:2379")

		DialTimeout: 5 * time.Second,            
		// Java: 等价于 .setDialTimeout(Duration.ofSeconds(5))
	})
	if err != nil {
		log.Fatal("connection to etcd failed",err)
	}

	defer cli.Close()
	

	// 启动 watcher，监听 /foo；用 goroutine 启动后台任务，类似 Java 的 new Thread(() -> ...)
	watcher.WatchKey(cli, "/foo", func(key, val string) {
		fmt.Printf("roger etcd events：%s = %s\n", key, val)
		c.Set(key,val)
		// 将来你可以在这里加 cache.Set(key, val)
	})
	
	
	go func() {
		for {
			v, ok := c.Get("/foo")
			if ok {
				fmt.Println(" read from cache /foo = ", v)
			} else {
				fmt.Println("cache is currently empty")
			}
			time.Sleep(2*time.Second)
		}
	} ()
	
	// 阻塞主线程，防止退出。select {} 是 Go 的“永久等待”，类似 Java 的 while (true) {}
	select {}
}