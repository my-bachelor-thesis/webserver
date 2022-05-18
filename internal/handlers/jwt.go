package handlers

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserId    int    `json:"user_id"`
	IsAdmin   bool   `json:"is_admin"`
	jwt.StandardClaims
}

func getClaimsFromRequest(c echo.Context) (*JwtCustomClaims, error) {
	if user, ok := c.Get("user").(*jwt.Token); ok {
		return user.Claims.(*JwtCustomClaims), nil
	}
	return nil, errors.New("couldn't convert to *jwt.Token")
}
