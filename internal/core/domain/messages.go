package domain

import (
	"encoding/json"
	uuid "github.com/google/uuid"
	"log"
)

type MetaInfo struct {
	ReqID string `json:"req_id"`
}

type Message struct {
	Payload   string `json:"payload"`
	ID        string `json:"id"`
	RespTopic string `json:"topic"`
	MetaInfo  MetaInfo
}

func NewResponseMessage(respTopicName string, payload string, reqID string) Message {
	id := uuid.New().String()

	return Message{
		Payload:   payload,
		ID:        id,
		RespTopic: respTopicName,
		MetaInfo: MetaInfo{
			ReqID: reqID,
		},
	}
}

func NewRequestMessage(respTopicName string, payload string) Message {
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

func (m *Message) Unmarshall(data []byte) error {
	return json.Unmarshal(data, m)
}
