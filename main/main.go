package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/kafka"
	"logagent/tailf"
)

var (
	typ      = "ini"
	filename = "../conf/logagent.conf"
)

func main() {
	// 1. 初始化配置文件
	err := loadConfig(typ, filename)
	if err != nil {
		fmt.Printf("配置文件加载失败: %v\n", err)
		return
	}

	// 2. 初始化日志
	err = initLogger()
	if err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		return
	}

	// 3. 初始化 tailf
	err = tailf.InitTail(appConfig.Collects)
	if err != nil {
		logs.Error("tailf 初始化失败: %v", err)
		return
	}

	// 测试
	// go func() {
	//	var count = 0
	//	for {
	//		count++
	//		logs.Debug("test for logger %d", count)
	//		time.Sleep(time.Second)
	//	}
	// }()

	// 4. 初始化 kafka
	err = kafka.InitKafka(appConfig.KafkaAddr)
	if err != nil {
		logs.Error("kafka 初始化失败: %v", err)
		return
	}

	// 启动业务逻辑(其实就是从 tailf 中不断读出 msg 然后塞给 kafka)
	err = run()
	if err != nil {
		logs.Error("启动业务逻辑失败: %v", err)
		return
	}
}
