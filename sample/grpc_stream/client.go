package main

import (
	"context"
	"io"
	"log"

	"fmt"

	"time"

	pb "./proto_stream"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	stream, err := client.SayHello(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Fatalf("close")
				//close(waitc)
				//return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			fmt.Println(in)
		}
	}()

	if err := stream.Send(&pb.HelloRequest{"request name"}); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
	for {
		time.Sleep(10 * time.Second)
	}
	stream.CloseSend()
	<-waitc
}
