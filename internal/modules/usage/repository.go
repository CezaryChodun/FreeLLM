package usage

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrModelResourcesNotFound = errors.New("model resources not found")

type ModelResources struct {
	ModelID               int       `db:"model_id" json:"model_id"`
	InputTokensPerMinute  int       `db:"input_tokens_per_minute" json:"input_tokens_per_minute"`
	OutputTokensPerMinute int       `db:"output_tokens_per_minute" json:"output_tokens_per_minute"`
	RequestsPerMinute     int       `db:"requests_per_minute" json:"requests_per_minute"`
	RequestsPerDay        int       `db:"requests_per_day" json:"requests_per_day"`
	LastUsed              time.Time `db:"last_used" json:"last_used"`
}

type ModelResourcesRepository struct {
	db *sqlx.DB
}

func NewModelResourcesRepository(db *sqlx.DB) *ModelResourcesRepository {
	return &ModelResourcesRepository{db: db}
}

func (r *ModelResourcesRepository) Create(resources *ModelResources) error {
	if resources.LastUsed.IsZero() {
		resources.LastUsed = time.Now()
	}
	_, err := r.db.NamedExec(`
		INSERT INTO usage_tracking (model_id, input_tokens_per_minute, output_tokens_per_minute, requests_per_minute, requests_per_day, last_used)
		VALUES (:model_id, :input_tokens_per_minute, :output_tokens_per_minute, :requests_per_minute, :requests_per_day, :last_used)
	`, resources)
	return err
}

func (r *ModelResourcesRepository) FindByModelID(modelID int) (*ModelResources, error) {
	var resources ModelResources
	err := r.db.Get(&resources, `
		SELECT model_id, input_tokens_per_minute, output_tokens_per_minute, requests_per_minute, requests_per_day, last_used
		FROM usage_tracking WHERE model_id = $1
	`, modelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModelResourcesNotFound
		}
		return nil, err
	}
	return &resources, nil
}

func (r *ModelResourcesRepository) Update(resources *ModelResources) error {
	if resources.LastUsed.IsZero() {
		resources.LastUsed = time.Now()
	}
	result, err := r.db.NamedExec(`
		UPDATE usage_tracking
		SET input_tokens_per_minute = :input_tokens_per_minute,
		    output_tokens_per_minute = :output_tokens_per_minute,
		    requests_per_minute = :requests_per_minute,
		    requests_per_day = :requests_per_day,
		    last_used = :last_used
		WHERE model_id = :model_id
	`, resources)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrModelResourcesNotFound
	}
	return nil
}

func (r *ModelResourcesRepository) List() ([]ModelResources, error) {
	var resources []ModelResources
	err := r.db.Select(&resources, `
		SELECT model_id, input_tokens_per_minute, output_tokens_per_minute, requests_per_minute, requests_per_day, last_used
		FROM usage_tracking ORDER BY model_id
	`)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
