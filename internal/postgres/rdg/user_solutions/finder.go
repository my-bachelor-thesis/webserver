package user_solutions

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetById(id int) (*UserSolution, error) {
	statement := fmt.Sprintf("select %s from user_solutions where id = $1", allFieldReplacedTimestamp)
	us := UserSolution{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, id), &us)
	return &us, err
}

func UpdateName(id int, name string) error {
	statement := "update user_solutions set name = $1 where id = $2"
	_, err := postgres.GetPool().Exec(postgres.GetCtx(), statement, name, id)
	return err
}

func load(qr pgx.Row, us *UserSolution) error {
	return qr.Scan(&us.Id, &us.UserId, &us.TaskId, &us.TestId, &us.LastModified,
		&us.Language, &us.Name, &us.Public, &us.Code, &us.ExitCode, &us.Output, &us.CompilationTime, &us.RealTime,
		&us.KernelTime, &us.UserTime, &us.MaxRamUsage, &us.BinarySize)
}
