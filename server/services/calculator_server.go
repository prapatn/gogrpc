package services

import (
	context "context"
	"fmt"
)

type calculatorServer struct {
}

func NewCalculatorServer() CalculatorServer {
	return calculatorServer{}
}

func (obj calculatorServer) mustEmbedUnimplementedCalculatorServer() {

}

func (obj calculatorServer) Hello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	result := fmt.Sprintf("Hello %v at %v", req.Name, req.CreatedDate.AsTime().Local())

	res := HelloResponse{
		Result: result,
	}

	return &res, nil
}
