package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"runtime"
	"webserver/internal/config"
	"webserver/internal/handlers"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/redis"
)

var frontendEndpoints = [...]string{"/", "add-task", "task", "login", "register", "logout", "about",
	"account-settings", "my-tasks", "approve", "statistic", "promote-user", "reset-password-request",
	"email-verification", "password-reset", "edit-task"}

// custom form validator

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

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
	logAndExitIfErr(e, redis.Connect())
	defer redis.CloseConnection()

	if config.GetInstance().IsProduction {
		e.Use(middleware.Recover())
		e.Use(middleware.Logger())
	}

	jwtConfig := middleware.JWTConfig{
		Claims:      &jwt.CustomClaims{},
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

	// enable custom form validator
	e.Validator = &CustomValidator{validator: validator.New()}

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

	// requests from editor
	e.GET("/editor/init-data/:id", handlers.InitDataForEditorGet)
	e.GET("/editor/solutions-tests/:id/:lang", handlers.SolutionsAndTestsGet)
	e.POST("/editor/test/:lang", handlers.OnlyTestPost)
	e.POST("/editor/test-and-save-solution/:lang", handlers.TestAndSaveSolutionPost)
	e.POST("/editor/test-and-save-test/:lang", handlers.TestAndSaveTestPost)
	e.POST("/editor/test-and-save-both/:lang", handlers.TestAndSaveBothPost)
	e.GET("/editor/code-of-test/:id", handlers.CodeOfTestGet)
	e.GET("/editor/code-of-solution/:id", handlers.CodeOfSolutionGet)
	e.POST("/editor/change-name-in-test", handlers.UpdateTestNamePost)
	e.POST("/editor/change-name-in-usersolution", handlers.UpdateUserSolutionNamePost)
	e.POST("/editor/change-testid-for-usersolution", handlers.UpdateTestIdForUserSolutionPost)
	e.POST("/editor/change-last-opened", handlers.UpdateLastOpenedPost)
	e.GET("/editor/get-last-opened/:task-id", handlers.LastOpenedGet)
	e.GET("/editor/get-solution-result/:user-solution-id/:test-id", handlers.UserSolutionsResultsGet)

	// from register
	e.GET("/register/is-valid-username/:username", handlers.IsValidUsername)
	e.GET("/register/is-valid-email/:email", handlers.IsValidEmail)
	e.POST("/register/form", handlers.RegisterPost)

	// from login
	e.POST("/login/form", handlers.LoginPost)
	e.POST("/login/reset-password", handlers.RequestResetPasswordPost)

	// from add-task
	e.POST("/add-task/form", handlers.AddTaskPost)

	// my tasks
	e.GET("/my-tasks/all", handlers.AllUsersTasksGet)
	e.POST("/my-tasks/publish", handlers.PublishTaskPost)
	e.POST("/my-tasks/unpublish", handlers.UnpublishTaskPost)
	e.POST("/my-tasks/delete", handlers.DeleteTaskPost)

	// approve
	e.GET("/not-approved/all", handlers.AllTasksUnapprovedGet)
	e.POST("/not-approved/approve", handlers.ApproveTaskPost)
	e.POST("/not-approved/deny", handlers.DenyTaskPost)

	// account setting
	e.POST("/account-setting/update-user-info", handlers.UpdateUserInfoPost)
	e.POST("/account-setting/update-password", handlers.UpdatePasswordPost)
	e.POST("/account-setting/update-email", handlers.UpdateEmailPost)

	// edit solution
	e.GET("/edit-task/get-saved/:task-id", handlers.UnpublishedSavedTaskGet)

	// task statistic
	e.GET("/statistic/:task-id", handlers.StatisticGet)

	// other
	e.POST("/do-password-reset", handlers.ResetPasswordPost)
	e.POST("/email-verification", handlers.EmailVerificationPost)
	e.POST("/promote-user-form", handlers.PromoteToAdminPost)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.GetInstance().Port)))
}

func logAndExitIfErr(e *echo.Echo, err error) {
	if err != nil {
		e.Logger.Fatal(err)
	}
}
