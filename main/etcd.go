package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	ec "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"logagent/tailf"
	"time"
)

type EtcdClient struct {
	Client *ec.Client
	Keys   []string
}

var (
	etcdClient *EtcdClient
	etcdKey    = "/logagent/config/"
	cocs       = make([]*tailf.CollectConfig, 0, 10)
)

func initEtcd(addr string) error {
	cli, err := ec.New(ec.Config{
		Endpoints:   []string{addr},
		DialTimeout: time.Second,
	})
	if err != nil {
		logs.Error(err)
		return err
	}
	fmt.Println("etcd 连接成功")

	etcdClient = &EtcdClient{
		Client: cli,
		Keys:   make([]string, 0, 10),
	}

	for _, ip := range localIPs {
		etcdClient.Keys = append(etcdClient.Keys, etcdKey+ip)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey+ip)
		if err != nil {
			continue
		}
		cancel()

		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey+ip {
				err = json.Unmarshal(v.Value, &cocs)
				if err != nil {
					logs.Error(err)
					continue
				}
				logs.Debug("log config is %v", *cocs[0])
			}
		}
	}

	initEtcdWatcher()
	return nil
}

func initEtcdWatcher() {
	for _, key := range etcdClient.Keys {
		go watchKey(key)
	}
}

func watchKey(key string) {
	cli, err := ec.New(ec.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("1:", err)
		return
	}

	for {
		rch := cli.Watch(context.Background(), key)
		var collectConfigs []*tailf.CollectConfig
		var getConfSucc = true

		for wresp := range rch {
			for _, ev := range wresp.Events { // 会有很多种类事件
				fmt.Printf("%s %q: %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &collectConfigs)
					if err != nil {
						logs.Error("key [%s], unmarshal error: %v", err)
						getConfSucc = false
						continue
					}
				}
			}

			if getConfSucc {
				tailf.UpdateConfig(collectConfigs)
			}
		}
	}
}
func GetCocs() []*tailf.CollectConfig {
	return cocs
}
