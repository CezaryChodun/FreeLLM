package proxy

import (
	"errors"
	"fmt"

	"github.com/cezarychodun/freellms/internal/modules/modelgroups"
	"github.com/cezarychodun/freellms/internal/modules/models"
	"github.com/cezarychodun/freellms/internal/modules/ratelimits"
	"github.com/cezarychodun/freellms/internal/modules/usage"
)

var ErrNoAvailableModel = errors.New("no available model with remaining usage")

type ModelSelector struct {
	modelRepo      *models.ModelRepository
	rateLimitRepo  *ratelimits.RateLimitRepository
	usageRepo      *usage.ModelResourcesRepository
	modelGroupRepo *modelgroups.ModelGroupRepository
}

func NewModelSelector(modelRepo *models.ModelRepository, rateLimitRepo *ratelimits.RateLimitRepository, usageRepo *usage.ModelResourcesRepository, modelGroupRepo *modelgroups.ModelGroupRepository) *ModelSelector {
	return &ModelSelector{
		modelRepo:      modelRepo,
		rateLimitRepo:  rateLimitRepo,
		usageRepo:      usageRepo,
		modelGroupRepo: modelGroupRepo,
	}
}

// Select picks the model furthest from hitting its rate limit.
// If groupName matches a model group, only models in that group are considered.
// Otherwise all models are considered.
func (s *ModelSelector) Select(groupName string) (models.Model, error) {
	return s.SelectExcluding(groupName, nil)
}

// SelectExcluding picks the best model excluding any model IDs in the blacklist.
func (s *ModelSelector) SelectExcluding(groupName string, excludedIDs []int) (models.Model, error) {
	candidates, err := s.getCandidates(groupName)
	if err != nil {
		return models.Model{}, err
	}

	filtered := candidates[:0]
	for _, c := range candidates {
		excluded := false
		for _, id := range excludedIDs {
			if c.model.ID == id {
				excluded = true
				break
			}
		}
		if !excluded {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		return models.Model{}, ErrNoAvailableModel
	}

	return chooseFurthestFromLimit(filtered), nil
}

type candidate struct {
	model          models.Model
	maxUtilization float64
}

func (s *ModelSelector) getCandidates(groupName string) ([]candidate, error) {
	var modelList []models.Model

	if groupName != "" {
		exists, err := s.modelGroupRepo.GroupExists(groupName)
		if err != nil {
			return nil, fmt.Errorf("checking group existence: %w", err)
		}
		if exists {
			modelList, err = s.modelGroupRepo.FindModelsByGroupName(groupName)
			if err != nil {
				return nil, fmt.Errorf("finding models by group: %w", err)
			}
		}
	}

	if modelList == nil {
		var err error
		modelList, err = s.modelRepo.List()
		if err != nil {
			return nil, fmt.Errorf("listing models: %w", err)
		}
	}

	var candidates []candidate
	for _, m := range modelList {
		limit, err := s.rateLimitRepo.FindByModel(m.Name, m.Provider)
		if err != nil {
			if errors.Is(err, ratelimits.ErrRateLimitNotFound) {
				continue
			}
			return nil, err
		}

		current, err := s.usageRepo.FindByModelID(m.ID)
		if err != nil {
			if errors.Is(err, usage.ErrModelResourcesNotFound) {
				candidates = append(candidates, candidate{model: m, maxUtilization: 0})
				continue
			}
			return nil, err
		}

		if !hasRemainingUsage(*limit, current) {
			continue
		}

		candidates = append(candidates, candidate{
			model:          m,
			maxUtilization: maxUtilizationPct(*limit, current),
		})
	}

	return candidates, nil
}

func chooseFurthestFromLimit(candidates []candidate) models.Model {
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.maxUtilization < best.maxUtilization {
			best = c
		}
	}
	return best.model
}

func maxUtilizationPct(limit ratelimits.RateLimit, current *usage.ModelResources) float64 {
	var maxPct float64
	if limit.InputTokensPerMinute > 0 {
		pct := float64(current.InputTokensPerMinute) / float64(limit.InputTokensPerMinute)
		if pct > maxPct {
			maxPct = pct
		}
	}
	if limit.RequestsPerMinute > 0 {
		pct := float64(current.RequestsPerMinute) / float64(limit.RequestsPerMinute)
		if pct > maxPct {
			maxPct = pct
		}
	}
	if limit.RequestsPerDay > 0 {
		pct := float64(current.RequestsPerDay) / float64(limit.RequestsPerDay)
		if pct > maxPct {
			maxPct = pct
		}
	}
	return maxPct
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
