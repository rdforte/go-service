// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// newMigration is responsible for setting up the migration from which we can migrate up/down
func newMigration(ctx context.Context, db *sqlx.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("newMigration: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://business/data/schema/sql",
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
