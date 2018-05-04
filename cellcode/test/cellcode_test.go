package test

import (
	"context"
	"testing"
	"time"

	"encoding/json"

	"fmt"

	"os"
	"os/exec"

	"github.com/it-chain/tesseract/cellcode/cell"
	"github.com/it-chain/tesseract/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

/* Mock Client
--------------------*/
type MockClient struct {
	address string
	conn    *grpc.ClientConn
	client  pb.DefaultServiceClient
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewMockClient(address string) (*MockClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	client := pb.NewDefaultServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	return &MockClient{
		address: address,
		conn:    conn,
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (c *MockClient) RunICode(request *pb.Request) (*pb.Response, error) {
	return c.client.RunICode(c.ctx, &pb.Request{Tx: request.Tx})
}

func (c *MockClient) Close() {
	c.conn.Close()
	c.cancel()
}

/* Run Server
--------------------*/
func before(port string) {
	GOPATH := os.Getenv("GOPATH")
	cmd := exec.Command("sh", GOPATH+"/src/github.com/it-chain/tesseract/cellcode/test/cellcode_test.sh "+port)
	cmd.Run()
	time.Sleep(3 * time.Second)
}

func after() {
	GOPATH := os.Getenv("GOPATH")
	os.RemoveAll(GOPATH + "/src/github.com/it-chain/tesseract/cellcode/test/wsdb")
	os.RemoveAll(GOPATH + "/src/github.com/it-chain/tesseract/cellcode/test/icode.so")
}

/* Test
--------------------*/
func TestQueryGetA(t *testing.T) {
	port := "50011"
	before(port)

	mc, _ := NewMockClient("127.0.0.1:" + port)
	tx, _ := json.Marshal(cell.TxInfo{
		Method: "query",
		ID:     "123",
		Params: cell.Params{
			Type:     1,
			Function: "getA",
			Args:     []string{""},
		},
	})

	res, err := mc.RunICode(&pb.Request{Tx: tx})
	assert.NoError(t, err)

	m := make(map[string]string)
	err = json.Unmarshal(res.Data, &m)
	assert.NoError(t, err)

	fmt.Println(res)

	after()
}

func TestInvokeInitA(t *testing.T) {
	port := "50011"
	before(port)

	mc, _ := NewMockClient("127.0.0.1:" + port)
	tx, _ := json.Marshal(cell.TxInfo{
		Method: "invoke",
		ID:     "124",
		Params: cell.Params{
			Type:     1,
			Function: "initA",
			Args:     []string{""},
		},
	})

	res, err := mc.RunICode(&pb.Request{Tx: tx})
	assert.NoError(t, err)

	fmt.Println(res)

	after()
}

func TestInvokeIncA(t *testing.T) {
	port := "50011"
	before(port)

	mc, _ := NewMockClient("127.0.0.1:" + port)
	tx, _ := json.Marshal(cell.TxInfo{
		Method: "invoke",
		ID:     "124",
		Params: cell.Params{
			Type:     1,
			Function: "incA",
			Args:     []string{""},
		},
	})

	res, err := mc.RunICode(&pb.Request{Tx: tx})
	assert.NoError(t, err)

	fmt.Println(res)

	after()
}

func TestAHundredTimesQuery(t *testing.T) {
	port := "50011"
	before(port)

	mc, _ := NewMockClient("127.0.0.1:" + port)
	tx, _ := json.Marshal(cell.TxInfo{
		Method: "query",
		ID:     "123",
		Params: cell.Params{
			Type:     1,
			Function: "getA",
			Args:     []string{""},
		},
	})

	startTime := time.Now()

	for i := 0; i < 100; i++ {
		_, err := mc.RunICode(&pb.Request{Tx: tx})
		assert.NoError(t, err)
	}
	endTime := time.Now()

	assert.WithinDuration(t, endTime, startTime, 2*time.Second)

	after()
}
