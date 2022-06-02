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
	myjwt "webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/tokens"
	"webserver/internal/postgres/rdg/users"
	"webserver/internal/postgres/transaction_scripts"
	"webserver/pkg/postgresutil"
)

func LoginPost(c echo.Context) error {
	var loginCredentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := bindAndValidate(c, &loginCredentials); err != nil {
		return err
	}

	// get user from DB
	user, err := users.GetByUsername(postgres.GetPool(), loginCredentials.Username)
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
	claims := &myjwt.CustomClaims{
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

	token := email_sender.GenerateToken()
	verificationToken := tokens.TokenForVerification{Token: token}

	if err := transaction_scripts.RegisterUser(user, &verificationToken); err != nil {
		return err
	}

	// send verification email
	return email_sender.SendVerificationToken(user.Email, token)
}

func IsValidUsername(c echo.Context) error {
	_, err := users.GetByUsername(postgres.GetPool(), c.Param("username"))
	return c.JSON(http.StatusOK, fmt.Sprintf("%t", !postgresutil.IsNoRowsInResultErr(err)))
}

func IsValidEmail(c echo.Context) error {
	_, err := users.GetByEmail(postgres.GetPool(), c.Param("email"))
	return c.JSON(http.StatusOK, fmt.Sprintf("%t", !postgresutil.IsNoRowsInResultErr(err)))
}

func UpdateUserInfoPost(c echo.Context) error {
	var request struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Username  string `json:"username" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	return transaction_scripts.UpdateUserInfo(c, request.Username, request.LastName, request.FirstName)
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

	encrypted, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = transaction_scripts.UpdateUserPassword(c, request.OldPassword, string(encrypted))
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	return err
}

func UpdateEmailPost(c echo.Context) error {
	var request struct {
		Email string `json:"new_email" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	user, token, err := transaction_scripts.UpdateUserEmail(c, request.Email)
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	if err != nil {
		return err
	}

	return email_sender.SendVerificationToken(user.Email, token)
}

func RequestResetPasswordPost(c echo.Context) error {
	var request struct {
		Email string `json:"registered_email" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	user, token, err := transaction_scripts.RequestResetUserPassword(request.Email)
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	if err != nil {
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

	encrypted, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = transaction_scripts.ResetUserPassword(request.Token, string(encrypted))
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	return err
}

func EmailVerificationPost(c echo.Context) error {
	var request struct {
		Token string `json:"token" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	err := transaction_scripts.UserEmailVerification(request.Token)
	if err, ok := err.(*transaction_scripts.BadRequestError); ok {
		return c.JSON(http.StatusBadRequest, err.Message)
	}
	return err
}

func PromoteToAdminPost(c echo.Context) error {
	var request struct {
		Username string `json:"username_to_promote" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	claims, err := myjwt.GetClaimsFromRequest(c)
	if err != nil {
		return err
	}

	if !claims.IsAdmin {
		return c.JSON(http.StatusForbidden, "not a admin")
	}

	return transaction_scripts.PromoteUserToAdmin(request.Username)
}
