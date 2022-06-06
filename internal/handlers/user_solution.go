package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/email_sender"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/transaction_scripts"
)

func UpdateUserSolutionNamePost(c echo.Context) error {
	err := transaction_scripts.UpdateUserSolutionName(c)
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	return err
}

func UserSolutionGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}
	if !claims.IsAdmin {
		return c.JSON(http.StatusMethodNotAllowed, "not a admin")
	}

	us, err := user_solutions.GetById(postgres.GetPool(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, us)
}

func DeleteUserSolutionPost(c echo.Context) error {
	var request struct {
		Id     int    `json:"id"`
		Reason string `json:"reason"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}
	if !claims.IsAdmin {
		return c.JSON(http.StatusMethodNotAllowed, "not a admin")
	}

	admin, user, us, err := transaction_scripts.DeleteUserSolution(request.Id, claims.UserId)
	if err != nil {
		return err
	}

	return email_sender.SendOnUserSolutionDeletion(user.Email, us.Code, admin.Username, admin.Email, request.Reason)
}
