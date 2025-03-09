package handler

import "github.com/nsqio/go-nsq"

type MessageHandler struct {
	handleFunc func(*nsq.Message) error
}

func NewMessageHandler(handleFunc func(*nsq.Message) error) *MessageHandler {
	return &MessageHandler{handleFunc: handleFunc}
}

func (h *MessageHandler) HandleMessage(message *nsq.Message) error {
	return h.handleFunc(message)
}
