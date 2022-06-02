package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres/transaction_scripts"
)

func UpdateUserSolutionNamePost(c echo.Context) error {
	return transaction_scripts.UpdateUserSolutionNamePost(c)
}
