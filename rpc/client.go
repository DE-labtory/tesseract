package rpc

import "github.com/it-chain/tesseract/pb"

type Client interface {
	RunICode(request *pb.Request) (*pb.Response, error)
	Close()
}
