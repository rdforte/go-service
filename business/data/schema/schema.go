// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rdforte/go-service/business/sys/database"
)

//go:embed seed/seed.sql
var seedDoc string

// newMigration is responsible for setting up the migration from which we can migrate up/down
func newMigration(ctx context.Context, db *sqlx.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("newMigration: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://business/data/schema/migrations",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("newMigration: %w", err)
	}

	return m, nil
}

// MigrateUp is responsible for migrating the database shema up
func MigrateUp(ctx context.Context, db *sqlx.DB) error {
	m, err := newMigration(ctx, db)
	if err != nil {
		return fmt.Errorf("MigrateUp: %w", err)
	}
	m.Steps(1)

	return nil
}

// MigrateDown is responsible for migrating the database shema down
func MigrateDown(ctx context.Context, db *sqlx.DB) error {
	m, err := newMigration(ctx, db)
	if err != nil {
		return fmt.Errorf("MigrateUp: %w", err)
	}
	m.Steps(-1)

	return nil
}

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seedDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
