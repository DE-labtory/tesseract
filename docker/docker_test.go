package docker_test

import (
	"context"
	"os"
	"testing"

	"path"
	"path/filepath"
	"runtime"

	dockerlib "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/it-chain/tesseract/docker"
	"github.com/stretchr/testify/assert"
	"github.com/it-chain/tesseract"
	"sync"
	"time"
)

type CleanFunc = func() error

func setup(t *testing.T, callback CleanFunc) CleanFunc {
	err := removeAllContainers()
	assert.NoError(t, err)

	return callback
}

func TestCreateVolume(t *testing.T) {
	// given
	name := "myvol"

	defer DeleteVolumeByName(name)

	// when
	vol, err := docker.CreateVolume(name)
	// then
	assert.NoError(t, err)

	// when
	result, err := docker.FindVolumeByName(vol.Name)
	// then
	assert.NoError(t, err)
	assert.Equal(t, result.Name, vol.Name)
	assert.Equal(t, result.Mountpoint, vol.Mountpoint)
}

func DeleteVolumeByName(name string) error {
	ctx := context.Background()
	cli, _ := dockerlib.NewEnvClient()
	defer cli.Close()

	return cli.VolumeRemove(ctx, name, true)
}

func TestCreateContainer(t *testing.T) {
	defer setup(t, removeAllContainers)()

	testGolangImg := tesseract.ContainerImage{
		Name: "golang",
		Tag:  "1.9",
	}

	docker.CreateVolume("testVolume")
	GOPATH := os.Getenv("GOPATH")
	// when
	res, err := docker.CreateContainer(tesseract.ContainerConfig{
			Name:           "container_mock",
			ContainerImage: testGolangImg,
			IP:             "127.0.0.1",
			Port:           "50001",
			StartCmd:       []string{"go","run","icode_1/icode.go","-p","50001"},
			Network:        nil,
			Volume:         nil,
			HostICodeRoot: path.Join(GOPATH,"src/github.com/it-chain/tesseract/mock"),
			ImgSrcRootPath: "/go/src",

		},
	)
	// then
	docker.StartContainer(res)
	a:=sync.WaitGroup{}
	go func() {
		time.Sleep(10 * time.Second)
		a.Done()
	}()
	a.Add(1)
	a.Wait()
	assert.NoError(t, err)

	// when
	containerName, err := getContainerName(res.ID)
	// then
	assert.NoError(t, err)
	assert.Equal(t, "/container_mock", containerName)
}

//func TestCreateContainer_WhenSameNamedContainerExist_RandomGenerateName(t *testing.T) {
//	defer setup(t, removeAllContainers)()
//
//	GOPATH := os.Getenv("GOPATH")
//	// when
//	res, err := docker.CreateContainer(
//		tesseract.GetDefaultImage(),
//		GOPATH + "/src/github.com/it-chain/tesseract/mock",
//		"github.com/mock",
//		"50005",
//	)
//	// then
//	assert.NoError(t, err)
//
//	// when
//	containerName, err := getContainerName(res.ID)
//	// then
//	assert.NoError(t, err)
//	assert.Equal(t, "/container_mock", containerName)
//
//	// when
//	res2, err := docker.CreateContainer(
//		tesseract.GetDefaultImage(),
//		GOPATH + "/src/github.com/it-chain/tesseract/mock",
//		"github.com/mock",
//		"50005",
//	)
//	// then
//	assert.NoError(t, err)
//
//	// when
//	randomGeneratedName, err := getContainerName(res2.ID)
//	// then
//	assert.NoError(t, err)
//	assert.NotEqual(t, "/container_mock", randomGeneratedName)
//}

//func TestIsContainerExist(t *testing.T) {
//	defer setup(t, removeAllContainers)()
//
//	GOPATH := os.Getenv("GOPATH")
//	// when
//	_, err := docker.CreateContainer(
//		tesseract.GetDefaultImage(),
//		GOPATH + "/src/github.com/it-chain/tesseract/mock",
//		"github.com/mock",
//		"50005",
//	)
//	// then
//	assert.NoError(t, err)
//
//	// when
//	exist := docker.IsContainerExist("container_mock")
//	// then
//	assert.Equal(t, true, exist)
//
//	// when
//	exist2 := docker.IsContainerExist("/strange_container_name")
//	// then
//	assert.Equal(t, false, exist2)
//}

func TestConvertToSlashedPath(t *testing.T) {
	if runtime.GOOS == "window" {
		GOPATH := os.Getenv("GOPATH")

		// when
		result := docker.ConvertToSlashedPath(GOPATH)
		// then
		assert.Equal(t, "/c", result[:2])
	}
}

func TestMakeICodeLogDir(t *testing.T) {
	currentPath, _ := filepath.Abs("./")
	defer os.RemoveAll(".tmp")

	targetPath := path.Join(currentPath, ".tmp", "dir1", "dir2")
	docker.MakeICodeLogDir(targetPath)

	_, err := os.Stat(path.Join(currentPath, ".tmp/dir1/dir2"))

	assert.Equal(t, false, os.IsNotExist(err))
}

//func TestStartContainer(t *testing.T) {
//	defer setup(t, removeAllContainers)()
//
//	//given
//	GOPATH := os.Getenv("GOPATH")
//	res, err := docker.CreateContainer(
//		tesseract.GetDefaultImage(),
//		GOPATH+"/src/github.com/it-chain/tesseract/mock",
//		"github.com/mock",
//		"50005",
//	)
//
//	// when
//	err = docker.StartContainer(res)
//	// then
//	assert.NoError(t, err)
//
//	// when
//	containerName, err := getContainerName(res.ID)
//	// then
//	assert.NoError(t, err)
//	assert.Equal(t, "/container_mock", containerName)
//}

func TestPullImage(t *testing.T) {
	err := docker.PullImage("golang:1.9")
	assert.NoError(t, err)
}

func TestHasImageWhenImageExist(t *testing.T) {
	testImg := "golang:1.9"
	//given
	image := testImg
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
	testImg := "golang:1.9"
	//given
	image := testImg
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
