package transaction_scripts

import (
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/user_solutions"
	"webserver/internal/postgres/rdg/users"
)

func UpdateUserSolutionName(c echo.Context) error {
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
		return NewBadRequestError("can't update this user solution")
	}

	us.Name = req.Name
	if err := us.UpdateName(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func DeleteUserSolution(solutionId, adminId int) (admin, user *users.User, us *user_solutions.UserSolution, err error) {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return nil, nil, nil, err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	us, err = user_solutions.GetById(tx, solutionId)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := us.HideFromStatistic(tx); err != nil {
		return nil, nil, nil, err
	}

	user, err = users.GetById(tx, us.UserId)
	if err != nil {
		return nil, nil, nil, err
	}

	admin, err = users.GetById(tx, adminId)
	if err != nil {
		return nil, nil, nil, err
	}

	return admin, user, us, tx.Commit(postgres.GetCtx())
}
