package article

import "github.com/grepsd/knowledge-database/internal/domain/model/entity"

type SaveArticle struct {
	url entity.ArticleURL
	title entity.ArticleTitle
	datetime entity.ArticleSavedDateTime
}

func (s SaveArticle) Url() entity.ArticleURL {
	return s.url
}

func (s SaveArticle) Title() entity.ArticleTitle {
	return s.title
}

func (s SaveArticle) Datetime() entity.ArticleSavedDateTime {
	return s.datetime
}

func NewSaveArticle(url entity.ArticleURL, title entity.ArticleTitle, datetime entity.ArticleSavedDateTime) SaveArticle {
	return SaveArticle{url: url, title: title, datetime: datetime}
}

