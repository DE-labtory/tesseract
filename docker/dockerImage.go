package docker

type DockerImage struct {
	Name string
	Tag  string
}

func (dc DockerImage) getFullName() string {
	return dc.Name + ":" + dc.Tag
}
