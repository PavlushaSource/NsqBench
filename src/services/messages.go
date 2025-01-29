package services

import (
	"encoding/json"
	uuid "github.com/google/uuid"
	"log"
)

type Message struct {
	Payload   string `json:"payload"`
	ID        string `json:"id"`
	RespTopic string `json:"topic"`
}

func NewMessage(respTopicName string, payload string) Message {
	id := uuid.New().String()

	return Message{
		Payload:   payload,
		ID:        id,
		RespTopic: respTopicName,
	}
}

func (m *Message) Marshall() []byte {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}
