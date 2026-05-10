package ratelimits

import (
	"fmt"
	"os"
	"strings"

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

func LoadConfig(repo *RateLimitRepository, configPath, defaultsDir string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	var cfg modelsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	if err := repo.Clear(); err != nil {
		return fmt.Errorf("clearing rate_limits: %w", err)
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

	for provider, models := range providerModels {
		defaults, err := loadDefaults(defaultsDir, provider)
		if err != nil {
			return fmt.Errorf("loading defaults for %s: %w", provider, err)
		}

		defaultsMap := make(map[string]defaultEntry)
		for _, d := range defaults {
			defaultsMap[d.Name] = d
		}

		for _, model := range models {
			d, ok := defaultsMap[model]
			if !ok {
				continue
			}
			rl := &RateLimit{
				Model:                 provider + "/" + model,
				InputTokensPerMinute:  d.TPM,
				OutputTokensPerMinute: 0,
				RequestsPerMinute:     d.RPM,
				RequestsPerDay:        d.RPD,
			}
			if err := repo.Create(rl); err != nil {
				return fmt.Errorf("inserting rate limit for %s: %w", rl.Model, err)
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
