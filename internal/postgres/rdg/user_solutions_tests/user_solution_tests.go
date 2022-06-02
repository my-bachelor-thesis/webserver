package user_solutions_tests

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	tableName = "user_solutions_tests"
	allFields = "user_solution_id, test_id, user_id"
)

type UserSolutionTest struct {
	UserSolutionId int `json:"user_solution_id"`
	TestId         int `json:"test_id"`
	UserId         int `json:"user_id"`
}

func (ust *UserSolutionTest) Insert(tx postgres.PoolInterface) error {
	statement := fmt.Sprintf(`
	insert into %s (%s)
	values (%s)`, tableName, allFields, postgresutil.GeneratePlaceholders(allFields))
	_, err := tx.Exec(postgres.GetCtx(), statement, getInsertFields(ust)...)
	return err
}

func (ust *UserSolutionTest) UpdateTestId(tx postgres.PoolInterface) error {
	statement := fmt.Sprintf("update %s set test_id = $1 where user_id = $2 and user_solution_id = $3", tableName)
	_, err := tx.Exec(postgres.GetCtx(), statement, ust.TestId, ust.UserId, ust.UserSolutionId)
	return err
}

func (ust *UserSolutionTest) Upsert(tx postgres.PoolInterface) error {
	_, err := GetByUserIdAndUserSolutionId(tx, ust.UserId, ust.UserSolutionId)
	if postgresutil.IsNoRowsInResultErr(err) {
		return ust.Insert(tx)
	}
	return ust.UpdateTestId(tx)
}

func getInsertFields(ust *UserSolutionTest) (res []interface{}) {
	res = append(res, ust.UserSolutionId, ust.TestId, ust.UserId)
	return
}
