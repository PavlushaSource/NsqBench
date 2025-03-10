package port

import (
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/nsqio/go-nsq"
)

type ConsumerFabric interface {
	NewSyncConsumer(topicName domain.Topic, channelName domain.Channel, nsqd string, handleFunc nsq.HandlerFunc) (*nsq.Consumer, error)
	Stop()
}
