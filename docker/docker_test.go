package docker_test

import (
	"context"
	"os"
	"testing"

	dockerlib "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/docker"
	"github.com/stretchr/testify/assert"
)

func TestCreateContainerWithCellCode(t *testing.T) {

	GOPATH := os.Getenv("GOPATH")
	res, err := docker.CreateContainer(
		tesseract.GetDefaultImage(),
		GOPATH+"/src/github.com/it-chain/tesseract/container/mock",
		"50005",
		tesseract.DefaultLogFileName,
	)
	defer func() {
		// Remove Docker Container
		err := docker.RemoveContainer(res.ID)
		assert.NoError(t, err)

		// need time to remove container
		//time.Sleep(20 * time.Second)
	}()
	assert.NoError(t, err)
}

func TestStartContainer(t *testing.T) {

	//given
	GOPATH := os.Getenv("GOPATH")
	res, err := docker.CreateContainer(
		tesseract.GetDefaultImage(),
		GOPATH+"/src/github.com/it-chain/tesseract/container/mock",
		"50005",
		tesseract.DefaultLogFileName,
	)

	err = docker.StartContainer(res)

	defer func() {
		// Remove Docker Container
		err := docker.KillContainer(res.ID)
		assert.NoError(t, err)

		err = docker.RemoveContainer(res.ID)
		assert.NoError(t, err)

		// need time to remove container
		//time.Sleep(20 * time.Second)
	}()
	assert.NoError(t, err)
}

func TestPullImage(t *testing.T) {

	err := docker.PullImage(tesseract.GetDefaultImage().GetFullName())
	assert.NoError(t, err)
}

func TestHasImageWhenImageExist(t *testing.T) {

	//given
	image := tesseract.GetDefaultImage().GetFullName()
	err := docker.PullImage(image)
	assert.NoError(t, err)

	//when
	flag, err := docker.HasImage(image)
	assert.NoError(t, err)

	//then
	assert.True(t, flag)

	defer func() {
		ctx := context.Background()
		cli, err := dockerlib.NewEnvClient()
		assert.NoError(t, err)

		_, err = cli.ImageRemove(ctx, image, types.ImageRemoveOptions{})
		assert.NoError(t, err)
	}()
}

func TestHasImageWhenImageDoesNotExist(t *testing.T) {

	//given
	image := tesseract.GetDefaultImage().GetFullName()
	removeImage(image)

	//when
	flag, err := docker.HasImage(image)
	assert.NoError(t, err)

	//then
	assert.False(t, flag)
}

func removeImage(image string) error {

	ctx := context.Background()
	cli, err := dockerlib.NewEnvClient()
	defer cli.Close()

	if err != nil {
		return err
	}

	_, err = cli.ImageRemove(ctx, image, types.ImageRemoveOptions{})

	if err != nil {
		return err
	}

	return nil
}
