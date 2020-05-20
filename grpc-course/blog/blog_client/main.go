package main

import (
	"context"
	"fmt"
	"io"
	"log"
	blogpb "proto/grpc-course/blog/blog_pb"
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
	c := blogpb.NewBlogServiceClient(cc)
	//Create MultipleBlog at a time
	fmt.Println("Adding Many Blogs")
	fmt.Println()
	requests := []*blogpb.AddManyRequest{
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu124",
				AuthorId: "Gopal124",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu125",
				AuthorId: "Gopal125",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu126",
				AuthorId: "Gopal126",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu127",
				AuthorId: "Gopal127",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu128",
				AuthorId: "Gopal128",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu129",
				AuthorId: "Gopal129",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
		&blogpb.AddManyRequest{
			Blog: &blogpb.Blog{
				Id:       "Venu130",
				AuthorId: "Gopal130",
				Title:    "My FIrst Blog Is Awesome",
				Content:  "Content of the first blog",
			},
		},
	}

	stream, addmanyErr := c.AddManyBlog(context.Background())

	if addmanyErr != nil {
		log.Fatalf("error while calling AddMany Blog %v", addmanyErr)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	AddManyres, admanyCloseError := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving for LongGreet: %v", admanyCloseError)
	}
	fmt.Printf("LongGreeting response : %v\n", AddManyres)
	fmt.Println()

	//Create SIngle Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		Id:       "Venu123",
		AuthorId: "Gopal123",
		Title:    "My FIrst Blog Is Awesome",
		Content:  "Content of the first blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error ")
	}
	fmt.Printf("Blog has been Created %v", res.GetBlog())
	fmt.Println()

	//read Blog
	fmt.Println()
	fmt.Println("Reading the Blog")
	readResp, err1 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		Id:       "Venu123",
		AuthorId: "Gopal123",
	})
	if err1 != nil {
		fmt.Printf("Error happened while reading: %v", err)
	}
	fmt.Printf("Book Found is : %v", readResp.GetBlog())
	fmt.Println()

	//Update Blog
	fmt.Println()
	fmt.Println("Updating the blog ...")
	updatedblog := &blogpb.Blog{
		Id:       "Venu123",
		AuthorId: "Gopal123",
		Title:    "My FIrst Blog Is Awesome got Updated",
		Content:  "Content of the first blog is updated to new Content",
	}
	updatedRes, updatEerr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: updatedblog})
	if err != nil {
		log.Fatalf("Unexpected error: %v ", updatEerr)
	}
	fmt.Printf("Blog has been Updated %v", updatedRes.GetBlog())

	//List Blog
	fmt.Println("Listing the blog ...")
	fmt.Println()
	listBlog := &blogpb.ListRequest{}
	listStream, listEerr := c.ListBlog(context.Background(), listBlog)
	if listEerr != nil {
		log.Fatalf("Unexpected error: %v ", listEerr)
	}
	i := 1
	for {
		myblog, err := listStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading : %v", err)
		}
		log.Printf("%d : Response from ListBlog:  %v", i, myblog.GetBlog())
		i++
	}
	fmt.Println()
	fmt.Println("Completed Listing the Blog")
	fmt.Println()

	//Delete Blog
	fmt.Println()
	fmt.Println("Deleting the blog ...")
	deleteblog := &blogpb.DeleteBlogRequest{
		Id:       "Venu123",
		AuthorId: "Gopal123",
	}
	deletedRes, deletedEerr := c.DeleteBlog(context.Background(), deleteblog)
	if err != nil {
		log.Fatalf("Unexpected error: %v ", deletedEerr)
	}
	fmt.Printf("Blog has been Deleted %v", deletedRes.GetId())

}
