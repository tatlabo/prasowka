package application

import (
	"html/template"
	"net/http"

	"prasowka/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *Application) Router() http.Handler {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"formatDate": utils.FormatDate,
		"not":        utils.Not,
		"equals":     utils.Equals,
		"notequals":  utils.Notequals,
	})

	r.LoadHTMLGlob("public/views/*")

	r.Static("static", "./static")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8099"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/hello", app.HelloWorldHandler)

	// r.GET("/health", app.healthHandler)

	r.GET("/", func(c *gin.Context) {
		app.HandleAllDaily(c)
	})

	r.GET("/news/:id", app.HandleProcessById)

	r.GET("/error", app.HandleError)

	r.GET("/json/news", app.HandleAllDailyJSON)

	r.GET("/news/raw/:id", app.HandleProcessById)

	return r
}

func (app *Application) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// func (app *Application) healthHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, app.DB.Health())
// }
