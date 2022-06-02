package users

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetById(tx postgres.PoolInterface, id int) (*User, error) {
	return getBySomething(tx, "id", strconv.Itoa(id))
}

func GetByUsername(tx postgres.PoolInterface, username string) (*User, error) {
	return getBySomething(tx, "username", username)
}

func GetByEmail(tx postgres.PoolInterface, email string) (*User, error) {
	return getBySomething(tx, "email", email)
}

func getBySomething(tx postgres.PoolInterface, fieldName, fieldValue string) (*User, error) {
	statement := fmt.Sprintf("select %s from users where %s = $1", allFields, fieldName)
	user := User{}
	err := load(tx.QueryRow(postgres.GetCtx(), statement, fieldValue), &user)
	if err == nil && user.Id == 0 {
		err = postgresutil.ErrNoRowsInResult
	}
	return &user, err
}

func load(qr pgx.Row, user *User) error {
	return qr.Scan(&user.Id, &user.IsAdmin, &user.FirstName, &user.LastName, &user.Username, &user.Email,
		&user.Password, &user.Activated)
}
