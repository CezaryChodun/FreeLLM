package proxy

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cezarychodun/freellms/internal/modules/ratelimits"
	"github.com/cezarychodun/freellms/internal/modules/usage"
)

var ErrNoAvailableModel = errors.New("no available model with remaining usage")

type ModelSelector struct {
	mu            sync.Mutex
	queue         []string
	rateLimitRepo *ratelimits.RateLimitRepository
	usageRepo     *usage.ModelResourcesRepository
}

func NewModelSelector(rateLimitRepo *ratelimits.RateLimitRepository, usageRepo *usage.ModelResourcesRepository) *ModelSelector {
	return &ModelSelector{
		rateLimitRepo: rateLimitRepo,
		usageRepo:     usageRepo,
	}
}

func (s *ModelSelector) Next() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.queue) == 0 {
		available, err := s.queryAvailable()
		if err != nil {
			return "", err
		}
		s.queue = available
	}

	if len(s.queue) == 0 {
		return "", ErrNoAvailableModel
	}

	model := s.queue[0]
	s.queue = s.queue[1:]
	return model, nil
}

func (s *ModelSelector) queryAvailable() ([]string, error) {
	limits, err := s.rateLimitRepo.List()
	if err != nil {
		return nil, fmt.Errorf("listing rate limits: %w", err)
	}

	var available []string
	for _, limit := range limits {
		modelName := stripProvider(limit.Model)

		current, err := s.usageRepo.FindByModel(modelName)
		if err != nil {
			if errors.Is(err, usage.ErrModelResourcesNotFound) {
				available = append(available, modelName)
				continue
			}
			return nil, fmt.Errorf("finding usage for %s: %w", modelName, err)
		}

		if hasRemainingUsage(limit, current) {
			available = append(available, modelName)
		}
	}

	return available, nil
}

func stripProvider(model string) string {
	if idx := strings.Index(model, "/"); idx != -1 {
		return model[idx+1:]
	}
	return model
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
