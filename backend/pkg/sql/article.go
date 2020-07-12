package sql

import (
	"errors"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
)

type ArticleRepository struct {
	db DBer
}

func NewArticleRepository(db DBer) ArticleRepository {
	return ArticleRepository{db: db}
}
func (r ArticleRepository) Create(a article.Article) error {
	query := `INSERT INTO articles (id, title, url, slug)
VALUES ($1, $2, $3, $4)`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return errors.New("failed to prepare query : " + err.Error())
	}
	_, err = stmt.Exec(a.ID.String(), a.Title, a.URL, a.Slug)
	if err != nil {
		return errors.New("failed to exec query : " + err.Error())
	}
	return nil
}

func (r ArticleRepository) GetAll() ([]article.Article, error) {
	var articles []article.Article
	query := `SELECT id, title, url, slug FROM articles ORDER BY title`
	results, err := r.db.Query(query)
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
		articles = append(articles, a)
	}
	return articles, nil
}

func (r ArticleRepository) GetOneById(id uuid.UUID) (art article.Article, err error) {
	const query = `SELECT id, title, url, slug FROM articles WHERE id = $1`
	result, err := r.db.Query(query, id)
	if err != nil {
		return art, errors.New("failed to exec query : " + err.Error())
	}
	var articleID string
	if !result.Next() {
		return art, errors.New("article not found")
	}
	err = result.Scan(&articleID, &art.Title, &art.URL, &art.Slug)
	if err != nil {
		return art, errors.New("failed to scan query result : " + err.Error())
	}
	art.ID = uuid.MustParse(articleID)
	return art, err
}

func (r ArticleRepository) GetOneBySlug(slug string) (art article.Article, err error) {
	const query = `SELECT id, title, url, slug FROM articles WHERE slug = $1`
	result, err := r.db.Query(query, slug)
	if err != nil {
		return art, errors.New("failed to exec query : " + err.Error())
	}
	var articleID string
	if !result.Next() {
		return art, errors.New("article not found")
	}
	err = result.Scan(&articleID, &art.Title, &art.URL, &art.Slug)
	if err != nil {
		return art, errors.New("failed to scan query result : " + err.Error())
	}
	art.ID = uuid.MustParse(articleID)
	return art, err
}

func (r ArticleRepository) Update(a article.Article) (err error) {
	query := `UPDATE articles SET title=$2, url=$3, slug=$4 WHERE id = $1`
	_, err = r.db.Exec(query, a.ID.String(), a.Title, a.URL, a.Slug)
	if err != nil {
		err = errors.New("failed to exec query : " + err.Error())
	}
	return err
}

func (r ArticleRepository) DeleteById(id uuid.UUID) error {
	query := `DELETE FROM articles WHERE id = $1`
	_, err := r.db.Exec(query, id.String())
	if err != nil {
		err = errors.New("failed to exec query : " + err.Error())
	}
	return err
}