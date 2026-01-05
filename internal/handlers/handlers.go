package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	a, err := ScrapArticle(w)

	if err != nil {
		ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "article", gin.H{"Article": a})
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
