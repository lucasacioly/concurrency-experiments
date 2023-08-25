package main

import (
	"context"
	"fmt"
	"net"

	"projeto6/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type primeServer struct {
	pb.UnimplementedPrimeServiceServer
}

func (s *primeServer) SeparatePrimeNumbers(ctx context.Context, req *pb.NumbersRequest) (*pb.NumbersResponse, error) {
	primeNumbers, nonPrimeNumbers := separatePrimeNumbers(req.Numbers)
	return &pb.NumbersResponse{
		PrimeNumbers:    primeNumbers,
		NonPrimeNumbers: nonPrimeNumbers,
	}, nil
}

func isPrime(num int32) bool {
	if num == 1 {
		return false
	}

	for i := 2; i*i <= int(num); i++ {
		if int(num)%i == 0 {
			return false
		}
	}

	return true
}

func separatePrimeNumbers(numbers []int32) ([]int32, []int32) {
	var primeNumbers []int32
	var nonPrimeNumbers []int32

	for _, num := range numbers {
		if isPrime(num) {
			primeNumbers = append(primeNumbers, num)
		} else {
			nonPrimeNumbers = append(nonPrimeNumbers, num)
		}
	}

	return primeNumbers, nonPrimeNumbers
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	server := grpc.NewServer()
	pb.RegisterPrimeServiceServer(server, &primeServer{})
	reflection.Register(server)
	fmt.Println("Starting gRPC server...")
	if err := server.Serve(lis); err != nil {
		fmt.Println("Failed to serve:", err)
	}
}
