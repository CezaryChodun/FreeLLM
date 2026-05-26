package ratelimits

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrRateLimitNotFound = errors.New("rate limit not found")

type RateLimitRepository struct {
	db *sqlx.DB
}

func NewRateLimitRepository(db *sqlx.DB) *RateLimitRepository {
	return &RateLimitRepository{db: db}
}

func (r *RateLimitRepository) Create(rl *RateLimit) error {
	_, err := r.db.NamedExec(`
		INSERT INTO rate_limits (model_name, model_provider, input_tokens_per_minute, output_tokens_per_minute, requests_per_minute, requests_per_day)
		VALUES (:model_name, :model_provider, :input_tokens_per_minute, :output_tokens_per_minute, :requests_per_minute, :requests_per_day)
	`, rl)
	return err
}

func (r *RateLimitRepository) FindByModel(name, provider string) (*RateLimit, error) {
	var rl RateLimit
	err := r.db.Get(&rl, `SELECT * FROM rate_limits WHERE model_name = $1 AND model_provider = $2`, name, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRateLimitNotFound
		}
		return nil, err
	}
	return &rl, nil
}

func (r *RateLimitRepository) List() ([]RateLimit, error) {
	var results []RateLimit
	err := r.db.Select(&results, `SELECT * FROM rate_limits ORDER BY model_provider, model_name`)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *RateLimitRepository) Clear() error {
	_, err := r.db.Exec(`DELETE FROM rate_limits`)
	return err
}
