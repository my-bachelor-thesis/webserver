package user_solutions

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetById(tx postgres.PoolInterface, id int) (*UserSolution, error) {
	statement := fmt.Sprintf("select %s from %s where id = $1", allFieldReplacedTimestamp, tableName)
	us := UserSolution{}
	err := load(tx.QueryRow(postgres.GetCtx(), statement, id), &us)
	return &us, err
}

func load(qr pgx.Row, us *UserSolution) error {
	return qr.Scan(&us.Id, &us.UserId, &us.TaskId, &us.LastModified,
		&us.Language, &us.Name, &us.Public, &us.HideInStatistic, &us.Code)
}
