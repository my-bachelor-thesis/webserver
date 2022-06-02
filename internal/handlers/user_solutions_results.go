package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/user_solutions_results"
	"webserver/pkg/postgresutil"
)

func UserSolutionsResultsGet(c echo.Context) error {
	userSolutionId, err := strconv.Atoi(c.Param("user-solution-id"))
	if err != nil {
		return err
	}

	testId, err := strconv.Atoi(c.Param("test-id"))
	if err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	usr, err := user_solutions_results.GetByUserIdUserSolutionIdAndTestId(postgres.GetPool(), claims.UserId, userSolutionId, testId)
	if err != nil {
		if postgresutil.IsNoRowsInResultErr(err) {
			return c.JSON(http.StatusOK, "not found")
		}
		return err
	}
	return c.JSON(http.StatusOK, usr)
}
