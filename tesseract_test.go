package tesseract_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"testing"

	"fmt"
	"time"

	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/docker"
	"github.com/it-chain/tesseract/pb"
	"github.com/stretchr/testify/assert"
)

func TestSetupContainer(t *testing.T) {

	//given
	GOPATH := os.Getenv("GOPATH")
	tests := map[string]struct {
		input  tesseract.ICodeInfo
		output string
		err    error
	}{
		"success": {
			input: tesseract.ICodeInfo{
				Directory: GOPATH + "/src/github.com/it-chain/tesseract/docker/mock",
			},
			output: "123",
			err:    nil,
		},
	}

	var setup = func() (*tesseract.Tesseract, func()) {
		te := tesseract.New()

		return te, func() {
			t.Log("container is closing")
			te.StopContainers()
		}
	}

	tesseract, tearDown := setup()

	defer tearDown()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		id, err := tesseract.SetupContainer(test.input)
		assert.Equal(t, err, test.err)

		defer docker.CloseContainer(id)
		for _, client := range tesseract.Clients {
			err := client.RunICode(&pb.Request{
				Uuid:         "1",
				Type:         "invoke",
				FunctionName: "initA",
				Args:         nil,
			}, func(response *pb.Response, err error) {
				if err != nil {
					fmt.Println("err in init A")
				}
				fmt.Println("res initA!")
			})
			err = client.RunICode(&pb.Request{
				Uuid:         "2",
				Type:         "invoke",
				FunctionName: "incA",
				Args:         nil,
			}, func(response *pb.Response, err error) {
				if err != nil {
					fmt.Println("err in inc A")
				}
				fmt.Println("res incA!")
			})
			err = client.RunICode(&pb.Request{
				Uuid:         "3",
				Type:         "invoke",
				FunctionName: "incA",
				Args:         nil,
			}, func(response *pb.Response, err error) {
				if err != nil {
					fmt.Println("err in inc A")
				}
				fmt.Println("res incA!")
			})
			fmt.Println("wait for 2 min")
			time.Sleep(120 * time.Second)
			err = client.RunICode(&pb.Request{
				Uuid:         "4",
				Type:         "query",
				FunctionName: "getA",
				Args:         nil,
			}, func(response *pb.Response, err error) {
				if err != nil {
					fmt.Println("err in get A")
				}
				fmt.Println("res getA! : ", string(response.Data))
			})
			time.Sleep(2 * time.Second)
			assert.NoError(t, err)
		}

	}
}

func clearContainer() {
	c1 := exec.Command("docker", "ps", "-a", "-f", "ancestor=golang:1.9", "-q")
	c2 := exec.Command("xargs", "-I", "{}", "docker", "rm", "{}")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
}
