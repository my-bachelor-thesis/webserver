package tests

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "user_id, task_id, test_id, last_modified, language, code, exit_code, output, compilation_time, real_time, kernel_time, user_time, max_ram_usage, binary_size"
	allFields          = "id, " + allFieldsWithoutId
)

var allFieldReplacedTimestamp = postgresutil.CallToCharOnTimestamp(allFields, "last_modified")

type UserSolution struct {
	Id              int     `json:"id"`
	UserId          int     `json:"user_id"`
	TaskId          int     `json:"task_id"`
	TestId          int     `json:"test_id"`
	LastModified    string  `json:"last_modified"`
	Language        string  `json:"language"`
	Code            string  `json:"code"`
	ExitCode        int     `json:"exit_code"`
	Output          string  `json:"output"`
	CompilationTime float32 `json:"compilation_time"`
	RealTime        float32 `json:"real_time"`
	KernelTime      float32 `json:"kernel_time"`
	UserTime        float32 `json:"user_time"`
	MaxRamUsage     float32 `json:"max_ram_usage"`
	BinarySize      float32 `json:"binary_size"`
}

func Insert(us *UserSolution) error {
	statement := fmt.Sprintf(`
	insert into user_solutions %s
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholder(allFieldsWithoutId))
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, us.UserId, us.TaskId, us.TestId, us.LastModified,
		us.Language, us.Code, us.ExitCode, us.Output, us.CompilationTime, us.RealTime, us.KernelTime, us.UserTime,
		us.MaxRamUsage, us.BinarySize).Scan(&us.Id)
}
