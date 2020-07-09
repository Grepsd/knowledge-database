package article

import (
	"errors"
	"github.com/grepsd/knowledge-database/internal/application/command/article"
	"github.com/grepsd/knowledge-database/internal/domain/model/entity"
	"github.com/grepsd/knowledge-database/internal/domain/repository"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewSaveArticle(t *testing.T) {

	type args struct {
		repository repository.Article
	}
	tests := []struct {
		name string
		args args
		want SaveArticle
	}{
		{"valid",
			args{
				repository: testArticleRepository{},
			},
			SaveArticle{repository: testArticleRepository{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSaveArticle(tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSaveArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testArticleRepository struct {
	saveReturnError bool
}

func (t testArticleRepository) GetOneById(id entity.ArticleID) (*repository.Article, error) {
	panic("implement me")
}

func (t testArticleRepository) GetByPageOrderedBySavedDateTime(pageNumber uint16) ([]*repository.Article, error) {
	panic("implement me")
}

func (t testArticleRepository) DeleteByID(id entity.ArticleID) error {
	panic("implement me")
}

func (t testArticleRepository) MarkAsRead(id entity.ArticleID) error {
	panic("implement me")
}

func (t testArticleRepository) AddTag(id entity.ArticleID, tagID entity.TagID) error {
	panic("implement me")
}

func (t testArticleRepository) RemoveTag(id entity.ArticleID, tagID entity.TagID) error {
	panic("implement me")
}

func (t testArticleRepository) Save(article *entity.Article) error {
	if t.saveReturnError {
		return errors.New("dummy error")
	}
	return nil
}

func TestSaveArticle_Handle(t *testing.T) {
	type fields struct {
		repository repository.Article
	}
	type args struct {
		command article.SaveArticle
	}
	now := time.Now()
	testArticleURL, _ := entity.NewArticleURL("http://go.com")
	invalidArticleTitle := entity.NewArticleTitle(strings.Repeat("d", entity.TitleMaxLength+1))
	cmd := article.NewSaveArticle(testArticleURL, entity.NewArticleTitle("valid"), entity.NewArticleSavedDateTime(now))
	invalidCmd := article.NewSaveArticle(testArticleURL, invalidArticleTitle, entity.NewArticleSavedDateTime(now))
	testArticle, err := entity.NewArticle(entity.NewArticleID(uuid.UUID{}), entity.NewArticleTitle("valid"), testArticleURL, entity.NewArticleReadDateTime(time.Time{}), entity.NewArticleSavedDateTime(now), entity.NewEmptyTags())
	if err != nil {
		t.Error("failed to create test article")
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Article
		wantErr bool
	}{
		{
			"valid",
			fields{testArticleRepository{}},
			args{cmd},
			testArticle,
			false,
		},
		{
			"valid",
			fields{testArticleRepository{}},
			args{invalidCmd},
			&entity.Article{},
			true,
		},
		{
			"valid",
			fields{testArticleRepository{true}},
			args{cmd},
			&entity.Article{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SaveArticle{
				repository: tt.fields.repository,
			}
			got, err := s.Handle(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("expected valid Article, got nil")
			}
			if tt.want.Title() != got.Title() || tt.want.Url() != got.Url() || tt.want.ReadDateTime() != got.ReadDateTime() {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}
