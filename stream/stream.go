package grpc

type GrpcServer interface {
	Start()
	Connect()
	SendMsg()
	ReceiveMsg()
}
