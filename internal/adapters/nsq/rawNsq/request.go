package rawNsq

import (
	"context"
	"errors"
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/handler"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/producer"
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"github.com/nsqio/go-nsq"
	"time"
)

type ServiceRequest struct {
	NsqLookupdAddr string
	NsqdAddr       string
	Producer       port.Producer
}

func (sr *ServiceRequest) Run(iterations int) error {
	ctx := context.Background()

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
	return nil
}

func (sr *ServiceRequest) Send(ctx context.Context, reqTopicName, respTopicName domain.Topic, msg string) error {
	msgToSend := domain.NewRequestMessage(string(respTopicName), msg)

	responseChannel := make(chan *domain.Message)

	channelName := fmt.Sprintf("%s-response#ephemeral", msgToSend.ID)

	messageHandler := handler.NewMessageHandler(func(message *nsq.Message) error {
		var m domain.Message
		if err := m.Unmarshall(message.Body); err != nil {
			return err
		}

		if msgToSend.ID == m.MetaInfo.ReqID {
			responseChannel <- &m
		}

		return nil
	})

	c, err := nsq.NewConsumer(string(respTopicName), channelName, nsq.NewConfig())
	if err != nil {
		fmt.Println("Err new consumer")
		return err
	}
	c.AddHandler(messageHandler)

	err = c.ConnectToNSQD(sr.NsqdAddr)
	if err != nil {
		fmt.Println("Err connect to NSQLookupAddr")
		return err
	}
	defer c.Stop()

	// Sleep for each synchronous message
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Message publish now - ", msgToSend.ID)

	err = sr.Producer.Publish(string(reqTopicName), msgToSend.Marshall())
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
		case msgReceive := <-responseChannel:
			fmt.Println("receive response, your message:", msgReceive.Payload)
			return nil
		}
	}
}

func NewServiceRequest(nsqLookupdAddr, nsqdAddr string) (port.RequestService, error) {
	p, err := producer.NewProducer(nsqdAddr)
	if err != nil {
		return nil, err
	}

	return &ServiceRequest{
		NsqLookupdAddr: nsqLookupdAddr,
		Producer:       p,
		NsqdAddr:       nsqdAddr,
	}, nil
}
