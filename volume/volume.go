package volume

// referenced go-docker repository: https://github.com/docker/go-docker/blob/master/api/types/volume.go
type Volume struct {
	CreatedAt  string
	Driver     string
	Mountpoint string
	Name       string
	Options    map[string]string
}

func NewVolume(createdAt, driver, mountpoint, name string, opts map[string]string) Volume {
	return Volume{
		CreatedAt:  createdAt,
		Driver:     driver,
		Mountpoint: mountpoint,
		Name:       name,
		Options:    opts,
	}
}

func (v Volume) GetID() string {
	return v.Name
}

func (v Volume) GetMountPoint() string {
	return v.Mountpoint
}
