package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Host:     "127.0.0.1",
			Port:     7432,
			Username: "guest",
			Password: "Guest0000!",
			Name:     "freellm",
			Charset:  "utf8",
		},
	}
}

/**

Dialect:  getEnv("DB_DIALECT", "postgres"),
Host:     getEnv("DB_HOST", ""),
Port:     getEnvAsInt("DB_PORT", 0),
Username: getEnv("DB_USERNAME", ""),
Password: getEnv("DB_PASSWORD", ""),
Name:     getEnv("DB_NAME", ""),
SSLMode:  getEnv("DB_SSLMODE", ""),

*/
