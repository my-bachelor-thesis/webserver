package tokens

import (
	"fmt"
	"webserver/internal/postgres"
)

func getByToken(token, tableName string) (userId int, err error) {
	statement := fmt.Sprintf("select user_id from %s where token = $1", tableName)
	err = postgres.GetPool().QueryRow(postgres.GetCtx(), statement, token).Scan(&userId)
	return
}

//func getByUserId(userId int, tableName string) (token string, err error) {
//	statement := fmt.Sprintf("select token from %s where user_id = $1", tableName)
//	err = postgres.GetPool().QueryRow(postgres.GetCtx(), statement, userId).Scan(&token)
//	return
//}

func GetTokenForPasswordResetByToken(token string) (*TokenForPasswordReset, error) {
	userId, err := getByToken(token, passwordResetTableName)
	return &TokenForPasswordReset{Token: token, UserId: userId}, err
}

func GetTokenForVerificationByToken(token string) (*TokenForVerification, error) {
	userId, err := getByToken(token, verificationTableName)
	return &TokenForVerification{Token: token, UserId: userId}, err
}

//func GetTokenForPasswordResetByUserId(userId int) (*TokenForPasswordReset, error) {
//	token, err := getByUserId(userId, passwordResetTableName)
//	return &TokenForPasswordReset{Token: token, UserId: userId}, err
//}
//
//func GetTokenForVerificationByUserId(userId int) (*TokenForVerification, error) {
//	token, err := getByUserId(userId, passwordResetTableName)
//	return &TokenForVerification{Token: token, UserId: userId}, err
//}
