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
		INSERT INTO rate_limits (name, input_tokens_per_minute, output_tokens_per_minute, requests_per_day)
		VALUES (:name, :input_tokens_per_minute, :output_tokens_per_minute, :requests_per_day)
	`, rl)
	return err
}

func (r *RateLimitRepository) FindByName(name string) (*RateLimit, error) {
	var rl RateLimit
	err := r.db.Get(&rl, `SELECT * FROM rate_limits WHERE name = $1`, name)
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
	err := r.db.Select(&results, `SELECT * FROM rate_limits ORDER BY name`)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *RateLimitRepository) Update(rl *RateLimit) error {
	result, err := r.db.NamedExec(`
		UPDATE rate_limits
		SET input_tokens_per_minute = :input_tokens_per_minute,
		    output_tokens_per_minute = :output_tokens_per_minute,
		    requests_per_day = :requests_per_day
		WHERE name = :name
	`, rl)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRateLimitNotFound
	}
	return nil
}

func (r *RateLimitRepository) Delete(name string) error {
	result, err := r.db.Exec(`DELETE FROM rate_limits WHERE name = $1`, name)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRateLimitNotFound
	}
	return nil
}
