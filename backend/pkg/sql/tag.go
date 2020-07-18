package sql

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/grepsd/knowledge-database/pkg/tag"
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

func (r *tagRepository) GetAll() ([]tag.Tag, error) {
	var tags []tag.Tag
	query := `SELECT id, name FROM tags`
	rows, err := r.db.Query(query)
	if err != nil {
		return tags, fmt.Errorf("failed to retrieve tags : %w", err)
	}
	for rows.Next() {
		var tag tag.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return tags, fmt.Errorf("failed to scan row : %w", err)
		}
		tags = append(tags, tag)
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
