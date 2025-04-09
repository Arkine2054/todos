package storer

import (
	"github.com/jmoiron/sqlx"
	logging "todos3/todos-api/pkg"
)

type PostgresStorer struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewPostgresStorer(db *sqlx.DB, logger logging.Logger) *PostgresStorer {
	return &PostgresStorer{
		db:     db,
		logger: logger,
	}
}
