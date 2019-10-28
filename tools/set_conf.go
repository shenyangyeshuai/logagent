package main

import (
	"context"
	"encoding/json"
	"fmt"
	ec "go.etcd.io/etcd/clientv3"
	"logagent/tailf"
	"time"
)

var (
	etcdKey = "/logagent/config/192.168.199.183"
)

func main() {
	cfg := ec.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: time.Second,
	}

	cli, err := ec.New(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()

	var logConfs []*tailf.CollectConfig

	logConfs = append(logConfs, &tailf.CollectConfig{
		LogPath: "/home/yeshuai/go-1.12/vendor/src/logagent/logs/logagent.log",
		Topic:   "nginx_log",
	})

	logConfs = append(logConfs, &tailf.CollectConfig{
		LogPath: "/home/yeshuai/go-1.12/vendor/src/logagent/logs/logagent2.log",
		Topic:   "nginx_log",
	})

	data, err := json.Marshal(logConfs)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// cli.Delete(ctx, etcdKey)
	// return

	_, err = cli.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println(err)
		return
	}
}
