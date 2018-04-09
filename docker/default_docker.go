package docker

import (
	"os"
	"strings"
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types"
)

const (
	imageName = "golang"
	imageTag = "1.9.2-alpine3.6"
)

func CreateContainer(smartcontractName string) error {
	GOPATH := os.Getenv("GOPATH")

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		// An error occurred while creating new Docker Client!
		return err
	}

	imageName_splited := strings.Split(imageName, "/")
	image := imageName_splited[len(imageName_splited)-1]

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{
			"go", "run",
			"/cellcode/cellcode.go",
			// Need to send smartcontract path
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds: []string{GOPATH + "/src/it-chain/tesseract/cellcode:/cellcode"},
	}, nil, "")
	if err != nil {
		// An error occurred while creating docker container!
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		// An error occurred while starting the container!
		return err
	}

	return nil
}

func StartContainer() {

}
