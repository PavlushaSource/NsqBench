package main

import (
	"github.com/PavlushaSource/NsqBench/src/services/domain"
	"github.com/PavlushaSource/NsqBench/src/services/requestAdapters/serviceRequestOpt"
	"log"
)

func main() {
	//sr, err := serviceRequestNow.NewServiceRequest("127.0.0.1:4161", "127.0.0.1:4150")
	sr, err := serviceRequestOpt.NewServiceRequestOpt("127.0.0.1:4161", "127.0.0.1:4150", domain.ResponseTopic, domain.ResponseChannel)
	if err != nil {
		log.Fatal(err)
	}

	defer sr.Close()

	if err = sr.Run(10); err != nil {
		log.Fatal(err)
	}
}
