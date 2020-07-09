package article

import (
	"github.com/grepsd/knowledge-database/internal/application/query"
	"github.com/grepsd/knowledge-database/internal/domain/repository"
)

type FindAll struct {
	repository repository.Article
}

func NewFindAll(repository repository.Article) *FindAll {
	return &FindAll{repository: repository}
}

func (f *FindAll) Handle(query *query.FindAll) {
	results, err :=  f.repository.FindAll()
	query.HandleResult(results, err)
}