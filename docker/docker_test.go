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

type CleanFunc = func() error

func setup(t *testing.T, callback CleanFunc) CleanFunc {
	err := removeAllContainers()
	assert.NoError(t, err)

	return callback
}

func TestCreateContainerWithCellCode(t *testing.T) {
	defer setup(t, removeAllContainers)()

	GOPATH := os.Getenv("GOPATH")
	// when
	res, err := docker.CreateContainer(
		tesseract.GetDefaultImage(),
		GOPATH+"/src/github.com/it-chain/tesseract/mock",
		"github.com/mock",
		"50005",
	)
	// then
	assert.NoError(t, err)

	// when
	containerName, err := getContainerName(res.ID)
	// then
	assert.NoError(t, err)
	assert.Equal(t, "/container_mock", containerName)
}

func TestCreateContainer_WhenSameNamedContainerExist_RandomGenerateName(t *testing.T) {
	defer setup(t, removeAllContainers)()

	GOPATH := os.Getenv("GOPATH")
	// when
	res, err := docker.CreateContainer(
		tesseract.GetDefaultImage(),
		GOPATH + "/src/github.com/it-chain/tesseract/mock",
		"github.com/mock",
		"50005",
	)
	// then
	assert.NoError(t, err)

	// when
	containerName, err := getContainerName(res.ID)
	// then
	assert.NoError(t, err)
	assert.Equal(t, "/container_mock", containerName)

	// when
	res2, err := docker.CreateContainer(
		tesseract.GetDefaultImage(),
		GOPATH + "/src/github.com/it-chain/tesseract/mock",
		"github.com/mock",
		"50005",
	)
	// then
	assert.NoError(t, err)

	// when
	randomGeneratedName, err := getContainerName(res2.ID)
	// then
	assert.NoError(t, err)
	assert.NotEqual(t, "/container_mock", randomGeneratedName)
}

func TestStartContainer(t *testing.T) {
	defer setup(t, removeAllContainers)()

	//given
	GOPATH := os.Getenv("GOPATH")
	res, err := docker.CreateContainer(
		tesseract.GetDefaultImage(),
		GOPATH+"/src/github.com/it-chain/tesseract/mock",
		"github.com/mock",
		"50005",
	)

	// when
	err = docker.StartContainer(res)
	// then
	assert.NoError(t, err)

	// when
	containerName, err := getContainerName(res.ID)
	// then
	assert.NoError(t, err)
	assert.Equal(t, "/container_mock", containerName)
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

func getContainerName(containerID string) (string, error) {
	ctx := context.Background()
	cli, err := dockerlib.NewEnvClient()
	defer cli.Close()

	if err != nil {
		return "", err
	}

	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}

	return containerJSON.Name, nil
}

func removeAllContainers() error {
	ctx := context.Background()
	cli, err := dockerlib.NewEnvClient()
	defer cli.Close()

	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, container := range containers {
		err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}

	return nil
}
