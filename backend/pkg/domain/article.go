package domain

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrCannotReadArticleBeforeItIsSaved = errors.New("cannot read article before it is saved")
var ErrCannotUnreadArticleBeforeItIsRead = errors.New("cannot unread article before it is read")

type ArticleRepository interface {
	Create(a Article) error
	GetOnById(id uuid.UUID) (Article, error)
	GetAll() ([]Article, error)
	GetAllByTags([]Tag) ([]Article, error)
	Delete(article Article) error
	Update(article Article) error
}

type Article struct {
	ID      uuid.UUID
	Title   string
	URL     string
	SavedOn time.Time
	ReadOn  time.Time
	Tags    []Tag
}

func NewArticle(id uuid.UUID, title string, url string, savedOn time.Time, readOn time.Time, tags []Tag) Article {
	return Article{ID: id, Title: title, URL: url, SavedOn: savedOn, ReadOn: readOn, Tags: tags}
}

func CreateArticle(id uuid.UUID, title string, url string, tags []Tag) Article {
	return NewArticle(id, title, url, time.Now(), time.Time{}, tags)
}

func (a *Article) Read(t time.Time) error {
	if t.Before(a.SavedOn) {
		return ErrCannotReadArticleBeforeItIsSaved
	}
	a.ReadOn = t
	return nil
}

func (a *Article) UnRead() error {
	if a.ReadOn.IsZero() {
		return ErrCannotUnreadArticleBeforeItIsRead
	}
	return nil
}

func (a *Article) Tag(t Tag) error {
	a.Tags = append(a.Tags, t)
	return nil
}

func (a *Article) UnTag(tag Tag) error {
	for index, t := range a.Tags {
		if t.id == tag.id {
			a.Tags[index] = a.Tags[index-1]
			a.Tags = a.Tags[:len(a.Tags)-1]
		}
	}
	return nil
}

func (a *Article) HasTag(tag Tag) bool {
	for _, t := range a.Tags {
		if t.id == tag.id {
			return true
		}
	}
	return false
}