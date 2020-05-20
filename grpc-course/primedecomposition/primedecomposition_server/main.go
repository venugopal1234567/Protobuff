package main

import (
	"fmt"
	"log"
	"net"
	"proto/grpc-course/primedecomposition/primedecompositionpb"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) PrimeDecomposeManyTimes(req *primedecompositionpb.PrimeDecomposeManyTimesRequest, stream primedecompositionpb.PrimeDecomposeService_PrimeDecomposeManyTimesServer) error {
	fmt.Printf("PrimeDecomposeManyTimes function got invoked %v", req)
	primeNumber := req.GetPrimedecompose().GetNum()
	var k int32
	k = 2
	for primeNumber > 1 {
		if primeNumber%k == 0 { // if k evenly divides into N
			res := &primedecompositionpb.PrimeDecomposeManyTimesResponse{
				Result: k,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			primeNumber = primeNumber / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %w", err)
	}

	s := grpc.NewServer()

	primedecompositionpb.RegisterPrimeDecomposeServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
