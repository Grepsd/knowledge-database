package article

import (
	"github.com/grepsd/knowledge-database/internal/application/command/article"
	"github.com/grepsd/knowledge-database/internal/domain/model/entity"
	"github.com/grepsd/knowledge-database/internal/domain/repository"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type SaveArticle struct {
	repository repository.Article
}

func NewSaveArticle(repository repository.Article) SaveArticle {
	return SaveArticle{repository: repository}
}

func (s *SaveArticle) Handle(command article.SaveArticle) (*entity.Article, error) {
	newArticle, err := entity.NewArticle(
		entity.NewArticleID(uuid.NewV4()),
		command.Title(),
		command.Url(),
		entity.NewArticleReadDateTime(time.Time{}),
		command.Datetime(),
		entity.NewEmptyTags())
	if err != nil {
		return &entity.Article{}, errors.Wrap(err, "command failed")
	}
	err = s.repository.Save(newArticle)
	if err != nil {
		return &entity.Article{}, errors.Wrap(err, "repository save failed")
	}
	return newArticle, nil
}
