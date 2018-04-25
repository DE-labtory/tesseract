package stream

import (
	"context"
	"log"

	"github.com/it-chain/tesseract/pb"
	"google.golang.org/grpc"
)

type DefaultClientStream struct {
	address string
	port    string

	client pb.StreamServiceClient
	stream pb.StreamService_StreamClient
}

func NewDefaultClientStream(address string, port string) *DefaultClientStream {
	conn, err := grpc.Dial(address+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewStreamServiceClient(conn)

	return &DefaultClientStream{
		address: address,
		port:    port,
		client:  client,
	}
}

func (c *DefaultClientStream) Connect() error {
	stream, err := c.client.Stream(context.Background())
	if err != nil {
		return err
	}

	c.stream = stream

	/*
		waitc := make(chan struct{})
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					log.Fatalf("close")
					//close(waitc)
					//return
				}
				if err != nil {
					log.Fatalf("Failed to receive a note : %v", err)
				}
				fmt.Println(in)
			}
		}()

		if err := stream.Send(&pb.Request{"request name"}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
		for {
			time.Sleep(10 * time.Second)
		}
		stream.CloseSend()
		<-waitc
	*/

	return nil
}

func (c *DefaultClientStream) SendRequest(request *pb.Request) error {
	return c.stream.Send(request)
}

func (c *DefaultClientStream) SendResponse() (response *pb.Response, err error) {
	return c.stream.Recv()
}
