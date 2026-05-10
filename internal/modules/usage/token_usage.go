package usage

import (
	"errors"
	"time"
)

func (r *ModelResourcesRepository) AddTokenUsage(model string, inputTokens int, outputTokens int, timestamp int) error {
	lastUsed := time.Unix(int64(timestamp), 0).UTC()

	resources, err := r.FindByModel(model)
	if err != nil {
		if errors.Is(err, ErrModelResourcesNotFound) {
			return r.Create(newModelResourcesFromUsage(model, inputTokens, outputTokens, lastUsed))
		}

		return err
	}

	updatedResources := applyTokenUsage(resources, inputTokens, outputTokens, lastUsed)

	return r.Update(updatedResources)
}

func newModelResourcesFromUsage(model string, inputTokens int, outputTokens int, lastUsed time.Time) *ModelResources {
	return &ModelResources{
		Model:                 model,
		InputTokensPerMinute:  inputTokens,
		OutputTokensPerMinute: outputTokens,
		RequestsPerMinute:     1,
		RequestsPerDay:        1,
		LastUsed:              lastUsed,
	}
}

func applyTokenUsage(resources *ModelResources, inputTokens int, outputTokens int, lastUsed time.Time) *ModelResources {
	updatedResources := *resources

	if isDifferentMinute(resources.LastUsed, lastUsed) {
		updatedResources.InputTokensPerMinute = inputTokens
		updatedResources.OutputTokensPerMinute = outputTokens
		updatedResources.RequestsPerMinute = 1
	} else {
		updatedResources.InputTokensPerMinute += inputTokens
		updatedResources.OutputTokensPerMinute += outputTokens
		updatedResources.RequestsPerMinute++
	}

	if isDifferentDay(resources.LastUsed, lastUsed) {
		updatedResources.RequestsPerDay = 1
	} else {
		updatedResources.RequestsPerDay++
	}

	updatedResources.LastUsed = lastUsed

	return &updatedResources
}

func isDifferentMinute(previous time.Time, current time.Time) bool {
	previous = previous.UTC()
	current = current.UTC()

	return previous.Year() != current.Year() ||
		previous.YearDay() != current.YearDay() ||
		previous.Hour() != current.Hour() ||
		previous.Minute() != current.Minute()
}

func isDifferentDay(previous time.Time, current time.Time) bool {
	previous = previous.UTC()
	current = current.UTC()

	return previous.Year() != current.Year() ||
		previous.YearDay() != current.YearDay()
}
