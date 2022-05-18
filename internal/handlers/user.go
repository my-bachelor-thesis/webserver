package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"webserver/internal/config"
	"webserver/internal/postgres/rdg/users"
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
	user, err := users.GetByUsername(loginCredentials.Username)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusBadRequest, "username doesn't exist")
	}
	if err != nil {
		return err
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

	if err = user.Insert(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func IsValidUsername(c echo.Context) error {
	_, err := users.GetByUsername(c.Param("username"))
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
		return c.JSON(http.StatusBadRequest, "bad password")
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
		Email string `json:"email" validate:"required"`
	}
	if err := bindAndValidate(c, &request); err != nil {
		return err
	}

	user, err := getUserFromJWTCookie(c)
	if err != nil {
		return err
	}

	user.Email = request.Email
	if err := user.UpdateEmail(); err != nil {
		return err
	}

	// TODO: insert token and logout
	return nil
}
