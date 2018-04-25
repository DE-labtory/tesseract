package docker

import (
	"context"
	"io"
	"os"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"

	"path/filepath"

	"github.com/it-chain/tesseract"
)

const (
	DefaultImageName = "golang"
	DefaultImageTag  = "1.9"
	GrpcGoImageName  = "grpc/go"
	GrpcGoImageTag   = "1.0"
)

func CreateContainerWithCellCode(dockerImage DockerImage, iCodeInfo tesseract.ICodeInfo, shPath string) (container.ContainerCreateCreatedBody, error) {
	GOPATH := os.Getenv("GOPATH")
	res := container.ContainerCreateCreatedBody{}
	image := dockerImage.Name + ":" + dockerImage.Tag

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
		Image: dockerImage.Name + ":" + dockerImage.Tag,
		Cmd: []string{
			"sh",
			"/sh/" + filepath.Base(shPath),
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		CapAdd: []string{"SYS_ADMIN"},
		Binds: []string{
			GOPATH + "/src:/go/src",
			iCodeInfo.Directory + ":/icode",
			filepath.Dir(shPath) + ":/sh"},
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
