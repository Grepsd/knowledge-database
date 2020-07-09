package repository

import (
	"database/sql"
	"github.com/grepsd/knowledge-database/internal/domain/model/entity"
	"github.com/grepsd/knowledge-database/internal/domain/repository"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type Article struct {
	db *sql.DB
}

func NewArticle() *Article {
	db, err := sql.Open("sqlite3", "file:knowledgedatabase.sqlite?cache=shared")
	if err != nil {
		panic("cannot connect to database : " + err.Error())
	}
	query := `create table articles
(
    id             varchar(70) not null
        constraint articles_pk
            primary key,
    title          varchar(255),
    url            varchar(255),
    saved_datetime datetime    not null,
    read_datetime  datetime default null
);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("cannot create table : " + err.Error())
	} else {
		log.Println("article table created")
	}
	return &Article{db}
}

func (a *Article) GetOneById(id entity.ArticleID) (*repository.Article, error) {
	panic("implement me")
}

func (a *Article) FindByPageOrderedBySavedDateTime(pageNumber uint16) ([]*repository.Article, error) {
	panic("implement me")
}

func (a *Article) FindAll() ([]*entity.Article, error) {
	rows, err := a.db.Query("SELECT id, title, url, saved_datetime, read_datetime FROM articles")
	if err != nil {
		return []*entity.Article{}, errors.Wrap(err, "query failed : ")
	}
	defer rows.Close()
	var articles []*entity.Article
	for rows.Next() {
		var id, title, url string
		var savedDateTime time.Time
		var readDateTime time.Time
		err := rows.Scan(&id, &title, &url, &savedDateTime, &readDateTime)
		if err != nil {
			log.Fatal(err)
		}
		u, err := uuid.FromString(id)
		if err != nil {
			log.Fatal("invalid uuid : " + err.Error())
		}
		articleId := entity.NewArticleID(u)
		articleTitle := entity.NewArticleTitle(title)
		articleURL, err := entity.NewArticleURL(url)
		if err != nil {
			log.Fatal("invalid URL : " + err.Error())
		}
		articleSavedDateTime := entity.NewArticleSavedDateTime(savedDateTime)
		articleReadDateTime := entity.NewArticleReadDateTime(readDateTime)
		article, err := entity.NewArticle(articleId, articleTitle, articleURL, articleReadDateTime, articleSavedDateTime, entity.NewEmptyTags())
		if err != nil {
			log.Fatal("failed to create article : " + err.Error())
		}
		articles = append(articles, article)

	}
	return articles, nil
}

func (a *Article) DeleteByID(id entity.ArticleID) error {
	panic("implement me")
}

func (a *Article) MarkAsRead(id entity.ArticleID) error {
	panic("implement me")
}

func (a *Article) AddTag(id entity.ArticleID, tagID entity.TagID) error {
	panic("implement me")
}

func (a *Article) RemoveTag(id entity.ArticleID, tagID entity.TagID) error {
	panic("implement me")
}

func (a *Article) Save(article *entity.Article) error {
	query := `INSERT INTO articles (id, title, url, saved_datetime, read_datetime) VALUES (?, ?, ?, ?, ?)`
	stmt, err := a.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(article.Id().String(), article.Title().String(), article.Url().String(), article.SavedDateTime().String(), time.Time{})
	if err != nil {
		log.Println("failed to save article to database : " + err.Error())
		return err
	}
	return nil
}
