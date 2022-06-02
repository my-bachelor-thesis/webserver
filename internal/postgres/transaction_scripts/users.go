package transaction_scripts

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"webserver/internal/email_sender"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tokens"
	"webserver/internal/postgres/rdg/users"
	"webserver/pkg/postgresutil"
)

func RegisterUser(user *users.User, token *tokens.TokenForVerification) error {
	conn, tx, err := getConnectionFromPoolAndStartRegularTrans()
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	if err = user.Insert(tx); err != nil {
		return err
	}

	if err := token.Insert(tx); err != nil {
		return err
	}

	return tx.Commit(postgres.GetCtx())
}

func UpdateUserEmail(c echo.Context, email string) (*users.User, string, error) {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	user, err := getUserFromJWTCookie(tx, c)
	if err != nil {
		return nil, "", err
	}

	if !user.Activated {
		return nil, "", errors.New("can't change, your previous email wasn't verified")
	}

	token := email_sender.GenerateToken()

	user.Email = email
	if err := user.UpdateEmailAndDeactivate(tx); err != nil {
		return nil, "", err
	}

	verificationToken := tokens.TokenForVerification{UserId: user.Id, Token: token}
	if err := verificationToken.Insert(tx); err != nil {
		return nil, "", err
	}

	return user, token, tx.Rollback(postgres.GetCtx())
}

func RequestResetUserPassword(email string) (*users.User, string, error) {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	user, err := users.GetByEmail(tx, email)
	if err != nil {
		return nil, "", err
	}

	token := email_sender.GenerateToken()
	resetToken := tokens.TokenForPasswordReset{Token: token, UserId: user.Id}
	if err := resetToken.Insert(tx); err != nil {
		if postgresutil.IsUniqueConstraintErr(err) {
			return nil, "", errors.New("reset email has already been sent")
		}
		return nil, "", err
	}

	return user, token, tx.Rollback(postgres.GetCtx())
}

func ResetUserPassword(token, encryptedPassword string) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	resetToken, err := tokens.GetTokenForPasswordResetByToken(tx, token)
	if postgresutil.IsNoRowsInResultErr(err) {
		return errors.New("invalid token")
	}
	if err != nil {
		return err
	}

	user, err := users.GetById(tx, resetToken.UserId)
	if err != nil {
		return err
	}
	user.Password = encryptedPassword

	if err := user.UpdatePassword(tx); err != nil {
		return err
	}

	if err := resetToken.Delete(tx); err != nil {
		return err
	}

	return tx.Rollback(postgres.GetCtx())
}

func UserEmailVerification(token string) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	verificationToken, err := tokens.GetTokenForVerificationByToken(tx, token)
	if postgresutil.IsNoRowsInResultErr(err) {
		return errors.New("invalid token")
	}
	if err != nil {
		return err
	}

	user, err := users.GetById(tx, verificationToken.UserId)
	if err != nil {
		return err
	}

	user.Activated = true
	if err := user.UpdateActivatedStatus(tx); err != nil {
		return err
	}

	if err := verificationToken.Delete(tx); err != nil {
		return err
	}

	return tx.Rollback(postgres.GetCtx())
}

func PromoteUserToAdmin(username string) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	user, err := users.GetByUsername(tx, username)
	if err != nil {
		return err
	}

	if err := user.PromoteToAdmin(tx); err != nil {
		return err
	}

	return tx.Rollback(postgres.GetCtx())
}

func UpdateUserInfo(c echo.Context, username, lastName, firstName string) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	user, err := getUserFromJWTCookie(tx, c)
	if err != nil {
		return err
	}

	user.Username = username
	user.LastName = lastName
	user.FirstName = firstName
	if err := user.UpdateFirstLastAndUsername(tx); err != nil {
		return err
	}

	return tx.Rollback(postgres.GetCtx())
}

func UpdateUserPassword(c echo.Context, oldPassword, encryptedNewPassword string) error {
	conn, tx, err := getConnectionFromPoolAndStartTrans(pgx.RepeatableRead)
	if err != nil {
		return err
	}
	defer tx.Rollback(postgres.GetCtx())
	defer conn.Release()

	user, err := getUserFromJWTCookie(tx, c)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("bad old password")
	}

	user.Password = encryptedNewPassword
	if err := user.UpdatePassword(tx); err != nil {
		return err
	}

	return tx.Rollback(postgres.GetCtx())
}
