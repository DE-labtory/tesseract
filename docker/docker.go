package docker

import "docker.io/go-docker/api/types/container"

type Docker interface {
	CreateContainerWithCellCode(dockerImage Image, dir string, shPath string, port string) (container.ContainerCreateCreatedBody, error)
	StartContainer(containerBody container.ContainerCreateCreatedBody) error
}
