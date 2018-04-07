package socket

type Socket interface {
	Start()
	Connect()
	SendMsg()
	ReceiveMsg()
}
