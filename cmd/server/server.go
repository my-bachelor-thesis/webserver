package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"runtime"
	"webserver/internal/config"
	"webserver/internal/handlers"
)

func init() {
	if runtime.GOOS != "linux" {
		log.Fatal("can only run on Linux")
	}
}

func main() {
	var e = echo.New()

	if config.GetInstance().IsProduction {
		e.Use(middleware.Logger())
	} else {
		e.Debug = true
	}

	jwtConfig := middleware.JWTConfig{
		Claims:     &handlers.JwtCustomClaims{},
		SigningKey: []byte(config.GetInstance().JWTSecret),
		TokenLookup: "cookie:auth",
	}

	// disable all CORS
	e.Use(middleware.CORS())

	allTemplates, err := getAllTemplates()
	logAndExitIfErr(e, err)
	parsedTemplates := &Template{
		templates: allTemplates,
	}
	e.Renderer = parsedTemplates

	// requests from website
	e.GET("/", handlers.IndexGet)
	e.GET("/login", handlers.LoginGet)
	e.POST("/login", handlers.LoginPost)
	e.GET("/register", handlers.RegisterGet)
	e.POST("/register", handlers.RegisterPost)
	e.GET("/add_task", handlers.AddTaskGet)
	e.POST("/add_task", handlers.AddTaskPost)
	e.GET("/task", handlers.TaskGet)

	// requests from editor
	e.GET("/task/:id/init", handlers.TaskInitGet)
	e.GET("/task/:id/solutions-tests", handlers.TaskSolutionsTestsGet)

	// experimental requests
	r := e.Group("/restricted")
	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("/restricted", handlers.RestrictedGet)

	// static
	e.Static("/static", config.GetInstance().PublicDir)

	e.Logger.Fatal(e.Start(config.GetInstance().Port))
}

func logAndExitIfErr(e *echo.Echo, err error) {
	if err != nil {
		e.Logger.Fatal(err)
	}
}
