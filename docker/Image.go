package docker

type Image struct {
	Name string
	Tag  string
}

func (dc Image) getFullName() string {
	return dc.Name + ":" + dc.Tag
}
