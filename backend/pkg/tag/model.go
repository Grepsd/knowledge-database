package tag

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewTag(id uuid.UUID, name string) *Tag {
	return &Tag{ID: id, Name: name}
}

type ReadRepositoryer interface {
	GetOneById(id uuid.UUID) (Tag, error)
	GetOneByName(name string) (Tag, error)
	GetAll(hasCategories bool) ([]Tag, error)
	FindTagsByArticleID(id []uuid.UUID) ([]Tag, error)
}
type WriteRepositoryer interface {
	Create(article Tag) error
	Update(article Tag) error
	DeleteById(id uuid.UUID) error
}

type ReadWriteRepositoryer interface {
	ReadRepositoryer
	WriteRepositoryer
}
