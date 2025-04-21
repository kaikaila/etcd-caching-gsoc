package watcher // 类似 Java 的 package 声明，比如 package com.example.main;

import (
	"context" // 类似 Java 的 java.util.concurrent.CancellationException + Future cancel 管理
	"fmt"     // 类似 Java 的 System.out.println

	// Java 的 java.time 包
	clientv3 "go.etcd.io/etcd/client/v3" // 导入 etcd 的 Go 客户端，类似 Java 的第三方依赖
)

// Go 中函数以小写字母开头表示 “包内可见”，类似 Java 的 package-private 函数
// 等价于：void watchKey(Client cli, String key)
func WatchKey(cli *clientv3.Client, key string, onPut func(string, string), onDelete func(string)) {
	// 相当于：Context ctx = new Context(); 用于控制取消、超时等
	ctx := context.Background()

	// 相当于：cli.watch(ctx, key)，返回一个异步的事件流（channel）
	rch := cli.Watch(ctx, key)

	// 类似 Java：System.out.println("Start watching key: " + key);
	fmt.Printf("Start watching key: %s\n", key)

	// Go 的 channel 可以 for 循环消费，类似 Java 的 while(true) + queue.take()
	for wresp := range rch {
		// 每个响应里可能有多个事件，比如 PUT、DELETE 等
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				onPut(string(ev.Kv.Key), string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				onDelete(string(ev.Kv.Key))
			}
		}
	}
}

// for memoryCache
func WatchKeySimple(cli *clientv3.Client, key string, onPut func(string, string), onDelete func(string)) {
	ch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	go func() {
		for resp := range ch {
			for _, ev := range resp.Events {
				k := string(ev.Kv.Key)
				v := string(ev.Kv.Value)
				switch ev.Type {
				case clientv3.EventTypePut:
					onPut(k, v)
				case clientv3.EventTypeDelete:
					onDelete(k)
				}
			}
		}
	}()
}

//for watchCache
func WatchKeyWithRevision(cli *clientv3.Client, key string, onPut func(string, string, int64), onDelete func(string, int64)) {
	ch := cli.Watch(context.Background(), key, clientv3.WithPrefix())
	go func() {
		for resp := range ch {
			for _, ev := range resp.Events {
				k := string(ev.Kv.Key)
				v := string(ev.Kv.Value)
				rev := ev.Kv.ModRevision
				switch ev.Type {
				case clientv3.EventTypePut:
					onPut(k, v, rev)
				case clientv3.EventTypeDelete:
					onDelete(k, rev)
				}
			}
		}
	}()
}
