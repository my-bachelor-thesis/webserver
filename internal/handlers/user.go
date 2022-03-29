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
	type loginCredentials struct {
		Username string
		Password string
	}
	var lc loginCredentials
	if err := c.Bind(&lc); err != nil {
		return err
	}

	// get user from DB
	user, err := users.GetByUsername(lc.Username)
	if postgresutil.IsNoRowsInResultErr(err) {
		return c.JSON(http.StatusBadRequest, "username doesn't exist")
	}
	if err != nil {
		return err
	}

	// Throws unauthorized error
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lc.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, "bad password")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user.Id,
		user.IsAdmin,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.GetInstance().JWTSecret))
	if err != nil {
		return err
	}

	jwtCookie := &http.Cookie{}
	jwtCookie.Name = "auth"
	jwtCookie.Value = t
	jwtCookie.Path = "/"
	jwtCookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(jwtCookie)

	return nil
}

func RegisterPost(c echo.Context) error {
	user := &users.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	// TODO: validate data here (https://echo.labstack.com/guide/request/)

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
	return c.JSON(http.StatusOK, fmt.Sprintf("%t",!postgresutil.IsNoRowsInResultErr(err)))
}
