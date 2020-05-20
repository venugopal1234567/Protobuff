package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"proto/grpc-course/greet/greetpb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("I am client")
	certFile := "ssl/ca.crt"
	creds, sslError := credentials.NewClientTLSFromFile(certFile, "")
	if sslError != nil {
		log.Fatalf("Error while loading CA trust certificate: %v", sslError)
		return
	}
	opts := grpc.WithTransportCredentials(creds)
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %w", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	//doBiDiStreaming(c)

	//doErrorsUnary(c)

	//doUnaryWithDeadline(c, 5*time.Second) //should complete

	//doUnaryWithDeadline(c, 2*time.Second) //should fail
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	//fmt.Printf("Created client: %f", c)
	fmt.Println("Starting Unary operation with DeadLine")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Venugopal",
			LastName:  "Hegde",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! DeadLine was exceeded")
			} else {
				fmt.Printf("unexpected error:%v", statusErr)
			}
		} else {
			log.Fatalf("error while calling Greer %v", err)
		}
		return
	}
	log.Printf("Response from Greet: %v", res.Message)
}

func doErrorsUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Do error Unary RPC ....")
	number := int32(10)
	//correct call
	doErrorCall(c, number)
	//do Error Call
	doErrorCall(c, -10)
}

func doErrorCall(c greetpb.GreetServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &greetpb.SquareRootRequest{Number: number})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent negative number")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot %v", err)
		}
	}
	fmt.Printf("Result of square root of %v :  %v", number, res.GetNumberRoot())
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Bidirectional Streaming RPC ....")

	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal1",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal2",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal3",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal4",
			},
		},
	}

	//Create stream by invoking client
	stream, err := c.GreetEveryOne(context.Background())
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
			fmt.Printf("Received : %v\n", res.GetMessage())
		}
		close(waitc)
	}()
	//block until everthing is done
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("CLient Streaming RPC ....")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Venugopal4",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())

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
	fmt.Printf("LongGreeting response : %v\n", res)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Created client: %f", c)
	fmt.Println("Server Streaming RPC ....")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Venugopal",
			LastName:  "Hegede",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
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
		log.Printf("Response from GreetManyTimes:  %v", message.GetMessage())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Printf("Created client: %f", c)
	fmt.Println("Starting Unary operation")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Venugopal",
			LastName:  "Hegde",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greer %v", err)
	}
	log.Printf("Response from Greet: %v", res.Message)
}
