package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/it-chain/tesseract/pb"
	"github.com/it-chain/yggdrasill/transaction"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

/* Mock Server
--------------------*/
type MockServer struct {
}

func (s *MockServer) RunICode(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return &pb.Response{[]byte("result")}, nil
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
	cs, _ := NewDefaultRpcClient("127.0.0.1:50001")
	fmt.Println(cs)
}

func TestRunICode(t *testing.T) {
	port := ":50002"

	ms := &MockServer{}
	server, lis := ListenMockServer(ms, port)

	defer func() {
		server.Stop()
		lis.Close()
	}()

	cs, err := NewDefaultRpcClient("127.0.0.1" + port)

	assert.NoError(t, err)

	tx, err := json.Marshal(transaction.DefaultTransaction{ID: "123"})

	res, err := cs.RunICode(&pb.Request{Tx: tx})

	log.Println(res)
	assert.NoError(t, err)

	assert.Equal(t, "test", string(res.Result))

	log.Println("Success")
}
