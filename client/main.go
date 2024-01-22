package main

import (
	"client/services"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	var cc *grpc.ClientConn
	var err error

	var creds credentials.TransportCredentials

	host := flag.String("host", "localhost:8000", "gRPC server host")
	tls := flag.Bool("tls", false, "use a secure TLS connection")
	flag.Parse()

	if *tls {
		certFile := "../tls/ca.crt"
		creds, err = credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		creds = insecure.NewCredentials()

	}

	cc, err = grpc.Dial(*host, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := services.NewCalculatorClient(cc)
	service := services.NewCalculatorService(client)

	// err = service.Hello("K")
	// err = service.Fibonacci(6)
	// err = service.Average(9, 2, 3, 4, 5)
	err = service.Sum(9, 2, 3, 4, 5)

	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			log.Printf("[%v] %v", grpcErr.Code(), grpcErr.Message())
		} else {
			log.Fatal(err)
		}
	}
}
