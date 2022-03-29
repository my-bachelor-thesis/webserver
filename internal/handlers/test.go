package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres/rdg/tests"
)

func UpdateTestNamePost(c echo.Context) error {
	req, test, err := bindAndFind(c, tests.GetById)
	if err != nil {
		return err
	}
	return test.UpdateName(req.Name)
}
