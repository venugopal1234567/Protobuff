package main

import (
	"context"
	"fmt"
	"io"
	"log"
	findmax "proto/grpc-course/findmax/findmaxpb"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I am client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %w", err)
	}
	defer cc.Close()
	c := findmax.NewFindMaxServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c findmax.FindMaxServiceClient) {
	fmt.Println("Bidirectional Streaming RPC ....")

	requests := []*findmax.FindMaxRequest{
		&findmax.FindMaxRequest{
			Number: &findmax.Number{
				Num: 1,
			},
		},
		&findmax.FindMaxRequest{
			Number: &findmax.Number{
				Num: 5,
			},
		},
		&findmax.FindMaxRequest{
			Number: &findmax.Number{
				Num: 3,
			},
		},
		&findmax.FindMaxRequest{
			Number: &findmax.Number{
				Num: 6,
			},
		},
		&findmax.FindMaxRequest{
			Number: &findmax.Number{
				Num: 2,
			},
		},
		&findmax.FindMaxRequest{
			Number: &findmax.Number{
				Num: 20,
			},
		},
	}

	//Create stream by invoking client
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error while creating : %v", err)
		return
	}

	waitc := make(chan struct{})
	//send buch of message to server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	//Receive bunch of message from server
	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received : %v\n", res.GetResult())
		}
		close(waitc)
	}()
	//block until everthing is done
	<-waitc
}
