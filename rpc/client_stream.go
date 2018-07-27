package rpc

import (
	"context"

	"github.com/it-chain/tesseract/pb"
	"google.golang.org/grpc"
)

type ClientStream struct {
	conn         *grpc.ClientConn
	client       pb.BistreamServiceClient
	clientStream pb.BistreamService_RunICodeClient
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewClientStream(address string) (*ClientStream, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	ctx, cf := context.WithCancel(context.Background())
	client := pb.NewBistreamServiceClient(conn)
	clientStream, err := client.RunICode(ctx)
	if err != nil {
		return nil, err
	}

	return &ClientStream{
		conn:         conn,
		client:       client,
		clientStream: clientStream,
		ctx:          ctx,
		cancel:       cf,
	}, nil
}

func (cs *ClientStream) RunIcode(request *pb.Request) error {
	return cs.clientStream.Send(request)
}
