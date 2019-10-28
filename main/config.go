package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"logagent/tailf"
)

var (
	appConfig *Config
)

type Config struct {
	LogLevel  string
	LogPath   string
	Collects  []*tailf.CollectConfig
	KafkaAddr string
	EtcdAddr  string
}

func NewConfig() *Config {
	return &Config{}
}

func loadConfig(typ, filename string) error {
	conf, err := config.NewConfig(typ, filename)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return err
	}

	appConfig = NewConfig()

	appConfig.LogLevel = conf.String("logs::level")
	appConfig.LogPath = conf.String("logs::path")

	collectConfig := tailf.NewCollectConfig()
	collectConfig.LogPath = conf.String("collect::path")
	collectConfig.Topic = conf.String("collect::topic")

	appConfig.Collects = append(appConfig.Collects, collectConfig)

	appConfig.KafkaAddr = conf.String("kafka::addr")
	appConfig.EtcdAddr = conf.String("etcd::addr")

	return nil
}
