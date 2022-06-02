package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/transaction_scripts"
)

func InitDataForEditorGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	data, err := transaction_scripts.GetInitDataForEditorByTaskId(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func SolutionsAndTestsGet(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}
	lang := c.Param("lang")
	data, err := transaction_scripts.GetUserSolutionsWithTestsByLanguage(lang, taskId, claims.UserId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func CodeOfTestGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	test, err := tests.GetById(postgres.GetPool(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, test)
}

func CodeOfSolutionGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	us, err := user_solutions.GetById(postgres.GetPool(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, us)
}
