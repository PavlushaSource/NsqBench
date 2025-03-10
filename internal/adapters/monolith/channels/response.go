package channels

import (
	"github.com/PavlushaSource/NsqBench/internal/core/domain"
	"github.com/PavlushaSource/NsqBench/internal/core/port"
)

type ServiceResponse struct {
	requestChannel  <-chan domain.Message
	responseChannel chan<- domain.Message
}

func NewServiceResponse(requestChannel <-chan domain.Message, responseChannel chan<- domain.Message) port.ResponseService {
	return &ServiceResponse{requestChannel, responseChannel}
}

func (s ServiceResponse) Run() error {
	for msg := range s.requestChannel {
		msgToSend := domain.NewResponseMessage(msg.RespTopic, "Hello too!", msg.ID)
		s.responseChannel <- msgToSend
	}
	return nil
}

func (s ServiceResponse) Close() error {
	close(s.responseChannel)
	return nil
}
