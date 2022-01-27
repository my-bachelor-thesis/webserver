package tests

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
)

func GetById(id int) (*Test, error) {
	statement := fmt.Sprintf("select %s from tests where id = $1", allFieldReplacedTimestamp)
	test := Test{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, id), &test)
	return &test, err
}

func load(qr pgx.Row, test *Test) error {
	return qr.Scan(&test.Id, &test.LastModified, &test.Final, &test.UserId,
		&test.TaskId, &test.Language, &test.Code)
}
