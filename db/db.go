package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase() (*Database, error) {
	db, err := sqlx.Open("postgres",
		"user=postgres password=postgres host=localhost port=5432 dbname=nikita3db sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
func (d *Database) GetDb() *sqlx.DB {
	return d.db
}
