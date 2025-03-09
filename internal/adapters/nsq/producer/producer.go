package producer

import (
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"github.com/nsqio/go-nsq"
)

func NewProducer(nsqdTCPaddr string) (port.Producer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdTCPaddr, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
