package main

import (
	"log"
	"net"
	"server/services"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	services.RegisterCalculatorServer(s, services.NewCalculatorServer())

	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
