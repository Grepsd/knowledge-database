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
	id            ArticleID
	title         ArticleTitle
	url           ArticleURL
	readDateTime  ArticleReadDateTime
	savedDateTime ArticleSavedDateTime
	tags          []*Tag
}
type ArticleID struct {
	value uuid.UUID
}

func (i ArticleID) String() string {
	return i.value.String()
}

type ArticleTitle struct {
	value string
}

func NewArticleTitle(value string) ArticleTitle {
	return ArticleTitle{value: value}
}

type ArticleURL struct {
	value string
}

func (u ArticleURL) String() string {
	return u.value
}

func NewArticleURL(value string) (ArticleURL, error) {
	if _, err := url.ParseRequestURI(value); err != nil {
		return ArticleURL{}, errors.Wrap(err, InvalidArticleURL.Error())
	}

	return ArticleURL{value: value}, nil
}

type ArticleReadDateTime struct {
	value time.Time
}
type ArticleSavedDateTime struct {
	value time.Time
}

func (t ArticleSavedDateTime) String() string {
	return t.value.String()
}

func (t ArticleReadDateTime) IsZero() bool {
	return t.value.IsZero()
}

func (t ArticleReadDateTime) Update(value time.Time) (ArticleReadDateTime, error) {
	t.value = value
	return t, nil
}

func (t *ArticleReadDateTime) String() string {
	return t.value.String()
}

func (a *Article) Id() ArticleID {
	return a.id
}

func (a *Article) Title() ArticleTitle {
	return a.title
}

func (a *Article) Url() ArticleURL {
	return a.url
}

func (a *Article) ReadDateTime() ArticleReadDateTime {
	return a.readDateTime
}

func (a Article) SavedDateTime() ArticleSavedDateTime {
	return a.savedDateTime
}

func (a *Article) Tags() []*Tag {
	return a.tags
}

func (at *ArticleTitle) Len() int {
	return len(at.value)
}

func (at ArticleTitle) String() string {
	return at.value
}

func NewArticleID(id uuid.UUID) ArticleID {
	return ArticleID{id}
}

func NewArticle(ID ArticleID, title ArticleTitle, URL ArticleURL, readDatetime ArticleReadDateTime, savedDatetime ArticleSavedDateTime, tags []*Tag) (*Article, error) {
	if title.Len() > TitleMaxLength {
		return nil, errors.WithMessage(ArticleTitleIsTooLong, fmt.Sprintf(" length : %d", title.Len()))
	}
	return &Article{id: ID, title: title, url: URL, readDateTime: readDatetime, savedDateTime: savedDatetime, tags: tags}, nil
}

func CreateArticle(title ArticleTitle, url ArticleURL) (*Article, error) {
	id := NewArticleID(uuid.NewV4())
	return NewArticle(id, title, url, NewArticleReadDateTime(time.Time{}), NewArticleSavedDateTime(time.Now()), NewEmptyTags())
}

func NewArticleSavedDateTime(t time.Time) ArticleSavedDateTime {
	return ArticleSavedDateTime{t}
}

func NewArticleReadDateTime(t time.Time) ArticleReadDateTime {
	return ArticleReadDateTime{t}
}

func NewEmptyTags() []*Tag {
	return NewArticleTags()
}

func NewArticleTags() []*Tag {
	return []*Tag{}
}

func (a *Article) HasTag(checkedTag *Tag) bool {
	for _, tag := range a.Tags() {
		if tag.ID() == checkedTag.ID() {
			return true
		}
	}
	return false
}

func (a *Article) AddTag(tag *Tag) error {
	if a.HasTag(tag) {
		return AlreadyHasTag
	}
	a.tags = append(a.Tags(), tag)
	return nil
}

func (a *Article) RemoveTag(toRemove *Tag) error {
	for index, tag := range a.tags {
		if tag.ID() == toRemove.ID() {
			a.tags[index] = a.tags[len(a.tags)-1]
			a.tags = a.tags[:len(a.tags)-1]
			return nil
		}
	}
	return ArticleDoesNotHaveTag
}

func (a *Article) CountTags() int {
	return len(a.tags)
}

func (a *Article) Read(time time.Time) error {
	if !a.readDateTime.IsZero() {
		return ArticleHasAlreadyBeenRead
	}
	a.UpdateReadDateTime(time)
	return nil
}

func (a *Article) HasBeenRead() bool {
	return !a.ReadDateTime().IsZero()
}

func (a *Article) UpdateReadDateTime(t time.Time) (ArticleReadDateTime) {
	dt := NewArticleReadDateTime(t)
	a.readDateTime = dt
	return a.readDateTime
}
