package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"prasowka/internal/filters"
	"prasowka/internal/models"
	"prasowka/internal/pagination"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrId  = fmt.Errorf("Błąd id")
	ErrInt = fmt.Errorf("wymagana liczba całkowita")
	ErrUrl = fmt.Errorf("Bład w adresie url")

	ErrRequired       = fmt.Errorf("wymagane")
	ErrReqMin         = fmt.Errorf("min.")
	ErrReqMax         = fmt.Errorf("max.")
	ErrEmailFormat    = fmt.Errorf("nieporawny format email")
	ErrEmailDuplacate = fmt.Errorf("adres email już zarejstrowany")

	ErrYear = fmt.Errorf("nieporawny rok")

	ErrIncorect = fmt.Errorf("niepoprawne")

	ErrUpdated = fmt.Errorf("Błąd aktualizacji. Nieaktualne dane. Spróbuj raz jeszcze.")

	ErrUser = fmt.Errorf("Nie odnaleziono takiego użytkownika.")

	ErrNotFound         = fmt.Errorf("nie odnaleziono danej pozycji/produktu w bazie danych")
	ErrDeadlineExceeded = fmt.Errorf("przekroczono limit czasu zapytania")
)

func (a *Application) ReadInt(c *gin.Context, key string, constructive int, max int) (int, error) {
	i := constructive

	if s, ok := c.GetQuery(key); ok {
		s = strings.TrimSpace(s)
		if s != "" {
			parsed, err := strconv.Atoi(s)
			if err != nil {
				return i, fmt.Errorf("%s — %w", key, ErrInt)
			}
			i = parsed
		}
	}

	if i <= 0 {
		return i, fmt.Errorf("%s — %w > 0", key, ErrInt)
	}

	if max > 0 && i > max {
		return i, fmt.Errorf("%s limited to %d", key, max)
	}

	return i, nil

}

func (app *Application) HandleAllDaily(c *gin.Context) {

	ValidationErrors := make(map[string]error)

	var pag pagination.Pagination
	pages := make([]int, 0)

	page, err := app.ReadInt(c, "page", 1, 0)
	if err != nil {
<<<<<<< HEAD
		ValidationErrors["Page"] = err
	}

	pageSize, err := app.ReadInt(c, "page_size", 17, 120)
	if err != nil {
		ValidationErrors["Page size"] = err
	}

	f := filters.Article{}
	f.Page = page
	f.PageSize = pageSize

	articles, total, err := models.SelectArticlesWhere(app.DB, f)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			ValidationErrors["NoteErr"] = fmt.Errorf("Nie udało się pobrać notatek, przekroczono limit czasu.")
		}
	}

	if total > 0 {
		pag = pagination.NewPagination(page, pageSize, total)
		for i := range pag.TotalPages {
			pages = append(pages, i+1)
		}
=======
		app.Logger.Error("Error: no database found.")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "Error Page", // for haeder title
			"Error": err,
		})
		return
>>>>>>> d83eca5fd0969efa2f39c913956192f7cc917ca8
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":      "Main website", // for haeder title
		"List":       articles,
		"Pagination": pag,
		"Pages":      pages,
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
