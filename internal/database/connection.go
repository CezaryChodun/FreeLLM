package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cezarychodun/freellms/internal/config"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

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

	if err := ensureSchemaMigrationsTable(db); err != nil {
		return err
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := isMigrationApplied(db, migration.Name)
		if err != nil {
			return err
		}

		if applied {
			continue
		}

		if err := applyMigration(db, migration); err != nil {
			return err
		}
	}

	return nil
}

type migration struct {
	Name string
	SQL  string
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)

	return err
}

func loadMigrations() ([]migration, error) {
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		names = append(names, name)
	}

	sort.Strings(names)

	migrations := make([]migration, 0, len(names))
	for _, name := range names {
		path := filepath.Join("migrations", name)

		content, err := migrationFiles.ReadFile(path)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration{
			Name: name,
			SQL:  string(content),
		})
	}

	return migrations, nil
}

func isMigrationApplied(db *sql.DB, name string) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM schema_migrations
			WHERE version = $1
		)
	`, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func applyMigration(db *sql.DB, migration migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(migration.SQL); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to apply migration %s: %w", migration.Name, err)
	}

	if _, err := tx.Exec(`
		INSERT INTO schema_migrations (version)
		VALUES ($1)
	`, migration.Name); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("Applied migration: %s\n", migration.Name)

	return nil
}
