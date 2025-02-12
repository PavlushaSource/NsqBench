package main

import (
	"fmt"
	"github.com/PavlushaSource/NsqBench/src/services/responseAdapters/serviceResponseNow"
	"log"
)

func main() {
	sr, err := serviceResponseNow.NewServiceResponse("127.0.0.1:4161", "127.0.0.1:4150")
	if err != nil {
		log.Fatal(err)
	}
	defer sr.Close()

	fmt.Println("Subscribed and start wait request!")
	if err = sr.Run(); err != nil {
		log.Fatal(err)
	}
}
