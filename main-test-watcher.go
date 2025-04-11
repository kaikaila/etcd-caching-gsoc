package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kaikaila/etcd-caching-gsoc/watcher"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main_temp() {
	// 连接本地 etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("❌ 连接 etcd 失败:", err)
	}
	defer cli.Close()

	// 用 WatchKey 监听 "/foo"
	watcher.WatchKey(cli, "/foo", func(key, val string) {
		fmt.Printf("✅ 收到 etcd 变更事件：key=%s, value=%s\n", key, val)
	})

	// 主线程挂起，等你在另一个终端执行 etcdctl
	select {}
}