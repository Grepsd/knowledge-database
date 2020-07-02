package entity

import (
	uuid "github.com/satori/go.uuid"
)

type Source struct {
	ID   uuid.UUID
	Name string
	URL  string
	Tags []Tag
}
