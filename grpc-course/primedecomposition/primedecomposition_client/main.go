package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"proto/grpc-course/primedecomposition/primedecompositionpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %w", err)
	}
	defer cc.Close()
	c := primedecompositionpb.NewPrimeDecomposeServiceClient(cc)

	doServerStreaming(c)
}

func doServerStreaming(c primedecompositionpb.PrimeDecomposeServiceClient) {
	fmt.Printf("Created client: %f", c)
	fmt.Println("Server Streaming RPC ....")
	req := &primedecompositionpb.PrimeDecomposeManyTimesRequest{
		Primedecompose: &primedecompositionpb.PrimeDecompose{
			Num: 12,
		},
	}
	resStream, err := c.PrimeDecomposeManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Server %v", err)
	}
	for {
		message, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading : %v", err)
		}
		log.Printf("Response from GreetManyTimes:  %v", message.GetResult())
	}
}
