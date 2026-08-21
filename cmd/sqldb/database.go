package sqldb

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
	DB *sql.DB
}

func (s *Sqlite) DbConn(dsn string) error {
	if s.DB != nil {
		return nil
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}

	// Sprawdza faktyczne połączenie
	if err := db.Ping(); err != nil {
		return err
	}

	// WAL improves concurrent read/write behavior for SQLite.
	if _, err := db.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		return fmt.Errorf("enable WAL mode: %w", err)
	}

	if _, err := db.Exec("PRAGMA synchronous = NORMAL;"); err != nil {
		return fmt.Errorf("set synchronous mode: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.Exec("PRAGMA busy_timeout = 5000;"); err != nil {
		return fmt.Errorf("set busy timeout: %w", err)
	}

	// SQLite allows one writer, but WAL supports concurrent readers.
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	s.DB = db
	return nil
}
