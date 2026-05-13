package proxy

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cezarychodun/freellms/internal/modules/models"
	"github.com/cezarychodun/freellms/internal/modules/ratelimits"
	"github.com/cezarychodun/freellms/internal/modules/usage"
)

var ErrNoAvailableModel = errors.New("no available model with remaining usage")

type ModelSelector struct {
	mu            sync.Mutex
	queue         []models.Model
	modelRepo     *models.ModelRepository
	rateLimitRepo *ratelimits.RateLimitRepository
	usageRepo     *usage.ModelResourcesRepository
}

func NewModelSelector(modelRepo *models.ModelRepository, rateLimitRepo *ratelimits.RateLimitRepository, usageRepo *usage.ModelResourcesRepository) *ModelSelector {
	return &ModelSelector{
		modelRepo:     modelRepo,
		rateLimitRepo: rateLimitRepo,
		usageRepo:     usageRepo,
	}
}

// Next returns the next available model for request routing.
func (s *ModelSelector) Next() (models.Model, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.queue) == 0 {
		available, err := s.queryAvailable()
		if err != nil {
			return models.Model{}, err
		}
		s.queue = available
	}

	if len(s.queue) == 0 {
		return models.Model{}, ErrNoAvailableModel
	}

	model := s.queue[0]
	s.queue = s.queue[1:]
	return model, nil
}

func (s *ModelSelector) queryAvailable() ([]models.Model, error) {
	allModels, err := s.modelRepo.List()
	if err != nil {
		return nil, fmt.Errorf("listing models: %w", err)
	}

	var available []models.Model
	for _, m := range allModels {
		limit, err := s.rateLimitRepo.FindByModel(m.Name, m.Provider)
		if err != nil {
			if errors.Is(err, ratelimits.ErrRateLimitNotFound) {
				continue
			}
			return nil, fmt.Errorf("finding rate limit for %s/%s: %w", m.Provider, m.Name, err)
		}

		current, err := s.usageRepo.FindByModelID(m.ID)
		if err != nil {
			if errors.Is(err, usage.ErrModelResourcesNotFound) {
				available = append(available, m)
				continue
			}
			return nil, fmt.Errorf("finding usage for model %d: %w", m.ID, err)
		}

		if hasRemainingUsage(*limit, current) {
			available = append(available, m)
		}
	}

	return available, nil
}

func hasRemainingUsage(limit ratelimits.RateLimit, current *usage.ModelResources) bool {
	if limit.InputTokensPerMinute > 0 && current.InputTokensPerMinute >= limit.InputTokensPerMinute {
		return false
	}
	if limit.RequestsPerMinute > 0 && current.RequestsPerMinute >= limit.RequestsPerMinute {
		return false
	}
	if limit.RequestsPerDay > 0 && current.RequestsPerDay >= limit.RequestsPerDay {
		return false
	}
	return true
}
