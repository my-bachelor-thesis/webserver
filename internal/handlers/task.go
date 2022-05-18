package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/postgres/rdg/tasks"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/pkg/postgresutil"
)

func AllTasksGet(c echo.Context) error {
	t, err := tasks.GetApprovedAndPublished()
	if err != nil {
		if postgresutil.IsNoRowsInResultErr(err) {
			return c.JSON(http.StatusOK, emptySliceResponse)
		}
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func AllTasksUnpublishedGet(c echo.Context) error {
	t, err := tasks.GetUnpublished()
	if err != nil {
		if postgresutil.IsNoRowsInResultErr(err) {
			return c.JSON(http.StatusOK, emptySliceResponse)
		}
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func PublishTaskPost(c echo.Context) error {
	_, task, err := bindAndFind(c, tasks.GetById)
	if err != nil {
		return err
	}
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}
	if claims.IsAdmin {
		task.ApproverId = claims.UserId
		return task.ApproveAndPublish()
	}
	return task.Publish()
}

func AllTasksUnapprovedGet(c echo.Context) error {
	t, err := tasks.GetUnapproved()
	if err != nil {
		if postgresutil.IsNoRowsInResultErr(err) {
			return c.JSON(http.StatusOK, emptySliceResponse)
		}
		return err
	}
	return c.JSON(http.StatusOK, t)
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

func AddPostPost(c echo.Context) error {
	type nameAndCode struct {
		Name     string `json:"name"`
		Code     string `json:"code" validate:"required"`
		Language string `json:"language" validate:"required"`
	}
	var request struct {
		Title           string         `json:"title" validate:"required"`
		Difficulty      string         `json:"difficulty" validate:"required"`
		Description     string         `json:"description" validate:"required"`
		FinalTests      []*nameAndCode `json:"final_tests" validate:"required"`
		PublicTests     []*nameAndCode `json:"pubic_tests"`
		PublicSolutions []*nameAndCode `json:"public_solutions"`
	}

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
	if err = task.Insert(); err != nil {
		return err
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

	fillTest := func(newTest *tests.Test, testFromRequest *nameAndCode) {
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
