package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kaikaila/etcd-caching-gsoc/watcher"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func testWatcherWithCtlOutput() {
	// Connect to etcd
	fmt.Println(">>> testWatcherWithCtlOutput started")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("❌ 连接 etcd 失败:", err)
	}
	defer cli.Close()

	// Use WatchKey to watch "/foo" (wrapped in a goroutine so it doesn't block)
	go watcher.WatchKey(cli, "/foo", func(key, val string) {
		fmt.Printf("✅ Received etcd event: key=%s, value=%s\n", key, val)
	}, func(key string) {
		fmt.Printf("✅ Deleted etcd event: key=%s\n", key)
	})

	// Automatically write to etcd to trigger watcher callback
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			_, err := cli.Put(context.Background(), "/foo", fmt.Sprintf("val-%d", i))
			if err != nil {
				log.Println("❌ Automatically writing to etcd failed:", err)
			} else {
				log.Printf("✍️ Automatically writing to etcd: /foo=val-%d\n", i)
			}
		}
	}()
	// Block the main thread (previously waited for etcdctl commands)
	// select {}
	time.Sleep(10 * time.Second)
	fmt.Println("🛑 TestWatcherWithCtlOutput finished. Proceeding to next test...")
}