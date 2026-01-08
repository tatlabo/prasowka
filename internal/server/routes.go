package server

import (
	"html/template"
	"net/http"

	"prasowka/internal/handlers"
	"prasowka/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
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
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/news", func(c *gin.Context) {
		handlers.HandleAllDaily(c)
	})

	r.GET("/news/:id", handlers.HandleProcessById)

	r.GET("/error", handlers.HandleError)

	r.GET("/json/news", handlers.HandleAllDailyJSON)

	r.GET("/news/raw/:id", handlers.HandleProcessById)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
