package main

import (
	"github.com/Shopify/sarama"
)

type KafkaProducer struct {
	asyncProducer sarama.AsyncProducer
	producerConf  *ProducerConf
}

type CallbackInfo struct {
	ProducerMessage ProducerMessage
	Err             error
}

type ProducerMessage struct {
	Topic     string
	Value     []byte
	Key       []byte
	Partition int32
	Offset    int64
}

type ProducerConf struct {
	Addrs       []string
	SuccessFunc func(*CallbackInfo)
	ErrorFunc   func(*CallbackInfo)
}

func InitProducer(pConf *ProducerConf) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_2_0_0

	p, err := sarama.NewAsyncProducer(pConf.Addrs, config)
	if err != nil {
		return nil, err
	}

	kafkaProducer := &KafkaProducer{p, pConf}

	go kafkaProducer.handleMsg()

	return kafkaProducer, nil
}

func (producer *KafkaProducer) PushMsg(msg *ProducerMessage) {
	producer.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic:     msg.Topic,
		Value:     sarama.ByteEncoder(msg.Value),
		Key:       sarama.ByteEncoder(msg.Key),
		Partition: msg.Partition,
	}
}

func (producer *KafkaProducer) handleMsg() {
	for {
		select {
		case info := <-producer.asyncProducer.Successes():
			var key []byte
			if info.Key != nil {
				key, _ = info.Key.Encode()
			}
			value, _ := info.Value.Encode()
			producer.producerConf.SuccessFunc(&CallbackInfo{
				ProducerMessage: ProducerMessage{
					Topic:     info.Topic,
					Value:     value,
					Partition: info.Partition,
					Key:       key,
					Offset:    info.Offset,
				},
				Err: nil,
			})
		case info := <-producer.asyncProducer.Errors():
			var key []byte
			if info.Msg.Key != nil {
				key, _ = info.Msg.Key.Encode()
			}
			value, _ := info.Msg.Value.Encode()
			producer.producerConf.ErrorFunc(&CallbackInfo{
				ProducerMessage: ProducerMessage{
					Topic:     info.Msg.Topic,
					Value:     value,
					Partition: info.Msg.Partition,
					Key:       key,
					Offset:    info.Msg.Offset,
				},
				Err: info.Err,
			})
		}
	}
}
