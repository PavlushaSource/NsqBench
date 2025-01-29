package main

import (
	"fmt"
	"github.com/PavlushaSource/NsqBench/src/services"
	"github.com/PavlushaSource/NsqBench/src/services/serviceResponse"
	"log"
)

func main() {
	sr, err := serviceResponse.NewServiceResponse("127.0.0.1:4161", "127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}

	err = sr.Subscribe(services.RequestTopic, "TestChannel")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Subscribed!")
	if err = sr.Run(); err != nil {
		log.Fatal(err)
	}
}
