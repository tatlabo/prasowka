package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	Table     string       `db:"table"`
}

type SqlInit struct {
	Create string
	Config []string
	Delete string
}

type WebsiteRender struct {
	Id        int    `db:"id" json:"id"`
	URL       string `db:"url" json:"url"`
	Title     string `db:"title" json:"title"`
	Body      string `db:"body" json:"body"`
	CreatedAt string `db:"created_at" json:"created_at"`
	Keywords  string `db:"keywords" json:"keywords"`
	Display   int    `db:"display" json:"display"`
	Done      int    `db:"done" json:"done"`
}

func HandleAllDaily(c *gin.Context) {
	// dbService := database.New()
	// defer dbService.Close()

	db, err := sql.Open("sqlite", "./db/websites.db")
	if err != nil {
		ErrorPage(c, err)
		return
	}
	defer db.Close()

	list, err := SelectAllArticles(db)
	if err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Main website", // for haeder title
		"List":  list,
	})
	// c.JSON(http.StatusOK, resp)
}

func HandleByID(c *gin.Context) {
	// dbService := database.New()
	// defer dbService.Close()
	id := c.Param("id")

	w := Website{}
	w.Id, _ = strconv.Atoi(id)

	db, err := sql.Open("sqlite", "./db/websites.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := w.SelectById(db); err != nil {
		panic(err)
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"Title":   "Main website", // for haeder title
		"Article": w,
	})
	// c.JSON(http.StatusOK, resp)
}

func HandleAllDailyJSON(c *gin.Context) {

	db, err := sql.Open("sqlite", "./db/websites.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	list, err := SelectAllArticles(db)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, list)
}

func HandleProcessById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorPage(c, fmt.Errorf("invalid id: %w", err))
		return
	}

	db, err := sql.Open("sqlite", "./db/websites.db")
	if err != nil {
		ErrorPage(c, err)
		return
	}

	defer db.Close()

	w := Website{Id: id}
	if err := w.SelectById(db); err != nil {
		ErrorPage(c, err)
		return
	}

	if w.Done != 1 {
		if err := w.ProcessWebsite(); err != nil {
			ErrorPage(c, err)
			return
		}

		if err := w.UpdateRaw(db); err != nil {
			ErrorPage(c, err)
			return
		}
	}
	s, err := ScrapArticle(w)

	if err != nil {
		ErrorPage(c, err)
		return
	}

	var articleTitle, articleLead string
	var articleContent []string

	if len(s) > 0 {
		articleTitle = s[0]
	}
	if len(s) > 1 {
		articleLead = s[1]
	}
	if len(s) > 2 {
		articleContent = s[2:]
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"Title":          "Main website", // for haeder title
		"Article":        w,
		"ArticleTitle":   articleTitle,
		"ArticleLead":    articleLead,
		"ArticleContent": articleContent,
		"CreatedAt":      w.CreatedAt.Format("2006-01-02 15:04:05"),
	})

}

func ErrorPage(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "error.html", gin.H{
		"Title": "Error Page", // for haeder title
		"Error": err,
	})
}

func HandleError(c *gin.Context) {
	err := fmt.Errorf("internal server error: something went wrong")
	c.HTML(http.StatusOK, "error.html", gin.H{
		"Title": "Error Page", // for haeder title
		"Error": err,
	})
}
