package main

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "test.db")+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func appliedVersions(t *testing.T, conn *sql.DB) []int {
	t.Helper()
	rows, err := conn.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var versions []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			t.Fatal(err)
		}
		versions = append(versions, v)
	}
	return versions
}

func TestEmbeddedMigrationsApplyAndAreIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.db")
	for i := 0; i < 2; i++ {
		db, err := NewDatabase(path)
		if err != nil {
			t.Fatalf("open #%d: %v", i+1, err)
		}
		migrations, _ := loadMigrations(migrationFiles, "migrations")
		if got := appliedVersions(t, db.conn); len(got) != len(migrations) {
			t.Fatalf("open #%d: applied %v, want %d migrations", i+1, got, len(migrations))
		}
		db.Close()
	}
}

func TestMigrateAppliesOnlyPendingMigrations(t *testing.T) {
	conn := openTestDB(t)
	ctx := context.Background()

	first := []migration{{1, "create_a", "CREATE TABLE a (id INTEGER);"}}
	if err := migrate(ctx, conn, first); err != nil {
		t.Fatal(err)
	}

	// Re-running 0001 would fail with "table a already exists".
	second := append(first, migration{2, "create_b", "CREATE TABLE b (id INTEGER);"})
	if err := migrate(ctx, conn, second); err != nil {
		t.Fatal(err)
	}

	if got := appliedVersions(t, conn); len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("applied %v, want [1 2]", got)
	}
}

func TestMigrateRollsBackOnFailure(t *testing.T) {
	conn := openTestDB(t)
	migrations := []migration{
		{1, "create_a", "CREATE TABLE a (id INTEGER);"},
		{2, "broken", "CREATE TABLE b (id INTEGER); NOT VALID SQL;"},
	}

	err := migrate(context.Background(), conn, migrations)
	if err == nil || !strings.Contains(err.Error(), "0002_broken") {
		t.Fatalf("expected failure naming 0002_broken, got %v", err)
	}

	var n int
	conn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE name IN ('a', 'b', 'schema_migrations')").Scan(&n)
	if n != 0 {
		t.Fatalf("expected full rollback, found %d objects", n)
	}
}

func TestMigrateRejectsNewerDatabase(t *testing.T) {
	conn := openTestDB(t)
	ctx := context.Background()
	two := []migration{{1, "a", "CREATE TABLE a (id INTEGER);"}, {2, "b", "CREATE TABLE b (id INTEGER);"}}
	if err := migrate(ctx, conn, two); err != nil {
		t.Fatal(err)
	}
	if err := migrate(ctx, conn, two[:1]); err == nil || !strings.Contains(err.Error(), "newer") {
		t.Fatalf("expected newer-schema error, got %v", err)
	}
}

func TestLoadMigrationsValidatesFiles(t *testing.T) {
	cases := map[string]struct {
		files   fstest.MapFS
		wantErr string
	}{
		"sorted": {files: fstest.MapFS{
			"m/0002_b.sql": {Data: []byte("B")},
			"m/0001_a.sql": {Data: []byte("A")},
		}},
		"gap": {files: fstest.MapFS{
			"m/0001_a.sql": {Data: []byte("A")},
			"m/0003_c.sql": {Data: []byte("C")},
		}, wantErr: "contiguous"},
		"duplicate": {files: fstest.MapFS{
			"m/0001_a.sql": {Data: []byte("A")},
			"m/0001_b.sql": {Data: []byte("B")},
		}, wantErr: "contiguous"},
		"bad name": {files: fstest.MapFS{
			"m/init.sql": {Data: []byte("A")},
		}, wantErr: "invalid migration file name"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := loadMigrations(tc.files, "m")
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 2 || got[0].sql != "A" || got[1].sql != "B" {
				t.Fatalf("unexpected order: %+v", got)
			}
		})
	}
}
