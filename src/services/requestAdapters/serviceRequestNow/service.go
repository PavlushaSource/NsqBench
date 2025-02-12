package serviceRequestNow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PavlushaSource/NsqBench/src/services"
	"github.com/PavlushaSource/NsqBench/src/services/domain"
	"github.com/nsqio/go-nsq"
	"time"
)

type ServiceRequest struct {
	NsqLookupdAddr string
	NsqdAddr       string
	Producer       services.NsqProducer
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
	msgToSend := services.NewMessage(string(respTopicName), msg)

	responseChannel := make(chan *nsq.Message)

	channelName := fmt.Sprintf("%s-response#ephemeral", msgToSend.ID)

	messageHandler := services.NewMessageHandler(func(message *nsq.Message) error {
		responseChannel <- message
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
			var m services.Message
			if err = json.Unmarshal(msgReceive.Body, &m); err != nil {
				return err
			}
			fmt.Println("receive response, your message:", m.Payload)
			return nil
		}
	}
}

func NewServiceRequest(nsqLookupdAddr, nsqdAddr string) (services.Requester, error) {
	producer, err := services.NewProducer(nsqdAddr)
	if err != nil {
		return nil, err
	}

	return &ServiceRequest{
		NsqLookupdAddr: nsqLookupdAddr,
		Producer:       producer,
		NsqdAddr:       nsqdAddr,
	}, nil
}
