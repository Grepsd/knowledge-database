package repository

import "github.com/grepsd/knowledge-database/internal/domain/model/entity"

type Tag interface {
	GetAll() ([]entity.Article, error)
	GetOneById(id entity.TagID) (entity.Tag, error)
	Save(tag entity.Tag) error
	DeleteById(id entity.TagID)
}
