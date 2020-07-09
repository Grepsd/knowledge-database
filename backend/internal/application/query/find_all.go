package query

import "github.com/grepsd/knowledge-database/internal/domain/model/entity"

type FindAll struct {
	resultHandler func([]*entity.Article, error)
}

func NewFindAll(resultHandler func([]*entity.Article, error)) *FindAll {
	return &FindAll{resultHandler: resultHandler}
}

func (f FindAll) HandleResult(result []*entity.Article, err error) {
	f.resultHandler(result, err)
}