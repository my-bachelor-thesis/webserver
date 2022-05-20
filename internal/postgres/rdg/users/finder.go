package users

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

func GetById(id int) (*User, error) {
	return getBySomething("id", strconv.Itoa(id))
}

func GetByUsername(username string) (*User, error) {
	return getBySomething("username", username)
}

func GetByEmail(email string) (*User, error) {
	return getBySomething("email", email)
}

func getBySomething(fieldName, fieldValue string) (*User, error) {
	statement := fmt.Sprintf("select %s from users where %s = $1", allFields, fieldName)
	user := User{}
	err := load(postgres.GetPool().QueryRow(postgres.GetCtx(), statement, fieldValue), &user)
	if err == nil && (user.Id == 0 || !user.Activated) {
		err = postgresutil.ErrNoRowsInResult
	}
	return &user, err
}

func load(qr pgx.Row, user *User) error {
	return qr.Scan(&user.Id, &user.IsAdmin, &user.FirstName, &user.LastName, &user.Username, &user.Email,
		&user.Password, &user.Activated)
}
