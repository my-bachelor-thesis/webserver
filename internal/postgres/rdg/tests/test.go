package tests

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "last_modified, language, code"
	allFields          = "id, " + allFieldsWithoutId
)

var allFieldReplacedTimestamp = postgresutil.CallToCharOnTimestamp(allFields, "last_modified")

type Test struct {
	Id           int    `json:"id"`
	LastModified string `json:"last_modified"`
	Language     string `json:"language"`
	Code         string `json:"code"`
}

func Insert(test *Test) error {
	statement := fmt.Sprintf(`
	insert into tests %s
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholder(allFieldsWithoutId))
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, test.LastModified, test.Language, test.Code).Scan(&test.Id)
}
