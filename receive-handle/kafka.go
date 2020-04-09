package main

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"kafka-client/receive-handle/models"
)

var kcs []*KafkaConsumer

func InitKafkaConsumers() error {
	for _, kc := range KafkaConsumers {
		var err error
		switch kc.Topic {
		case models.TopicBizReport:
			err = setBizReportConsumer(kc)
		case models.TopicLsmReport:
			err = setLsmReportConsumer(kc)
		case models.TopicVVSession:
			err = setVVSessionConsumer(kc)
		default:
			logs.Error("undefined or error topic :", kc.Topic)
			return errors.New("undefined or error topic")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func setBizReportConsumer(kc Consumer) error {
	conf := setDefaultConsumerConf(kc)
	conf.SuccessFunc = func(msg *ConsumerMessage) {
		logs.Info("kafka receive msg :topic:%s, Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
	return initConsumer(conf)
}

func setLsmReportConsumer(kc Consumer) error {
	conf := setDefaultConsumerConf(kc)
	conf.SuccessFunc = func(msg *ConsumerMessage) {
		logs.Info("kafka receive msg :topic:%s, Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
	return initConsumer(conf)
}

func setVVSessionConsumer(kc Consumer) error {
	conf := setDefaultConsumerConf(kc)
	conf.SuccessFunc = func(msg *ConsumerMessage) {
		logs.Info("kafka receive msg :topic:%s, Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
	return initConsumer(conf)
}

func setDefaultConsumerConf(kc Consumer) *ConsumerConf {
	return &ConsumerConf{
		Topic:     kc.Topic,
		Addrs:     kc.Addrs,
		Partition: int32(kc.Partitions[0]),
		GroupId:   kc.GroupId,
		SuccessFunc: func(msg *ConsumerMessage) {
			logs.Info("kafka receive msg :topic:%s, Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		},
		ErrorFunc: func(msg *ConsumerError) {
			logs.Error("kafka receive msg err, topic:%s, Partition:%d, err:%v", msg.Topic, msg.Partition, msg.Err)
		},
	}

}

func initConsumer(conf *ConsumerConf) error {
	if kc, err := NewKafkaConsumer(conf); err != nil {
		logs.Error("failed to init kafka Consumer : %s", err.Error())
		return err
	} else {
		kcs = append(kcs, kc)
		logs.Info("success to init Consumer, topic:%s,addrs:%v,len(addrs):%v,partition:%d,groupId:%s ", conf.Topic, conf.Addrs, len(conf.Addrs), conf.Partition, conf.GroupId)
	}
	return nil
}

func MarkOffsetConsumers(sig string) {
	for _, kc := range kcs {
		kc.CloseConsumer()
	}
}

func handleLsmReportMsg(data []byte) {
	lsmReport := new(models.LsmReport)
	if err := json.Unmarshal(data, lsmReport); err != nil {
		logs.Warn("invalid LsmReport: %s", err.Error())
		return
	}
}

func handleVVSessionMsg(data []byte) {
	vvSession := new(models.VVSession)
	if err := json.Unmarshal(data, vvSession); err != nil {
		logs.Warn("invalid VVSession: %s", err.Error())
		return
	}
}
