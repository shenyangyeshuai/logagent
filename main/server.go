package main

import (
	"github.com/astaxie/beego/logs"
	"logagent/kafka"
	"logagent/tailf"
	"time"
)

func run() error {
	for {
		msg := tailf.GetOneLine()
		err := sendToKafka(msg)
		if err != nil {
			logs.Error("消息发送给 kafka 失败")
			time.Sleep(time.Second)
			continue
		}
	}
	return nil
}

func sendToKafka(msg *tailf.TextMsg) error {
	err := kafka.SendToKafka(msg.Msg, msg.Topic)
	if err != nil {
		logs.Error("发送给 kafka 失败")
		return err
	}

	return nil
}
