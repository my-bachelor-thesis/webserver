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
		Name string `json:"name"`
		Code string `json:"code"`
	}
	type request struct {
		Title           string         `json:"title"`
		Difficulty      string         `json:"difficulty"`
		Description     string         `json:"description"`
		GoFinalTest     string         `json:"go_final_test"`
		GoTests         []*nameAndCode `json:"go_tests"`
		GoSolutions     []*nameAndCode `json:"go_solutions"`
		PythonFinalTest string         `json:"python_final_test"`
		PythonTests     []*nameAndCode `json:"python_tests"`
		PythonSolutions []*nameAndCode `json:"python_solutions"`
	}

	// todo: all in one transaction

	req := request{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	// insert task
	task := tasks.Task{
		AuthorId:   claims.UserId,
		Title:      req.Title,
		Difficulty: req.Difficulty,
		Text:       req.Description,
	}
	if err = task.Insert(); err != nil {
		return err
	}

	makeSolutionsSlice := func(solutions []*nameAndCode, language string) (res []*user_solutions.UserSolution) {
		for _, solution := range solutions {
			var u user_solutions.UserSolution
			u.UserId = claims.UserId
			u.TaskId = task.Id
			u.Language = language
			u.Name = solution.Name
			u.Public = true
			u.Code = solution.Code
			res = append(res, &u)
		}
		return
	}

	makeTestsSlice := func(ts []*nameAndCode, language string) (res []*tests.Test) {
		for _, test := range ts {
			var t tests.Test
			t.Name = test.Name
			t.Code = test.Code
			t.Language = "go"
			t.UserId = claims.UserId
			t.Public = true
			t.TaskId = task.Id
			res = append(res, &t)
		}
		return
	}

	insertFinal := func(code, language string) error {
		var test tests.Test
		test.Name = "final"
		test.Code = code
		test.Language = language
		test.UserId = claims.UserId
		test.TaskId = task.Id
		if err = test.Insert(); err != nil {
			return err
		}
		return nil
	}

	if req.GoFinalTest != "" {
		if err = user_solutions.InsertMany(makeSolutionsSlice(req.GoSolutions, "go")); err != nil {
			return err
		}
		if err = tests.InsertMany(makeTestsSlice(req.GoTests, "go")); err != nil {
			return err
		}
		if err = insertFinal(req.GoFinalTest, "go"); err != nil {
			return err
		}
	}

	if req.PythonFinalTest != "" {
		if err = user_solutions.InsertMany(makeSolutionsSlice(req.PythonSolutions, "python")); err != nil {
			return err
		}
		if err = tests.InsertMany(makeTestsSlice(req.PythonTests, "python")); err != nil {
			return err
		}
		if err = insertFinal(req.PythonFinalTest, "python"); err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, nil)
}
