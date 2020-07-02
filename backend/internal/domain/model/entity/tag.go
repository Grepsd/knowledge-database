package entity

import (
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const NameMaxLength = 25
const NameMinLength = 1

var NameIsTooLong = fmt.Errorf("name is too long, max : %d", NameMaxLength)
var NameIsTooShort = fmt.Errorf("name is too short, min : %d", NameMinLength)

type Tag struct {
	ID   uuid.UUID
	Name string
}

func NewTag(ID uuid.UUID, name string) (*Tag, error) {
	if len(name) > NameMaxLength {
		return nil, errors.Wrap(NameIsTooLong, fmt.Sprintf("length : %d", len(name)))
	}
	if len(name) < NameMinLength {
		return nil, errors.Wrap(NameIsTooShort, fmt.Sprintf("length : %d", len(name)))
	}
 	return &Tag{ID: ID, Name: name}, nil
}

func CreateTag(name string) (*Tag, error) {
	return NewTag(uuid.NewV4(), name)
}
