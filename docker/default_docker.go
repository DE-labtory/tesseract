package docker

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"fmt"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

const (
	DefaultImageName = "golang"
	DefaultImageTag  = "1.9"
)

func CreateContainerWithCellCode(dockerImage Image, dir string, shPath string, port string) (container.ContainerCreateCreatedBody, error) {

	GOPATH := os.Getenv("GOPATH")
	res := container.ContainerCreateCreatedBody{}
	image := dockerImage.GetFullName()

	exist, err := HasImage(image)

	if err != nil {
		return res, err
	}

	if !exist {
		fmt.Println("Image not exit. try to pull")
		if err := PullImage(image); err != nil {
			return res, err
		}
		fmt.Println("finish to pull image")
	}

	ctx := context.Background()
	cli, err := docker.NewEnvClient()

	if err != nil {
		return res, err
	}
	fmt.Println("docer env success")
	portBindings := nat.PortMap{
		nat.Port(port + "/tcp"): []nat.PortBinding{{
			HostIP:   "0.0.0.0",
			HostPort: port,
		}},
	}

	res, err = cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd: []string{
			"sh",
			"/sh/" + filepath.Base(shPath),
			port,
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
			GOPATH + "/src:/go/src",
			dir + ":/icode",
			filepath.Dir(shPath) + ":/sh"},
	}, nil, "")

	log.Printf(GOPATH + "/src:/go/src")
	log.Println(res)
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

func GetLocalIPAddressFromContainer(containerID string) (string, error) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()

	inspectBody, err := cli.ContainerInspect(ctx, containerID)

	if err != nil {
		// An error occurred while starting the container!
		return "", err
	}

	return inspectBody.NetworkSettings.IPAddress, nil
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

		if len(image.RepoTags) == 0 {
			continue
		}
		if name == image.RepoTags[0] {
			return true, nil
		}
	}

	return false, nil
}

func CloseContainer(id string) error {

	ctx := context.Background()
	cli, _ := docker.NewEnvClient()

	err := cli.ContainerKill(ctx, id, "9")
	if err != nil {
		return err
	}

	cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{})

	return nil
}

// todo : create test case
func GetUsingPorts() ([]Port, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	portList := make([]Port, 0)
	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})

	if err != nil {
		return portList, err
	}

	for _, oneContainer := range containerList {
		for _, containerPort := range oneContainer.Ports {
			portInfo := Port{
				IP:          containerPort.IP,
				PrivatePort: int(containerPort.PrivatePort),
				PublicPort:  int(containerPort.PublicPort),
			}
			portList = append(portList, portInfo)
		}
	}

	return portList, nil
}
