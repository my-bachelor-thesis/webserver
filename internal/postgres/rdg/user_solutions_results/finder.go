package user_solutions_results

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetByUserIdUserSolutionIdAndTestId(userId, userSolutionId, testId int) (*UserSolutionsResults, error) {
	statement := fmt.Sprintf("select %s from %s where user_id = $1 and user_solution_id = $2 and test_id = $3", allFields, tableName)
	usr := UserSolutionsResults{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, userId, userSolutionId, testId), &usr)
	return &usr, err
}

func load(qr pgx.Row, usr *UserSolutionsResults) error {
	return qr.Scan(&usr.UserSolutionId, &usr.TestId, &usr.UserId, &usr.ExitCode, &usr.Output,
		&usr.CompilationTime, &usr.RealTime, &usr.KernelTime, &usr.UserTime, &usr.MaxRamUsage, &usr.BinarySize)
}
