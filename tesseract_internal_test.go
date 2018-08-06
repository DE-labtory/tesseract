package tesseract

import (
	"os"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestGetAvailablePort(t *testing.T) {
	GOPATH := os.Getenv("GOPATH")
	var setup = func() (*Tesseract, func()) {
		te := New()

		return te, func() {

			t.Log("container is closing")
			te.StopContainers()
		}
	}

	tesseract, tearDown := setup()

	defer tearDown()
	code1 := ICodeInfo{
		Directory: GOPATH + "/src/github.com/it-chain/tesseract/docker/mock",
	}
	code2 := ICodeInfo{
		Directory: GOPATH + "/src/github.com/it-chain/tesseract/docker/mock",
	}
	code3 := ICodeInfo{
		Directory: GOPATH + "/src/github.com/it-chain/tesseract/docker/mock",
	}
	fmt.Println("code 1 test")
	_, err := tesseract.SetupContainer(code1)
	assert.NoError(t, err)
	fmt.Println("code 2 test")
	_, err = tesseract.SetupContainer(code2)
	assert.NoError(t, err)
	fmt.Println("code 3 test")
	_, err = tesseract.SetupContainer(code3)
	assert.NoError(t, err)
}
