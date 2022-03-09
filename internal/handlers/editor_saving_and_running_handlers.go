package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
)

type RequestForTesting struct {
	Solution   string `json:"solution"`
	SolutionId int    `json:"solution_id"`
	Test       string `json:"test"`
	TestId     int    `json:"test_id"`
	TaskId     int    `json:"task_id"`
}

type ResultFromTesting struct {
	Solution         *user_solutions.UserSolution `json:"solution"`
	TestId           int                          `json:"test_id"`            // id of a test that was inserted into the db
	TestLastModified string                       `json:"test_last_modified"` // last_modified of a test that was inserted into the db
}

func runSolutionsOnTesterApi(c echo.Context) (*user_solutions.UserSolution, *RequestForTesting, error) {
	incoming := &RequestForTesting{}
	if err := c.Bind(incoming); err != nil {
		return nil, nil, err
	}

	postData, err := json.Marshal(incoming)
	if err != nil {
		return nil, nil, err
	}

	lang := c.Param("lang")
	resp, err := http.Post(fmt.Sprintf("http://localhost:4000/%s", lang), "application/json", bytes.NewBuffer(postData))
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	us := &user_solutions.UserSolution{}
	if err = json.NewDecoder(resp.Body).Decode(us); err != nil {
		return nil, nil, err
	}
	return us, incoming, nil
}

func OnlyTestPost(c echo.Context) error {
	us, _, err := runSolutionsOnTesterApi(c)
	if err != nil {
		return err
	}
	out := &ResultFromTesting{Solution: us}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveSolutionPost(c echo.Context) error {
	us, incoming, err := runSolutionsOnTesterApi(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	fillUserSolution(c, us, incoming, userId)

	if err = us.Insert(); err != nil {
		fmt.Println(err)
		return err
	}
	us.Code = ""
	out := &ResultFromTesting{Solution: us}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveTestPost(c echo.Context) error {
	us, incoming, err := runSolutionsOnTesterApi(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	test := &tests.Test{}
	fillTest(c, test, incoming, userId)

	if err = test.Insert(); err != nil {
		return err
	}
	// don't send code back
	us.Code = ""
	out := &ResultFromTesting{TestId: test.Id, TestLastModified: test.LastModified, Solution: us}
	return c.JSON(http.StatusOK, out)
}

func TestAndSaveBothPost(c echo.Context) error {
	us, incoming, err := runSolutionsOnTesterApi(c)
	if err != nil {
		return err
	}

	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	fillUserSolution(c, us, incoming, userId)

	if err = us.Insert(); err != nil {
		return err
	}

	test := &tests.Test{}
	fillTest(c, test, incoming, userId)

	if err = test.Insert(); err != nil {
		return err
	}

	test.Code = ""
	us.Code = ""
	out := &ResultFromTesting{TestId: test.Id, TestLastModified: test.LastModified, Solution: us}
	return c.JSON(http.StatusOK, out)
}

func fillUserSolution(c echo.Context, us *user_solutions.UserSolution, incoming *RequestForTesting, userId int) {
	us.UserId = userId
	us.TaskId = incoming.TaskId
	us.TestId = incoming.TestId
	us.Code = incoming.Solution
	us.Language = c.Param("lang")
}

func fillTest(c echo.Context, test *tests.Test, incoming *RequestForTesting, userId int) {
	test.UserId = userId
	test.Code = incoming.Test
	test.TaskId = incoming.TaskId
	test.Language = c.Param("lang")
}
