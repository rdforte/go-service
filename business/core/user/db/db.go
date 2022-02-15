package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rdforte/go-service/business/sys/database"
	"go.uber.org/zap"
)

// Store manages the set of API's for user access.
type Store struct {
	log    *zap.SugaredLogger
	sqlxDB *sqlx.DB
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Store {
	return Store{
		log:    log,
		sqlxDB: sqlxDB,
	}
}

// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, usr User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, usr); err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}
