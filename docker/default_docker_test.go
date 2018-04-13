package docker

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
	//"docker.io/go-docker/api/types/container"
	"os"
	"bytes"
	"os/exec"
	"io"
	"time"
)

func TestCreateContainerWithCellCode(t *testing.T) {
	GOPATH := os.Getenv("GOPATH")
	res, err := CreateContainerWithCellCode(ICodeInfo{"icode", GOPATH + "/src/github.com/it-chain/tesseract/test/icode_test"}, "/Users/hackurity/go/src/github.com/it-chain/tesseract/cellcode")
	assert.NoError(t, err)

	log.Print(res)
}

func TestStartContainer(t *testing.T) {
	GOPATH := os.Getenv("GOPATH")
	res, err := CreateContainerWithCellCode(ICodeInfo{"icode", GOPATH + "/src/github.com/it-chain/tesseract/test/icode_test"}, "/Users/hackurity/go/src/github.com/it-chain/tesseract/cellcode")

	err = StartContainer(res)
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	_, err = os.Stat("../cellcode/query")
	assert.NoError(t, err)


	defer func() {
		// Remove Test Docker Container
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

		// Remove Success File(Query) that created by icode
		os.Remove("../cellcode/query")
	}()
}
