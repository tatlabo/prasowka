package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"prasowka/internal/models"
	"slices"
	"strings"
	"time"
)

var (
	ErrScaningPage = fmt.Errorf("Error in scanning page.")
	ErrConnDb      = fmt.Errorf("Error connection to db.")
)

func main() {

	start := time.Now()
	defer func() {
		fmt.Printf("Scan completed in %s\n", time.Since(start))
	}()

	var path string
	w := models.Website{}

	if len(os.Args) == 1 {
		fmt.Println("Please provide a path")
		os.Exit(1)
	}

	path = os.Args[1]
	path = strings.TrimSpace(path)
	path = strings.TrimRight(path, "/")

	w.URL = template.URL(path)

	if err := w.ProcessWebsite(); err != nil {
		fmt.Println(fmt.Errorf("Error processing website: %w", err))
		os.Exit(1)
	}

	err := ScrapPage(w)
	if err != nil {
		fmt.Printf("%v", ErrScaningPage)
		os.Exit(1)
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

func RefreshSource(w *models.Website, db *sql.DB) ([]models.Website, error) {

	models.CreateSourceTable(db)
	models.CreateArticleTable(db)

	ctx := context.Background()
	// get Body from source URL
	if err := w.ProcessWebsite(); err != nil {
		return []models.Website{}, fmt.Errorf("failed to process source website: %w", err)
	}
	// Insert source website to db, get source ID
	w.CreatedAt = time.Now()
	if err := w.SourceToDb(ctx, db); err != nil {
		return []models.Website{}, fmt.Errorf("failed to insert source website to db: %w", err)
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
			return []models.Website{}, err
		}

		l := len(subpages)
		if l == 0 {
			log.Println("No subpages found in source body")
			return []models.Website{}, nil
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

	return []models.Website{}, nil

}
