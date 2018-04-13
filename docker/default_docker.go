package docker

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types"
)

const (
	imageName = "golang"
	imageTag = "1.9"
)

type ICodeInfo struct {
	Name string
	Directory string
}

func CreateContainerWithCellCode(iCodeInfo ICodeInfo, cellCodeDir string) (container.ContainerCreateCreatedBody, error) {
	res := container.ContainerCreateCreatedBody{}
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		// An error occurred while creating new Docker Client!
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
		// An error occurred while creating docker container!
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
