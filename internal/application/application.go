package application

import (
	"database/sql"
)

type Application struct {
	// Config         config.Config
	// Template       *template.Template
	DB     *sql.DB
	Logger func(msg string)

	// Producter      Producter
	// Servicer       models.Servicer
	// UserModel      models.UserModel
	// User           models.User
	// Noter          models.Noter
	// Note           models.Note
	// sessionManager *scs.SessionManager
	// Role           role
	// TemplateData   templateData
}

func (app *Application) New(db *sql.DB) *Application {

	app.DB = db

	return app
}
