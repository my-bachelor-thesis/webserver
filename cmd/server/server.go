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
		Claims:      &handlers.JwtCustomClaims{},
		SigningKey:  []byte(config.GetInstance().JWTSecret),
		TokenLookup: "cookie:auth",
		Skipper: func(c echo.Context) bool {
			authCookie, err := c.Cookie("auth")
			if err != nil {
				return true
			}
			if authCookie.Value == "" {
				return true
			}
			return false
		},
	}

	// disable all CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.JWTWithConfig(jwtConfig))

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
	e.GET("/logout", handlers.Logout)

	// requests from editor
	e.GET("/init-data/:id", handlers.InitDataForEditorGet)
	e.GET("/solutions-tests/:id/:lang", handlers.SolutionsAndTestsGet)
	e.POST("/test/:lang", handlers.OnlyTestPost)
	e.POST("/test-and-save-solution/:lang", handlers.TestAndSaveSolutionPost)
	e.POST("/test-and-save-test/:lang", handlers.TestAndSaveTestPost)
	e.POST("/test-and-save-both/:lang", handlers.TestAndSaveBothPost)
	e.GET("/code-of-test/:id", handlers.CodeOfTestGet)
	e.GET("/code-of-solution/:id", handlers.CodeOfSolutionGet)

	// static
	e.Static("/static", config.GetInstance().PublicDir)

	e.Logger.Fatal(e.Start(config.GetInstance().Port))
}

func logAndExitIfErr(e *echo.Echo, err error) {
	if err != nil {
		e.Logger.Fatal(err)
	}
}
