package user_solutions_test_ids

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetByUserIdAndUserSolutionId(userId, userSolutionId int) (*UserSolutionsTestIds, error) {
	statement := fmt.Sprintf("select %s from user_solutions_test_ids where user_id = $1 and user_solution_id = $2", allFields)
	usti := UserSolutionsTestIds{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, userId, userSolutionId), &usti)
	return &usti, err
}

func load(qr pgx.Row, usti *UserSolutionsTestIds) error {
	return qr.Scan(&usti.UserSolutionId, &usti.TestId, &usti.UserId)
}
