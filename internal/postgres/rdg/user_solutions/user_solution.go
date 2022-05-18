package user_solutions

import (
	"fmt"
	"strings"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "user_id, task_id, last_modified, language, name, public, code"
	allFields          = "id, " + allFieldsWithoutId
	tableName          = "user_solutions"
)

var (
	allFieldReplacedTimestamp             = postgresutil.CallToCharOnTimestamp(allFields, "last_modified")
	placeHoldersWithTimestampAndWithoutId = postgresutil.GeneratePlaceholdersAndReplace(allFieldsWithoutId, map[int]string{2: "CURRENT_TIMESTAMP"})
)

type UserSolution struct {
	Id           int    `json:"id"`
	UserId       int    `json:"user_id"`
	TaskId       int    `json:"task_id"`
	LastModified string `json:"last_modified"`
	Language     string `json:"language"`
	Name         string `json:"name"`
	Public       bool   `json:"public"`
	Code         string `json:"code"`
}

func (us *UserSolution) Insert() error {
	statement := fmt.Sprintf(`
	insert into %s (%s)
	values (%s)
	returning id, to_char(last_modified, 'DD.MM.YY, HH24:MI:SS')`, tableName, allFieldsWithoutId, placeHoldersWithTimestampAndWithoutId)
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, getInsertFields(us)...).Scan(&us.Id, &us.LastModified)
}

func (us *UserSolution) UpdateName() error {
	statement := fmt.Sprintf("update %s set name = $1 where id = $2", tableName)
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, us.Name, us.Id)
	return err
}

func InsertMany(us []*UserSolution) error {
	statement := fmt.Sprintf(`
	insert into %s (%s)
	values `, tableName, allFieldsWithoutId)
	var vals []interface{}
	var palaceholders string
	numberOfFields := strings.Count(allFieldsWithoutId, ",")
	placeholderIndex := 1
	for _, row := range us {
		palaceholders = postgresutil.GeneratePlaceholdersAndReplaceFromIndex(allFieldsWithoutId, map[int]string{2: "CURRENT_TIMESTAMP"}, placeholderIndex)
		placeholderIndex += numberOfFields
		statement += fmt.Sprintf("(%s),", palaceholders)
		vals = append(vals, getInsertFields(row)...)
	}
	statement = statement[0 : len(statement)-1]
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, vals...)
	return err
}

func getInsertFields(us *UserSolution) (res []interface{}) {
	res = append(res, us.UserId, us.TaskId,
		us.Language, us.Name, us.Public, us.Code)
	return
}
