package legacy_stream

/*

import (
	"context"
	"log"

	"fmt"

	"github.com/it-chain/tesseract/pb"
	"google.golang.org/grpc"
)

type DefaultClientStream struct {
	address string
	port    string
	//pb := stream_pb

	client pb.StreamServiceClient
	stream pb.StreamService_StreamClient
}

func NewDefaultClientStream(address string, port string) *DefaultClientStream {
	fmt.Println(address + ":" + port)
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	return &DefaultClientStream{
		address: address,
		port:    port,
		client:  client,
	}
}

func (c *DefaultClientStream) Connect() error {
	ctx, _ := context.WithCancel(context.Background())
	stream, err := c.client.Stream(ctx)
	if err != nil {
		return err
	}

	c.stream = stream


	err = stream.Send(&pb.Request{"testset"})

	return err
}

func (c *DefaultClientStream) SendRequest(request *pb.Request) error {
	return c.stream.Send(request)
}

func (c *DefaultClientStream) SendResponse() (response *pb.Response, err error) {
	return c.stream.Recv()
}
*/
