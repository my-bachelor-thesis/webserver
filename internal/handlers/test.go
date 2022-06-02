package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres/transaction_scripts"
)

func UpdateTestNamePost(c echo.Context) error {
	return transaction_scripts.UpdateTestName(c)
}
