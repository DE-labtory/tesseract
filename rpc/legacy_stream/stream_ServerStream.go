package legacy_stream

/*
import (
	"fmt"
	"io"

	"github.com/it-chain/tesseract/pb"
)

type DefaultServerStream struct {
	Port    string
	Handler func()
}

func NewDefaultServerStream(port string, handler func()) *DefaultServerStream {
	return &DefaultServerStream{
		Port:    port,
		Handler: handler,
	}
}

func (s *DefaultServerStream) Stream(stream pb.StreamService_StreamServer) error {
	fmt.Println("in Stream")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(in)

		s.Handler()
	}

	return nil
}
*/
