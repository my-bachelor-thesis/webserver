package transaction_scripts

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tests"
)

func UpdateTestName(c echo.Context) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	req, test, err := bindAndFind(tx, c, tests.GetById)
	if err != nil {
		return err
	}

	if test.Final || test.Public {
		return errors.New( "can't update this test")
	}

	test.Name = req.Name
	if err := test.UpdateName(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}
