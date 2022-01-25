package users

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "is_admin, first_name, last_name, username, email, password"
	allFields          = "id, " + allFieldsWithoutId
)

type User struct {
	Id        int    `json:"id"`
	IsAdmin   bool   `json:"is_admin"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func Insert(user *User) error {
	statement := fmt.Sprintf(`
	insert into users %s
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholder(allFieldsWithoutId))
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, user.IsAdmin, user.FirstName, user.LastName, user.Email, user.Password).Scan(&user.Id)
}
