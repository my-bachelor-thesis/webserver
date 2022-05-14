package user_solutions_results

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	tableName = "user_solutions_results"
	allFields = "user_solution_id, test_id, user_id, exit_code, output, compilation_time, real_time, kernel_time, user_time, max_ram_usage, binary_size"
)

type UserSolutionsResults struct {
	UserSolutionId  int     `json:"user_solution_id"`
	TestId          int     `json:"test_id"`
	UserId          int     `json:"user_id"`
	ExitCode        int     `json:"exit_code"`
	Output          string  `json:"output"`
	CompilationTime float32 `json:"compilation_time"`
	RealTime        float32 `json:"real_time"`
	KernelTime      float32 `json:"kernel_time"`
	UserTime        float32 `json:"user_time"`
	MaxRamUsage     float32 `json:"max_ram_usage"`
	BinarySize      float32 `json:"binary_size"`
}

func (usr *UserSolutionsResults) Insert() error {
	statement := fmt.Sprintf(`
	insert into %s (%s)
	values (%s)`, tableName, allFields, postgresutil.GeneratePlaceholders(allFields))
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, getInsertFields(usr)...)
	return err
}

func getInsertFields(usr *UserSolutionsResults) (res []interface{}) {
	res = append(res, usr.UserSolutionId, usr.TestId, usr.UserId, usr.ExitCode, usr.Output, usr.CompilationTime,
		usr.RealTime, usr.KernelTime, usr.UserTime, usr.MaxRamUsage, usr.BinarySize)
	return
}
