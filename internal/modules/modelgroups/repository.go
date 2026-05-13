package modelgroups

import "github.com/jmoiron/sqlx"

type ModelGroupRepository struct {
	db *sqlx.DB
}

func NewModelGroupRepository(db *sqlx.DB) *ModelGroupRepository {
	return &ModelGroupRepository{db: db}
}

func (r *ModelGroupRepository) Create(name string) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO model_groups (name) VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, name).Scan(&id)
	return id, err
}

func (r *ModelGroupRepository) AddMember(groupID, modelID int) error {
	_, err := r.db.Exec(`
		INSERT INTO model_group_members (model_group_id, model_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, groupID, modelID)
	return err
}

func (r *ModelGroupRepository) Clear() error {
	_, err := r.db.Exec(`DELETE FROM model_groups`)
	return err
}
