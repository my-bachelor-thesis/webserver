package tokens

import (
	"fmt"
	"webserver/internal/postgres"
)

func getByToken(tx postgres.PoolInterface, token, tableName string) (userId int, err error) {
	statement := fmt.Sprintf("select user_id from %s where token = $1", tableName)
	err = tx.QueryRow(postgres.GetCtx(), statement, token).Scan(&userId)
	return
}

func GetTokenForPasswordResetByToken(tx postgres.PoolInterface, token string) (*TokenForPasswordReset, error) {
	userId, err := getByToken(tx, token, passwordResetTableName)
	return &TokenForPasswordReset{Token: token, UserId: userId}, err
}

func GetTokenForVerificationByToken(tx postgres.PoolInterface, token string) (*TokenForVerification, error) {
	userId, err := getByToken(tx, token, verificationTableName)
	return &TokenForVerification{Token: token, UserId: userId}, err
}
