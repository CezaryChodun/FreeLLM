package usage

import (
	"errors"
	"time"
)

// IncrementRequestCount increments RPM and RPD. Call this before sending the request.
func (r *ModelResourcesRepository) IncrementRequestCount(modelID int) error {
	now := time.Now().UTC()

	resources, err := r.FindByModelID(modelID)
	if err != nil {
		if errors.Is(err, ErrModelResourcesNotFound) {
			return r.Create(&ModelResources{
				ModelID:           modelID,
				RequestsPerMinute: 1,
				RequestsPerDay:    1,
				LastUsed:          now,
			})
		}
		return err
	}

	updated := *resources
	if isDifferentMinute(resources.LastUsed, now) {
		updated.InputTokensPerMinute = 0
		updated.OutputTokensPerMinute = 0
		updated.RequestsPerMinute = 1
	} else {
		updated.RequestsPerMinute++
	}
	if isDifferentDay(resources.LastUsed, now) {
		updated.RequestsPerDay = 1
	} else {
		updated.RequestsPerDay++
	}
	updated.LastUsed = now
	return r.Update(&updated)
}

// AddTokenUsage adds input/output token counts after a successful request.
func (r *ModelResourcesRepository) AddTokenUsage(modelID int, inputTokens int, outputTokens int, timestamp int) error {
	lastUsed := time.Unix(int64(timestamp), 0).UTC()

	resources, err := r.FindByModelID(modelID)
	if err != nil {
		if errors.Is(err, ErrModelResourcesNotFound) {
			return r.Create(&ModelResources{
				ModelID:               modelID,
				InputTokensPerMinute:  inputTokens,
				OutputTokensPerMinute: outputTokens,
				LastUsed:              lastUsed,
			})
		}
		return err
	}

	updated := *resources
	if isDifferentMinute(resources.LastUsed, lastUsed) {
		updated.InputTokensPerMinute = inputTokens
		updated.OutputTokensPerMinute = outputTokens
	} else {
		updated.InputTokensPerMinute += inputTokens
		updated.OutputTokensPerMinute += outputTokens
	}
	updated.LastUsed = lastUsed
	return r.Update(&updated)
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
