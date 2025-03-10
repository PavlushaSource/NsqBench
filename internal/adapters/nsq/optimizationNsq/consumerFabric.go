package optimizationNsq

import (
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/adapters/nsq/handler"
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

var ConsumerFabricInstance = NewConsumerFabric()

type ConsumerFabric struct {
	consumers map[string]*nsq.Consumer

	sync.RWMutex
	Router *ConsumerRouter
}

func (c *ConsumerFabric) get(key string) (*nsq.Consumer, bool) {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.consumers[key]
	return v, ok
}

func (c *ConsumerFabric) register(topicName domain.Topic, channelName domain.Channel, nsqd string, handlerFunc nsq.Handler) (*nsq.Consumer, error) {
	newC, err := nsq.NewConsumer(string(topicName), string(channelName), nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("nsq new consumer failed: %w", err)
	}

	newC.SetLoggerLevel(nsq.LogLevelError)
	newC.AddHandler(handlerFunc)

	err = newC.ConnectToNSQD(nsqd)
	if err != nil {
		return nil, fmt.Errorf("nsq connect to nsqd failed: %w", err)
	}

	c.Lock()
	defer c.Unlock()
	c.consumers[GenerateMapKey(topicName, channelName)] = newC

	return newC, nil
}

func (c *ConsumerFabric) Stop() {
	c.Lock()
	defer c.Unlock()

	for k, v := range c.consumers {
		v.Stop()
		delete(c.consumers, k)
	}
}

func (c *ConsumerFabric) NewSyncConsumer(topicName domain.Topic, channelName domain.Channel, nsqd string, handleFunc nsq.HandlerFunc) (*nsq.Consumer, error) {

	key := GenerateMapKey(topicName, channelName)
	consumer, ok := c.get(key)
	//fmt.Println("Register new consumer for topic:", topicName, "channel:", channelName)

	if !ok {
		fmt.Println("Register new consumer for topic:", topicName, "channel:", channelName)
		hndl := handler.NewMessageHandler(handleFunc)

		var err error
		consumer, err = c.register(topicName, channelName, nsqd, hndl)
		if err != nil {
			return nil, err
		}

		// Time sleep only for new consumer register
		time.Sleep(100 * time.Millisecond)

	}

	return consumer, nil
}

func NewConsumerFabric() port.ConsumerFabric {
	return &ConsumerFabric{
		Router:    NewConsumerRouter(),
		consumers: make(map[string]*nsq.Consumer),
	}
}

func GenerateMapKey(topicName domain.Topic, channelName domain.Channel) string {
	return fmt.Sprintf("%s-%s", topicName, channelName)
}

var Router = NewConsumerRouter()

type ConsumerRouter struct {
	sync.RWMutex
	responseChannels map[string]chan *domain.Message
}

func NewConsumerRouter() *ConsumerRouter {
	return &ConsumerRouter{
		responseChannels: make(map[string]chan *domain.Message),
	}
}

func (c *ConsumerRouter) RegisterChannel(UUID string) chan *domain.Message {
	c.Lock()
	defer c.Unlock()

	ch := make(chan *domain.Message)
	c.responseChannels[UUID] = ch

	return ch
}

func (c *ConsumerRouter) GetChannel(UUID string) (chan *domain.Message, bool) {
	c.RLock()
	defer c.RUnlock()

	ch, ok := c.responseChannels[UUID]
	return ch, ok
}
