package config

type Model struct {
	Constraint *ModelConstraints
}

type ModelConstraints struct {
	Name              string `yaml:"name"`
	TokensInPerMinute int    `yaml:"TPM"`
	RequestPerMinute  int    `yaml:"RPM"`
	RequestsPerDay    int    `yaml:"RPD"`
}
