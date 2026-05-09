package usage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrModelResourcesNotFound = errors.New("model resources not found")

type ModelResources struct {
	Model                 string    `db:"model" json:"model"`
	InputTokensPerMinute  int       `db:"input_tokens_per_minute" json:"input_tokens_per_minute"`
	OutputTokensPerMinute int       `db:"output_tokens_per_minute" json:"output_tokens_per_minute"`
	RequestsPerDay        int       `db:"requests_per_day" json:"requests_per_day"`
	LastUsed              time.Time `db:"last_used" json:"last_used"`
}

type ModelResourcesRepository struct {
	db *sqlx.DB
}

func NewModelResourcesRepository(db *sqlx.DB) *ModelResourcesRepository {
	return &ModelResourcesRepository{
		db: db,
	}
}

func (r *ModelResourcesRepository) Create(resources *ModelResources) error {
	if resources.LastUsed.IsZero() {
		resources.LastUsed = time.Now()
	}

	_, err := r.db.NamedExec(`
		INSERT INTO remaining_resources (
			model,
			input_tokens_per_minute,
			output_tokens_per_minute,
			requests_per_day,
			last_used
		)
		VALUES (
			:model,
			:input_tokens_per_minute,
			:output_tokens_per_minute,
			:requests_per_day,
			:last_used
		)
	`, resources)

	return err
}

func (r *ModelResourcesRepository) FindByModel(model string) (*ModelResources, error) {
	var resources ModelResources

	err := r.db.Get(&resources, `
		SELECT
			model,
			input_tokens_per_minute,
			output_tokens_per_minute,
			requests_per_day,
			last_used
		FROM remaining_resources
		WHERE model = $1
	`, model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrModelResourcesNotFound
		}

		return nil, err
	}

	return &resources, nil
}

func (r *ModelResourcesRepository) List() ([]ModelResources, error) {
	resources := make([]ModelResources, 0)

	err := r.db.Select(&resources, `
		SELECT
			model,
			input_tokens_per_minute,
			output_tokens_per_minute,
			requests_per_day,
			last_used
		FROM remaining_resources
		ORDER BY model
	`)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *ModelResourcesRepository) Update(resources *ModelResources) error {
	if resources.LastUsed.IsZero() {
		resources.LastUsed = time.Now()
	}

	result, err := r.db.NamedExec(`
		UPDATE remaining_resources
		SET
			input_tokens_per_minute = :input_tokens_per_minute,
			output_tokens_per_minute = :output_tokens_per_minute,
			requests_per_day = :requests_per_day,
			last_used = :last_used
		WHERE model = :model
	`, resources)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrModelResourcesNotFound
	}

	return nil
}

func (r *ModelResourcesRepository) Delete(model string) error {
	result, err := r.db.Exec(`
		DELETE FROM remaining_resources
		WHERE model = $1
	`, model)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrModelResourcesNotFound
	}

	return nil
}

func (r *ModelResourcesRepository) AddTokenUsage(model string, inputTokens int, outputTokens int, timestamp int) error {
	lastUsed := time.Unix(int64(timestamp), 0).UTC()

	resources, err := r.FindByModel(model)
	if err != nil {
		if errors.Is(err, ErrModelResourcesNotFound) {
			return r.Create(&ModelResources{
				Model:                 model,
				InputTokensPerMinute:  inputTokens,
				OutputTokensPerMinute: outputTokens,
				RequestsPerDay:        1,
				LastUsed:              lastUsed,
			})
		}

		return err
	}

	inputTokensPerMinute := resources.InputTokensPerMinute
	outputTokensPerMinute := resources.OutputTokensPerMinute
	requestsPerDay := resources.RequestsPerDay

	if isDifferentMinute(resources.LastUsed, lastUsed) {
		inputTokensPerMinute = inputTokens
		outputTokensPerMinute = outputTokens
	} else {
		inputTokensPerMinute += inputTokens
		outputTokensPerMinute += outputTokens
	}

	if isDifferentDay(resources.LastUsed, lastUsed) {
		requestsPerDay = 1
	} else {
		requestsPerDay++
	}

	return r.Update(&ModelResources{
		Model:                 model,
		InputTokensPerMinute:  inputTokensPerMinute,
		OutputTokensPerMinute: outputTokensPerMinute,
		RequestsPerDay:        requestsPerDay,
		LastUsed:              lastUsed,
	})
}

func isDifferentMinute(previous time.Time, current time.Time) bool {
	previous = previous.UTC()
	current = current.UTC()

	out := previous.Year() != current.Year() ||
		previous.YearDay() != current.YearDay() ||
		previous.Hour() != current.Hour() ||
		previous.Minute() != current.Minute()

	fmt.Printf("Out: %v\n, prev: %v, curr: %v\n", out, previous, current)

	return out
}

func isDifferentDay(previous time.Time, current time.Time) bool {
	previous = previous.UTC()
	current = current.UTC()

	return previous.Year() != current.Year() ||
		previous.YearDay() != current.YearDay()
}
