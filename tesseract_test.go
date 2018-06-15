package tesseract_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/cellcode/cell"
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
				Directory: GOPATH + "/src/github.com/it-chain/tesseract/cellcode/mock/icode/",
			},
			output: "123",
			err:    nil,
		},
	}

	var setup = func(config tesseract.Config) (*tesseract.Tesseract, func()) {
		te := tesseract.New(config)

		return te, func() {

			t.Log("container is closing")
			te.StopContainer()
		}
	}

	tesseract, tearDown := setup(tesseract.Config{ShPath: GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"})

	defer tearDown()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		id, err := tesseract.SetupContainer(test.input)
		assert.Equal(t, err, test.err)

		defer docker.CloseContainer(id)

		for _, client := range tesseract.Clients {

			tx, _ := json.Marshal(cell.TxInfo{
				Method: "invoke",
				ID:     "124",
				Params: cell.Params{
					Type:     1,
					Function: "initA",
					Args:     []string{""},
				},
			})

			_, err := client.RunICode(&pb.Request{Tx: tx})
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
