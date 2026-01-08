package handlers

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func (w *Website) LastSourceWebsite(db *sql.DB) error {

	sql := `SELECT id, url, body, created_at, keywords, display FROM source WHERE url = ? ORDER BY created_at DESC LIMIT 1;`

	timeStr := ""
	err := db.QueryRow(sql, w.URL).Scan(&w.Id, &w.URL, &w.Body, &timeStr, &w.Keywords, &w.Display)
	if err != nil {
		return err
	}

	w.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", timeStr)

	return nil

}

func (w *Website) SelectByURL(db *sql.DB) error {

	sql := `SELECT id, title, url, body, created_at, keywords, display FROM daily WHERE url = ? ORDER BY created_at DESC LIMIT 1;`
	timeStr := ""
	err := db.QueryRow(sql, w.URL).Scan(&w.Id, &w.Title, &w.URL, &w.Body, &timeStr, &w.Keywords, &w.Display)
	if err != nil {
		return err
	}

	w.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", timeStr)

	return nil

}

func (w *Website) SelectById(db *sql.DB) error {

	table := "daily"

	sql := fmt.Sprintf(`SELECT daily.id, CONCAT(source.url, daily.url) as url, 
	daily.title, daily.body, daily.created_at, daily.keywords, daily.display, daily.done 
	FROM %s JOIN source ON daily.source_id = source.id WHERE daily.id = ?;`, table)

	timeStr := ""
	err := db.QueryRow(sql, w.Id).Scan(&w.Id, &w.URL, &w.Title, &w.Body, &timeStr, &w.Keywords, &w.Display, &w.Done)
	if err != nil {
		return err
	}

	w.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", timeStr)

	return nil

}

func (w *Website) UpdateRaw(db *sql.DB) error {

	sql := `UPDATE daily SET body=?, done=? WHERE id = ?;`
	w.Done = 1
	_, err := db.Exec(sql, w.Body, w.Done, w.Id)
	if err != nil {
		return err
	}

	return nil

}

func (w *Website) AddWebsite(db *sql.DB) error {

	stmt := `INSERT OR IGNORE INTO daily (source_id, url, body, title, created_at, keywords, display) 
	VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id;`

	err := db.QueryRow(stmt, w.SourceId, w.URL, w.Body, w.Title, w.CreatedAt.Format("2006-01-02 15:04:05"), w.Keywords, w.Display).Scan(&w.Id)
	if err != nil {
		return err
	}

	return nil
}

func (w *Website) ProcessWebsite() error {
	www := string(w.URL)
	res, err := http.Get(www)
	if err != nil {
		return fmt.Errorf("Error getting website: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Error res.StatusCode != http.StatusOK %w", err)
	}

	scanner := bufio.NewScanner(res.Body)
	doc := ""
	for scanner.Scan() {
		doc += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning website: %w", err)
	}

	w.Body = doc
	return nil
}
