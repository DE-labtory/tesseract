package docker

type Docker interface {
	CreateContainer()
	StartContainer()
}
