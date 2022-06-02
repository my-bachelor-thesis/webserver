package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/task_statistics"
)

func StatisticGet(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("task-id"))
	if err != nil {
		return err
	}
	ts, err := task_statistics.GetByTaskId(postgres.GetPool(), taskId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ts)
}
