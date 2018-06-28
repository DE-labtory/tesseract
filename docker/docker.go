package docker

import "docker.io/go-docker/api/types/container"

type Docker interface {
	CreateContainerWithCellCode(dockerImage Image, dir string, shPath string, port string) (container.ContainerCreateCreatedBody, error)
	StartContainer(containerBody container.ContainerCreateCreatedBody) error
	GetUsingPorts() ([]Port, error)
}

type Port struct {
	// Host IP address that the container's port is mapped to
	IP string

	// Port on the container
	PrivatePort int

	// Port exposed on the host
	PublicPort int
}
