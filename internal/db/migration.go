package db

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func getProjectRootPath() (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b)) // go back from /db to /internal
	migrationsPath := filepath.Join(basePath, "../migration")
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return "", fmt.Errorf("failed to get abs path: %w", err)
	}

	return absPath, nil
}

func RunMigration(dsn string) error {
	absPath, _ := getProjectRootPath()
	m, err := migrate.New(
		"file://"+absPath,
		fmt.Sprintf("pgx5://%s", dsn),
	)

	if err != nil {
		return fmt.Errorf("migration init failed: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
