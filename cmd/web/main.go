package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"prasowka/cmd/conf"
	"prasowka/cmd/server"
	"prasowka/cmd/sqldb"
	"prasowka/internal/application"
	"prasowka/scan"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	var app application.Application

	var cfg conf.Config
	cfg.New()

	var database sqldb.Sqlite
	err := database.DbConn(cfg.DNS)
	if err != nil {
		slog.Info("Error connection to db")
		os.Exit(1)
	}

	app.New(database.DB)

	router := app.Router()

	slog.Info("Starting", "port:", cfg.Port)
	s := server.NewServer(router)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// refresh at 10 past every hour
	RunEveryHour(func() {
		scan.ReadSource(database.DB, "https://www.rmf24.pl/")
	}, 01)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(s, done)

	err = s.ListenAndServe()
	//
	//
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")

}

func RunEveryHour(fn func(), n int) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), n, 0, 0, now.Location())

			// If it's already past 03 minutes this hour, schedule for next hour
			if !next.After(now) {
				next = next.Add(time.Hour)
			}

			time.Sleep(time.Until(next))

			fn()
		}
	}()
}
