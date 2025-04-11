#!/bin/bash

echo "🔧 停止并移除 etcd 容器..."

docker stop etcd >/dev/null 2>&1
docker rm etcd >/dev/null 2>&1

echo "✅ etcd 容器已安全停止并移除。"