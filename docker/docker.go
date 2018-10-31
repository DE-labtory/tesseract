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
	"github.com/docker/go-connections/nat"
	"github.com/it-chain/iLogger"
	"github.com/it-chain/tesseract"
	"docker.io/go-docker/api/types/network"
)

func CreateContainer(config tesseract.ContainerConfig) (container.ContainerCreateCreatedBody, error) {

	GOPATH := os.Getenv("GOPATH")
	res := container.ContainerCreateCreatedBody{}

	if GOPATH == "" {
		return res, errors.New("invalid GOPATH. check your GOPATH")
	}

	imageName := config.ContainerImage.GetFullName()
	setUpImage(imageName)

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	defer cli.Close()

	if err != nil {
		return res, err
	}

	networkName := ""
	portBinding := nat.PortMap{}
	exposedPort := nat.PortSet{}
	var networkConfig *network.NetworkingConfig

	if config.Network != nil {
		networkName = config.Network.Name
		endpointSetting := make(map[string]*network.EndpointSettings)
		endpointSetting[config.Network.Name] = &network.EndpointSettings{
			IPAddress:           config.ContainerIp,
		}
		networkConfig = &network.NetworkingConfig{
			EndpointsConfig: endpointSetting,
		}
	}else{
		exposedPort = nat.PortSet{
			nat.Port(config.Port + "/tcp"): struct{}{},
		}
		portBinding = nat.PortMap{
			nat.Port(config.Port + "/tcp"): []nat.PortBinding{{
				HostIP:   config.HostIp,
				HostPort: config.Port,
			}},
		}

	}

	if err != nil {
		return res, err
	}

	containerName := config.Name
	if IsContainerExist(containerName) {
		iLogger.Infof(nil, "[Tesseract] The Container name exists, creating new name - ContainerName: [%s]", containerName)
		containerName = ""
	}

	res, err = cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		Cmd:          config.StartCmd,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: exposedPort,

	}, &container.HostConfig{
		CapAdd:       []string{"SYS_ADMIN"},
		PortBindings: portBinding,
		Binds:        config.Mount,
		NetworkMode:  container.NetworkMode(networkName),
	}, networkConfig, containerName)

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

func StartContainer(containerBody container.ContainerCreateCreatedBody) (types.ContainerJSON, error) {

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	defer cli.Close()

	err = cli.ContainerStart(ctx, containerBody.ID, types.ContainerStartOptions{})
	if err != nil {
		// An error occurred while starting the container!
		return types.ContainerJSON{}, err
	}

	return cli.ContainerInspect(ctx, containerBody.ID)
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

func CreateVolume(name string) (tesseract.Volume, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	vol, err := FindVolumeByName(name)
	if err != nil {
		return tesseract.Volume{}, err
	}

	if !isVolumeEmpty(vol) {
		return tesseract.NewVolume(vol.CreatedAt, vol.Driver, vol.Mountpoint, vol.Name, vol.Options), nil
	}

	res, err := cli.VolumeCreate(ctx, convToVolumesCreateBody(name))
	if err != nil {
		return tesseract.Volume{}, err
	}

	return tesseract.NewVolume(res.CreatedAt, res.Driver, res.Mountpoint, res.Name, res.Options), nil
}

func RemoveVolume(name string, force bool) error {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	return cli.VolumeRemove(ctx, name, force)
}

func CreateNetwork(name string, subnet string) (tesseract.Network, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	findNetwork, err := FindNetworkByName(name)
	if err != nil {
		return tesseract.Network{}, err
	}

	if !isNetworkEmpty(findNetwork) {
		return tesseract.Network{}, errors.New("already exist network name")
	}

	res, err := cli.NetworkCreate(ctx, name, types.NetworkCreate{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Options: nil,
			Config: []network.IPAMConfig{
				{Subnet: subnet},
			},
		},
	})

	if err != nil {
		return tesseract.Network{}, err
	}

	return tesseract.Network{
		ID:   res.ID,
		Name: name,
	}, nil
}

func FindNetworkByName(name string) (tesseract.Network, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	networkList, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return tesseract.Network{}, err
	}
	for _, net := range networkList {
		if net.Name == name {
			return tesseract.Network{
				ID:     net.ID,
				Name:   net.Name,
			}, nil
		}
	}
	return tesseract.Network{}, nil
}

func FindVolumeByName(name string) (tesseract.Volume, error) {
	ctx := context.Background()
	cli, _ := docker.NewEnvClient()
	defer cli.Close()

	listBody, err := cli.VolumeList(ctx, filters.Args{})
	if err != nil {
		return tesseract.Volume{}, err
	}

	for _, vol := range listBody.Volumes {
		if vol.Name == name {
			return tesseract.Volume{
				CreatedAt:  vol.CreatedAt,
				Driver:     vol.Driver,
				Mountpoint: vol.Mountpoint,
				Name:       vol.Name,
				Options:    vol.Options,
			}, nil
		}
	}

	return tesseract.Volume{}, nil
}

func isVolumeEmpty(vol tesseract.Volume) bool {
	return reflect.DeepEqual(vol, tesseract.Volume{})
}

func isNetworkEmpty(vol tesseract.Network) bool {
	return reflect.DeepEqual(vol, tesseract.Network{})
}

func convToVolumesCreateBody(name string) volume.VolumesCreateBody {
	return volume.VolumesCreateBody{
		Driver:     "local",
		DriverOpts: map[string]string{},
		Labels:     map[string]string{},
		Name:       name,
	}
}
