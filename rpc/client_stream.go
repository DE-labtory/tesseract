package rpc

import (
	"context"

	"log"
	"time"

	"fmt"
	"io"

	"github.com/it-chain/tesseract/pb"
	"google.golang.org/grpc"
)

const (
	defaultDialTimeout = 3 * time.Second
)

type ClientStream struct {
	conn         *grpc.ClientConn
	client       pb.BistreamServiceClient
	clientStream pb.BistreamService_RunICodeClient
	ctx          context.Context
	cancel       context.CancelFunc
	handler      func(response *pb.Response, err error)
}

func NewClientStream(address string) (*ClientStream, error) {
	dialContext, _ := context.WithTimeout(context.Background(), defaultDialTimeout)

	conn, err := grpc.DialContext(dialContext, address, grpc.WithInsecure())
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

func (cs *ClientStream) SetHandler(handler func(response *pb.Response, err error)) {
	cs.handler = handler
}

func (cs *ClientStream) StartHandle() {
	go func() {
		for {
			res, err := cs.clientStream.Recv()
			if err == io.EOF {
				fmt.Println("io.EOF handle finish.")
				return
			}
			if cs.handler == nil {
				log.Fatal("error in start handle. there is no handle")
				return
			}
			cs.handler(res, err)
		}
	}()
}

func (cs *ClientStream) RunICode(request *pb.Request) error {
	return cs.clientStream.Send(request)
}
func (c *ClientStream) Ping() (*pb.Empty, error) {
	return c.client.Ping(c.ctx, &pb.Empty{})
}
