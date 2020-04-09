package models

import (
	"github.com/Shopify/sarama"
	"time"
)

type PeerInfo struct {
	Pid       string
	PidPrefix string
	Version   string
	PublibIP  string
	IPId      int
	NatType   int8
	IsLogin   int8
}

const (
	TopicBizReport   = "biz_report"
	TopicPeerSession = "peer_session"
	TopicLsmReport   = "lsm_report"
	TopicVVSession   = "vv_session"
)

type BizReportULFlow struct {
	Timestamp int64 `json:"ts"`
	Duration  int   `json:"period"`
	Bytes     int64 `json:"bytes"`
}

type BizReportDLFlow struct {
	Timestamp int64 `json:"ts"`
	Duration  int   `json:"period"`
	AppBytes  int64 `json:"app"`
	CdnBytes  int64 `json:"cdn"`
	P2pBytes  int64 `json:"p2p"`
}

type BizReportUpload struct {
	Flows []BizReportULFlow `json:"flow"`
}

type BizReportDownload struct {
	Url            string            `json:"url"`
	Vvid           string            `json:"vid"`
	Type           string            `json:"type"`
	FirstPlayTime  float64           `json:"boottime"`
	DelayTime      float64           `json:"delaytime"`
	BufferingCount int               `json:"bufferingcount"`
	Flows          []BizReportDLFlow `json:"flow"`
}

type BizReportDistribute struct {
	Url   string            `json:"url"`
	Type  string            `json:"type"`
	Flows []BizReportULFlow `json:"flow"`
}

type BizReport struct {
	SerialId    string                `json:"id"`
	PeerId      string                `json:"pid"`
	Duration    int                   `json:"period"`
	Uploads     BizReportUpload       `json:"upload"`
	Downloads   []BizReportDownload   `json:"download"`
	Distributes []BizReportDistribute `json:"push"`
}

type PeerSession struct {
	Pid       string    `json:"pid"`
	Version   string    `json:"version"`
	PublicIP  string    `json:"pub_ip"`
	NatType   int       `json:"nat"`
	Timestamp time.Time `json:"ts"`
	IsLogin   bool      `json:"is_login"`
}

type LsmReport struct {
	Pid       string    `json:"pid"`
	LsmFree   int64     `json:"lsm_free"`
	LsmTotal  int64     `json:"lsm_total"`
	DiskFree  int64     `json:"disk_free"`
	DiskTotal int64     `json:"disk_total"`
	Timestamp time.Time `json:"ts"`
}

type VVSession struct {
	Pid       string    `json:"pid"`
	Url       string    `json:"url"`
	Timestamp time.Time `json:"ts"`
	IsStart   bool      `json:"is_start"`
}

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
