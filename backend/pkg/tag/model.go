package tag

type Tag struct {
	ID   string
	Name string
}

func NewTag(id string, name string) *Tag {
	return &Tag{ID: id, Name: name}
}
