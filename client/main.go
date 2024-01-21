package main

import (
	"client/services"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	creds := insecure.NewCredentials()
	cc, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := services.NewCalculatorClient(cc)
	service := services.NewCalculatorService(client)

	// err = service.Hello("Go")
	err = service.Fibonacci(6)
	if err != nil {
		log.Fatal(err)
	}
}
