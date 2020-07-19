package sql

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/article"
	"github.com/grepsd/knowledge-database/pkg/tag"
	"strings"
)

type articleRepository struct {
	db DBer
}

func NewArticleRepository(db DBer) *articleRepository {
	return &articleRepository{db: db}
}

var ErrDuplicateKey = errors.New("duplicate key value violates unique constraint")

func (r articleRepository) Create(a article.Article) error {
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

func (r articleRepository) GetAll() ([]article.Article, error) {
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

func (r articleRepository) GetOneById(id uuid.UUID) (art article.Article, err error) {
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

func (r articleRepository) GetOneBySlug(slug string) (art article.Article, err error) {
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

func (r articleRepository) Update(a article.Article) (err error) {
	query := `UPDATE articles SET title=$2, url=$3, slug=$4 WHERE id = $1`
	_, err = r.db.Exec(query, a.ID.String(), a.Title, a.URL, a.Slug)
	if err != nil {
		err = errors.New("failed to exec query : " + err.Error())
	}
	return err
}

func (r articleRepository) DeleteById(id uuid.UUID) error {
	query := `DELETE FROM articles WHERE id = $1`
	_, err := r.db.Exec(query, id.String())
	if err != nil {
		err = errors.New("failed to exec query : " + err.Error())
	}
	return err
}

func (r *articleRepository) GetArticleTags(id uuid.UUID) ([]*tag.Tag, error) {
	var tags []*tag.Tag
	query := `SELECT tags.id, tags.name FROM article_tags at INNER JOIN tags USING (tag_id) WHERE at.tag_id = ?`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return tags, fmt.Errorf("failed to retrieve article tags : %w", err)
	}
	for rows.Next() {
		var tag *tag.Tag
		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return tags, fmt.Errorf("failed to scan results : %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *articleRepository) AssignTagToArticle(articleID uuid.UUID, tagID uuid.UUID) error {
	query := `INSERT INTO articles_tags (article_id, tag_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, articleID, tagID)
	if err != nil {
		if strings.Contains("duplicate key value violates unique constraint", err.Error()) {
			return ErrDuplicateKey
		}
		return errors.New("failed to exec query : " + err.Error())
	}
	return nil
}

func (r *articleRepository) GetArticleCategories(id uuid.UUID) ([]uuid.UUID, error) {
	var articlesID []uuid.UUID
	query := `SELECT tag_id FROM articles_tags WHERE article_id = $1`

	rows, err := r.db.Query(query, id.String())

	if err != nil {
		return []uuid.UUID{}, errors.New("failed to exec GetArticleCategories query : " + err.Error())
	}

	for rows.Next() {
		var tagID uuid.UUID
		rows.Scan(&tagID)
		articlesID = append(articlesID, tagID)
	}

	return articlesID, nil
}
