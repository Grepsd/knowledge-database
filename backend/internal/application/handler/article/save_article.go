package article

import (
	"github.com/grepsd/knowledge-database/internal/application/command/article"
	"github.com/grepsd/knowledge-database/internal/domain/model/entity"
	"github.com/grepsd/knowledge-database/internal/domain/repository"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type SaveArticle struct {
	command    article.SaveArticle
	repository repository.Article
}

func NewSaveArticle(command article.SaveArticle, repository repository.Article) SaveArticle {
	return SaveArticle{command: command, repository: repository}
}

func (s *SaveArticle) Execute() (*entity.Article, error) {
	article, err := entity.NewArticle(
		entity.NewArticleID(uuid.NewV4()),
		s.command.Title(),
		s.command.Url(),
		entity.NewArticleReadDateTime(),
		s.command.Datetime(),
		entity.NewEmptyTags())
	if err != nil {
		return &entity.Article{}, errors.Wrap(err, "command failed")
	}
	err = s.repository.Save(article)
	if err != nil {
		return &entity.Article{}, errors.Wrap(err, "repository save failed")
	}
	return article, nil
}
