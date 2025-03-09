package optNsq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/handler"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/producer"
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"github.com/nsqio/go-nsq"
	"time"
)

type ServiceRequestOpt struct {
	NsqLookupdAddr  string
	NsqdAddr        string
	Producer        port.Producer
	Consumer        *nsq.Consumer
	ResponseChannel chan *nsq.Message
}

func (sr *ServiceRequestOpt) Run(iterations int) error {
	ctx := context.Background()

	for i := 0; i < iterations; i++ {
		err := sr.Send(ctx, domain.RequestTopic, domain.ResponseTopic, "hello")
		if err != nil {
			return fmt.Errorf("not send request (%d): %w", i, err)
		}
	}
	return nil
}

func (sr *ServiceRequestOpt) Close() error {
	sr.Producer.Stop()
	sr.Consumer.Stop()
	return nil
}

func (sr *ServiceRequestOpt) Send(ctx context.Context, reqTopicName, respTopicName domain.Topic, msg string) error {
	msgToSend := domain.NewRequestMessage(string(respTopicName), msg)

	//channelName := fmt.Sprintf("%s-response#ephemeral", msgToSend.ID)

	fmt.Println("Message publish now - ", msgToSend.ID)

	err := sr.Producer.Publish(string(reqTopicName), msgToSend.Marshall())
	if err != nil {
		return err
	}

	fmt.Println("Start wait for response")
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			if errors.Is(context.DeadlineExceeded, err) {
				return fmt.Errorf("error: Context timeout")
			}

			return err
		case msgReceive := <-sr.ResponseChannel:
			var m domain.Message
			if err = json.Unmarshal(msgReceive.Body, &m); err != nil {
				return err
			}
			fmt.Println("receive response, your message:", m.Payload)
			return nil
		}
	}
}

func (sr *ServiceRequestOpt) NewConsumer(consumerTopicName domain.Topic, consumerChannelName domain.Channel) error {
	c, err := nsq.NewConsumer(string(consumerTopicName), string(consumerChannelName), nsq.NewConfig())
	if err != nil {
		fmt.Println("Err new consumer")
		return err
	}

	messageHandler := handler.NewMessageHandler(func(message *nsq.Message) error {
		sr.ResponseChannel <- message
		return nil
	})
	c.AddHandler(messageHandler)

	now := time.Now()
	err = c.ConnectToNSQD(sr.NsqdAddr)
	fmt.Println("Time connect to NSQLookupAddr: ", time.Since(now))
	if err != nil {
		fmt.Println("Err connect to NSQLookupAddr")
		return err
	}

	sr.Consumer = c
	return nil
}

func NewServiceRequestOpt(nsqLookupdAddr, nsqdAddr string, consumerTopicName domain.Topic, consumerChannelName domain.Channel) (port.RequestService, error) {
	p, err := producer.NewProducer(nsqdAddr)
	if err != nil {
		return nil, err
	}

	requester := &ServiceRequestOpt{
		NsqLookupdAddr:  nsqLookupdAddr,
		Producer:        p,
		NsqdAddr:        nsqdAddr,
		ResponseChannel: make(chan *nsq.Message),
	}

	err = requester.NewConsumer(consumerTopicName, consumerChannelName)
	if err != nil {
		return nil, err
	}

	return requester, nil
}
