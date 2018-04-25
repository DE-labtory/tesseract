package stream

type Stream interface {
	Send()
	Receive()
}
