package user_solutions_tests

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetByUserIdAndUserSolutionId(userId, userSolutionId int) (*UserSolutionTest, error) {
	statement := fmt.Sprintf("select %s from %s where user_id = $1 and user_solution_id = $2", allFields, tableName)
	ust := UserSolutionTest{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, userId, userSolutionId), &ust)
	return &ust, err
}

func load(qr pgx.Row, ust *UserSolutionTest) error {
	return qr.Scan(&ust.UserSolutionId, &ust.TestId, &ust.UserId)
}
