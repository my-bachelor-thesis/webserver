package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/email_sender"
	"webserver/internal/postgres/rdg/task_with_solutions_and_tests"
	"webserver/internal/postgres/rdg/tasks"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/users"
	"webserver/internal/postgres/transaction_scripts"
	"webserver/pkg/postgresutil"
)

func AllTasksGet(c echo.Context) error {
	search, date, name := GetSearchDateName(c)
	tsks, err := tasks.GetApprovedAndPublishedByFilter(search, date, name)

	return returnEmptySliceIfNoRows(c, tsks, err)
}

func AllUsersTasksGet(c echo.Context) error {
	user, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	search, date, name := GetSearchDateName(c)

	t, err := tasks.GetByAuthorIdAndFilter(user.UserId, search, date, name)
	return returnEmptySliceIfNoRows(c, t, err)
}

func PublishTaskPost(c echo.Context) error {
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	_, task, err := bindAndFindWithUserId(c, tasks.GetByIdAndAuthorId, claims.UserId)
	if err != nil {
		return err
	}

	if claims.IsAdmin {
		task.ApproverId = claims.UserId
		return task.ApproveAndPublish()
	}
	return task.Publish()
}

func UnpublishTaskPost(c echo.Context) error {
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	_, task, err := bindAndFindWithUserId(c, tasks.GetByIdAndAuthorId, claims.UserId)
	if err != nil {
		return err
	}
	return task.Unpublish()
}

func DeleteTaskPost(c echo.Context) error {
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	_, task, err := bindAndFindWithUserId(c, tasks.GetByIdAndAuthorId, claims.UserId)
	if err != nil {
		return err
	}

	return task.Delete()
}

func AllTasksUnapprovedGet(c echo.Context) error {
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}
	if !claims.IsAdmin {
		return c.JSON(http.StatusForbidden, "not admin, forbidden")
	}

	search, date, name := GetSearchDateName(c)

	t, err := tasks.GetUnapproved(search, date, name)
	return returnEmptySliceIfNoRows(c, t, err)
}

func ApproveTaskPost(c echo.Context) error {
	_, task, err := bindAndFind(c, tasks.GetById)
	if err != nil {
		return err
	}
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}
	task.ApproverId = claims.UserId
	return task.Approve()
}

func DenyTaskPost(c echo.Context) error {
	var request struct {
		Reason   string `json:"reason" validate:"required"`
		AuthorId int    `json:"author_id" validate:"required"`
		TaskId   int    `json:"task_id" validate:"required"`
	}

	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	task, err := tasks.GetById(request.TaskId)
	if err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	admin, err := users.GetById(claims.UserId)
	if err != nil {
		return err
	}

	user, err := users.GetById(request.AuthorId)
	if err != nil {
		return err
	}

	if err := task.Unapprove(); err != nil {
		return err
	}

	return email_sender.SendDenial(user.Email, task.Title, admin.Username, admin.Email, request.Reason)
}

func AddPostPost(c echo.Context) error {
	// TODO: move to transaction script
	var request task_with_solutions_and_tests.TaskWithSolutionsAndTests

	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	// todo: all in one transaction

	// insert task
	task := tasks.Task{
		AuthorId:   claims.UserId,
		Title:      request.Title,
		Difficulty: request.Difficulty,
		Text:       request.Description,
	}

	// if updating
	if request.TaskId != 0 {
		task.Id = request.TaskId
		if err = task.UpdateTitleDifficultyAndText(); err != nil {
			return err
		}
		if err := transaction_scripts.DeleteAllPublicOrFinal(task.Id); err != nil {
			return err
		}
	} else {
		if err = task.Insert(); err != nil {
			return err
		}
	}

	if len(request.PublicSolutions) > 0 {
		var publicSolutions []*user_solutions.UserSolution
		for _, solution := range request.PublicSolutions {
			var u user_solutions.UserSolution
			u.UserId = claims.UserId
			u.TaskId = task.Id
			u.Language = solution.Language
			u.Name = solution.Name
			u.Public = true
			u.Code = solution.Code
			publicSolutions = append(publicSolutions, &u)
		}
		if err = user_solutions.InsertMany(publicSolutions); err != nil {
			return err
		}
	}

	fillTest := func(newTest *tests.Test, testFromRequest *task_with_solutions_and_tests.NameAndCode) {
		newTest.Name = testFromRequest.Name
		newTest.Public = true
		newTest.UserId = claims.UserId
		newTest.TaskId = task.Id
		newTest.Language = testFromRequest.Language
		newTest.Code = testFromRequest.Code
	}

	if len(request.PublicTests) > 0 {
		var publicTests []*tests.Test
		for _, test := range request.PublicTests {
			var t tests.Test
			fillTest(&t, test)
			publicTests = append(publicTests, &t)
		}
		if err = tests.InsertMany(publicTests); err != nil {
			return err
		}
	}

	var finalTests []*tests.Test
	for _, test := range request.FinalTests {
		var t tests.Test
		fillTest(&t, test)
		t.Final = true
		t.Name = "Final"
		finalTests = append(finalTests, &t)
	}
	return tests.InsertMany(finalTests)
}

func UnpublishedSavedTaskGet(c echo.Context) error {
	solutionId, err := strconv.Atoi(c.Param("task-id"))
	if err != nil {
		return err
	}

	user, err := getUserFromJWTCookie(c)
	if err != nil {
		return err
	}

	task, err := task_with_solutions_and_tests.GetByTaskId(solutionId, user.Id)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusMethodNotAllowed, "Task not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, task)
}

func GetSearchDateName(c echo.Context) (search, date, name string) {
	search = c.QueryParam("search")
	date = c.QueryParam("date")
	name = c.QueryParam("name")
	return search, date, name
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
