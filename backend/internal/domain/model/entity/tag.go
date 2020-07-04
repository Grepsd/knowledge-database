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
	id   *TagID
	name *TagName
}
type TagID struct {
	value uuid.UUID
}

func (i *TagID) String() string {
	return i.value.String()
}

type TagName struct {
	value string
}

func (t *Tag) ID() *TagID {
	return t.id
}

func (t *Tag) Name() *TagName {
	return t.name
}

func (t *TagName) Len() int {
	return len(t.value)
}

func (t *TagName) String() string {
	return t.value
}

func NewTagID(value uuid.UUID) *TagID {
	return &TagID{value: value}
}

func NewTagName(value string) *TagName {
	return &TagName{value: value}
}

func NewTag(id *TagID, name *TagName) (*Tag, error) {
	if name.Len() > NameMaxLength {
		return nil, errors.Wrap(NameIsTooLong, fmt.Sprintf("length : %d", name.Len()))
	}
	if name.Len() < NameMinLength {
		return nil, errors.Wrap(NameIsTooShort, fmt.Sprintf("length : %d", name.Len()))
	}
	return &Tag{id, name}, nil
}
