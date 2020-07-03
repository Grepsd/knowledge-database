package repository

import "github.com/grepsd/knowledge-database/internal/domain/model/entity"

type ArticleRepository interface {
	GetAll() ([]entity.Article, error)
	GetAllNotRead() ([]entity.Article, error)
	GetAllByTags([]entity.Tag) ([]entity.Article, error)
	GetOneById() ([]entity.Article, error)
}