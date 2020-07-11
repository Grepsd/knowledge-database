package article

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"regexp"
)

type Article struct {
	ID    uuid.UUID `json "id"`
	Title string    `json "title"`
	URL   string    `json "url"`
	Slug  string    `json "slug"`
}

func NewArticle(id uuid.UUID, title string, url string, slug string) *Article {
	return &Article{ID: id, Title: title, URL: url, Slug:slug}
}

func GenerateSlugFromTitle(title string) (slug string, err error) {
	re, err := regexp.Compile("[^a-zA-Z0-9-_]+")
	if err != nil {
		// "failed to compile regexp : "+err.Error()
		return slug, errors.As()
	}
	res := bytes.ToLower(re.ReplaceAll([]byte(key), []byte("_")))

}