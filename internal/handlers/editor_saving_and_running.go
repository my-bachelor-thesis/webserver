package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/config"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/user_solutions_results"
	"webserver/internal/postgres/transaction_scripts"
)

type RequestForTesting struct {
	Solution   string `json:"solution" validate:"required"`
	SolutionId int    `json:"solution_id"`
	Test       string `json:"test" validate:"required"`
	TestId     int    `json:"test_id" validate:"required"`
	TaskId     int    `json:"task_id" validate:"required"`
	HashId     string `json:"hash_id" validate:"required"`
}

type ResultFromTesting struct {
	Result           *user_solutions_results.UserSolutionResult `json:"result"`
	InsertedTest     *transaction_scripts.InsertedTest          `json:"inserted_test"`
	InsertedSolution *transaction_scripts.InsertedSolution      `json:"inserted_solution"`
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

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	usr.UserId = claims.UserId
	usr.TestId = req.TestId
	usr.UserSolutionId = req.SolutionId

	if err := usr.Insert(postgres.GetPool()); err != nil {
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

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	us := user_solutions.UserSolution{}
	fillUserSolution(c, &us, req, claims.UserId)

	usr.UserId = claims.UserId
	usr.TestId = req.TestId

	inserted, err := transaction_scripts.EditorSaveSolution(&us, usr)
	if err != nil {
		return err
	}

	out := &ResultFromTesting{InsertedSolution: inserted, Result: usr}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveTestPost(c echo.Context) error {
	usr, req, err := bindRequestRunAndSaveResult(c)
	if err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	test := tests.Test{}
	fillTest(c, &test, req, claims.UserId)

	usr.UserId = claims.UserId
	usr.UserSolutionId = req.SolutionId

	inserted, err := transaction_scripts.EditorSaveTest(&test, usr)
	if err != nil {
		return err
	}

	out := &ResultFromTesting{InsertedTest: inserted, Result: usr}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveBothPost(c echo.Context) error {
	usr, req, err := bindRequestRunAndSaveResult(c)
	if err != nil {
		return err
	}

	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	test := tests.Test{}
	fillTest(c, &test, req, claims.UserId)

	us := user_solutions.UserSolution{}
	fillUserSolution(c, &us, req, claims.UserId)

	usr.UserId = claims.UserId

	insertedSolution, insertedTest, err := transaction_scripts.EditorSaveTestAndSolution(&test, &us, usr)
	if err != nil {
		return err
	}

	out := &ResultFromTesting{InsertedTest: insertedTest, InsertedSolution: insertedSolution, Result: usr}
	return c.JSON(http.StatusOK, out)
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
