package serviceResponseNow

import (
	"encoding/json"
	"fmt"
	"github.com/PavlushaSource/NsqBench/src/services"
	"github.com/PavlushaSource/NsqBench/src/services/domain"
	"github.com/nsqio/go-nsq"
)

type ServiceResponse struct {
	NsqLookupdAddr  string
	NsqdAddr        string
	Producer        services.NsqProducer
	Consumer        *nsq.Consumer
	Messages        chan *nsq.Message
	CurrentReceived int
}

func (sr *ServiceResponse) Subscribe(topic domain.Topic, channel domain.Channel, iterations int) error {
	c, err := nsq.NewConsumer(string(topic), string(channel), nsq.NewConfig())

	if err != nil {
		fmt.Println("Err new consumer")
		return err
	}

	messageHandler := services.NewMessageHandler(func(message *nsq.Message) error {
		sr.CurrentReceived++
		sr.Messages <- message

		if iterations <= sr.CurrentReceived {
			close(sr.Messages)
		}
		return nil
	})

	c.AddHandler(messageHandler)
	err = c.ConnectToNSQD(sr.NsqdAddr)
	if err != nil {
		fmt.Println("Err connect to NSQLookupAddr")
		return err
	}
	sr.Consumer = c
	return nil
}

func NewServiceResponse(nsqLookupdAddr, nsqdAddr string, iterations int) (services.Responser, error) {
	producer, err := services.NewProducer(nsqdAddr)
	if err != nil {
		return nil, err
	}

	responser := &ServiceResponse{
		NsqLookupdAddr: nsqLookupdAddr,
		NsqdAddr:       nsqdAddr,
		Messages:       make(chan *nsq.Message),
		Producer:       producer,
	}

	err = responser.Subscribe(domain.RequestTopic, domain.RequestChannel, iterations)
	if err != nil {
		return nil, err
	}

	return responser, nil
}

func (sr *ServiceResponse) Run() error {
	for msgReceive := range sr.Messages {
		var m services.Message
		if err := json.Unmarshal(msgReceive.Body, &m); err != nil {
			return err
		}
		fmt.Println("Received your request, message: ", m.Payload)
		msgToSend := services.NewMessage(m.RespTopic, "ServiceResponse accept your request")
		err := sr.Producer.Publish(m.RespTopic, msgToSend.Marshall())
		if err != nil {
			return err
		}
	}
	return nil
}

func (sr *ServiceResponse) Close() error {
	sr.Consumer.Stop()
	sr.Producer.Stop()
	return nil
}
