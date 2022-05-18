package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/postgres/rdg/last_opened"
	"webserver/pkg/postgresutil"
)

func UpdateLastOpenedPost(c echo.Context) error {
	var req last_opened.LastOpened
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}
	req.UserId = claims.UserId
	err = req.Upsert()
	return err
}

func LastOpenedGet(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("task-id"))
	if err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	lo, err := last_opened.GetByUserIdAndTaskId(claims.UserId, taskId)
	if err != nil {
		if postgresutil.IsNoRowsInResultErr(err) {
			return c.JSON(http.StatusOK, "not found")
		}
		return err
	}
	return c.JSON(http.StatusOK, lo)
}
