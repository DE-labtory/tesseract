package docker

type Docker interface {
	CreateContainerWithCellCode()
	StartContainer()
}
