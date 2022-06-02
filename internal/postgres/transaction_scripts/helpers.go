package transaction_scripts

import (
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/users"
)

func getConnectionFromPoolAndStartTrans(isolationLevel pgx.TxIsoLevel) (*pgxpool.Conn, pgx.Tx, error) {
	conn, err := postgres.GetPool().Acquire(postgres.GetCtx())
	if err != nil {
		return nil, nil, err
	}

	tx, err := conn.BeginTx(postgres.GetCtx(), pgx.TxOptions{IsoLevel: isolationLevel})
	return conn, tx, err
}

func getConnectionFromPoolAndStartRegularTrans() (*pgxpool.Conn, pgx.Tx, error) {
	conn, err := postgres.GetPool().Acquire(postgres.GetCtx())
	if err != nil {
		return nil, nil, err
	}

	tx, err := conn.Begin(postgres.GetCtx())
	return conn, tx, err
}

func getUserFromJWTCookie(tx postgres.PoolInterface, c echo.Context) (*users.User, error) {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return nil, err
	}
	return users.GetById(tx, claims.UserId)
}

type BadRequestError struct {
	Message string
}

func (b *BadRequestError) Error() string {
	return b.Message
}

func NewBadRequestError(msg string) *BadRequestError {
	return &BadRequestError{Message: msg}
}
