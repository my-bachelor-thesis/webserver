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
	Id        int    `json:"id"` // default is 0
	IsAdmin   bool   `json:"is_admin" form:"is_admin"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Username  string `json:"username" form:"username"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
}

func (user *User) Insert() error {
	statement := fmt.Sprintf(`
	insert into users (%s)
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholders(allFieldsWithoutId))
	return postgres.GetPool().QueryRow(postgres.GetCtx(), statement, user.IsAdmin, user.FirstName, user.LastName,
		user.Username, user.Email, user.Password).Scan(&user.Id)
}
