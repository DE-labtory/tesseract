package container

// referenced go-docker repository: https://github.com/docker/go-docker/blob/master/api/types/volume.go
type Volume struct {
	CreatedAt string
	Driver string
	Mountpoint string
	Name string
	Options map[string]string
}
