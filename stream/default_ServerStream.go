package stream

import (
	"context"
	"log"
	"net"

	pb "../message"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DefaultGrpcServer struct {
	port string
}

func NewDefaultGrpcServer(port string) *DefaultGrpcServer {

	if string(port[0]) != ":" {
		port = ":" + port
	}
	return &DefaultGrpcServer{
		port: port,
	}
}

func (s *DefaultGrpcServer) Listen() {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterGrpcMessageServer(server, &DefaultGrpcServer{})
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *DefaultGrpcServer) Connect(ctx context.Context, in *pb.ConnectionRequest) (*pb.ConnectionReply, error) {
	return &pb.ConnectionReply{"test"}, nil
}
