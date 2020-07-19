package sql

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/tag"
	"strings"
)

type tagRepository struct {
	db DBer
}

func (r *tagRepository) GetOneById(id uuid.UUID) (tag.Tag, error) {
	panic("implement me")
}

func (r *tagRepository) GetOneByName(name string) (tag.Tag, error) {
	panic("implement me")
}

func (r *tagRepository) GetAll(hasCategories bool) ([]tag.Tag, error) {
	var tags []tag.Tag
	query := `SELECT tags.id, tags.name FROM tags`
	if hasCategories {
		query = query + ` INNER JOIN articles_tags ON tags.id = articles_tags.tag_id`
	}
	query = query + ` ORDER BY name`
	fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return tags, fmt.Errorf("failed to retrieve tags : %w", err)
	}
	for rows.Next() {
		var t tag.Tag
		err := rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return tags, fmt.Errorf("failed to scan row : %w", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *tagRepository) Create(t tag.Tag) error {
	query := `INSERT INTO tags (id, name) VALUES($1, $2)`
	_, err := r.db.Exec(query, t.ID, t.Name)
	if err != nil {
		return fmt.Errorf("failed to create article : %w", err)
	}
	return nil
}

func (r *tagRepository) Update(article tag.Tag) error {
	panic("implement me")
}

func (r *tagRepository) DeleteById(id uuid.UUID) error {
	panic("implement me")
}

func NewTagRepository(db DBer) *tagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindTagsByArticleID(tagsID []uuid.UUID) ([]tag.Tag, error) {
	var tags []tag.Tag

	if len(tagsID) == 0 {
		return []tag.Tag{}, nil
	}

	tagsIDInterface := make([]interface{}, len(tagsID))
	for i, id := range tagsID {
		tagsIDInterface[i] = id
	}

	addPlaceHolders := func(cnt int) string {
		placeholders := make([]string, cnt)
		for i := 1; i <= cnt; i++ {
			placeholders[i-1] = fmt.Sprintf(`$%d`, i)
		}
		return strings.Join(placeholders, `, `)
	}

	query := `SELECT id, name FROM tags WHERE id IN (` + addPlaceHolders(len(tagsID)) + `)`
	rows, err := r.db.Query(query, tagsIDInterface...)
	if err != nil {
		return []tag.Tag{}, fmt.Errorf("failed to exec FindTagsByArticleID query : %w", err)
	}

	for rows.Next() {
		var t tag.Tag
		rows.Scan(&t.ID, &t.Name)
		tags = append(tags, t)
	}

	return tags, nil
}
