package channels

import (
	"context"
	"errors"
	"fmt"
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
)

type ServiceRequest struct {
	responseChannel <-chan domain.Message
	requestChannel  chan<- domain.Message
}

func NewServiceRequest(responseChannel <-chan domain.Message, requestChannel chan<- domain.Message) port.RequestService {
	return &ServiceRequest{responseChannel, requestChannel}
}

func (s *ServiceRequest) Send(ctx context.Context, reqTopicName domain.Topic, msg string) error {
	msgToSend := domain.NewRequestMessage(string(reqTopicName), msg)

	// send request with ctx
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(context.DeadlineExceeded, err) {
			return fmt.Errorf("error: Context timeout")
		}
		return err
	case s.requestChannel <- msgToSend:
		//fmt.Println("Message publish now - ", msgToSend.Payload)
	}

	// receive response with ctx
	select {
	case <-ctx.Done():
		err := ctx.Err()
		if errors.Is(context.DeadlineExceeded, err) {
			return fmt.Errorf("error: Context timeout")
		}

		return err
	case msgReceive := <-s.responseChannel:
		//fmt.Println("receive response, your message:", msgReceive.Payload)
		_ = msgReceive
		return nil
	}
}

func (s *ServiceRequest) Run(ctx context.Context, iterations int) error {
	for i := 0; i < iterations; i++ {
		err := s.Send(ctx, domain.RequestTopic, "hello")
		if err != nil {
			return fmt.Errorf("not send request (%d): %w", i, err)
		}
	}
	return nil
}

func (s *ServiceRequest) Close() error {
	close(s.requestChannel)
	return nil
}
