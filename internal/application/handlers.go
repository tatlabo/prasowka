package application

import (
	"fmt"
	"net/http"
	"prasowka/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *Application) HandleAllDaily(c *gin.Context) {

	list, err := models.SelectAllArticles(app.DB)
	if err != nil {
		app.Logger.Error("Error: no database found.")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "Error Page", // for haeder title
			"Error": err,
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Main website", // for haeder title
		"List":  list,
	})
	// c.JSON(http.StatusOK, resp)
}

func (app *Application) HandleByID(c *gin.Context) {

	id := c.Param("id")

	w := models.Website{}
	w.Id, _ = strconv.Atoi(id)

	if err := w.SelectById(app.DB); err != nil {
		app.Logger.Error("Error: no database for single article")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "Error Page", // for haeder title
			"Error": err,
		})
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"Title":   "Main website", // for haeder title
		"Article": w,
	})
	// c.JSON(http.StatusOK, resp)
}

func (app *Application) HandleAllDailyJSON(c *gin.Context) {

	list, err := models.SelectAllArticles(app.DB)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, list)
}

func (app *Application) HandleProcessById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.ErrorPage(c, fmt.Errorf("invalid id: %w", err))
		return
	}

	w := models.Website{Id: id}
	if err := w.SelectById(app.DB); err != nil {
		app.ErrorPage(c, err)
		return
	}

	if w.Done != 1 {
		if err := w.ProcessWebsite(); err != nil {
			app.ErrorPage(c, err)
			return
		}

		if err := w.UpdateRaw(app.DB); err != nil {
			app.ErrorPage(c, err)
			return
		}
	}
	a, err := models.ScrapArticle(w)

	if err != nil {
		app.ErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "article", gin.H{"Article": a})
}

func (app *Application) ErrorPage(c *gin.Context, err error) {
	c.HTML(http.StatusOK, "error.html", gin.H{
		"Title": "Error Page", // for haeder title
		"Error": err,
	})
}

func (app *Application) HandleError(c *gin.Context) {
	err := fmt.Errorf("internal server error: something went wrong")
	c.HTML(http.StatusOK, "error.html", gin.H{
		"Title": "Error Page", // for haeder title
		"Error": err,
	})
}
