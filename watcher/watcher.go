package watcher // 类似 Java 的 package 声明，比如 package com.example.main;

import (
	"context" // 类似 Java 的 java.util.concurrent.CancellationException + Future cancel 管理
	"fmt"     // 类似 Java 的 System.out.println

	// Java 的 java.time 包
	clientv3 "go.etcd.io/etcd/client/v3" // 导入 etcd 的 Go 客户端，类似 Java 的第三方依赖
)

// Go 中函数以小写字母开头表示 “包内可见”，类似 Java 的 package-private 函数
// 等价于：void watchKey(Client cli, String key)
func WatchKey(cli *clientv3.Client, key string) {
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
			// 类似 Java：System.out.printf("Type: %s Key: %s Value: %s\n", ev.getType(), ev.getKey(), ev.getValue());
			fmt.Printf("Type: %s Key: %s Value: %s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}