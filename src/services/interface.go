package services

import (
	nsq "github.com/nsqio/go-nsq"
)

type Status bool

const (
	OK     Status = true
	Reject Status = false
)

type NsqProducer interface {
	// Publish synchronously publishes a message body to the specified topic, returning
	// an error if publish failed
	Publish(topicName string, body []byte) error

	// Stop initiates a graceful stop of the Producer (permanent)
	Stop()
}

func NewProducer(nsqdTCPaddr string) (NsqProducer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdTCPaddr, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
