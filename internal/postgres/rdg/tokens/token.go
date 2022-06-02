package tokens

import (
	"fmt"
	"webserver/internal/postgres"
	"webserver/pkg/postgresutil"
)

const (
	allFields              = "user_id, token"
	passwordResetTableName = "tokens_for_password_reset"
	verificationTableName  = "tokens_for_verification"
)

var (
	allFieldsPlaceholders = postgresutil.GeneratePlaceholders(allFields)
)

type TokenForPasswordReset struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type TokenForVerification struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

func insertFunc(tx postgres.PoolInterface, tableName, token string, userId int) error {
	statement := fmt.Sprintf(`
	insert into %s (%s)
	values (%s)`, tableName, allFields, allFieldsPlaceholders)
	_, err := tx.Exec(postgres.GetCtx(), statement, userId, token)
	return err
}

func (tfpr *TokenForPasswordReset) Insert(tx postgres.PoolInterface) error {
	return insertFunc(tx, passwordResetTableName, tfpr.Token, tfpr.UserId)
}

func (tfr *TokenForVerification) Insert(tx postgres.PoolInterface) error {
	return insertFunc(tx, verificationTableName, tfr.Token, tfr.UserId)
}

func deleteFunc(tx postgres.PoolInterface, tableName string, userId int) error {
	statement := fmt.Sprintf("delete from %s where user_id = $1", tableName)
	_, err := tx.Exec(postgres.GetCtx(), statement, userId)
	return err
}

func (tfpr *TokenForPasswordReset) Delete(tx postgres.PoolInterface) error {
	return deleteFunc(tx, passwordResetTableName, tfpr.UserId)
}

func (tfr *TokenForVerification) Delete(tx postgres.PoolInterface) error {
	return deleteFunc(tx, verificationTableName, tfr.UserId)
}
