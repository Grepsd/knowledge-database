package article

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/tag"
	"regexp"
)

type Article struct {
	ID    uuid.UUID `json "id"`
	Title string    `json "title"`
	URL   string    `json "url"`
	Slug  string    `json "slug"`
	Tags  []*tag.Tag    `json "tags"`
}

type ReadRepositoryer interface {
	GetOneById(id uuid.UUID) (Article, error)
	GetOneBySlug(slug string) (Article, error)
	GetAll() ([]Article, error)
}
type WriteRepositoryer interface {
	Create(article Article) error
	Update(article Article) error
	DeleteById(id uuid.UUID) error
	AssignTagToArticle(articleID uuid.UUID, tagID uuid.UUID) error
}

type ReadWriteRepositoryer interface {
	ReadRepositoryer
	WriteRepositoryer
}

func NewArticle(id uuid.UUID, title string, url string, slug string, tags []*tag.Tag) *Article {
	return &Article{ID: id, Title: title, URL: url, Slug: slug, Tags: tags}
}

func GenerateSlugFromTitle(title string) (string, error) {
	re, err := regexp.Compile("[^a-zA-Z0-9-_]+")
	if err != nil {
		return "", fmt.Errorf("failed to compile regexp : %w", err)
	}
	slug := string(bytes.ToLower(re.ReplaceAll([]byte(title), []byte("_"))))
	return slug, err
}
