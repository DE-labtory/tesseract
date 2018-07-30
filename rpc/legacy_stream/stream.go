package legacy_stream

type Stream interface {
	Send()
	Receive()
}
