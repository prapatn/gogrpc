package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalculatorService interface {
	Hello(name string) error
	Fibonacci(n uint32) error
	Average(number ...float64) error
	Sum(numbers ...int32) error
}

type calculatorService struct {
	calculatorClient CalculatorClient
}

func NewCalculatorService(calculatorClient CalculatorClient) CalculatorService {
	return calculatorService{calculatorClient: calculatorClient}
}

func (base calculatorService) Hello(name string) error {
	req := HelloRequest{
		Name:        name,
		CreatedDate: timestamppb.Now(),
	}
	res, err := base.calculatorClient.Hello(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: Hello\n")
	fmt.Printf("Request: %v\n", req.Name)
	fmt.Printf("Response: %v\n", res.Result)

	return nil
}

func (base calculatorService) Fibonacci(n uint32) error {
	req := FibonacciRequest{
		N: n,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := base.calculatorClient.Fibonacci(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service: Fibonacci\n")
	fmt.Printf("Request: %v\n", req.N)
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Printf("Response: %v\n", response.Result)
	}

	return nil
}

func (base calculatorService) Average(numbers ...float64) error {

	stream, err := base.calculatorClient.Average(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Service: Average\n")

	for _, n := range numbers {
		req := AverageRequest{
			Number: n,
		}
		stream.Send(&req)
		time.Sleep(time.Second)
		fmt.Printf("Request: %v\n", n)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	fmt.Printf("Response: %v\n", res.Result)
	return nil
}

func (base calculatorService) Sum(numbers ...int32) error {
	stream, err := base.calculatorClient.Sum(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Service: Sum\n")

	//Request
	go func() {
		for _, n := range numbers {
			req := SumRequest{
				Number: n,
			}
			stream.Send(&req)
			fmt.Printf("Request: %v\n", n)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	done := make(chan bool)
	errs := make(chan error)
	//Response
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errs <- err
			}
			fmt.Printf("Response: %v\n", res.Result)
		}
		done <- true

	}()

	select {
	case <-done:
		return nil
	case <-errs:
		return err
	}
}
