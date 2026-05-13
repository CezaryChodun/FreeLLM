package ratelimits

import (
	"fmt"
	"os"
	"strings"

	"github.com/cezarychodun/freellms/internal/modules/models"
	"gopkg.in/yaml.v3"
)

type modelEntry struct {
	Model string `yaml:"model"`
}

type modelsConfig struct {
	Models []modelEntry `yaml:"models"`
}

type defaultEntry struct {
	Name string `yaml:"name"`
	TPM  int    `yaml:"TPM"`
	RPM  int    `yaml:"RPM"`
	RPD  int    `yaml:"RPD"`
}

func LoadConfig(modelRepo *models.ModelRepository, rateLimitRepo *RateLimitRepository, configPath, defaultsDir string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	var cfg modelsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	if err := rateLimitRepo.Clear(); err != nil {
		return fmt.Errorf("clearing rate_limits: %w", err)
	}
	if err := modelRepo.Clear(); err != nil {
		return fmt.Errorf("clearing models: %w", err)
	}

	// Group models by provider
	providerModels := make(map[string][]string)
	for _, m := range cfg.Models {
		parts := strings.SplitN(m.Model, "/", 2)
		if len(parts) != 2 {
			continue
		}
		providerModels[parts[0]] = append(providerModels[parts[0]], parts[1])
	}

	for provider, modelNames := range providerModels {
		defaults, err := loadDefaults(defaultsDir, provider)
		if err != nil {
			return fmt.Errorf("loading defaults for %s: %w", provider, err)
		}

		defaultsMap := make(map[string]defaultEntry)
		for _, d := range defaults {
			defaultsMap[d.Name] = d
		}

		for _, name := range modelNames {
			d, ok := defaultsMap[name]
			if !ok {
				continue
			}

			// Create model entry (instance 1)
			m := &models.Model{Name: name, Provider: provider, Instance: 1}
			if _, err := modelRepo.Create(m); err != nil {
				return fmt.Errorf("inserting model %s/%s: %w", provider, name, err)
			}

			// Create rate limit (idempotent per name+provider)
			_, err := rateLimitRepo.FindByModel(name, provider)
			if err == nil {
				continue // already exists
			}
			rl := &RateLimit{
				ModelName:             name,
				ModelProvider:         provider,
				InputTokensPerMinute:  d.TPM,
				OutputTokensPerMinute: 0,
				RequestsPerMinute:     d.RPM,
				RequestsPerDay:        d.RPD,
			}
			if err := rateLimitRepo.Create(rl); err != nil {
				return fmt.Errorf("inserting rate limit for %s/%s: %w", provider, name, err)
			}
		}
	}

	return nil
}

func loadDefaults(dir, provider string) ([]defaultEntry, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s.yml", dir, provider))
	if err != nil {
		return nil, err
	}

	var entries []defaultEntry
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
