package scan

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"prasowka/internal/models"
	"slices"
	"time"
)

var (
	ErrScaningPage = fmt.Errorf("Error in scanning page.")
	ErrConnDb      = fmt.Errorf("Error connection to db.")
)

func ReadSource(path string) {

	start := time.Now()
	defer func() {
		slog.Info("Scan completed", "duration", time.Since(start).String())
	}()

	w := models.Website{}

	w.URL = template.URL(path)

	if err := w.ProcessWebsite(); err != nil {
		slog.Info("Error processing website: ", "time", time.Now().Format("2006-01-02 15:04:05"))
	}

	err := ScrapPage(w)
	if err != nil {
		slog.Info("Error scaning page: ", "with", ErrScaningPage)
	}

}

func ScrapPage(w models.Website) error {

	db, err := sql.Open("sqlite", "../db/websites.db")
	if err != nil {
		return fmt.Errorf("%v,  %s", ErrConnDb, err)
	}
	defer db.Close()

	err = models.CreateSourceTable(db)
	if err != nil {
		return fmt.Errorf("CreateSourceTable failed; %w", err)
	}

	err = models.CreateArticleTable(db)
	if err != nil {
		return fmt.Errorf("CreateArticleTable failed; %w", err)
	}

	articles, err := models.RefreshSource(&w, db)
	if err != nil {
		return fmt.Errorf("RefreshSource failed; %w", err)
	}

	var list []models.Website

	for i := range articles {
		n := models.Website{}
		n.SourceId = w.Id
		n.URL = articles[i].URL
		n.Title = articles[i].Title
		n.CreatedAt = time.Now()

		list = append(list, n)
		fmt.Printf("Found article: %s - %s\n", n.Title, n.URL)
	}

	err = models.AddWebsiteList(list, db)
	if err != nil {
		return fmt.Errorf("AddWebsiteList failed, %w", err)
	}

	return nil

}

func RefreshSource(w *models.Website, db *sql.DB) (news []models.Website, err error) {

	models.CreateSourceTable(db)
	models.CreateArticleTable(db)

	ctx := context.Background()
	// get Body from source URL
	if err := w.ProcessWebsite(); err != nil {
		return news, fmt.Errorf("failed to process source website: %w", err)
	}
	// Insert source website to db, get source ID
	w.CreatedAt = time.Now()
	if err := w.SourceToDb(ctx, db); err != nil {
		return news, fmt.Errorf("failed to insert source website to db: %w", err)
	}

	// get ALL existing articles urls from db
	existing, err := models.ExistingURL(w, db)
	if err != nil {
		log.Fatal(err)
	}

	newArticles := []models.Website{}

	if w.Body != "" {
		subpages, err := models.ParseSourceBody(w)
		if err != nil {
			return news, err
		}

		l := len(subpages)
		if l == 0 {
			log.Println("No subpages found in source body")
			return news, nil
		}
		//compare existing articles URL with new subpages
		for i := range subpages {
			currentTitle := subpages[i].URL
			if slices.Contains(existing, currentTitle) {
				continue
			}
			newArticles = append(newArticles, subpages[i])
		}

		return newArticles, nil

	}

	return news, nil

}
