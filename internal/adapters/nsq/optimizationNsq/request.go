package optimizationNsq

import (
	"context"
	"errors"
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/producer"
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"github.com/nsqio/go-nsq"
)

type ServiceRequest struct {
	NsqLookupdAddr  string
	NsqdAddr        string
	Producer        port.Producer
	ResponseChannel chan *nsq.Message
}

func (sr *ServiceRequest) Run(ctx context.Context, iterations int) error {
	//ctx := context.Background()

	for i := 0; i < iterations; i++ {
		err := sr.Send(ctx, domain.RequestTopic, domain.ResponseTopic, "hello")
		if err != nil {
			return fmt.Errorf("not send request (%d): %w", i, err)
		}
	}
	return nil
}

func (sr *ServiceRequest) Close() error {
	sr.Producer.Stop()
	ConsumerFabricInstance.Stop()
	return nil
}

func (sr *ServiceRequest) Send(ctx context.Context, reqTopicName, respTopicName domain.Topic, msg string) error {
	msgToSend := domain.NewRequestMessage(string(respTopicName), msg)

	respCh := Router.RegisterChannel(msgToSend.ID)

	handleFunc := func(message *nsq.Message) error {
		var m domain.Message
		if err := m.Unmarshall(message.Body); err != nil {
			return fmt.Errorf("unmarshal message failed: %w", err)
		}

		if ch, ok := Router.GetChannel(m.MetaInfo.ReqID); ok {
			//fmt.Println("receive response, your message:", m.Payload)
			ch <- &m
		}

		return nil
	}

	c, err := ConsumerFabricInstance.NewSyncConsumer(
		respTopicName, domain.ResponseChannel, sr.NsqdAddr, handleFunc)

	// consumer instance not needed for sync request
	_ = c
	if err != nil {
		return fmt.Errorf("new sync consumer failed: %w", err)
	}

	err = sr.Producer.Publish(string(reqTopicName), msgToSend.Marshall())
	if err != nil {
		return err
	}
	//fmt.Println("Message publish now - ", msgToSend.ID)

	//fmt.Println("Start wait response msg")
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			if errors.Is(context.DeadlineExceeded, err) {
				return fmt.Errorf("error: Context timeout")
			}

			return nil
		case <-respCh:
			return nil
		}
	}
}

func NewServiceRequest(nsqLookupdAddr, nsqdAddr string) (port.RequestService, error) {
	p, err := producer.NewProducer(nsqdAddr)
	if err != nil {
		return nil, err
	}

	requester := &ServiceRequest{
		NsqLookupdAddr:  nsqLookupdAddr,
		Producer:        p,
		NsqdAddr:        nsqdAddr,
		ResponseChannel: make(chan *nsq.Message),
	}

	return requester, nil
}
