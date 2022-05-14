package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres/rdg/user_solutions_tests"
)

func UpdateTestIdForUserSolutionPost(c echo.Context) error {
	type request struct {
		UserSolutionId int `json:"user_solution_id"`
		TestId         int `json:"test_id"`
	}
	var req request
	if err := c.Bind(&req); err != nil {
		return err
	}
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}
	ust := user_solutions_tests.UserSolutionsTests{TestId: req.TestId, UserSolutionId: req.UserSolutionId, UserId: claims.UserId}
	return ust.Upsert()
}
