package pkg

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	db *sql.DB
}

func NewDB(dsn string) *DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to db : " + err.Error())
	}
	return &DB{db}
}

func (d *DB) CreateArticle(a *article.Article) error {
	query := `INSERT INTO articles (id, title, url, slug)
VALUES ($1, $2, $3, $4)`
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return errors.New("failed to prepare query : " + err.Error())
	}
	_, err = stmt.Exec(a.ID.String(), a.Title, a.URL, a.Slug)
	if err != nil {
		return errors.New("failed to exec query : " + err.Error())
	}
	return nil
}

func (d *DB) ListArticles() ([]*article.Article, error) {
	var articles []*article.Article
	query := `SELECT id, title, url, slug FROM articles ORDER BY title`
	results, err := d.db.Query(query)
	if err != nil {
		return articles, errors.New("failed to exec query : " + err.Error())
	}
	for results.Next() {
		var a article.Article
		var articleID string
		err := results.Scan(&articleID, &a.Title, &a.URL, &a.Slug)
		if err != nil {
			return articles, errors.New("failed to scan query result : " + err.Error())
		}
		a.ID = uuid.MustParse(articleID)
		articles = append(articles, &a)
	}
	return articles, nil
}

func (d *DB) GetOneById(id uuid.UUID) (art article.Article, err error) {
	const query = `SELECT id, title, url, slug FROM articles WHERE id = $1`
	result, err := d.db.Query(query, id)
	if err != nil {
		return art, errors.New("failed to exec query : " + err.Error())
	}
	var articleID string
	if !result.Next() {
		return art, errors.New("no article found")
	}
	err = result.Scan(&articleID, &art.Title, &art.URL, &art.Slug)
	if err != nil {
		return art, errors.New("failed to scan query result : " + err.Error())
	}
	art.ID = uuid.MustParse(articleID)
	return art, err
}


func (d *DB) GetOneBySlug(slug string) (art article.Article, err error) {
	const query = `SELECT id, title, url, slug FROM articles WHERE slug = $1`
	result, err := d.db.Query(query, slug)
	if err != nil {
		return art, errors.New("failed to exec query : " + err.Error())
	}
	var articleID string
	if !result.Next() {
		return art, errors.New("no article found")
	}
	err = result.Scan(&articleID, &art.Title, &art.URL, &art.Slug)
	if err != nil {
		return art, errors.New("failed to scan query result : " + err.Error())
	}
	art.ID = uuid.MustParse(articleID)
	return art, err
}

func (d *DB) UpdateArticle(a *article.Article) (err error) {
	query := `UPDATE articles SET title=$2, url=$3, slug=$4 WHERE id = $1`
	_, err = d.db.Exec(query, a.ID.String(), a.Title, a.URL, a.Slug)
	if err != nil {
		err = errors.New("failed to exec query : " + err.Error())
	}
	return err
}
