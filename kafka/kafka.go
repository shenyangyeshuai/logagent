package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

func InitKafka(addr string) error {
	// kafka 配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 创建一个生产者客户端
	var err error
	client, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Error("kafka 生产者启动失败")
		return err
	}
	// defer client.Close() // NMB

	logs.Debug("kafka 启动成功")

	return nil
}

func SendToKafka(data, topic string) error {
	// 生产一条消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	// 交给 kafka
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		logs.Error("(生产)发送给 kafka 的消息失败")
		return err
	}
	logs.Debug("发送成功! pid: %v, offset: %v, topic: %v\n", pid, offset, topic)

	return nil
}
