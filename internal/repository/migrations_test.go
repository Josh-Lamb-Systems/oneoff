package repository

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/meysam81/oneoff/internal/domain"
)

func TestRunMigrationsFromEmbeddedFiles(t *testing.T) {
	ctx := context.Background()
	// A fresh directory with no migrations/ folder, like a release install
	dbPath := filepath.Join(t.TempDir(), "oneoff.db")

	if err := RunMigrations(ctx, dbPath, "up"); err != nil {
		t.Fatalf("migrate up: %v", err)
	}
	// Re-running must be a no-op
	if err := RunMigrations(ctx, dbPath, "up"); err != nil {
		t.Fatalf("second migrate up: %v", err)
	}

	repo, err := NewSQLiteRepository(ctx, dbPath)
	if err != nil {
		t.Fatalf("open repository: %v", err)
	}
	defer func() { _ = repo.Close() }()

	if err := repo.Ping(ctx); err != nil {
		t.Fatalf("ping: %v", err)
	}
	if _, err := repo.ListJobs(ctx, domain.JobFilter{Limit: 10}); err != nil {
		t.Fatalf("jobs table missing after migrations: %v", err)
	}
	if _, err := repo.GetProject(ctx, "default"); err != nil {
		t.Fatalf("default project missing after migrations: %v", err)
	}

	if err := RunMigrations(ctx, dbPath, "down"); err != nil {
		t.Fatalf("migrate down: %v", err)
	}
}
