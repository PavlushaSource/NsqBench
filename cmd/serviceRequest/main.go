package main

import (
	"github.com/PavlushaSource/NsqBench/src/services/requestAdapters/serviceRequestNow"
	"log"
)

func main() {
	sr, err := serviceRequestNow.NewServiceRequest("127.0.0.1:4161", "127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}

	defer sr.Close()

	if err = sr.Run(10); err != nil {
		log.Fatal(err)
	}
}
