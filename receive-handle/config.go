package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
)

var KafkaConsumers []Consumer

type KafkaConfig struct {
	Consumers []Consumer `json:"consumers"`
}

type Consumer struct {
	Topic      string   `json:"topic"`
	Addrs      []string `json:"addrs"`
	Partitions []int32  `json:"partitions"`
	GroupId    string   `json:"group_id"`
}

func ConfigInit() error {
	data, err := ioutil.ReadFile("conf/kafka_configure.json")
	if err == nil {
		config := KafkaConfig{}
		err := json.Unmarshal(data, &config)
		if err == nil {
			logs.Debug("[config] kafka_configure.json load ok.")
			KafkaConsumers = config.Consumers
		} else {
			logs.Debug("[config] kafka_configure.json unmarshal failed, err=%s.", err.Error())
			return err
		}
	} else {
		logs.Debug("[config] kafka_configure.json load failed, err=%s.", err.Error())
		return err
	}
	logs.Debug("[config] Kafka Consumers count  : %d", len(KafkaConsumers))

	for i, kc := range KafkaConsumers {
		logs.Debug("[config] Consumer %d	: topic: %s, addrs: %v , partitions: %v , groupid:%s", i, kc.Topic, kc.Addrs, kc.Partitions, kc.GroupId)
	}
	return nil
}
