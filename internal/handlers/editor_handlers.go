package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/postgres/rdg/initial_data_for_editor"
	"webserver/internal/postgres/rdg/user_solutions_with_tests"
)

func TaskInitGet(c echo.Context) error {
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

func TaskSolutionsTestsGet(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	lang := c.QueryParam("lang")
	data, err := user_solutions_with_tests.GetByLanguage(lang, taskId, 1)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}
