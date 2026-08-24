package models

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"prasowka/internal/filters"
	"slices"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	_ "modernc.org/sqlite"
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

type IndexRender struct {
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

type ArticleRender struct {
	Website
	Lead    string   `db:"lead" json:"lead"`
	Content []string `db:"content" json:"content"`
}

type SqlInit struct {
	Create string
	Config []string
	Delete string
}

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

	var d strings.Builder

	for scanner.Scan() {
		d.WriteString(scanner.Text())
		d.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning website: %w", err)
	}

	w.Body = d.String()
	return nil
}

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

	limit := 25
	offset := 0

	sql := fmt.Sprintf(`SELECT daily.id, CONCAT(source.url, daily.url) as url, 
	daily.title, daily.body, daily.created_at, daily.keywords, daily.display, daily.done 
	FROM daily JOIN source ON daily.source_id = source.id ORDER BY daily.created_at DESC LIMIT %d OFFSET %d;
	`, limit, offset)

	sqlCount := `SELECT COUNT(*) FROM daily;`

	var count int
	err = db.QueryRow(sqlCount).Scan(&count)
	if err != nil {
		return []Website{}, err
	}

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

func SelectArticlesWhere(db *sql.DB, f filters.Article) (l []Website, count int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT COUNT(*) OVER () AS count, daily.id, CONCAT(source.url, daily.url) as url, 
	daily.title, daily.body, daily.created_at, daily.keywords, daily.display, daily.done 
	FROM daily JOIN source ON daily.source_id = source.id ORDER BY daily.created_at DESC
	LIMIT :limit OFFSET :offset;`

	args := []any{
		sql.Named("offset", f.PageSize*(f.Page-1)),
		sql.Named("limit", f.PageSize),
	}

	rows, err := db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return []Website{}, 0, err
	}
	defer rows.Close()

	timeStr := ""
	for rows.Next() {
		next := Website{}
		err := rows.Scan(&count, &next.Id, &next.URL, &next.Title, &next.Body, &timeStr, &next.Keywords, &next.Display, &next.Done)
		if err != nil {
			return []Website{}, 0, err
		}

		next.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", timeStr)
		l = append(l, next)
	}

	return l, count, nil

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
	list := htmlquery.Find(doc, "//div/div/article/h3")

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

func ScrapArticle(w Website) (a ArticleRender, err error) {

	doc, err := htmlquery.Parse(strings.NewReader(w.Body))
	if err != nil {
		return a, err
	}

	a.Website = w

	title := htmlquery.FindOne(doc, "//h1[@class='article-title']")
	lead := htmlquery.FindOne(doc, "//p[@class='article-page__lead']")
	rest := htmlquery.Find(doc, "//div[@class='article-page__content article_speakable']/p")
	dateEllement := htmlquery.FindOne(doc, "//div[@class='article-date']/meta")

	var dateString string
	if dateEllement != nil {
		dateString = htmlquery.SelectAttr(dateEllement, "content")
	}

	a.Website.CreatedAt, err = time.Parse("2006-01-02T15:04:05", dateString)
	if err != nil {
		a.Website.CreatedAt = w.CreatedAt
	}

	if title != nil {
		a.Website.Title = htmlquery.InnerText(title)
	}
	if lead != nil {
		a.Lead = htmlquery.InnerText(lead)
	}

	for i := range rest {
		text := htmlquery.InnerText(rest[i])
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		a.Content = append(a.Content, text)
	}
	a.Body = ""

	return a, nil
}
