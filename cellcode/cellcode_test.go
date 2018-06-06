package main

import (
	"context"
	"testing"
	"time"

	"encoding/json"

	"fmt"

	"os"
	"os/exec"

	"syscall"

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
//cmd로 cellcode를 실행시킬 경우 2개의 process가 생성된다. grpc에서 process를 독립적으로 만들기 때문에 cellcode를 종료하더라도 grpc는 남아있다.
//따라서 pgid로 kill해야 모두 종료된다.
func SetupTest(t *testing.T, port string) func() {

	t.Log("before")

	GOPATH := os.Getenv("GOPATH")
	cmd := exec.Command("go", "build", "-buildmode=plugin",
		"-o", GOPATH+"/src/github.com/it-chain/tesseract/cellcode/mock/tmp/icode.so",
		GOPATH+"/src/github.com/it-chain/tesseract/cellcode/mock/icode/icode.go")
	err := cmd.Run()
	assert.NoError(t, err)

	cmd2 := exec.Command("go", "run",
		GOPATH+"/src/github.com/it-chain/tesseract/cellcode/cellcode.go",
		GOPATH+"/src/github.com/it-chain/tesseract/cellcode/mock/tmp/icode.so", port)

	//pgid를 set(한번에 kill하기 위해)
	cmd2.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err = cmd2.Start()
	assert.NoError(t, err)

	//sh를 이용한 build시간 확보
	time.Sleep(3 * time.Second)

	return func() {
		// t is from the outer SetupTest scope
		t.Log("after")

		//pgid를 사용해서 kill
		if err := syscall.Kill(-cmd2.Process.Pid, syscall.SIGINT); err != nil {
			t.Fatal("failed to kill process: ", err)
		}

		os.RemoveAll("./wsdb")
	}
}

/* Test
--------------------*/
func TestInvokeInitA(t *testing.T) {
	port := "50011"
	defer SetupTest(t, port)()

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
	assert.Equal(t, res.Result, "Success")
}

func TestQueryGetA(t *testing.T) {
	port := "50011"
	defer SetupTest(t, port)()

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

	tx, _ = json.Marshal(cell.TxInfo{
		Method: "query",
		ID:     "123",
		Params: cell.Params{
			Type:     1,
			Function: "getA",
			Args:     []string{""},
		},
	})

	res, err = mc.RunICode(&pb.Request{Tx: tx})
	assert.NoError(t, err)

	m := make(map[string]string)
	err = json.Unmarshal(res.Data, &m)
	assert.NoError(t, err)
	assert.Equal(t, m["A"], "0")
}

func TestInvokeIncA(t *testing.T) {

	port := "50011"
	defer SetupTest(t, port)()

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
}

func TestAHundredTimesQuery(t *testing.T) {
	port := "50011"
	defer SetupTest(t, port)()

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
}
