package users

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFieldsWithoutId = "is_admin, first_name, last_name, username, email, password, activated"
	allFields          = "id, " + allFieldsWithoutId
)

type User struct {
	Id        int    `json:"id"` // default is 0
	IsAdmin   bool   `json:"is_admin" form:"is_admin"`
	FirstName string `json:"first_name" form:"first_name" validate:"required"`
	LastName  string `json:"last_name" form:"last_name" validate:"required"`
	Username  string `json:"username" form:"username" validate:"required"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required"`
	Activated bool   `json:"activated" form:"activated"`
}

func (user *User) Insert(tx postgres.PoolInterface) error {
	statement := fmt.Sprintf(`
	insert into users (%s)
	values (%s)
	returning id`, allFieldsWithoutId, postgresutil.GeneratePlaceholders(allFieldsWithoutId))
	return tx.QueryRow(postgres.GetCtx(), statement, user.IsAdmin, user.FirstName, user.LastName,
		user.Username, user.Email, user.Password, user.Activated).Scan(&user.Id)
}

func (user *User) UpdateFirstLastAndUsername(tx postgres.PoolInterface) error {
	statement := "update users set first_name = $1, last_name = $2, username = $3 where id = $4"
	_, err := tx.Exec(postgres.GetCtx(), statement, user.FirstName, user.LastName, user.Username, user.Id)
	return err
}

func (user *User) UpdateEmailAndDeactivate(tx postgres.PoolInterface) error {
	statement := "update users set email = $1, activated = false where id = $2"
	_, err := tx.Exec(postgres.GetCtx(), statement, user.Email, user.Id)
	user.Activated = false
	return err
}

func (user *User) UpdatePassword(tx postgres.PoolInterface) error {
	statement := "update users set password = $1 where id = $2"
	_, err := tx.Exec(postgres.GetCtx(), statement, user.Password, user.Id)
	return err
}

func (user *User) UpdateActivatedStatus(tx postgres.PoolInterface) error {
	statement := "update users set activated = $1 where id = $2"
	_, err := tx.Exec(postgres.GetCtx(), statement, user.Activated, user.Id)
	return err
}

func (user *User) PromoteToAdmin(tx postgres.PoolInterface) error {
	statement := "update users set is_admin = true where id = $1"
	_, err := tx.Exec(postgres.GetCtx(), statement, user.Id)
	return err
}
