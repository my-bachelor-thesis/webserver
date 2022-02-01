package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/postgres/rdg/initial_data_for_editor"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/user_solutions_with_tests"
)

func InitDataForEditorGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	data, err := initial_data_for_editor.GetByTaskId(id)
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
	userId, err := getUserId(c)
	if err != nil {
		return err
	}
	lang := c.Param("lang")
	data, err := user_solutions_with_tests.GetByLanguage(lang, taskId, userId)
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
	test, err := tests.GetById(id)
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
	us, err := user_solutions.GetById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, us)
}