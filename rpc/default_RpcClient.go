package rpc

import (
	"context"
	"log"
	"time"

	"fmt"

	"github.com/it-chain/tesseract/pb"
	"google.golang.org/grpc"
)

type DefaultRpcClient struct {
	address string
	conn    *grpc.ClientConn
	client  pb.DefaultServiceClient
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewDefaultRpcClient(address string) *DefaultRpcClient {
	fmt.Println(address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewDefaultServiceClient(conn)

	return &DefaultRpcClient{
		address: address,
		conn:    conn,
		client:  client,
	}
}

func (c *DefaultRpcClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	c.ctx = ctx
	c.cancel = cancel

	return nil
}

func (c *DefaultRpcClient) RunICode(request *pb.Request) (*pb.Response, error) {
	return c.client.RunICode(c.ctx, &pb.Request{request.Test})
}

func (c *DefaultRpcClient) Close() {
	c.conn.Close()
	c.cancel()
}
