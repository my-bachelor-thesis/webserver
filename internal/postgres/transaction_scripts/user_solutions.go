package transaction_scripts

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/user_solutions"
)

func UpdateUserSolutionNamePost(c echo.Context) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	req, us, err := bindAndFind(tx, c, user_solutions.GetById)
	if err != nil {
		return err
	}

	if us.Public {
		return errors.New("can't update this user solution")
	}

	us.Name = req.Name
	if err := us.UpdateName(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}
