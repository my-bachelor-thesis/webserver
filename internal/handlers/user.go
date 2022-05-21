package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"webserver/internal/config"
	"webserver/internal/email_sender"
	"webserver/internal/postgres/rdg/tokens"
	"webserver/internal/postgres/rdg/users"
	"webserver/pkg/postgresutil"
)

// Tested
func LoginPost(c echo.Context) error {
	var loginCredentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := bindAndValidate(c, &loginCredentials); err != nil {
		return err
	}

	// get user from DB
	user, err := users.GetByUsername(loginCredentials.Username)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusBadRequest, "username doesn't exist")
	}
	if err != nil {
		return err
	}
	if !user.Activated {
		return c.JSON(http.StatusMethodNotAllowed, "user with unconfirmed email")
	}

	// Throws unauthorized error
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredentials.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, "bad password")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user.Id,
		user.IsAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and
	t, err := token.SignedString([]byte(config.GetInstance().JWTSecret))
	if err != nil {
		return err
	}

	// set cookies
	cookieTimeout := time.Now().Add(24 * time.Hour)
	c.SetCookie(&http.Cookie{Name: "auth", Value: t, Path: "/", Expires: cookieTimeout})
	c.SetCookie(&http.Cookie{Name: "username", Value: user.Username, Path: "/", Expires: cookieTimeout})
	c.SetCookie(&http.Cookie{Name: "first_name", Value: user.FirstName, Path: "/", Expires: cookieTimeout})
	c.SetCookie(&http.Cookie{Name: "last_name", Value: user.LastName, Path: "/", Expires: cookieTimeout})
	c.SetCookie(&http.Cookie{Name: "email", Value: user.Email, Path: "/", Expires: cookieTimeout})

	return nil
}

// Tested
func RegisterPost(c echo.Context) error {
	user := &users.User{}
	if err := bindAndValidate(c, user); err != nil {
		return err
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)
	user.Activated = false

	// start transaction
	//tx, err := postgres.GetPool().Begin(postgres.GetCtx())
	//if err != nil {
	//	return err
	//}
	//defer tx.Rollback(postgres.GetCtx())

	if err = user.Insert(); err != nil {
		return err
	}

	token := email_sender.GenerateToken()
	verificationToken := tokens.TokenForVerification{UserId: user.Id, Token: token}
	if err := verificationToken.Insert(); err != nil {
		return err
	}

	// send verification email
	return email_sender.SendVerificationToken(user.Email, token)
}

func IsValidUsername(c echo.Context) error {
	_, err := users.GetByUsername(c.Param("username"))
	return c.JSON(http.StatusOK, fmt.Sprintf("%t", !postgresutil.IsNoRowsInResultErr(err)))
}

func IsValidEmail(c echo.Context) error {
	_, err := users.GetByEmail(c.Param("email"))
	return c.JSON(http.StatusOK, fmt.Sprintf("%t", !postgresutil.IsNoRowsInResultErr(err)))
}

// odskusane
func UpdateUserInfoPost(c echo.Context) error {
	var request struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	user, err := getUserFromJWTCookie(c)
	if err != nil {
		return err
	}

	user.Username = request.Username
	user.LastName = request.LastName
	user.FirstName = request.FirstName
	return user.UpdateFirstLastAndUsername()
}

func UpdatePasswordPost(c echo.Context) error {
	var request struct {
		OldPassword         string `json:"old_password" validate:"required"`
		NewPassword         string `json:"new_password" validate:"required"`
		NewPasswordRepeated string `json:"repeat_password" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	if request.NewPasswordRepeated != request.NewPassword {
		return c.JSON(http.StatusBadRequest, "passwords don't match")
	}

	user, err := getUserFromJWTCookie(c)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		return c.JSON(http.StatusBadRequest, "bad old password")
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)
	return user.UpdatePassword()
}

func UpdateEmailPost(c echo.Context) error {
	var request struct {
		Email string `json:"new_email" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	user, err := getUserFromJWTCookie(c)
	if err != nil {
		return err
	}

	if !user.Activated {
		return c.JSON(http.StatusMethodNotAllowed, "can't change, your previous email wasn't verified")
	}

	user.Email = request.Email
	if err := user.UpdateEmailAndDeactivate(); err != nil {
		return err
	}

	token := email_sender.GenerateToken()
	verificationToken := tokens.TokenForVerification{UserId: user.Id, Token: token}
	if err := verificationToken.Insert(); err != nil {
		return err
	}

	return email_sender.SendVerificationToken(user.Email, token)
}

// Tested
func RequestResetPasswordPost(c echo.Context) error {
	var request struct {
		Email string `json:"registered_email" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	user, err := users.GetByEmail(request.Email)
	if err != nil {
		return err
	}

	token := email_sender.GenerateToken()
	resetToken := tokens.TokenForPasswordReset{Token: token, UserId: user.Id}
	if err := resetToken.Insert(); err != nil {
		if postgresutil.IsUniqueConstraintErr(err) {
			return c.JSON(http.StatusMethodNotAllowed, "reset email has already been sent")
		}
		return err
	}

	return email_sender.SendResetToken(user.Email, token)
}

func ResetPasswordPost(c echo.Context) error {
	var request struct {
		Token            string `json:"token" validate:"required"`
		Password         string `json:"new_password" validate:"required"`
		RepeatedPassword string `json:"repeat_password" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	if request.RepeatedPassword != request.Password {
		return c.JSON(http.StatusBadRequest, "passwords don't match")
	}

	resetToken, err := tokens.GetTokenForPasswordResetByToken(request.Token)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusMethodNotAllowed, "invalid token")
	}
	if err != nil {
		return err
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := users.GetById(resetToken.UserId)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)

	if err := user.UpdatePassword(); err != nil {
		return err
	}

	return resetToken.Delete()
}

func EmailVerificationPost(c echo.Context) error {
	var request struct {
		Token string `json:"token" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	verificationToken, err := tokens.GetTokenForVerificationByToken(request.Token)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusMethodNotAllowed, "invalid token")
	}
	if err != nil {
		return err
	}

	user, err := users.GetById(verificationToken.UserId)
	if err != nil {
		return err
	}

	// TODO: in trans

	user.Activated = true
	if err := user.UpdateActivatedStatus(); err != nil {
		return err
	}

	return verificationToken.Delete()
}
