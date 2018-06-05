package docker

type Image struct {
	Name string
	Tag  string
}

func (dc Image) GetFullName() string {
	return dc.Name + ":" + dc.Tag
}
