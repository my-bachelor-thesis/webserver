package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/user_solutions_tests"
)

func UpdateTestIdForUserSolutionPost(c echo.Context) error {
	var request struct {
		UserSolutionId int `json:"user_solution_id" validate:"required"`
		TestId         int `json:"test_id" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}
	ust := user_solutions_tests.UserSolutionTest{TestId: request.TestId, UserSolutionId: request.UserSolutionId,
		UserId: claims.UserId}
	return ust.Upsert(postgres.GetPool())
}
