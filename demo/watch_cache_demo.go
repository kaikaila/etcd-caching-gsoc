package main

import (
	"fmt"
	"time"

	"github.com/kaikaila/etcd-caching-gsoc/cache"
	"github.com/kaikaila/etcd-caching-gsoc/watcher"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func testWatcherWithWatchCache() {
	wc_temp := cache.NewMemoryCache()
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	defer cli.Close()
	wc := wc_temp.(cache.CacheWithSink)
	watcher.WatchKey(cli, "/foo", wc.HandlePut, wc.HandleDelete)

	go func() {
		for {
			v, ok := wc.Get("/foo")
			if ok {
				fmt.Println(" watchcache → /foo = ", v)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	select {}
}