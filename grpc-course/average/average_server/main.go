package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"proto/grpc-course/average/averagepb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Average(stream averagepb.AverageService_AverageServer) error {
	fmt.Println("Average function got invoked with stream")
	var result int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&averagepb.AverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error While reading client: %w", err)
		}
		number := req.GetAverage().GetNumber()
		result += number
	}
}
func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	averagepb.RegisterAverageServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
