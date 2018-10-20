package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"docker.io/go-docker/api/types/filters"
	"docker.io/go-docker/api/types/volume"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	git"github.com/docker/go-connections/nat"
	"github.com/it-chain/iLogger"
	"github.com/it-chain/tesseract"
)

func CreateContainer(containerImage tesseract.ContainerImage, srcPath string, destPath string, IP string, port string) (container.ContainerCreateCreatedBody, error) {

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
			HostIP:   GetHostIpAddress(),
			HostPort: port,
		}},
	}

	logDirPath := makeICodeLogPath(srcPath)
	err = MakeICodeLogDir(logDirPath)
	if err != nil {
		return res, err
	}

	containerName := makeICodeContainerName(srcPath)
	if IsContainerExist(containerName) {
		iLogger.Infof(nil, "[tesseract] container name \"%s\" exist, container name now random generated", containerName)
		containerName = ""
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
			"it-chain-volume" + ":/go/src/" + "github.com/junbeomlee",
		},

	}, nil, "")

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

func IsContainerExist(name string) bool {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	containerList, _ := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	for _, container := range containerList {
		if container.Names[0] == fmt.Sprintf("/%s", name) {
			return true
		}
	}

	return false
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

func GetHostIpAddress() string {
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	host, _ := docker.ParseHostURL(cli.DaemonHost())
	return strings.Split(host.Host, ":")[0]
}

func MakeICodeLogDir(logDirPath string) error {
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
	icodeLogPath := srcPath
	logDir := fmt.Sprintf("icode_%s", filepath.Base(srcPath))
	return path.Join(icodeLogPath, "../../icode-logs", logDir)
}

func makeICodePath(srcPath string) string {
	return ConvertToSlashedPath(srcPath)
}

func ConvertToSlashedPath(srcPath string) string {
	splited := strings.Split(srcPath, ":")

	if len(splited) <= 1 {
		return srcPath
	}

	driveName := strings.ToLower(splited[0])
	return strings.Replace("/"+driveName+splited[1], "\\", "/", -1)

}

func makeICodeContainerName(srcPath string) string {
	icodeName := filepath.Base(srcPath)
	return fmt.Sprintf("container_%s", icodeName)
}

func CreateVolume(name string) (types.Volume, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	vol, err := FindVolumeByName(name)
	if err != nil {
		return types.Volume{}, err
	}
	if !isVolumeEmpty(vol) {
		return types.Volume{}, errors.New(fmt.Sprintf("volume [name: %s] already exist", vol.Name, vol.Mountpoint))
	}

	res, err := cli.VolumeCreate(ctx, convToVolumesCreateBody(name))
	if err != nil {
		return types.Volume{}, err
	}

	return res, nil
}

func FindVolumeByName(name string) (types.Volume, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	listBody, err := cli.VolumeList(ctx, filters.Args{})
	if err != nil {
		return types.Volume{}, err
	}

	for _, vol := range listBody.Volumes {
		if vol.Name == name {
			return *vol, nil
		}
	}

	return types.Volume{}, nil
}

func isVolumeEmpty(vol types.Volume) bool {
	return reflect.DeepEqual(vol, types.Volume{})
}

func convToVolumesCreateBody(name string) volume.VolumesCreateBody {
	return volume.VolumesCreateBody{
		Driver:     "local",
		DriverOpts: map[string]string{},
		Labels:     map[string]string{},
		Name:       name,
	}
}
