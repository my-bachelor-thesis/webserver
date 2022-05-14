package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"runtime"
	"webserver/internal/config"
	"webserver/internal/handlers"
	"webserver/internal/postgres"
)

var frontendEndpoints = [...]string{"/", "add-task", "task", "login", "register", "logout", "about", "account-settings", "not-published", "approve"}

func init() {
	if runtime.GOOS != "linux" {
		log.Fatal("can only run on Linux")
	}
}

func main() {
	var e = echo.New()

	logAndExitIfErr(e, config.LoadConfig())
	logAndExitIfErr(e, postgres.CreateDbPool())
	defer postgres.ClosePool()

	if config.GetInstance().IsProduction {
		e.Use(middleware.Recover())
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
			if err != nil || authCookie.Value == "" {
				return true
			}
			return false
		},
	}
	e.Use(middleware.JWTWithConfig(jwtConfig))

	// disable all CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	for _, endpoint := range frontendEndpoints {
		e.File(endpoint, "web/frontend/public/index.html")
	}
	e.Static("/", "web/frontend/public/")

	// requests from home
	e.GET("/home/all-tasks", handlers.AllTasksGet)

	// requests from editor TODO: all prefix editor
	e.GET("/init-data/:id", handlers.InitDataForEditorGet)
	e.GET("/solutions-tests/:id/:lang", handlers.SolutionsAndTestsGet)
	e.POST("/test/:lang", handlers.OnlyTestPost)
	e.POST("/test-and-save-solution/:lang", handlers.TestAndSaveSolutionPost)
	e.POST("/test-and-save-test/:lang", handlers.TestAndSaveTestPost)
	e.POST("/test-and-save-both/:lang", handlers.TestAndSaveBothPost)
	e.GET("/code-of-test/:id", handlers.CodeOfTestGet)
	e.GET("/code-of-solution/:id", handlers.CodeOfSolutionGet)
	e.POST("/editor/change-name-in-test", handlers.UpdateTestNamePost)
	e.POST("/editor/change-name-in-usersolution", handlers.UpdateUserSolutionNamePost)
	e.POST("/editor/change-testid-for-usersolution", handlers.UpdateTestIdForUserSolutionPost)
	e.POST("/editor/change-last-opened", handlers.UpdateLastOpenedPost)
	e.GET("/editor/get-last-opened/:task-id", handlers.LastOpenedGet)
	e.GET("/editor/get-solution-result/:user-solution-id/:test-id", handlers.UserSolutionsResultsGet)

	// from register
	e.GET("/register/is-valid-username/:username", handlers.IsValidUsername)
	e.POST("/register/form", handlers.RegisterPost)

	// from login
	e.POST("/login/form", handlers.LoginPost)

	// TODO: change to path/action
	// from add-task
	e.POST("/add-task/form", handlers.AddPostPost)

	// publish
	e.GET("/not-published/all", handlers.AllTasksUnpublishedGet)
	e.POST("/not-published/publish", handlers.PublishTaskPost)

	// approve
	e.GET("/not-approved/all", handlers.AllTasksUnapprovedGet)
	e.POST("/not-approved/approve", handlers.ApproveTaskPost)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.GetInstance().Port)))
}

func logAndExitIfErr(e *echo.Echo, err error) {
	if err != nil {
		e.Logger.Fatal(err)
	}
}
