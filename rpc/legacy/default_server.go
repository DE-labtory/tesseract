package legacy

import (
	"context"

	"encoding/json"

	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
)

type DefaultRpcServer struct {
	Port    string
	Handler func(*cell.TxInfo) pb.Response
}

func NewDefaultRpcServer(port string, handler func(*cell.TxInfo) pb.Response) *DefaultRpcServer {
	return &DefaultRpcServer{
		Port:    port,
		Handler: handler,
	}
}

func (s *DefaultRpcServer) RunICode(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	tx := cell.TxInfo{}
	json.Unmarshal(request.Tx, &tx)
	res := s.Handler(&tx)
	return &res, nil
}

func (s *DefaultRpcServer) Ping(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
