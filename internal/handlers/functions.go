package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "modernc.org/sqlite"
)

func ExistingURL(w *Website, db *sql.DB) (l []template.URL, err error) {

	sql := `SELECT id, title, url, created_at FROM daily;`

	timeStr := ""
	rows, err := db.Query(sql)
	if err != nil {
		return l, err
	}

	for rows.Next() {
		next := Website{}
		err := rows.Scan(&next.Id, &next.Title, &next.URL, &timeStr)
		if err != nil {
			return l, err
		}

		l = append(l, next.URL)
	}

	return l, nil

}

func SelectAllArticles(db *sql.DB) (l []Website, err error) {

	sql := `SELECT daily.id, CONCAT(source.url, daily.url) as url, 
	daily.title, daily.body, daily.created_at, daily.keywords, daily.display, daily.done 
	FROM daily JOIN source ON daily.source_id = source.id ORDER BY daily.created_at DESC;
	`

	rows, err := db.Query(sql)
	if err != nil {
		return []Website{}, err
	}

	timeStr := ""
	for rows.Next() {
		next := Website{}
		err := rows.Scan(&next.Id, &next.URL, &next.Title, &next.Body, &timeStr, &next.Keywords, &next.Display, &next.Done)
		if err != nil {
			return []Website{}, err
		}

		next.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", timeStr)
		l = append(l, next)
	}

	return l, nil

}

func PragmaConfig(db *sql.DB) error {

	config := [3]string{
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA foreign_keys = ON;`,
		`PRAGMA busy_timeout = 5000;`,
	}

	for _, pragma := range config {
		_, err := db.Exec(pragma)
		if err != nil {
			return err
		}
	}

	return nil

}

func ReadFromDbSource(w *Website, db *sql.DB) (subpages []Website, err error) {

	err = w.LastSourceWebsite(db)
	if err != nil {
		return []Website{}, err
	}

	if w.Body != "" {
		subpages, err := ParseSourceBody(w)
		if err != nil {
			log.Fatal(err)
		}

		return subpages, nil
	}

	return []Website{}, nil
}

func ParseSourceBody(w *Website) ([]Website, error) {

	doc, err := htmlquery.Parse(strings.NewReader(w.Body))
	if err != nil {
		return nil, err
	}
	list := htmlquery.Find(doc, "//div/div/h3/span")

	subpages := []Website{}

	for _, node := range list {
		subpage := Website{}

		a := htmlquery.FindOne(node, "//a")
		title := htmlquery.InnerText(a)
		title = strings.TrimSpace(title)
		link := htmlquery.SelectAttr(a, "href")

		subpage.Title = title
		subpage.URL = template.URL(link)
		subpage.CreatedAt = time.Now()

		subpage.SourceId = w.Id
		subpages = append(subpages, subpage)

	}

	return subpages, nil
}

func AddWebsite(ctx context.Context, db *sql.DB, w *Website) error {

	if err := w.ProcessWebsite(); err != nil {
		log.Fatal(err)
	}

	if err := w.AddWebsite(db); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized and data inserted successfully.")
	return nil
}

func AddWebsiteList(w []Website, db *sql.DB) error {

	batchSize := 100
	totalInserted := 0

	c := len(w)

	for i := 0; i < c; i += batchSize {
		end := min(i+batchSize, c)
		batch := w[i:end]

		// Begin transaction for this batch
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error beginning transaction: %w", err)
		}

		// Prepare statement for batch
		stmt, err := tx.Prepare(`INSERT OR IGNORE INTO daily (source_id, url, body, title, created_at, keywords, display) 
		VALUES (?, ?, ?, ?, ?, ?, ?);`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error preparing statement: %w", err)
		}

		// Insert batch
		for i := range batch {
			var (
				source_id  = batch[i].SourceId
				url        = batch[i].URL
				body       = batch[i].Body
				title      = batch[i].Title
				created_at = batch[i].CreatedAt.Format("2006-01-02 15:04:05")
				keywords   = batch[i].Keywords
				display    = batch[i].Display
			)
			_, err = stmt.Exec(source_id, url, body, title, created_at, keywords, display)
			if err != nil {
				stmt.Close()
				tx.Rollback()
				return fmt.Errorf("error inserting file: %w", err)
			}
		}

		stmt.Close()
		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing transaction: %w", err)
		}

		totalInserted += len(batch)
	}

	return nil
}

func (w *Website) SourceToDb(ctx context.Context, db *sql.DB) error {

	stmt := `INSERT INTO source (url, body, created_at, keywords, display) VALUES (?, ?, ?, ?, ?) RETURNING id;`

	err := db.QueryRow(stmt, w.URL, w.Body, w.CreatedAt.Format("2006-01-02 15:04:05"), w.Keywords, w.Display).Scan(&w.Id)
	if err != nil {
		return err
	}

	return nil
}

func (w *Website) ArticelToDb(ctx context.Context, db *sql.DB) error {

	stmt := `INSERT INTO source (url, body, created_at, keywords, display) VALUES (?, ?, ?, ?, ?) RETURNING id;`

	err := db.QueryRow(stmt, w.URL, w.Body, w.CreatedAt.Format("2006-01-02 15:04:05"), w.Keywords, w.Display).Scan(&w.Id)
	if err != nil {
		return err
	}

	return nil
}

// type Website struct {
// 	Id        int          `db:"id"`
// 	SourceId  int          `db:"source_id"`
// 	URL       template.URL `db:"url" json:"url"`
// 	Title     string       `db:"title" json:"title"`
// 	Body      string       `db:"body" json:"body"`
// 	Blob      []byte       `db:"raw"`
// 	CreatedAt time.Time    `db:"created_at" json:"created_at"`
// 	Keywords  string       `db:"keywords" json:"keywords"`
// 	Display   int          `db:"display" json:"display"`
// 	Done      int          `db:"done" json:"done"`
// 	MD5       string       `db:"md5"`
// }

type ArticleRender struct {
	Website
	Lead    string   `db:"lead" json:"lead"`
	Content []string `db:"content" json:"content"`
}

func ScrapArticle(w Website) (a ArticleRender, err error) {

	doc, err := htmlquery.Parse(strings.NewReader(w.Body))
	if err != nil {
		return a, err
	}

	a.Website = w

	title := htmlquery.FindOne(doc, "//h1[@class='article-title']")
	lead := htmlquery.FindOne(doc, "//p[@class='article-lead']")
	rest := htmlquery.Find(doc, "//div[@class='articleContent']/p")
	dateEllement := htmlquery.FindOne(doc, "//div[@class='article-date']/meta")
	dateString := htmlquery.SelectAttr(dateEllement, "content")

	a.Website.CreatedAt, err = time.Parse("2006-01-02T15:04:05", dateString)
	if err != nil {
		a.Website.CreatedAt = w.CreatedAt
	}

	a.Website.Title = htmlquery.InnerText(title)
	a.Lead = htmlquery.InnerText(lead)

	for i := range rest {
		text := htmlquery.InnerText(rest[i])
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		a.Content = append(a.Content, text)
	}

	// for _, node := range list {
	// 	subpage := Website{}

	// 	a := htmlquery.FindOne(node, "//a")
	// 	title := htmlquery.InnerText(a)
	// 	title = strings.TrimSpace(title)
	// 	link := htmlquery.SelectAttr(a, "href")

	// 	subpage.Title = title
	// 	subpage.URL = template.URL(link)
	// 	subpage.CreatedAt = time.Now()

	// 	subpage.SourceId = w.Id
	// 	subpages = append(subpages, subpage)

	// }

	return a, nil
}
