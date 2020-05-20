package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	blogpb "proto/grpc-course/blog/blog_pb"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var dynamosess *session.Session

type server struct{}

type blogItem struct {
	Id       string
	AuthorId string
	Content  string
	Title    string
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Processing Create blog request : %v", req)
	svc := dynamodb.New(dynamosess)
	blog := req.GetBlog()
	myblog := blogItem{
		Id:       blog.GetId(),
		AuthorId: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}
	av, err := dynamodbattribute.MarshalMap(myblog)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}
	tableName := "Blog"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err1 := svc.PutItem(input)
	if err1 != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err1))
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blog.GetId(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Processing Read blog request : %v", req)
	blog := blogItem{}
	tableName := "Blog"
	Id := req.GetId()
	AuthorId := req.GetAuthorId()
	svc := dynamodb.New(dynamosess)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(Id),
			},
			"AuthorId": {
				S: aws.String(AuthorId),
			},
		},
	})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &blog)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	if blog.Id == "" {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("NotFound error"))
	}
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blog.Id,
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Processing Update blog request : %v", req)
	tableName := "Blog"
	svc := dynamodb.New(dynamosess)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(req.GetBlog().GetTitle()),
			},
			":c": {
				S: aws.String(req.GetBlog().GetContent()),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(req.GetBlog().GetId()),
			},
			"AuthorId": {
				S: aws.String(req.GetBlog().GetAuthorId()),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Title = :t , Content = :c"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}
	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       req.GetBlog().GetId(),
			AuthorId: req.GetBlog().GetAuthorId(),
			Title:    req.GetBlog().GetTitle(),
			Content:  req.GetBlog().GetContent(),
		},
	}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Processing Delete blog request : %v", req)
	tableName := "Blog"
	svc := dynamodb.New(dynamosess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(req.GetId()),
			},
			"AuthorId": {
				S: aws.String(req.GetAuthorId()),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}
	return &blogpb.DeleteBlogResponse{
		Id:       req.GetId(),
		AuthorId: req.GetAuthorId(),
	}, nil
}

func (*server) ListBlog(req *blogpb.ListRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Printf("BlogService function got invoked with stream: %v", req)
	svc := dynamodb.New(dynamosess)
	tableName := "Blog"
	proj := expression.NamesList(expression.Name("Id"), expression.Name("AuthorId"), expression.Name("Title"), expression.Name("Content"))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}
	result, err := svc.Scan(params)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	for _, i := range result.Items {
		blog := blogItem{}

		err = dynamodbattribute.UnmarshalMap(i, &blog)

		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error: %v", err))
		}
		res := &blogpb.ListBlogResponse{
			Blog: &blogpb.Blog{
				Id:       blog.Id,
				AuthorId: blog.AuthorId,
				Title:    blog.Title,
				Content:  blog.Content,
			},
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) AddManyBlog(stream blogpb.BlogService_AddManyBlogServer) error {
	fmt.Println("AddManyBlog function got invoked with stream")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&blogpb.AddManyResponse{
				Message: "Added Successfully",
			})
		}
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error: %v", err))
		}

		svc := dynamodb.New(dynamosess)
		blog := req.GetBlog()
		myblog := blogItem{
			Id:       blog.GetId(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		}
		av, err := dynamodbattribute.MarshalMap(myblog)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error: %v", err))
		}
		tableName := "Blog"

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		_, err1 := svc.PutItem(input)
		if err1 != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error: %v", err1))
		}

	}
	return nil
}

func main() {
	//If we crash the go code , we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Hello BLog")

	fmt.Println("Connecting to DynamoDB")
	sess, err := session.NewSession(&aws.Config{
		Endpoint:                      aws.String("http://127.0.0.1:8000"),
		Region:                        aws.String("eu-west-2"),
		CredentialsChainVerboseErrors: aws.Bool(true),
		Credentials:                   credentials.NewStaticCredentials("123", "123", ""),
	})
	dynamosess = sess
	if err != nil {
		fmt.Println("Error connecting to client")
	}
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %w", err)
	}

	// opts := []grpc.ServerOption{}
	// tls := false
	// if tls {
	// 	certFile := "ssl/server.crt"
	// 	keyFile := "ssl/server.pem"
	// 	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	// 	if sslErr != nil {
	// 		log.Fatalf("Failed loading certificates: %v", sslErr)
	// 		return
	// 	}
	// 	opts = append(opts, grpc.Creds(creds))
	// }
	//s := grpc.NewServer(opts...)
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s, &server{})
	reflection.Register(s)
	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	//Wait for Control C
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	//Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listner")
	lis.Close()
	fmt.Println("End of Program")
}
