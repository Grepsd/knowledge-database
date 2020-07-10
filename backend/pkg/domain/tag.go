package domain

import (
	"github.com/google/uuid"
	"time"
)
type TagRepository interface {
	Create(tag Tag) error
	Delete(tag Tag) error
	GetAll() ([]Tag, error)
}
type Tag struct {
	id        uuid.UUID
	name      string
	createdOn time.Time
}

func (t Tag) CreatedOn() time.Time {
	return t.createdOn
}

func (t Tag) Name() string {
	return t.name
}

func (t Tag) Id() uuid.UUID {
	return t.id
}

func NewTag(id uuid.UUID, name string, createdOn time.Time) Tag {
	return Tag{id: id, name: name, createdOn: createdOn}
}

func CreateTag(id uuid.UUID, name string) Tag {
	return NewTag(id, name, time.Now())
}
