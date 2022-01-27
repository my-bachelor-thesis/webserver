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

func Insert(test *Test) error {
	statement := fmt.Sprintf(`
	insert into tests %s
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholder(allFieldsWithoutId))
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, test.LastModified, test.Final, test.UserId,
	test.TaskId, test.Language, test.Code).Scan(&test.Id)
}
