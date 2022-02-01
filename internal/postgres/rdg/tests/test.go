package tests

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "last_modified, final, user_id, task_id, language, code"
	allFields          = "id, " + allFieldsWithoutId
)

var allFieldReplacedTimestamp = postgresutil.CallToCharOnTimestamp(allFields, "last_modified")

type Test struct {
	Id           int    `json:"id"`
	LastModified string `json:"last_modified"`
	Final        bool   `json:"final"`
	UserId       int    `json:"user_id"`
	TaskId       int    `json:"task_id"`
	Language     string `json:"language"`
	Code         string `json:"code"`
}

func (test *Test) Insert() error {
	placeholders := postgresutil.GeneratePlaceholdersAndReplace(allFieldsWithoutId, map[int]string{0: "CURRENT_TIMESTAMP"})
	statement := fmt.Sprintf(`
	insert into tests (%s)
	values (%s)
	returning id, to_char(last_modified, 'DD.MM.YY, HH24:MI:SS')`, allFieldsWithoutId, placeholders)
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, test.Final, test.UserId,
		test.TaskId, test.Language, test.Code).Scan(&test.Id, &test.LastModified)
}
