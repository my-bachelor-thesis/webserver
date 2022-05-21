package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/config"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/user_solutions_results"
)

type RequestForTesting struct {
	Solution   string `json:"solution" validate:"required"`
	SolutionId int    `json:"solution_id" validate:"required"`
	Test       string `json:"test" validate:"required"`
	TestId     int    `json:"test_id" validate:"required"`
	TaskId     int    `json:"task_id" validate:"required"`
	HashId     string `json:"hash_id" validate:"required"`
}

type InsertedSolution struct {
	Id           int    `json:"id"`            // id of a solution that was inserted into the db
	LastModified string `json:"last_modified"` // last_modified of a solution that was inserted into the db
}

type InsertedTest struct {
	Id           int    `json:"id"`            // id of a test that was inserted into the db
	LastModified string `json:"last_modified"` // last_modified of a test that was inserted into the db
}

type ResultFromTesting struct {
	Result           *user_solutions_results.UserSolutionResult `json:"result"`
	InsertedTest     *InsertedTest                              `json:"inserted_test"`
	InsertedSolution *InsertedSolution                            `json:"inserted_solution"`
}

func bindRequestRunAndSaveResult(c echo.Context) (*user_solutions_results.UserSolutionResult, *RequestForTesting, error) {
	req := &RequestForTesting{}
	if err := bindAndValidate(c, req); err != nil {
		return nil, nil, err
	}

	postData, err := json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}

	lang := c.Param("lang")
	resp, err := http.Post(fmt.Sprintf("%s/%s", config.GetInstance().TesterApiURL, lang),
		"application/json", bytes.NewBuffer(postData))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(fmt.Sprintf("got status code %d from the Testrer API", resp.StatusCode))
	}

	usr := &user_solutions_results.UserSolutionResult{}
	if err = json.NewDecoder(resp.Body).Decode(usr); err != nil {
		return nil, nil, err
	}

	return usr, req, nil
}

func OnlyTestPost(c echo.Context) error {
	usr, req, err := bindRequestRunAndSaveResult(c)
	if err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	if err := insertUserSolutionResult(usr, claims.UserId, req.TestId, req.SolutionId); err != nil {
		return err
	}

	out := &ResultFromTesting{Result: usr}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveSolutionPost(c echo.Context) error {
	usr, req, err := bindRequestRunAndSaveResult(c)
	if err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	insertedSolution, err := insertSolution(c, req, claims.UserId)
	if err != nil {
		return err
	}

	if err := insertUserSolutionResult(usr, claims.UserId, req.TestId, insertedSolution.Id); err != nil {
		return err
	}

	out := &ResultFromTesting{InsertedSolution: insertedSolution, Result: usr}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveTestPost(c echo.Context) error {
	usr, req, err := bindRequestRunAndSaveResult(c)
	if err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	insertedTest, err := insertTest(c, req, claims.UserId)
	if err != nil {
		return err
	}

	if err := insertUserSolutionResult(usr, claims.UserId, insertedTest.Id, req.SolutionId); err != nil {
		return err
	}

	out := &ResultFromTesting{InsertedTest: insertedTest, Result: usr}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveBothPost(c echo.Context) error {
	usr, req, err := bindRequestRunAndSaveResult(c)
	if err != nil {
		return err
	}

	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return err
	}

	insertedSolution, err := insertSolution(c, req, claims.UserId)
	if err != nil {
		return err
	}

	insertedTest, err := insertTest(c, req, claims.UserId)
	if err != nil {
		return err
	}

	if err := insertUserSolutionResult(usr, claims.UserId, insertedTest.Id, insertedSolution.Id); err != nil {
		return err
	}

	out := &ResultFromTesting{InsertedTest: insertedTest, InsertedSolution: insertedSolution, Result: usr}
	return c.JSON(http.StatusOK, out)
}

func insertTest(c echo.Context, req *RequestForTesting, userId int) (*InsertedTest, error) {
	test := &tests.Test{}
	fillTest(c, test, req, userId)

	err := test.Insert()
	return &InsertedTest{Id: test.Id, LastModified: test.LastModified}, err
}

func insertSolution(c echo.Context, req *RequestForTesting, userId int) (*InsertedSolution, error) {
	us := user_solutions.UserSolution{}
	fillUserSolution(c, &us, req, userId)

	err := us.Insert()
	return &InsertedSolution{Id: us.Id, LastModified: us.LastModified}, err
}

func insertUserSolutionResult(usr *user_solutions_results.UserSolutionResult, userId, testId, solutionId int) error {
	usr.UserId = userId
	usr.TestId = testId
	usr.UserSolutionId = solutionId
	return usr.Insert()
}

func fillUserSolution(c echo.Context, us *user_solutions.UserSolution, req *RequestForTesting, userId int) {
	us.UserId = userId
	us.TaskId = req.TaskId
	us.Code = req.Solution
	us.Language = c.Param("lang")
}

func fillTest(c echo.Context, test *tests.Test, req *RequestForTesting, userId int) {
	test.UserId = userId
	test.Code = req.Test
	test.TaskId = req.TaskId
	test.Language = c.Param("lang")
}
