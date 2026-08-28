// Package application provides the application's shared dependencies.
package application

import (
	"database/sql"
)

type Application struct {
	DB *sql.DB
}

func New(db *sql.DB) *Application {
	return &Application{DB: db}
}
