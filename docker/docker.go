package docker

import (
	"context"
	"errors"
	"io"
	"os"

	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/it-chain/tesseract"
)

func CreateContainer(containerImage tesseract.ContainerImage, srcPath string, destPath string, port string) (container.ContainerCreateCreatedBody, error) {

	GOPATH := os.Getenv("GOPATH")
	res := container.ContainerCreateCreatedBody{}

	if GOPATH == "" {
		return res, errors.New("invalid GOPATH. check your GOPATH")
	}

	imageName := containerImage.GetFullName()
	setUpImage(imageName)

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	defer cli.Close()

	if err != nil {
		return res, err
	}

	portBindings := nat.PortMap{
		nat.Port(port + "/tcp"): []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: port,
		}},
	}

	err = makeICodeLogDir(srcPath)
	if err != nil {
		return res, err
	}

	res, err = cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd: []string{
			"go",
			"run",
			path.Join("/go/src", destPath, "icode.go"),
			"-p" + port,
		},
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: nat.PortSet{
			nat.Port(port + "/tcp"): struct{}{},
		},
	}, &container.HostConfig{
		CapAdd:       []string{"SYS_ADMIN"},
		PortBindings: portBindings,
		Binds: []string{
			srcPath + ":/go/src/" + destPath,
		},
		Mounts: []mount.Mount{
			{
				ReadOnly: false,
				Type:     mount.TypeBind,
				Source:   makeICodeLogPath(srcPath),
				Target:   "/go/log",
			},
		},
	}, nil, makeICodeContainerName(srcPath))

	if err != nil {
		return res, err
	}

	return res, nil
}

func setUpImage(imageName string) error {

	exist, err := HasImage(imageName)

	if err != nil {
		return err
	}

	if !exist {
		if err := PullImage(imageName); err != nil {
			return err
		}
	}

	return nil
}

func StartContainer(containerBody container.ContainerCreateCreatedBody) error {

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	defer cli.Close()

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
	defer cli.Close()

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
	defer cli.Close()

	imageList, err := cli.ImageList(ctx, types.ImageListOptions{})

	if err != nil {
		return false, err
	}

	for _, image := range imageList {

		if len(image.RepoTags) == 0 {
			continue
		}
		if name == image.RepoTags[0] {
			return true, nil
		}
	}

	return false, nil
}

func RemoveContainer(id string) error {

	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	return cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})
}

func KillContainer(id string) error {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	return cli.ContainerKill(ctx, id, "9")
}

// todo : create test case
func GetPorts() ([]types.Port, error) {

	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	portList := make([]types.Port, 0)
	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})

	if err != nil {
		return portList, err
	}

	for _, container := range containerList {
		for _, port := range container.Ports {
			portInfo := types.Port{
				IP:          port.IP,
				PrivatePort: port.PrivatePort,
				PublicPort:  port.PublicPort,
			}
			portList = append(portList, portInfo)
		}
	}

	return portList, nil
}

func makeICodeLogDir(srcPath string) error {
	logDirPath := makeICodeLogPath(srcPath)

	_, err := os.Stat(logDirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(logDirPath, 0755)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}

func makeICodeLogPath(srcPath string) string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	return path.Join(basePath, "../logs", fmt.Sprintf("icode_%s", filepath.Base(srcPath)))
}

func makeICodeContainerName(srcPath string) string {
	icodeName := filepath.Base(srcPath)
	return fmt.Sprintf("container_%s", icodeName)
}
