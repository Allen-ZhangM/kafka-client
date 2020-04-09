package main

import (
	"kafka-client/kago"
	"sync/atomic"
	"time"
)

type KafkaConsumer struct {
	consumer      *kago.PartitionConsumer
	offsetManager *kago.PartitionOffsetManager
	consumerConf  *ConsumerConf
	offset        int64
}

type ConsumerConf struct {
	Topic       string
	Addrs       []string
	Partition   int32
	GroupId     string
	SuccessFunc func(*ConsumerMessage)
	ErrorFunc   func(*ConsumerError)
}

type ConsumerMessage struct {
	Topic     string
	Value     []byte
	Key       []byte
	Partition int32
	Offset    int64
}

type ConsumerError struct {
	Topic     string
	Partition int32
	Err       error
}

func NewKafkaConsumer(conf *ConsumerConf) (*KafkaConsumer, error) {
	//partition Consumer
	config := kago.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.OffsetLocalOrServer = 0 //0,local  1,server  2,newest

	kago.InitOffsetFile() //Initialize the offset file and execute globally once

	//According to the config.OffsetLocalOrServe to pconsumer consumption from the specified offset
	pconsumer, err := kago.InitPartitionConsumer(conf.Addrs, conf.Topic, conf.Partition, conf.GroupId, config)
	if err != nil {
		return nil, err
	}

	pOffsetManager, err := kago.InitPartitionOffsetManager(conf.Addrs, conf.Topic, conf.GroupId, conf.Partition, config)
	if err != nil {
		return nil, err
	}
	kc := &KafkaConsumer{
		consumer:      pconsumer,
		offsetManager: pOffsetManager,
		consumerConf:  conf,
		offset:        pOffsetManager.NextOffset(),
	}

	go kc.handleMsg()

	return kc, nil

}

func (consumer *KafkaConsumer) CloseConsumer() {
	consumer.offsetManager.MarkOffset(consumer.consumerConf.Topic, consumer.consumerConf.Partition, consumer.offset, consumer.consumerConf.GroupId, true)
	consumer.consumer.Close()
}

func (consumer *KafkaConsumer) handleMsg() {
	for {
		select {
		case info := <-consumer.consumer.Recv():
			consumer.consumerConf.SuccessFunc(&ConsumerMessage{
				Topic:     info.Topic,
				Value:     info.Value,
				Partition: info.Partition,
				Key:       info.Key,
				Offset:    info.Offset,
			})
			//Offset is recorded after receiving the message first, and the mode is at least one consumption;  IfExactOnce: false does not log to a local file
			consumer.offsetManager.MarkOffset(info.Topic, info.Partition, info.Offset, consumer.consumerConf.GroupId, false)
			atomic.SwapInt64(&consumer.offset, info.Offset)
		case info := <-consumer.consumer.Errors():
			consumer.consumerConf.ErrorFunc(&ConsumerError{
				Topic:     info.Topic,
				Partition: info.Partition,
				Err:       info.Err,
			})
		}
	}
}
