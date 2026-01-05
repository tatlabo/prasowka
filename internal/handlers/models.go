package handlers

import (
	"html/template"
	"time"
)

type Website struct {
	Id        int          `db:"id"`
	SourceId  int          `db:"source_id"`
	URL       template.URL `db:"url" json:"url"`
	Title     string       `db:"title" json:"title"`
	Body      string       `db:"body" json:"body"`
	Blob      []byte       `db:"raw"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	Keywords  string       `db:"keywords" json:"keywords"`
	Display   int          `db:"display" json:"display"`
	Done      int          `db:"done" json:"done"`
	MD5       string       `db:"md5"`
}

type SqlInit struct {
	Create string
	Config []string
	Delete string
}
