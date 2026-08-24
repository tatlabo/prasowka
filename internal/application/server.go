package application

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   *sql.DB
}

func (app *Application) NewServer() *http.Server {
	// port, _ := strconv.Atoi(os.Getenv("PORT"))
	port := 8888
	NewServer := &Server{
		port: port,
		db:   app.DB,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] Listening and serving HTTP on %s\n", server.Addr)

	return server
}
