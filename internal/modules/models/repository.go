package models

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrModelNotFound = errors.New("model not found")

type ModelRepository struct {
	db *sqlx.DB
}

func NewModelRepository(db *sqlx.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) Create(m *Model) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO models (name, provider, instance)
		VALUES ($1, $2, $3)
		RETURNING id
	`, m.Name, m.Provider, m.Instance).Scan(&id)
	return id, err
}

func (r *ModelRepository) FindByID(id int) (*Model, error) {
	var m Model
	err := r.db.Get(&m, `SELECT id, name, provider, instance FROM models WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModelNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *ModelRepository) FindByNameProviderInstance(name, provider string, instance int) (*Model, error) {
	var m Model
	err := r.db.Get(&m, `SELECT id, name, provider, instance FROM models WHERE name = $1 AND provider = $2 AND instance = $3`, name, provider, instance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModelNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *ModelRepository) List() ([]Model, error) {
	var results []Model
	err := r.db.Select(&results, `SELECT id, name, provider, instance FROM models ORDER BY id`)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *ModelRepository) Clear() error {
	_, err := r.db.Exec(`DELETE FROM models`)
	return err
}
