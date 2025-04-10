package main

// Java: 相当于 public class Main，main 函数的入口类必须是这个名字

import (
	"context"
	// Java: 类似于 ScheduledExecutorService 的 cancel/timeout 控制器（上下文）

	"fmt"
	// Java: 相当于 System.out.println、String.format

	// Java: 相当于 import com.etcd.client.v3.EtcdClient （引入三方库）

	"time"
	// Java: 类似 java.time.Duration，用来控制 sleep、超时

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Java: 相当于 public static void main(String[] args) {}

	// 创建一个 etcd 客户端连接（配置 IP、端口、超时时间）
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, 
		// Java: 等价于 new EtcdClient.Builder().setEndpoints("localhost:2379")

		DialTimeout: 5 * time.Second,            
		// Java: 等价于 .setDialTimeout(Duration.ofSeconds(5))
	})
	if err != nil {
		// Java: 相当于 if (err != null) throw new RuntimeException(err)
		panic(err)
	}

	defer cli.Close()
	// Java: 相当于 try-finally 块，在函数结束时自动执行 cli.close()

	// 写入 key-value
	_, err = cli.Put(context.Background(), "foo", "bar")
	// Java: 相当于 client.put("foo", "bar")（异步方法）

	if err != nil {
		panic(err)
	}
	fmt.Println("Put foo=bar") 
	// Java: System.out.println("Put foo=bar")

	// 读取 key 的值
	resp, err := cli.Get(context.Background(), "foo")
	// Java: EtcdClient.GetResponse resp = client.get("foo")

	if err != nil {
		panic(err)
	}
	for _, kv := range resp.Kvs {
		// Java: for (KV kv : resp.getKvs())
		fmt.Printf("Get %s=%s\n", kv.Key, kv.Value)
		// Java: System.out.printf("Get %s=%s", kv.getKey(), kv.getValue())
	}
}