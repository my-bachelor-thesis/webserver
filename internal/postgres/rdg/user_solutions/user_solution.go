package user_solutions

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "user_id, task_id, last_modified, language, name, public, code, exit_code, output, compilation_time, real_time, kernel_time, user_time, max_ram_usage, binary_size"
	allFields          = "id, " + allFieldsWithoutId
)

var (
	allFieldReplacedTimestamp             = postgresutil.CallToCharOnTimestamp(allFields, "last_modified")
	placeHoldersWithTimestampAndWithoutId = postgresutil.GeneratePlaceholdersAndReplace(allFieldsWithoutId, map[int]string{2: "CURRENT_TIMESTAMP"})
)

type UserSolution struct {
	Id              int     `json:"id"`
	UserId          int     `json:"user_id"`
	TaskId          int     `json:"task_id"`
	LastModified    string  `json:"last_modified"`
	Language        string  `json:"language"`
	Name            string  `json:"name"`
	Public          bool    `json:"public"`
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

func (us *UserSolution) Insert() error {
	statement := fmt.Sprintf(`
	insert into user_solutions (%s)
	values (%s)
	returning id, to_char(last_modified, 'DD.MM.YY, HH24:MI:SS')`, allFieldsWithoutId, placeHoldersWithTimestampAndWithoutId)
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, getInsertFields(us)...).Scan(&us.Id, &us.LastModified)
}

func (us *UserSolution) UpdateName() error {
	statement := "update user_solutions set name = $1 where id = $2"
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, us.Name, us.Id)
	return err
}

func InsertMany(us []*UserSolution) error {
	statement := fmt.Sprintf(`
	insert into user_solutions (%s)
	values `, allFieldsWithoutId)
	var vals []interface{}
	for _, row := range us {
		statement += fmt.Sprintf("(%s),", placeHoldersWithTimestampAndWithoutId)
		vals = append(vals, getInsertFields(row)...)
	}
	statement = statement[0 : len(statement)-1]
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, vals...)
	return err
}

func getInsertFields(us *UserSolution) (res []interface{}) {
	res = append(res, us.UserId, us.TaskId,
		us.Language, us.Name, us.Public, us.Code, us.ExitCode, us.Output, us.CompilationTime, us.RealTime,
		us.KernelTime, us.UserTime, us.MaxRamUsage, us.BinarySize)
	return
}
