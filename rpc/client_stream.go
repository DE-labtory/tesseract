package rpc

import (
	"context"
	"io"
	"time"

	"github.com/DE-labtory/iLogger"
	"github.com/DE-labtory/tesseract"
	"github.com/DE-labtory/tesseract/pb"
	"google.golang.org/grpc"
)

const (
	defaultDialTimeout = 3 * time.Second
)

type ClientStream struct {
	conn         *grpc.ClientConn
	client       pb.BistreamServiceClient
	clientStream pb.BistreamService_RunICodeClient
	ctx          context.Context
	cancel       context.CancelFunc
	Handler      *DefaultHandler
}

func NewClientStream(address string) (*ClientStream, error) {
	dialContext, _ := context.WithTimeout(context.Background(), defaultDialTimeout)

	conn, err := grpc.DialContext(dialContext, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	ctx, cf := context.WithCancel(context.Background())
	client := pb.NewBistreamServiceClient(conn)
	clientStream, err := client.RunICode(ctx)
	if err != nil {
		//conn.Close()
		//cf()
		return nil, err
	}

	return &ClientStream{
		conn:         conn,
		client:       client,
		clientStream: clientStream,
		ctx:          ctx,
		cancel:       cf,
	}, nil
}

func (cs *ClientStream) SetHandler(handler *DefaultHandler) {
	cs.Handler = handler
}

func (cs *ClientStream) StartHandle() {
	go func() {
		for {
			res, err := cs.clientStream.Recv()
			if err == io.EOF || res == nil {
				iLogger.Info(nil, "[Tesseract] client stream finish")
				return
			}
			if cs.Handler == nil {
				iLogger.Fatal(nil, "[Tesseract] error in start handle. there is no handle")
				return
			}
			cs.Handler.Handle(res, err)
		}
	}()
}

func (cs *ClientStream) RunICode(request *pb.Request, callBack tesseract.CallBack) error {
	cs.Handler.AddCallback(request.Uuid, callBack)
	return cs.clientStream.Send(request)
}

func (c *ClientStream) Ping() (*pb.Empty, error) {
	return c.client.Ping(c.ctx, &pb.Empty{})
}

func (c *ClientStream) Close() {

	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		c.conn.Close()
	}
}

type DefaultHandler struct {
	callBacks map[string]tesseract.CallBack
}

func NewDefaultHandler() *DefaultHandler {

	return &DefaultHandler{
		callBacks: make(map[string]tesseract.CallBack),
	}
}

func (d *DefaultHandler) Handle(response *pb.Response, err error) {
	callbackFunc := d.callBacks[response.Uuid]

	if callbackFunc == nil {
		iLogger.Panicf(nil, "[Tesseract] error in handle uuid : %s", response.Uuid)
	}

	callbackFunc(response, err)
	delete(d.callBacks, response.Uuid)
}

func (d *DefaultHandler) AddCallback(uuid string, callback tesseract.CallBack) {
	d.callBacks[uuid] = callback
}
