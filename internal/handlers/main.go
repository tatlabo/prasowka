package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"
	"time"
)

func RefreshSource(w *Website, db *sql.DB) ([]Website, error) {

	CreateSourceTable(db)
	CreateArticleTable(db)

	ctx := context.Background()
	// get Body from source URL
	if err := w.ProcessWebsite(); err != nil {
		return []Website{}, fmt.Errorf("failed to process source website: %w", err)
	}
	// Insert source website to db, get source ID
	w.CreatedAt = time.Now()
	if err := w.SourceToDb(ctx, db); err != nil {
		return []Website{}, fmt.Errorf("failed to insert source website to db: %w", err)
	}
	// get ALL existing articles urls from db
	existing, err := ExistingURL(w, db)
	if err != nil {
		log.Fatal(err)
	}

	newArticles := []Website{}

	if w.Body != "" {
		subpages, err := ParseSourceBody(w)
		if err != nil {
			return []Website{}, err
		}

		l := len(subpages)
		if l == 0 {
			log.Println("No subpages found in source body")
			return []Website{}, nil
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

	return []Website{}, nil

}
