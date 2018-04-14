package docker

import (
	"context"
	"io"
	"os"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"

	"github.com/it-chain/tesseract"
)

const (
	imageName = "golang"
	imageTag  = "1.9"
)

func CreateContainerWithCellCode(iCodeInfo tesseract.ICodeInfo, cellCodeDir string) (container.ContainerCreateCreatedBody, error) {

	res := container.ContainerCreateCreatedBody{}
	image := imageName + ":" + imageTag

	exist, err := HasImage(image)

	if err != nil {
		return res, err
	}

	if !exist {
		if err := PullImage(image); err != nil {
			return res, err
		}
	}

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		return res, err
	}

	res, err = cli.ContainerCreate(ctx, &container.Config{
		Image: imageName + ":" + imageTag,
		Cmd: []string{
			"sh",
			"/cellcode/setup.sh",
			// Need to send smartcontract path
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: []string{cellCodeDir + ":/cellcode", iCodeInfo.Directory + ":/icode"},
	}, nil, "")

	if err != nil {
		return res, err
	}

	return res, nil
}

func StartContainer(containerBody container.ContainerCreateCreatedBody) error {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()

	err = cli.ContainerStart(ctx, containerBody.ID, types.ContainerStartOptions{})
	if err != nil {
		// An error occurred while starting the container!
		return err
	}

	return nil
}

func PullImage(imageName string) error {

	ctx := context.Background()
	cli, err := docker.NewEnvClient()

	if err != nil {
		return err
	}

	resp, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	defer resp.Close()

	io.Copy(os.Stdout, resp)

	if err != nil {
		return err
	}

	return nil
}

func HasImage(name string) (bool, error) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()

	imageList, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		return false, err
	}

	for _, image := range imageList {
		if name == image.RepoTags[0] {
			return true, nil
		}
	}

	return false, nil
}
