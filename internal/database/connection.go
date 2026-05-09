package database

import (
	"database/sql"
	"fmt"

	"github.com/cezarychodun/freellms/internal/config"
	_ "github.com/lib/pq"
)

func Open(dbConfig *config.DBConfig) (*sql.DB, error) {
	switch dbConfig.Dialect {
	case "postgres":
		return openPostgres(dbConfig)
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", dbConfig.Dialect)
	}
}

func openPostgres(dbConfig *config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Name,
		"disable",
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	fmt.Println("Migrating database...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS remaining_resources (
			model TEXT PRIMARY KEY,
			input_tokens_per_minute INTEGER NOT NULL,
			output_tokens_per_minute INTEGER NOT NULL,
			requests_per_day INTEGER NOT NULL,
			last_used TIMESTAMP NOT NULL
		)
	`)
	return err
}
