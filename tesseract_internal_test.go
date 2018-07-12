package tesseract

import (
	"os"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestGetAvailablePort(t *testing.T) {
	GOPATH := os.Getenv("GOPATH")
	var setup = func(config Config) (*Tesseract, func()) {
		te := New(config)

		return te, func() {

			t.Log("container is closing")
			te.StopContainer()
		}
	}

	tesseract, tearDown := setup(Config{ShPath: GOPATH + "/src/github.com/it-chain/tesseract/sh/default_setup.sh"})

	defer tearDown()
	code1 := ICodeInfo{
		Directory: GOPATH + "/src/github.com/it-chain/tesseract/cellcode/mock/icode/",
	}
	code2 := ICodeInfo{
		Directory: GOPATH + "/src/github.com/it-chain/tesseract/cellcode/mock/icode/",
	}
	code3 := ICodeInfo{
		Directory: GOPATH + "/src/github.com/it-chain/tesseract/cellcode/mock/icode/",
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
