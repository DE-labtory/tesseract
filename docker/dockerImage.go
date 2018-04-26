package docker

type DockerImage struct {
	Name string
	Tag  string
}

func (dc DockerImage) getName() string {
	return dc.Name + ":" + dc.Tag
}
