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

type TestIncomingJson struct {
	Solution   string `json:"solution"`
	SolutionId int    `json:"solution_id"`
	Test       string `json:"test"`
	TestId     int    `json:"test_id"`
	TaskId     int    `json:"task_id"`
}

type TestOutGoingJson struct {
	Solution *user_solutions.UserSolution `json:"solution"`
	Test     *tests.Test                  `json:"test"`
}

func runSolutionsOnTesterApi(c echo.Context) (*user_solutions.UserSolution, *TestIncomingJson, error) {
	incoming := &TestIncomingJson{}
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
	out := &TestOutGoingJson{Solution: us}
	return c.JSON(http.StatusOK, out)
}
func fillUserSolution(c echo.Context, us *user_solutions.UserSolution, incoming *TestIncomingJson, userId int) {
	us.UserId = userId
	us.TaskId = incoming.TaskId
	us.TestId = incoming.TestId
	us.Code = incoming.Solution
	us.Language = c.Param("lang")
}

func fillTest(c echo.Context, test *tests.Test, incoming *TestIncomingJson, userId int) {
	test.UserId = userId
	test.Code = incoming.Test
	test.TaskId = incoming.TaskId
	test.Language = c.Param("lang")
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
	out := &TestOutGoingJson{Solution: us}
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
	us.Code = ""
	test.Code = ""
	out := &TestOutGoingJson{Test: test, Solution: us}
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
	out := &TestOutGoingJson{Test: test, Solution: us}
	return c.JSON(http.StatusOK, out)
}
