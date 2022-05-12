package user_solutions_test_ids

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const allFields = "user_solution_id, test_id, user_id"

type UserSolutionsTestIds struct {
	UserSolutionId int `json:"user_solution_id"`
	TestId         int `json:"test_id"`
	UserId         int `json:"user_id"`
}

func (usti *UserSolutionsTestIds) Insert() error {
	statement := fmt.Sprintf(`
	insert into user_solutions_test_ids (%s)
	values (%s)`, allFields, postgresutil.GeneratePlaceholders(allFields))
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, usti.UserSolutionId, usti.TestId, usti.UserId)
	return err
}

func (usti *UserSolutionsTestIds) UpdateTestId() error {
	statement := "update user_solutions_test_ids set test_id = $1 where user_id = $2 and user_solution_id = $3"
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, usti.TestId, usti.UserId, usti.UserSolutionId)
	return err
}

func (usti *UserSolutionsTestIds) Upsert() error {
	_, err := GetByUserIdAndUserSolutionId(usti.UserId, usti.UserSolutionId)
	if postgresutil.IsNoRowsInResultErr(err) {
		return usti.Insert()
	}
	return usti.UpdateTestId()
}
