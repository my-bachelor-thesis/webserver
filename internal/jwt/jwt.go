package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	UserId    int    `json:"user_id"`
	IsAdmin   bool   `json:"is_admin"`
	jwt.StandardClaims
}

func GetClaimsFromRequest(c echo.Context) (*CustomClaims, error) {
	if user, ok := c.Get("user").(*jwt.Token); ok {
		return user.Claims.(*CustomClaims), nil
	}
	return nil, errors.New("couldn't convert to *jwt.Token")
}