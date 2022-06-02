package tests

import (
	"fmt"
	"strings"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "last_modified, final, name, public, user_id, task_id, language, code"
	allFields          = "id, " + allFieldsWithoutId
)

var (
	allFieldReplacedTimestamp             = postgresutil.CallToCharOnTimestamp(allFields, "last_modified")
	placeHoldersWithTimestampAndWithoutId = postgresutil.GeneratePlaceholdersAndReplace(allFieldsWithoutId, map[int]string{0: "CURRENT_TIMESTAMP"})
)

type Test struct {
	Id           int    `json:"id"` // 0 is default
	LastModified string `json:"last_modified"`
	Final        bool   `json:"final"`
	Name         string `json:"name"`
	Public       bool   `json:"public"`
	UserId       int    `json:"user_id"`
	TaskId       int    `json:"task_id"`
	Language     string `json:"language"`
	Code         string `json:"code"`
}

func (test *Test) Insert(tx postgres.PoolInterface) error {
	statement := fmt.Sprintf(`
	insert into tests (%s)
	values (%s)
	returning id, to_char(last_modified, 'DD.MM.YY, HH24:MI:SS')`, allFieldsWithoutId, placeHoldersWithTimestampAndWithoutId)
	return tx.QueryRow(postgres.GetCtx(), statement, getInsertFields(test)...).Scan(&test.Id, &test.LastModified)
}

func (test *Test) UpdateName(tx postgres.PoolInterface) error {
	statement := "update tests set name = $1 where id = $2"
	_, err := tx.Exec(postgres.GetCtx(), statement, test.Name, test.Id)
	return err
}

func InsertMany(tx postgres.PoolInterface, tests []*Test) error {
	statement := fmt.Sprintf(`
	insert into tests (%s)
	values `, allFieldsWithoutId)
	var vals []interface{}
	var palaceholders string
	numberOfFields := strings.Count(allFieldsWithoutId, ",")
	placeholderIndex := 1
	for _, row := range tests {
		palaceholders = postgresutil.GeneratePlaceholdersAndReplaceFromIndex(allFieldsWithoutId, map[int]string{0: "CURRENT_TIMESTAMP"}, placeholderIndex)
		placeholderIndex += numberOfFields
		statement += fmt.Sprintf("(%s),", palaceholders)
		vals = append(vals, getInsertFields(row)...)
	}
	statement = statement[0 : len(statement)-1]
	_, err := tx.Exec(postgres.GetCtx(), statement, vals...)
	return err
}

func getInsertFields(test *Test) (res []interface{}) {
	res = append(res, test.Final, test.Name, test.Public,
		test.UserId, test.TaskId, test.Language, test.Code)
	return
}
