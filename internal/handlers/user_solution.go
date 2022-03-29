package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres/rdg/user_solutions"
)

func UpdateUserSolutionNamePost(c echo.Context) error {
	req, us, err := bindAndFind(c, user_solutions.GetById)
	if err != nil {
		return err
	}
	return us.UpdateName(req.Name)
}
