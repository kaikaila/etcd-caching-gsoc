#!/bin/bash

docker rm -f etcd >/dev/null 2>&1

docker run -d -p 2379:2379 --name etcd quay.io/coreos/etcd:v3.5.11 \
  etcd --advertise-client-urls http://localhost:2379 \
       --listen-client-urls http://0.0.0.0:2379
