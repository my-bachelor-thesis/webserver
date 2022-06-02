package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/postgres/transaction_scripts"
)

func UpdateTestNamePost(c echo.Context) error {
	err := transaction_scripts.UpdateTestName(c)
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	return err
}
