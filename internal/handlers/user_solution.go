package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/postgres/rdg/user_solutions"
)

func UpdateUserSolutionNamePost(c echo.Context) error {
	req, us, err := bindAndFind(c, user_solutions.GetById)
	if err != nil {
		return err
	}
	if us.Public {
		return c.JSON(http.StatusForbidden, "can't update this user solution")
	}
	return us.UpdateName(req.Name)
}
