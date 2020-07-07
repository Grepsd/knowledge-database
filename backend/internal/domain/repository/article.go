package repository

import "github.com/grepsd/knowledge-database/internal/domain/model/entity"

type Article interface {
	GetOneById(id entity.ArticleID) (*Article, error)
	GetByPageOrderedBySavedDateTime(pageNumber uint16) ([]*Article, error)
	DeleteByID(id entity.ArticleID) error
	MarkAsRead(id entity.ArticleID) error
	AddTag(id entity.ArticleID, tagID entity.TagID) error
	RemoveTag(id entity.ArticleID, tagID entity.TagID) error
	Save(article *entity.Article) error
}
