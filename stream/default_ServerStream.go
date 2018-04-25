package stream

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/it-chain/tesseract/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type DefaultServerStream struct {
	port    string
	handler func()
}

func NewDefaultServerStream(port string) *DefaultServerStream {

	if string(port[0]) != ":" {
		port = ":" + port
	}
	return &DefaultServerStream{
		port: port,
	}
}

func (s *DefaultServerStream) Listen(handler func()) {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &DefaultServerStream{})
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *DefaultServerStream) Stream(stream pb.StreamService_StreamServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(in)

		s.handler()
	}

	return nil
}
