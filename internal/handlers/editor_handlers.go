package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"webserver/internal/postgres/rdg/initial_data_for_editor"
	"webserver/internal/postgres/rdg/tests"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/user_solutions_with_tests"
)

func InitDataForEditorGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	data, err := initial_data_for_editor.GetByTaskId(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func SolutionsAndTestsGet(c echo.Context) error {
	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	lang := c.Param("lang")
	data, err := user_solutions_with_tests.GetByLanguage(lang, taskId, 1)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func TestSolutionPost(c echo.Context) error {
	incoming := &SolutionIncomingJson{}
	if err := c.Bind(incoming); err != nil {
		return err
	}

	postData, err := json.Marshal(incoming)
	if err != nil {
		return err
	}

	lang := c.Param("lang")
	resp, err := http.Post(fmt.Sprintf("http://localhost:4000/%s", lang), "application/json", bytes.NewBuffer(postData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var outgoing SolutionOutgoingJson
	if err = json.NewDecoder(resp.Body).Decode(&outgoing); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, outgoing)
}

func CodeOfTestGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	test, err := tests.GetById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, test)
}

func CodeOfSolutionGet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	us, err := user_solutions.GetById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, us)
}

type SolutionIncomingJson struct {
	Solution string `json:"solution"`
	Test     string `json:"test"`
}

type SolutionOutgoingJson struct {
	ExitCode        int     `json:"exit_code"`
	Out             string  `json:"out"`
	CompilationTime float32 `json:"compilation_time"` // in s
	RealTime        float32 `json:"real_time"`        // in s
	KernelTime      float32 `json:"kernel_time"`      // in s
	UserTime        float32 `json:"user_time"`        // in s
	MaxRamUsage     float32 `json:"max_ram_usage"`    // in MB
	BinarySize      float32 `json:"binary_size"`      // in MB
}
