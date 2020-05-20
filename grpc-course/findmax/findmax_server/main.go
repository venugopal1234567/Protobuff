package main

import (
	"fmt"
	"io"
	"log"
	"net"
	findmax "proto/grpc-course/findmax/findmaxpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) FindMax(stream findmax.FindMaxService_FindMaxServer) error {
	fmt.Println("FindMax function got invoked with stream")
	var prevNumber int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("There is an errror while receiving: %w", err)
		}
		number := req.GetNumber().GetNum()
		if number > prevNumber {
			sendErr := stream.Send(&findmax.FindMaxResponse{
				Result: number,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending to client: %v", err)
				return err
			}
		}
		prevNumber = number
	}
}

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	findmax.RegisterFindMaxServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
