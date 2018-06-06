package rpc

import (
	"context"
	"time"

	"encoding/json"

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
	txBytes, err := json.Marshal(request.Tx)
	if err != nil {
		return nil, err
	}
	return c.client.RunICode(c.ctx, &pb.Request{Tx: txBytes})
}

func (c *DefaultRpcClient) Close() {
	c.conn.Close()
	c.cancel()
}
