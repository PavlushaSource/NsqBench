package main

import (
	"context"
	"github.com/PavlushaSource/NsqBench/src/services"
	"github.com/PavlushaSource/NsqBench/src/services/serviceRequest"
	"log"
	"time"
)

func main() {
	sr, err := serviceRequest.NewServiceRequest("127.0.0.1:4161", "127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = sr.Send(ctx, services.RequestTopic, services.ResponseTopic, "bro, вышли мне ответ пж")
	if err != nil {
		log.Fatal(err)
	}
}
