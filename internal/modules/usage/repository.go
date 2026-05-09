package usage

import (
	"database/sql"
	"errors"
	"time"
)

type RemainingResources struct {
	Model                 string
	InputTokensPerMinute  int
	OutputTokensPerMinute int
	RequestsPerDay        int
	LastUsed              time.Time
}

type UsageService struct {
	db *sql.DB
}

var (
	service *UsageService
)

func NewUsageService(db *sql.DB) *UsageService {
	service = &UsageService{
		db: db,
	}
	return service
}

func GetUSageService() *UsageService {
	return service
}

func (s *UsageService) GetModelResources(model string) (*RemainingResources, error) {
	resources, err := s.findModelResources(model)
	if err == nil {
		return resources, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if err := s.CreateModelResources(model, 0, 0, 0); err != nil {
		return nil, err
	}

	return s.findModelResources(model)
}

func (s *UsageService) UpdateModelResources(
	model string,
	inputTokensPerMinute int,
	outputTokensPerMinute int,
	requestsPerDay int,
) (*RemainingResources, error) {
	_, err := s.GetModelResources(model)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(`
		UPDATE remaining_resources
		SET
			input_tokens_per_minute = $2,
			output_tokens_per_minute = $3,
			requests_per_day = $4,
			last_used = $5
		WHERE model = $1
	`,
		model,
		inputTokensPerMinute,
		outputTokensPerMinute,
		requestsPerDay,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return s.findModelResources(model)
}

func (s *UsageService) CreateModelResources(
	model string,
	inputTokensPerMinute int,
	outputTokensPerMinute int,
	requestsPerDay int,
) error {
	_, err := s.db.Exec(`
		INSERT INTO remaining_resources (
			model,
			input_tokens_per_minute,
			output_tokens_per_minute,
			requests_per_day,
			last_used
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (model) DO NOTHING
	`,
		model,
		inputTokensPerMinute,
		outputTokensPerMinute,
		requestsPerDay,
		time.Now(),
	)

	return err
}

func (s *UsageService) findModelResources(model string) (*RemainingResources, error) {
	var resources RemainingResources

	err := s.db.QueryRow(`
		SELECT
			model,
			input_tokens_per_minute,
			output_tokens_per_minute,
			requests_per_day,
			last_used
		FROM remaining_resources
		WHERE model = $1
	`, model).Scan(
		&resources.Model,
		&resources.InputTokensPerMinute,
		&resources.OutputTokensPerMinute,
		&resources.RequestsPerDay,
		&resources.LastUsed,
	)
	if err != nil {
		return nil, err
	}

	return &resources, nil
}
