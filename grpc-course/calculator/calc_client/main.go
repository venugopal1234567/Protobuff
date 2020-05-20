package main

import (
	"context"
	"fmt"
	"log"
	"proto/grpc-course/calculator/calcpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %w", err)
	}
	defer cc.Close()
	c := calcpb.NewCalcServiceClient(cc)

	doUnary(c)
}

func doUnary(c calcpb.CalcServiceClient) {
	fmt.Printf("Created client: %f", c)
	fmt.Println("Starting Unary operation")
	req := &calcpb.SumRequest{
		Sum: &calcpb.Sum{
			A: 20,
			B: 30,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greer %v", err)
	}
	log.Printf("Response from Calculator: %v", res.Result)
}
