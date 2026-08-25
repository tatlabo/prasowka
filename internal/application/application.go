package application

import (
	"database/sql"
)

type Application struct {
	DB *sql.DB
}

func (app *Application) New(db *sql.DB) *Application {
	app.DB = db
	return app
}
