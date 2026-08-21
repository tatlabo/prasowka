package application

import (
	"database/sql"

	"go.uber.org/zap"
)

type Application struct {
	// Config         config.Config
	// Template       *template.Template
	Logger *zap.SugaredLogger
	DB     *sql.DB
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

func (app *Application) New(logger *zap.SugaredLogger, db *sql.DB) *Application {

	app.Logger = logger

	app.DB = db

	return app
}
