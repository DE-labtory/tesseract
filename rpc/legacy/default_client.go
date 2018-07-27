package rpc

import (
	"context"
	"time"

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

func Connect(address string) (*DefaultRpcClient, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	client := pb.NewDefaultServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	return &DefaultRpcClient{
		address: address,
		conn:    conn,
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

//todo test request -> request
func (c *DefaultRpcClient) RunICode(request *pb.Request) (*pb.Response, error) {
	return c.client.RunICode(c.ctx, request)
}

func (c *DefaultRpcClient) Ping() (*pb.Empty, error) {
	return c.client.Ping(c.ctx, &pb.Empty{})
}

func (c *DefaultRpcClient) Close() {
	c.conn.Close()
	c.cancel()
}
