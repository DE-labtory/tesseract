package rpc

type Stream interface {
	Send()
	Receive()
}
