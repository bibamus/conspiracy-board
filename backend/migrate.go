package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path"
	"regexp"
	"sort"
	"strconv"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type migration struct {
	version int
	name    string
	sql     string
}

var migrationFileName = regexp.MustCompile(`^(\d+)_([a-z0-9_]+)\.sql$`)

// loadMigrations reads NNNN_name.sql files from dir, sorted by version.
// Versions must start at 1 and be contiguous so a missing file is caught.
func loadMigrations(fsys fs.FS, dir string) ([]migration, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations: %w", err)
	}

	var migrations []migration
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		match := migrationFileName.FindStringSubmatch(entry.Name())
		if match == nil {
			return nil, fmt.Errorf("invalid migration file name %q (expected NNNN_name.sql)", entry.Name())
		}
		version, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("invalid migration version in %q: %w", entry.Name(), err)
		}
		body, err := fs.ReadFile(fsys, path.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration %q: %w", entry.Name(), err)
		}
		migrations = append(migrations, migration{version: version, name: match[2], sql: string(body)})
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	for i, m := range migrations {
		if m.version != i+1 {
			return nil, fmt.Errorf("migration versions must be contiguous from 1: expected %d, found %d_%s", i+1, m.version, m.name)
		}
	}
	return migrations, nil
}

// migrate applies all pending migrations and records each one in
// schema_migrations. Everything runs in one BEGIN IMMEDIATE transaction, so
// concurrent processes serialize on the write lock and a failed migration
// leaves the database untouched.
func migrate(ctx context.Context, db *sql.DB, migrations []migration) error {
	// A dedicated connection is required so BEGIN/COMMIT and the statements
	// in between run on the same SQLite connection.
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		return fmt.Errorf("failed to begin migration transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			conn.ExecContext(context.Background(), "ROLLBACK")
		}
	}()

	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}

	var current int
	if err := conn.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&current); err != nil {
		return fmt.Errorf("failed to read schema version: %w", err)
	}
	if current > len(migrations) {
		return fmt.Errorf("database schema version %d is newer than this binary supports (%d)", current, len(migrations))
	}

	for _, m := range migrations[current:] {
		if _, err := conn.ExecContext(ctx, m.sql); err != nil {
			return fmt.Errorf("migration %04d_%s failed: %w", m.version, m.name, err)
		}
		if _, err := conn.ExecContext(ctx, "INSERT INTO schema_migrations (version, name) VALUES (?, ?)", m.version, m.name); err != nil {
			return fmt.Errorf("failed to record migration %04d_%s: %w", m.version, m.name, err)
		}
		log.Printf("applied migration %04d_%s", m.version, m.name)
	}

	if _, err := conn.ExecContext(ctx, "COMMIT"); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}
	committed = true
	return nil
}
