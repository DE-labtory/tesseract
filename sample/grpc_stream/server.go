package main

import (
	"log"
	"net"

	pb "./proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(stream *pb.HelloRequest) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)
		... // look for notes to be sent to client
		for _, note := range s.routeNotes[key] {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func (s *server) SayHelloAgain(stream *pb.HelloRequest) error {
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
