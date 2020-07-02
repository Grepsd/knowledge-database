package entity

import (
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/url"
	"time"
)

const TitleMaxLength = 200

var ArticleTitleIsTooLong = fmt.Errorf("article title is too long, max %d", TitleMaxLength)
var InvalidArticleURL = errors.New("invalid url")
var AlreadyHasTag = errors.New("already has tag")
var ArticleDoesNotHaveTag = errors.New("does not have this tag")
var ArticleHasAlreadyBeenRead = errors.New("article has already been read")

type Article struct {
	id            uuid.UUID
	title         string
	url           string
	readDateTime  time.Time
	savedDateTime time.Time
	tags          []*Tag
}

func (a *Article) Id() uuid.UUID {
	return a.id
}

func (a *Article) Title() string {
	return a.title
}

func (a *Article) Url() string {
	return a.url
}

func (a *Article) ReadDateTime() time.Time {
	return a.readDateTime
}

func (a *Article) SavedDateTime() time.Time {
	return a.savedDateTime
}

func (a *Article) Tags() []*Tag {
	return a.tags
}

func NewArticle(ID uuid.UUID, title string, URL string, readDatetime time.Time, savedDatetime time.Time, tags []*Tag) (*Article, error) {
	if _, err := url.ParseRequestURI(URL); err != nil {
		return nil, errors.Wrap(err, InvalidArticleURL.Error())
	}
	if len(title) > TitleMaxLength {
		return nil, errors.WithMessage(ArticleTitleIsTooLong, fmt.Sprintf(" length : %d", len(title)))
	}
	return &Article{id: ID, title: title, url: URL, readDateTime: readDatetime, savedDateTime: savedDatetime, tags: tags}, nil
}

func CreateArticle(title string, url string) (*Article, error) {
	return NewArticle(uuid.NewV4(), title, url, time.Time{}, time.Now(), NewEmptyTags())
}

func NewEmptyTags() []*Tag {
	return []*Tag{}
}

func (a *Article) HasTag(checkedTag *Tag) bool {
	for _, tag := range a.tags {
		if tag.ID == checkedTag.ID {
			return true
		}
	}
	return false
}

func (a *Article) AddTag(tag *Tag) error {
	if a.HasTag(tag) {
		return AlreadyHasTag
	}
	a.tags = append(a.tags, tag)
	return nil
}

func (a *Article) RemoveTag(toRemove *Tag) error {
	for index, tag := range a.tags {
		if tag.ID == toRemove.ID {
			a.tags[index] = a.tags[len(a.tags)-1]
			//a.tags[len(a.tags)-1] = &Tag{}
			a.tags = a.tags[:len(a.tags)-1]
			return nil
		}
	}
	return ArticleDoesNotHaveTag
}

func (a *Article) CountTags() int {
	return len(a.tags)
}

func (a *Article) Read() error {
	if !a.readDateTime.IsZero() {
		return ArticleHasAlreadyBeenRead
	}
	a.readDateTime = time.Now()
	return nil
}

func (a *Article) HasBeenRead() bool {
	return !a.readDateTime.IsZero()
}
