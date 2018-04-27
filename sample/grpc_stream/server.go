package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "./proto_stream"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	Name string
}

func (s *server) SayHello(stream pb.Greeter_SayHelloServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(in.Name)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
