package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"prasowka/internal/handlers"
	"strings"
	"time"
)

func main() {

	var path string
	w := handlers.Website{}

	switch len(os.Args) {
	case 1:
		fmt.Println("Please provide a path")
		os.Exit(1)
	case 2:
		path = os.Args[1]
		path = strings.TrimSpace(path)
		path = strings.TrimRight(path, "/")

		w.URL = template.URL(path)

		if err := w.ProcessWebsite(); err != nil {
			fmt.Println(fmt.Errorf("Error processing website: %w", err))
			os.Exit(1)
		}

	}

	ScrapPage(w)

}

func ScrapPage(w handlers.Website) {

	db, err := sql.Open("sqlite", "../db/websites.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	articles, err := handlers.RefreshSource(&w, db)
	if err != nil {
		log.Fatal(err)
	}

	var list []handlers.Website

	for i := range articles {
		n := handlers.Website{}
		n.SourceId = w.Id
		n.URL = articles[i].URL
		n.Title = articles[i].Title
		n.CreatedAt = time.Now()

		list = append(list, n)
		fmt.Printf("Found article: %s - %s\n", n.Title, n.URL)
	}

	err = handlers.AddWebsiteList(list, db)
	if err != nil {
		log.Fatal(err)
	}

}
