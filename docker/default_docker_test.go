package docker

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"

	"github.com/it-chain/tesseract"
	"github.com/stretchr/testify/assert"
)

func TestCreateContainerWithCellCode(t *testing.T) {
	GOPATH := os.Getenv("GOPATH")
	res, err := CreateContainerWithCellCode(
		DockerImage{DefaultImageName, DefaultImageTag},
		tesseract.ICodeInfo{"icode", GOPATH + "/src/github.com/it-chain/tesseract/test/icode_test"},
		GOPATH+"/src/github.com/it-chain/tesseract/sh/default_setup.sh",
		"50001",
	)
	assert.NoError(t, err)

	log.Print(res)
}

func TestStartContainer(t *testing.T) {
	GOPATH := os.Getenv("GOPATH")
	res, err := CreateContainerWithCellCode(
		DockerImage{DefaultImageName, DefaultImageTag},
		tesseract.ICodeInfo{"icode", GOPATH + "/src/github.com/it-chain/tesseract/test/icode_test"},
		GOPATH+"/src/github.com/it-chain/tesseract/sh/default_setup.sh",
		"50003",
	)

	err = StartContainer(res)
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	//_, err = os.Stat("../cellcode/query")
	//assert.NoError(t, err)

	/*	defer func() {
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
	}()*/
}

func TestPullImage(t *testing.T) {
	err := PullImage(DefaultImageName + ":" + DefaultImageTag)
	assert.NoError(t, err)
}

func TestHasImageWhenImageExist(t *testing.T) {

	//given
	image := DefaultImageName + ":" + DefaultImageTag
	err := PullImage(image)
	assert.NoError(t, err)

	//when
	flag, err := HasImage(image)
	assert.NoError(t, err)

	//then
	assert.True(t, flag)

	defer func() {
		ctx := context.Background()
		cli, err := docker.NewEnvClient()
		assert.NoError(t, err)
		_, err = cli.ImageRemove(ctx, image, types.ImageRemoveOptions{})
		assert.NoError(t, err)
	}()
}

func TestHasImageWhenImageDoesNotExist(t *testing.T) {

	//given
	image := DefaultImageName + ":" + DefaultImageTag
	removeImage(image)

	//when
	flag, err := HasImage(image)
	assert.NoError(t, err)

	//then
	assert.False(t, flag)
}

func removeImage(image string) error {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()

	if err != nil {
		return err
	}

	_, err = cli.ImageRemove(ctx, image, types.ImageRemoveOptions{})

	if err != nil {
		return err
	}

	return nil
}
