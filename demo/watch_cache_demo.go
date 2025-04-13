package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kaikaila/etcd-caching-gsoc/cache"
	"github.com/kaikaila/etcd-caching-gsoc/watcher"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func testWatcherWithWatchCache() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	wc := cache.NewWatchCache(nil)

	watcher.WatchKeyWithRevision(cli, "/foo", wc.HandlePut, wc.HandleDelete)

	// 持续读取 WatchCache 并打印状态与 revision
	go func() {
		for {
			val, ok := wc.Get("/foo")
			if ok {
				fmt.Printf("read from watchcache /foo = %s [rev = %d]\n", val, wc.Revision())
			} else {
				fmt.Printf("read from watchcache /foo = (not found) [rev = %d]\n", wc.Revision())
			}
			time.Sleep(2 * time.Second)
		}
	}()

	select {} // 阻塞主线程
}