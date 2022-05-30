package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/postgres/rdg/tests"
)

func UpdateTestNamePost(c echo.Context) error {
	req, test, err := bindAndFind(c, tests.GetById)
	if err != nil {
		return err
	}

	if test.Final || test.Public {
		return c.JSON(http.StatusForbidden, "can't update this test")
	}

	test.Name = req.Name
	return test.UpdateName()
}
