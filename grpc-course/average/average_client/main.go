package main

import (
	"context"
	"fmt"
	"log"
	"proto/grpc-course/average/averagepb"
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
	c := averagepb.NewAverageServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	doClientStreaming(c)
}

func doClientStreaming(c averagepb.AverageServiceClient) {
	fmt.Println("CLient Streaming RPC ....")

	requests := []*averagepb.AverageRequest{
		&averagepb.AverageRequest{
			Average: &averagepb.Average{
				Number: 2,
			},
		},
		&averagepb.AverageRequest{
			Average: &averagepb.Average{
				Number: 3,
			},
		},
		&averagepb.AverageRequest{
			Average: &averagepb.Average{
				Number: 4,
			},
		},
		&averagepb.AverageRequest{
			Average: &averagepb.Average{
				Number: 5,
			},
		},
		&averagepb.AverageRequest{
			Average: &averagepb.Average{
				Number: 6,
			},
		},
	}
	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("error while calling Long Greet %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving for LongGreet: %v", err)
	}
	fmt.Printf("The average is : %v\n", res)

}
