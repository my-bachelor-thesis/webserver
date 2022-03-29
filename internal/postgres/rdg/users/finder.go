package users

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetById(id int) (*User, error) {
	if id == 0 {
		return nil, postgresutil.ErrNoRowsInResult
	}
	statement := fmt.Sprintf("select %s from users where id = $1", allFields)
	user := User{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, id), &user)
	return &user, err
}

func GetByUsername(username string) (*User, error) {
	statement := fmt.Sprintf("select %s from users where username = $1", allFields)
	user := User{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, username), &user)
	if err == nil && user.Id == 0 {
		err = postgresutil.ErrNoRowsInResult
	}
	return &user, err
}

func load(qr pgx.Row, user *User) error {
	return qr.Scan(&user.Id, &user.IsAdmin, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password)
}
