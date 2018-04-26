package rpc

import (
	"fmt"
	"testing"

	"context"
	"log"
	"net"

	"github.com/it-chain/tesseract/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

/* Mock Server
--------------------*/
type MockServer struct {
}

func (s *MockServer) RunICode(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return &pb.Response{request.Test}, nil
}

func ListenMockServer(ms *MockServer, port string) (*grpc.Server, net.Listener) {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterDefaultServiceServer(server, ms)
	reflection.Register(server)

	go func() {
		fmt.Println("[MockServer] Listen")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return server, lis
}

/* Test
--------------------*/
func TestNewDefaultRpcClient(t *testing.T) {
	cs := NewDefaultRpcClient("127.0.0.1:50001")
	fmt.Println(cs)
}

func TestConnect(t *testing.T) {
	port := ":50001"

	ms := &MockServer{}
	server, lis := ListenMockServer(ms, port)

	defer func() {
		server.Stop()
		lis.Close()
	}()

	cs := NewDefaultRpcClient("127.0.0.1" + port)
	err := cs.Connect()

	assert.NoError(t, err)
}

func TestRunICode(t *testing.T) {
	port := ":50002"

	ms := &MockServer{}
	server, lis := ListenMockServer(ms, port)

	defer func() {
		server.Stop()
		lis.Close()
	}()

	cs := NewDefaultRpcClient("127.0.0.1" + port)
	err := cs.Connect()
	log.Println(cs)

	assert.NoError(t, err)

	res, err := cs.RunICode(&pb.Request{"TestRunICode"})

	log.Println(res)
	assert.NoError(t, err)

	assert.Equal(t, "TestRunICode", res.Test)

	log.Println("Success")
}
