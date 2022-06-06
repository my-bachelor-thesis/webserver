package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/email_sender"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/task_with_solutions_and_tests"
	"webserver/internal/postgres/rdg/tasks"
	"webserver/internal/postgres/transaction_scripts"
	"webserver/pkg/postgresutil"
)

type DenyRequest struct {
	Reason   string `json:"reason" validate:"required"`
	AuthorId int    `json:"author_id" validate:"required"`
	TaskId   int    `json:"task_id" validate:"required"`
}

func AllTasksGet(c echo.Context) error {
	by, err := getFilterParams(c)
	if err != nil {
		return err
	}
	tsks, err := tasks.GetApprovedAndPublishedByFilter(postgres.GetPool(), by)

	return returnEmptySliceIfNoRows(c, tsks, err)
}

func AllUsersTasksGet(c echo.Context) error {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	by, err := getFilterParams(c)
	if err != nil {
		return err
	}

	t, err := tasks.GetByAuthorIdAndFilter(postgres.GetPool(), claims.UserId, by)
	return returnEmptySliceIfNoRows(c, t, err)
}

func PublishTaskPost(c echo.Context) error {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	return transaction_scripts.PublishTask(c, claims)
}

func UnpublishTaskPost(c echo.Context) error {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	return transaction_scripts.UnpublishTask(c, claims)
}

func DeleteTaskPost(c echo.Context) error {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	return transaction_scripts.DeleteTask(c, claims)
}

func AllTasksUnapprovedGet(c echo.Context) error {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}
	if !claims.IsAdmin {
		return c.JSON(http.StatusForbidden, "not admin, forbidden")
	}

	by, err := getFilterParams(c)
	if err != nil {
		return err
	}

	t, err := tasks.GetUnapproved(postgres.GetPool(), by)
	return returnEmptySliceIfNoRows(c, t, err)
}

func ApproveTaskPost(c echo.Context) error {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	task, user, admin, err := transaction_scripts.ApproveTask(c, claims)
	if err != nil {
		return err
	}

	return email_sender.SendOnTaskApproval(user.Email, admin.Username, admin.Email, task.Title)
}

func DenyTaskPost(c echo.Context) error {
	var request DenyRequest

	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	task, user, admin, err := transaction_scripts.DenyTask(claims, request.TaskId, request.AuthorId)
	if err != nil {
		return err
	}

	return email_sender.SendDenial(user.Email, task.Title, admin.Username, admin.Email, request.Reason)
}

func AddTaskPost(c echo.Context) error {
	var request task_with_solutions_and_tests.TaskWithSolutionsAndTests

	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	return transaction_scripts.AddTask(claims, &request)
}

func UnpublishedSavedTaskGet(c echo.Context) error {
	solutionId, err := strconv.Atoi(c.Param("task-id"))
	if err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	task, err := transaction_scripts.GetTaskWithSolutionsAndTasksByTaskId(solutionId, claims.UserId)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusMethodNotAllowed, "Task not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

func getFilterParams(c echo.Context) (*tasks.FilterBy, error) {
	by := &tasks.FilterBy{}
	var err error
	by.Search = c.QueryParam("search")
	by.Date = c.QueryParam("date")
	by.Difficulty = c.QueryParam("difficulty")
	by.Page, err = strconv.Atoi(c.QueryParam("page"))
	by.NotPublished = c.QueryParam("not-published")
	if err == nil && by.Page < 1 {
		err = errors.New("wrong page")
	}
	return by, err
}

func returnEmptySliceIfNoRows(c echo.Context, result interface{}, err error) error {
	if err != nil {
		if postgresutil.IsNoRowsInResultErr(err) {
			return c.JSON(http.StatusOK, emptySliceResponse)
		}
		return err
	}
	return c.JSON(http.StatusOK, result)
}
