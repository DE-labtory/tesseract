package rpc

import (
	"context"

	"github.com/it-chain/tesseract/pb"
)

type DefaultRpcServer struct {
	Port    string
	Handler func()
}

func NewDefaultRpcServer(port string, handler func()) *DefaultRpcServer {
	return &DefaultRpcServer{
		Port:    port,
		Handler: handler,
	}
}

func (s *DefaultRpcServer) RunICode(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return &pb.Response{request.Test}, nil
}
